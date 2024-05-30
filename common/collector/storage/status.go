package storage

import (
	"context"
)

// PutAllRequestWithTTL is a data structure used to pass a data to a storage
// implementation (DAO).
// Storage implementation should parse this struct to extract necessary fields.
type PutAllRequestWithTTL struct {
	Key     string
	Value   []byte
	TTL     int64
	LeaseID int64
}

// PutOpOptionWithTTL is an opts used in PutGetterWithOpts interface's
// PutForceWithOpts method.
type PutOpOptionWithTTL struct {
	// LeaseID is not a TTL, instead it is an ID of a database key lease.
	// Database key lease is kind of a ticket, that is valid for TTL period,
	// once TTL period is over, ticket is not valid.
	LeaseID string

	// TTL is an additional parameter that could be used to set brand-new
	// database key lease with an expiration time of TTL.
	// TTL unit is a second.
	TTL string
}

// PutGetterWithOpts is an extension of the fundamental PutGetter
// interface. Implementing these methods is completely optional,
// because features like Lease (value with TTL) may be enabled in one
// storage implementation, but be absent in another.
type PutGetterWithOpts interface {
	// GetWithLease should get a key from a storage along with a
	// TTL value assigned to that key (0 if not assigned)
	GetWithLease(ctx context.Context, key string) ([]byte, string, error)

	// PutForceWithOpts not only puts value into key without any
	// conditional checks, e.g. "insert only if absent, etc.),
	// but also allows to pass custom options opts to the underlying
	// implementation, which may parse it according to its own rules.
	//
	// As an example, you can pass TTL to the opts.
	PutForceWithOpts(ctx context.Context, key string, value []byte, opts interface{}) error

	// Delete a key from a storage permanently.
	Delete(ctx context.Context, key string) error
}

// SessionStatusRepository is an interface that all Status services should
// implement!
//
// To recall, Status service is a service that just reports
// something. That something depends on the Status server implementation,
// e.g. some Status server can report SessionID lifetime, others can
// report OS used in a system, etc.
type SessionStatusRepository interface {
	PutGetterWithOpts
}

// SessionStatusRepository is a factory method to return a
// session status interface for the repository interaction.
func (c *Client) SessionStatusRepository() SessionStatusRepository {
	return c.prot.(SessionStatusRepository)
}
