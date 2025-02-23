/*
The verification service serves stored votes and qualifying properties for
requested vote identifiers.
*/
package main

import (
	"os"
	"strings"
	"time"

	"ivxv.ee/common/collector/command"
	"ivxv.ee/common/collector/command/exit"
	"ivxv.ee/common/collector/conf"
	"ivxv.ee/common/collector/container"
	"ivxv.ee/common/collector/errors"
	"ivxv.ee/common/collector/log"
	"ivxv.ee/common/collector/q11n"
	"ivxv.ee/common/collector/server"
	"ivxv.ee/common/collector/status/client"
	status "ivxv.ee/common/collector/status/client/rpc"
	"ivxv.ee/common/collector/storage"
	internal "ivxv.ee/verification/internal/sessionstatus/rpc"
	//ivxv:modules common/collector/container
	//ivxv:modules common/collector/storage
)

// RPC is the handler for verification service calls.
type RPC struct {
	status               client.Verifier
	election             *conf.Election
	qps                  []q11n.Protocol // Qualifying properties expected for each vote.
	storage              *storage.Client
	foreignCode          string // Administrative unit code for foreign voters.
	predefinedDistrictID string
}

// Args are the arguments provided to a call of RPC.Verify.
type Args struct {
	server.Header
	VoteID []byte `size:"16"` // Vote identifier of the vote to return.
}

// Response is the response returned by RPC.Verify.
type Response struct {
	server.Header
	Type          container.Type  // Type of the stored vote container.
	Vote          []byte          // The stored vote.
	Qualification q11n.Properties // Qualifying properties for the vote.
	ChoicesList   []byte          // Choices list
}

// Verify is the remote procedure call performed by clients to retrieve votes
// from the collector for verification.
func (r *RPC) Verify(args Args, resp *Response) error {
	log.Log(args.Ctx, VerifyReq{VoteID: args.VoteID})
	now := time.Now()

	// Build up VerifyReq for session status service
	verifyReq := status.NewVerifyReqBuilder().
		WithServiceMethod(internal.Verify).
		WithRequest(args.Header).
		Build()

	// SessionID security check
	ok, err := r.status.Verify(&verifyReq)
	if err != nil {
		log.Error(args.Ctx, VerifySessionIDError{Err: err})
		return server.ErrBadRequest
	}
	if !ok {
		log.Error(args.Ctx, VerifyUpdateSessionIDError{})
		return server.ErrBadRequest
	}

	// If the verification count changed between retrieving it and
	// attempting to update, then that means that there was another
	// concurrent verification session for this ID that got there first. If
	// this happens, then try from the beginning, but still use the time
	// the request originally came in.
	for {
		// Get the verification statistics for the vote ID.
		count, at, latest, choicesList, err := r.storage.GetVerificationStats(
			args.Ctx, args.VoteID, r.foreignCode, r.predefinedDistrictID)
		if err != nil {
			if errors.CausedBy(err, new(storage.NotExistError)) != nil {
				log.Error(args.Ctx, BadVoteIDError{Err: err})
				return server.ErrBadRequest
			}
			log.Error(args.Ctx, GetVerificationStatsError{Err: log.Alert(err)})
			return server.ErrInternal
		}

		// Check that we have not reached the count limit. We want this
		// to look exactly like a bad vote identifier was submitted, so
		// return ErrBadRequest instead of a more descriptive error.
		if limit := r.election.Verification.Count; limit > 0 && count >= limit {
			log.Error(args.Ctx, VerificationCountError{Count: count})
			return server.ErrBadRequest
		}

		// Check that we are inside the time limit.
		if limit := r.election.Verification.Minutes; limit > 0 &&
			now.After(at.Add(time.Duration(limit)*time.Minute)) {

			log.Error(args.Ctx, VerificationTimeError{At: at})
			return server.ErrBadRequest
		}

		// Check if this is the latest vote.
		if r.election.Verification.LatestOnly && !latest {
			log.Error(args.Ctx, VerificationNotLatestError{})
			return server.ErrBadRequest
		}

		// Retrieve the stored vote tied to the requested vote
		// identifier and increase the verification count, as long as
		// it has not changed.
		vote, err := r.storage.GetVerification(args.Ctx, args.VoteID, count, r.qps...)
		if err != nil {
			if errors.CausedBy(err, new(storage.UnexpectedValueError)) != nil {
				log.Log(args.Ctx, ConcurrentVerificationWarning{Err: err})
				continue
			}
			log.Error(args.Ctx, GetVerificationError{Err: log.Alert(err)})
			return server.ErrInternal
		}

		resp.Type = vote.VoteType
		resp.Vote = vote.Vote
		resp.Qualification = vote.Qualification
		resp.ChoicesList = choicesList

		// Convert vote.Qualification:map[q11n.Protocol][]byte to
		// logq11n:map[q11n.Protocol]log.Sensitive for logging.
		logq11n := make(map[q11n.Protocol]log.Sensitive)
		for p, b := range vote.Qualification {
			logq11n[p] = b
		}
		log.Log(args.Ctx, VerifyResp{
			Type:          resp.Type,
			Vote:          log.Sensitive(resp.Vote),
			Qualification: logq11n,
		})
		return nil
	}
}

