package ee.ivxv.key.util;

import ee.ivxv.common.model.CandidateList;
import ee.ivxv.common.model.DistrictList;
import ee.ivxv.common.model.Proof;
import ee.ivxv.common.service.console.Progress;
import ee.ivxv.common.service.report.Reporter;
import ee.ivxv.common.util.I18nConsole;
import ee.ivxv.common.util.Json;
import ee.ivxv.common.util.Util;
import ee.ivxv.key.model.Invalid;
import ee.ivxv.key.model.PlainBallotBox;
import ee.ivxv.key.model.Tally;
import ee.ivxv.key.model.Vote;
import ee.ivxv.key.protocol.SigningProtocol;

import java.io.ByteArrayOutputStream;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.HashMap;
import java.util.Map;
import java.util.concurrent.BlockingQueue;
import java.util.concurrent.Callable;
import java.util.concurrent.LinkedBlockingQueue;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

/**
 * Helper class for processing election result.
 */
public class ElectionResult {
    static final Logger log = LoggerFactory.getLogger(ElectionResult.class);

    private static final String INVALID_VOTE_PATH_TMPL = "invalid";
    private static final String PROOF_PATH_TMPL = "proof";
    private static final String PROOF_INVALID_PATH_TMPL = "proof-invalid";
    private static final String TALLY_SUFFIX = ".tally";
    private static final String PLAIN_SUFFIX = ".plain";
    private static final String TALLY_SIG_SUFFIX = TALLY_SUFFIX + ".signature";
    private static final String PLAIN_SIG_SUFFIX = PLAIN_SUFFIX + ".signature";

    private final String electionName;
    private final CandidateList candidates;
    private final DistrictList districts;
    private final boolean withProof;
    private final boolean proveInvalid;
    private final BlockingQueue<Object> votes = new LinkedBlockingQueue<>();
    private final Map<String, Tally> tallySet;
    private final Map<String, PlainBallotBox> plainSet;
    private final Proof proof;
    private final Proof proofInvalid;
    private final Invalid invalid;

    /**
     * Initialize using values.
     *
     * @param electionName Election identifier
     * @param candidates   Candidates list
     * @param districts    Districts list
     * @param withProof    Boolean indicating if encrypted ballots should be decrypted with proof of
     *                     correct decryption.
     * @param proveInvalid Boolean indicating whether proofs of correct decryption should be output
     *                     for decrypted ballots deemed invalid (e.g., due to incorrect padding).
     */
    public ElectionResult(String electionName, CandidateList candidates, DistrictList districts,
                          boolean withProof, boolean proveInvalid) {
        this.electionName = electionName;
        this.candidates = candidates;
        this.districts = districts;
        this.withProof = withProof;
        this.proveInvalid = withProof && proveInvalid;
        this.tallySet = new HashMap<>();
        this.plainSet = new HashMap<>();
        this.proof = withProof ? new Proof(electionName) : null;
        this.proofInvalid = this.proveInvalid ? new Proof(electionName) : null;
        this.invalid = new Invalid(electionName);
    }

    /**
     * Get the worker for computing the tally.
     * <p>
     * The worker works in parallel to decryptions and continuously computes the tally. It separates
     * votes which are invalid or are given for invalid candidates
     *
     * @param voteCount
     * @param console
     * @param reporter
     * @return
     */
    public ResultWorker getResultWorker(int voteCount, I18nConsole console, Reporter reporter) {
        return new ResultWorker(voteCount, console, reporter);
    }

    /**
     * Set the end of incoming votes.
     * <p>
     * After the last vote has been added, set the end marker indicating the worker thread to stop
     * waiting for incoming votes.
     */
    public void setEot() {
        votes.add(Util.EOT);
    }

    /**
     * Add a vote for worker to process.
     *
     * @param vote
     */
    public void addVote(Vote vote) {
        votes.add(vote);
    }

    /**
     * Output tally and the corresponding signature.
     *
     * @param outDir Output directory to store the tally and signature.
     * @param signer Protocol for constructing the signature for the tally.
     * @throws Exception When writing or communication with card tokens fail.
     */
    public void outputTally(Path outDir, SigningProtocol signer) throws Exception {
        for (Map.Entry<String, Tally> tally : tallySet.entrySet()) {
            ByteArrayOutputStream in = new ByteArrayOutputStream();
            Json.write(tally.getValue(), in);
            byte[] signature = signer.sign(in.toByteArray());
            Files.write(outDir.resolve(Paths.get(tally.getKey() + TALLY_SUFFIX)), in.toByteArray());
            Files.write(outDir.resolve(Paths.get(tally.getKey() + TALLY_SIG_SUFFIX)), signature);
        }
    }

    /**
     * Output the decrypted ballot box and the corresponding signature.
     *
     * @param outDir Output directory to store the decrypted BB and signature.
     * @param signer Protocol for constructing the signature for the decrypted BB.
     * @throws Exception When writing or communication with card tokens fail.
     */
    public void outputPlainBB(Path outDir, SigningProtocol signer) throws Exception {
        for (Map.Entry<String, PlainBallotBox> plain: plainSet.entrySet()) {
            ByteArrayOutputStream in = new ByteArrayOutputStream();
            Json.write(plain.getValue(), in);
            byte[] signature = signer.sign(in.toByteArray());
            Files.write(outDir.resolve(Paths.get(plain.getKey() + PLAIN_SUFFIX)), in.toByteArray());
            Files.write(outDir.resolve(Paths.get(plain.getKey() + PLAIN_SIG_SUFFIX)), signature);
        }
    }

    /**
     * Output proofs of correct decryptions.
     *
     * @param outDir Output directory to store the proofs file.
     * @throws Exception When writing the proofs file fails.
     */
    public void outputProof(Path outDir) throws Exception {
        if (!withProof) {
            return;
        }
        Json.write(proof, outDir.resolve(Util.prefixedPath(electionName, PROOF_PATH_TMPL)));
    }

    /**
     * Output proofs of correct decryption for invalid votes.
     *
     * @param outDir Output directory to store the proofs file.
     * @throws Exception When writing the proofs file fails.
     */
    public void outputInvalidProof(Path outDir) throws Exception {
        if (!proveInvalid) {
            return;
        }
        Json.write(proofInvalid, outDir.resolve(Util.prefixedPath(electionName, PROOF_INVALID_PATH_TMPL)));
    }

    /**
     * Output invalid votes.
     *
     * @param outDir Output directory to store the invalid votes.
     * @throws Exception When writing the file fails.
     */
    public void outputInvalid(Path outDir) throws Exception {
        Json.write(invalid,
                outDir.resolve(Util.prefixedPath(electionName, INVALID_VOTE_PATH_TMPL)));
    }

    private class ResultWorker implements Callable<Void> {
        private final int voteCount;
        private I18nConsole console;

        public ResultWorker(int voteCount, I18nConsole console, Reporter reporter) {
            this.voteCount = voteCount;
            this.console = console;
        }

        @Override
        public Void call() throws Exception {
            Progress progress = console.startProgress(voteCount);
            Object obj;
            while ((obj = votes.take()) != Util.EOT) {
                progress.increase(1);
                Vote vote;
                if (obj instanceof Vote) {
                    vote = (Vote) obj;
                } else {
                    throw new IllegalArgumentException(
                            "Unexpected decryption result type: " + obj.getClass());
                }
                String choice = Tally.INVALID_VOTE_ID;
                if (vote.getProof() != null) {
                    // the message has been decrypted
                    if (isValidChoice(vote)) {
                        // the message contains a valid choice string.
                        choice = getCandidateNumber(vote);
                    } else {
                        // the message was correctly decrypted, but this does not represent a valid
                        // choice string
                        log.warn("Choice is not correctly encoded: invalid vote");
                    }
                } else {
                    // the message has not been decrypted. This can be
                    // caused by incorrect group elements etc. It was not decrypted
                    // to prevent any leaks about the key.
                    log.warn("Ciphertext not correctly encoded: invalid vote");
                }
                if (withProof && !choice.equals(Tally.INVALID_VOTE_ID)) {
                    // output proof of correct decryption only if it is requested and the choice
                    // string is valid. As the decryption proof also contains the whole encrypted
                    // message, then this may leak identifiable information
                    proof.addProof(vote.getProof());
                }
                if (choice.equals(Tally.INVALID_VOTE_ID)) {
                    invalid.getInvalid().add(vote);
                    if (proveInvalid) {
                        // Only if so configured, output the proof of correct decryption also for
                        // votes with invalid choice strings. This includes votes with invalid padding,
                        // since then the choice string is meaningless (it cannot be reliably extracted).
                        proofInvalid.addProof(vote.getProof());
                    }
                }
                addVoteToTally(vote, choice);
            }
            progress.finish();
            return null;
        }

        private String getCandidateNumber(Vote vote) {
            String canidateNumber = Tally.INVALID_VOTE_ID;
            try {
                String message = vote.getProof().getDecrypted().stripPadding().getUTF8DecodedMessage();
                canidateNumber = message.split(Util.UNIT_SEPARATOR)[0];
            } catch (IllegalArgumentException ignored) {
            }
            return canidateNumber;
        }

        private boolean isValidChoice(Vote vote) {
            String voteStr;
            try {
                voteStr = vote.getProof().getDecrypted().stripPadding().getUTF8DecodedMessage();
            } catch (IllegalArgumentException ignored) {
                return false;
            }
            String[] voteParts = voteStr.split(Util.UNIT_SEPARATOR, 3);
            if (voteParts.length != 3) {
                return false;
            }
            Map<String, Map<String, Map<String, String>>> ds = candidates.getCandidates();
            if (!ds.containsKey(vote.getDistrict())) {
                return false;
            }
            Map<String, Map<String, String>> ps = ds.get(vote.getDistrict());
            if (!ps.containsKey(voteParts[1])) {
                return false;
            }
            Map<String, String> ids = ps.get(voteParts[1]);
            if (!ids.containsKey(voteParts[0])) {
                return false;
            }
            String name = ids.get(voteParts[0]);
            if (!name.equals(voteParts[2])) {
                return false;
            }
            return true;
        }

        private void addVoteToTally(Vote vote, String choice) {
            tallySet.computeIfAbsent(vote.getQuestion(),
                            q -> new Tally(electionName, candidates, districts)).getByParish()
                    .get(vote.getDistrict()).get(vote.getStation())
                    .compute(choice, (c, count) -> count + 1);
            plainSet.computeIfAbsent(vote.getQuestion(),
                    q -> new PlainBallotBox(electionName, candidates, districts)).getByParish()
                    .get(vote.getDistrict()).get(vote.getStation())
                    .add(choice);
        }

    }
}