func main() {
	// Call verifymain in a separate function so that it can set up defers
	// and have them trigger before returning with a non-zero exit code.
	os.Exit(verifymain())
}

func verifymain() (code int) {
	c := command.New("ivxv-verification", "")
	defer func() {
		code = c.Cleanup(code)
	}()

	// Configure session status client
	statusClient, errCode := internal.NewClient(c)
	if statusClient == nil || errCode != 0 {
		return errCode
	}

	// Create new RPC instance with the election configuration and storage
	// client.
	rpc := &RPC{election: c.Conf.Election, storage: c.Storage, status: statusClient}

	var start, stop time.Time
	var err error
	if c.Conf.Election != nil {
		// Check election configuration time values.
		if start, err = c.Conf.Election.ServiceStartTime(); err != nil {
			return c.Error(exit.Config, StartTimeError{Err: err},
				"bad service start time:", err)
		}

		if stop, err = c.Conf.Election.VerificationStopTime(); err != nil {
			return c.Error(exit.Config, StopTimeError{Err: err},
				"bad service stop time:", err)
		}

		// Initialize the list of qualifying properties which must
		// exist for each vote.
		for _, q := range c.Conf.Election.Qualification {
			rpc.qps = append(rpc.qps, q.Protocol)
		}

		// If not an empty string then ignore voterlist, use predefined district ID for all voters
		rpc.predefinedDistrictID = strings.TrimSpace(c.Conf.Election.IgnoreVoterList)

		// Get foreign code (voterforeignehak) from election config
		rpc.foreignCode = strings.TrimSpace(c.Conf.Election.VoterForeignEHAKDefault())

		// No authentication is used during verification.
	}

	var s *server.S
	if c.Conf.Technical != nil {
		// Configure a new server with the service instance
		// configuration and the RPC handler instance.
		cert, key := conf.TLS(conf.Sensitive(c.Service.ID))
		if s, err = server.New(&server.Conf{
			CertPath: cert,
			KeyPath:  key,
			Address:  c.Service.Address,
			End:      stop,
			Filter:   &c.Conf.Technical.Filter,
			Version:  &c.Conf.Version,
		}, rpc); err != nil {
			return c.Error(exit.Config, ServerConfError{Err: err},
				"failed to configure server:", err)
		}
	}

	// Start listening for incoming connections during the voting period.
	if c.Until >= command.Execute {
		if err = s.ServeAt(c.Ctx, start); err != nil {
			return c.Error(exit.Unavailable, ServeError{Err: err},
				"failed to serve verification service:", err)
		}
	}
	return exit.OK
}
