package ee.ivxv.audit.tools;

import ee.ivxv.audit.AuditContext;
import ee.ivxv.audit.Msg;
import ee.ivxv.audit.tools.IntegrityTool.IntegrityArgs;
import ee.ivxv.audit.util.RawBallotWithDigest;
import ee.ivxv.common.cli.Arg;
import ee.ivxv.common.cli.Args;
import ee.ivxv.common.cli.Tool;
import ee.ivxv.common.crypto.hash.HashType;
import ee.ivxv.common.model.AnonymousBallotBox;
import ee.ivxv.common.service.bbox.impl.FileSource;
import ee.ivxv.common.service.bbox.impl.ZipSource;
import ee.ivxv.common.service.bbox.impl.ZipSourceRaw;
import ee.ivxv.common.service.console.Progress;
import ee.ivxv.common.util.I18nConsole;
import ee.ivxv.common.util.Json;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.apache.commons.collections4.Bag;
import org.apache.commons.collections4.bag.TreeBag;

import java.io.BufferedReader;
import java.io.FileReader;
import java.io.IOException;
import java.math.BigInteger;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.*;

public class IntegrityTool implements Tool.Runner<IntegrityArgs> {

    private static final String BDOC_EXTENSION = ".bdoc";
    private static final String ASICE_EXTENSION = ".asice";
    private static final String RAW_BALLOT_EXTENSION = ".ballot";
    private static final Character DIR_SEP = '/';
    private static final String EXPECTED_BALLOT_REF = "<voter-id>" + DIR_SEP + "<ballot-id>+XXXX";

    private static final String CSV_FIELD_SEPARATOR = "\t";
    private static final int LOG_HASH_ENTRY_FIELDS = 5;
    private static final int BB_ERRORS_ENTRY_FIELDS = 3;

    private static final String BB_ERRORS_BALLOT_REGEX = "^\\d{11}/\\d+\\+\\d{4}$";

    private final Logger log = LoggerFactory.getLogger(IntegrityTool.class);
    private final AuditContext ctx;
    private final I18nConsole console;

    public IntegrityTool(AuditContext ctx) {
        this.ctx = ctx;
        this.console = new I18nConsole(ctx.i.console, ctx.i.i18n);
    }

    @Override
    public boolean run(IntegrityArgs args) throws Exception {
        // Verify whether all provided files are found on the filesystem.
        boolean allFilesExist = verifyInputFiles(args);
        if (!allFilesExist) {
            console.println();
            return false;
        }

        boolean abortEarly = args.abortEarly.value();

        console.println();

        console.println(Msg.m_bb_loading, args.bb.value());
        List<RawBallotWithDigest> ballots = getBallots(args);
        console.println(Msg.m_bb_loaded);
        BigInteger checksumAll = sumBbHashes(ballots);

        console.println();

        // Checksums of ballots in the processed ballot box.
        console.println(Msg.m_anon_loading, args.anonBoxPath.value());
        AnonymousBallotBox anonBb = Json.read(args.anonBoxPath.value(), AnonymousBallotBox.class);
        AnonBbDigests anonDigests = getAnonBbDigests(anonBb);
        console.println(Msg.m_anon_loaded);
        console.println();

        console.println(Msg.m_log_accepted, args.acceptedLogPath.value());
        console.println(Msg.m_log_squashed, args.squashedLogPath.value());
        console.println(Msg.m_log_revoked, args.revokedLogPath.value());
        console.println(Msg.m_log_anonymised, args.anonLogPath.value());
        console.println(Msg.m_bb_errors, args.bbErrorsPath.value());

        console.println();

        // Checksums of all accepted (valid received) ballots.
        List<byte[]> acceptedDigests = getBallotDigestsFromLog(args.acceptedLogPath.value());
        BigInteger acceptedDigestsSum = sumHashes(acceptedDigests);

        Set<String> erroredBallotNames = getErroredBallotsNames(args.bbErrorsPath.value());
        BbDigests bbDigests = getValidInvalidSums(ballots, acceptedDigests, erroredBallotNames);
        if (bbDigests == null) {
            // This should happen only if ballots in the original ballot box appear in neither
            // the error log nor the accepted log.
            console.println(Msg.m_integrity_bb_consistent, getYesNoMessage(false));
            if (abortEarly) {
                console.println();
                return true;
            }
        } else {
            // This check should never fail but is here for clarity.
            // If true: accepted ballots in the original box are consistent with the log of accepted ballots.
            boolean validChecksumsMatch = bbDigests.validSum.equals(acceptedDigestsSum);

            // This check should never fail but is here for clarity.
            // If true: accepted ballots + refused ballots = original box.
            boolean sumIsConsistent = checksumAll.equals(bbDigests.invalidSum.add(bbDigests.validSum));

            boolean bbIsConsistent = validChecksumsMatch && sumIsConsistent;
            console.println(Msg.m_integrity_bb_consistent, getYesNoMessage(bbIsConsistent));
            if (!bbIsConsistent) {

                if (abortEarly) {
                    console.println();
                    return true;
                }

                // Check if accepted ballots are included in the original box.
                // This provides more information on the inconsistency failure.
                boolean acceptedMatch = matchHashesToBb(acceptedDigests, ballots);
                console.println(Msg.m_integrity_match_valid, getYesNoMessage(acceptedMatch));
            }
        }

        // Checksums of accepted ballots that were squashed.
        List<byte[]> squashedDigests = getBallotDigestsFromLog(args.squashedLogPath.value());
        BigInteger squashedDigestsSum = sumHashes(squashedDigests);

        // Checksums of accepted ballots that were revoked.
        List<byte[]> revokedDigests = getBallotDigestsFromLog(args.revokedLogPath.value());
        BigInteger revokedDigestsSum = sumHashes(revokedDigests);

        // Checksums of accepted ballots forwarded for tallying (remaining after squash + revoke).
        List<byte[]> anonLogDigests = getBallotDigestsFromLog(args.anonLogPath.value());
        BigInteger anonLogDigestsSum = sumHashes(anonLogDigests);

        // If true: anon ballots are consistent with the log of anon ballots.
        boolean anonDigestMatches = anonDigests.checksumSum.equals(anonLogDigestsSum);
        console.println(Msg.m_integrity_match_anon_logs, getYesNoMessage(anonDigestMatches));
        if (!anonDigestMatches && abortEarly) {
            console.println();
            return true;
        }

        BigInteger postprocessDigestSum = squashedDigestsSum.add(revokedDigestsSum.add(anonLogDigestsSum));

        // If true: anon + revoked + squashed = all accepted ballots.
        boolean acceptedDigestMatches = acceptedDigestsSum.equals(postprocessDigestSum);
        console.println(Msg.m_integrity_match_logs, getYesNoMessage(acceptedDigestMatches));

        // If all checks succeed, the following audit trail is verified:
        // - anon BB matches with anon log
        // - invalid ballots (log) + valid ballots (log) = original BB
        // - anon ballots (log) + squashed ballots (log) + revoked ballots (log) = valid ballots (log)
        // Inferred: anon ballots are part of the original ballot box.

        console.println();

        return true;
    }

    private Msg getYesNoMessage(boolean success) {
        if (success) return Msg.m_yes;
        return Msg.m_no;
    }

    private List<RawBallotWithDigest> getBallots(IntegrityArgs args) {
        List<RawBallotWithDigest> ballots = new ArrayList<>();

        // Read ballot box
        FileSource zipContainer = new ZipSource(args.bb.value());

        // Start progress, counting from 0 up to the amount of ballots in a ballot box
        int ballotCount = zipContainer.countFilesWithSuffix(Arrays.asList(BDOC_EXTENSION, ASICE_EXTENSION));
        Progress progress = console.startProgress(ballotCount);

        // Read ballots from a ballot box one by one
        zipContainer.processFiles((fileName, fileContent) -> {

            // Only look for voter signed ballots
            if (fileName.contains(BDOC_EXTENSION) || fileName.contains(ASICE_EXTENSION)) {
                try {
                    // Read voter signed container as byte stream
                    FileSource signedContainer = new ZipSourceRaw(fileContent.readAllBytes());

                    // Signed container is a zip file, that contains files inside
                    signedContainer.processFiles((ballotName, ballotContent) -> {

                        // Only look for a ballot file
                        if (ballotName.contains(RAW_BALLOT_EXTENSION)) {
                            try {
                                byte[] ballot = ballotContent.readAllBytes();
                                String ballotId = getBallotContainerId(fileName);
                                ballots.add(new RawBallotWithDigest(ballot, ballotId));
                            } catch (Exception e) {
                                throw new RuntimeException(e);
                            }
                        }
                    });
                } catch (IOException e) {
                    throw new RuntimeException(e);
                }
                progress.increase(1);
            }
        });

        progress.finish();

        return ballots;
    }

    private static String getBallotContainerId(String fileName) throws Exception {
        // Getting the container path in the BB is necessary
        // for ballot identification through the error log.
        int i = fileName.lastIndexOf(DIR_SEP);
        if (i < 0) {
            // It's OK to throw since we can reasonably expect a valid ballot box.
            throw new Exception("Expected name " + EXPECTED_BALLOT_REF + " but got " + fileName);
        }
        int j = fileName.lastIndexOf(DIR_SEP, i - 1);
        int extensionLength = fileName.contains(BDOC_EXTENSION) ?
                BDOC_EXTENSION.length() :
                ASICE_EXTENSION.length();
        String ballotContainerId = fileName.substring(j + 1, fileName.length() - extensionLength);
        if (!ballotContainerId.matches(BB_ERRORS_BALLOT_REGEX)) {
            throw new Exception("Expected name " + EXPECTED_BALLOT_REF + " but got " + fileName);
        }
        return ballotContainerId;
    }

    private BbDigests getValidInvalidSums(List<RawBallotWithDigest> bb,
                                          List<byte[]> acceptedDigests, Set<String> invalidBallotIds) {
        BigInteger invalidBallotChecksumSum = BigInteger.ZERO;
        BigInteger validBallotChecksumSum = BigInteger.ZERO;

        // Use a tree bag since comparing arrays of bytes does not work in a regular bag.
        // We must use a bag instead of a set since it might happen (statistically unlikely unless intentional)
        // that there are ballots with the same ciphertext. Therefore, a set would not correctly represent the
        // state of the accepted ballots.
        // Use a tree bag since comparing arrays of bytes does not work in a regular bag.
        Bag<byte[]> digests = new TreeBag<>(Arrays::compare);
        digests.addAll(acceptedDigests);

        // Notify in case of multiple matching ciphertexts so that the occurrence(s) can be investigated further.
        int recurringCts = digests.size() - digests.uniqueSet().size();
        if (recurringCts != 0) {
            console.println(Msg.m_integrity_recurring_ct);
            log.warn("There are {} ciphertext recurrences among the accepted ballots", recurringCts);

            Set<byte[]> recurringDigests = new HashSet<>();
            for (byte[] c : digests) if (digests.getCount(c) != 1) recurringDigests.add(c);
            for (byte[] c : recurringDigests) log.warn("Recurring ballot: {}", Base64.getEncoder().encodeToString(c));
            console.println();
        }

        // Avoid mutating the input set. These are IDs so they must be unique.
        Set<String> invalids = new HashSet<>(invalidBallotIds);

        boolean errored = false;

        for (RawBallotWithDigest ballot : bb) {
            // Make sure that each ballot is referenced in exactly one file.
            int seenCount = 0;

            if (invalids.remove(ballot.getId())) {
                invalidBallotChecksumSum =
                        invalidBallotChecksumSum.add(new BigInteger(1, ballot.getRawDigest()));
                ++seenCount;
            }

            if (digests.remove(ballot.getRawDigest(), 1)) {
                validBallotChecksumSum =
                        validBallotChecksumSum.add(new BigInteger(1, ballot.getRawDigest()));
                ++seenCount;
            }

            if (seenCount == 0) {
                log.warn("Ballot '{}' not found in the acceptance/rejection logs", ballot.getId());
                errored = true;
            }

            if (seenCount > 1) {
                log.warn("Ballot '{}' present in both the acceptance and rejection logs", ballot.getId());
                errored = true;
            }
        }

        if (errored) return null;

        // Make sure that there are not more log entries than there are ballots in the BB.
        if (!digests.isEmpty() || !invalids.isEmpty()) {
            for (String invalid : invalids)
                log.warn("Ballot '{}' is missing from the original ballot box", invalid);
            for (byte[] digest : digests)
                log.warn("Missing from the original ballot box: {}", Base64.getEncoder().encodeToString(digest));
            return null;
        }

        return new BbDigests(invalidBallotChecksumSum, validBallotChecksumSum);
    }

    /**
     * Reads the digests from a processor logfile.
     *
     * @param path the logfile path
     * @return a list of digests extracted from the logfile
     * @throws Exception If the logfile does not exist or cannot be properly read
     */
    private List<byte[]> getBallotDigestsFromLog(Path path) throws Exception {
        List<byte[]> hashes = new ArrayList<>();

        BufferedReader br = new BufferedReader(new FileReader(path.toFile()));
        String line;
        while ((line = br.readLine()) != null) {
            String[] values = line.split(CSV_FIELD_SEPARATOR);
            if (values.length != LOG_HASH_ENTRY_FIELDS) continue;

            byte[] hash = Base64.getDecoder().decode(values[1]);
            hashes.add(hash);
        }

        return hashes;
    }

    /**
     * Reads ballot identifiers from the ballot validation error report.
     *
     * @param path the report path
     * @return a set of ballot identifiers extracted from the report
     * @throws Exception If the report cannot be found be read
     */
    private Set<String> getErroredBallotsNames(Path path) throws Exception {
        // Use a set since there may be multiple errors pertaining to the
        // same ballot in the log.
        Set<String> erroredBallotNames = new HashSet<>();

        BufferedReader br = new BufferedReader(new FileReader(path.toFile()));
        String line;
        while ((line = br.readLine()) != null) {
            String[] values = line.split(CSV_FIELD_SEPARATOR);
            if (values.length != BB_ERRORS_ENTRY_FIELDS) continue;
            if (!values[0].matches(BB_ERRORS_BALLOT_REGEX)) continue;

            erroredBallotNames.add(values[0]);
        }

        return erroredBallotNames;
    }

    private static BigInteger sumHashes(List<byte[]> hashes) {
        BigInteger sum = BigInteger.ZERO;

        for (byte[] hash : hashes) {
            sum = sum.add(new BigInteger(1, hash));
        }

        return sum;
    }

    private static BigInteger sumBbHashes(List<RawBallotWithDigest> bb) {
        BigInteger sum = BigInteger.ZERO;

        for (RawBallotWithDigest ballot : bb) {
            sum = sum.add(new BigInteger(1, ballot.getRawDigest()));
        }

        return sum;
    }

    private AnonBbDigests getAnonBbDigests(AnonymousBallotBox bb) {
        // Also extract checksums in a list to avoid looping again later.
        List<byte[]> digests = new ArrayList<>();
        BigInteger sum = BigInteger.ZERO;

        for (Map<String, Map<String, List<byte[]>>> smap : bb.getDistricts().values()) {
            for (Map<String, List<byte[]>> qmap : smap.values()) {
                for (List<byte[]> clist : qmap.values()) {
                    for (byte[] c : clist) {
                        byte[] hash = HashType.SHA256.getFunction().digest(c);
                        digests.add(hash);
                        sum = sum.add(new BigInteger(1, hash));
                    }
                }
            }
        }

        return new AnonBbDigests(digests, sum);
    }

    /**
     * Verifies whether digests correspond to ballots in a ballot box.
     *
     * @param digests the digests of some set of ballots
     * @param bb      the ballot box to match against
     * @return whether all ballots were present in the box
     */
    private boolean matchHashesToBb(List<byte[]> digests, List<RawBallotWithDigest> bb) {
        boolean allOk = true;
        // Potential perf improvement: would it be faster to iterate over the full ballot box
        // and try to remove them from a set of digests? The return value is then whether
        // the set is empty. During normal execution this should never be invoked in any case.
        for (byte[] digest : digests) {
            boolean matches = bb.stream().anyMatch(b -> Arrays.equals(digest, b.getRawDigest()));
            if (!matches) {
                log.warn("Missing from the original ballot box: {}", Base64.getEncoder().encodeToString(digest));
                allOk = false;
            }
        }
        return allOk;
    }

    record BbDigests(BigInteger invalidSum, BigInteger validSum) {
    }

    record AnonBbDigests(List<byte[]> digests, BigInteger checksumSum) {
    }

    private boolean verifyInputFiles(IntegrityArgs args) {
        List<Path> paths = new ArrayList<>(List.of(
                args.bb.value(),
                args.anonBoxPath.value(),
                args.acceptedLogPath.value(),
                args.squashedLogPath.value(),
                args.revokedLogPath.value(),
                args.anonLogPath.value(),
                args.bbErrorsPath.value()));

        for (Path path : paths) {
            if (Files.notExists(path)) {
                console.println();
                console.println(Msg.e_file_missing, path);
                return false;
            }
        }

        return true;
    }


    public static class IntegrityArgs extends Args {
        Arg<Path> bb = Arg.aPath(Msg.arg_ballotbox, true, false);
        Arg<Path> anonBoxPath = Arg.aPath(Msg.arg_anon_bb);
        Arg<Path> acceptedLogPath = Arg.aPath(Msg.arg_log_accepted);
        Arg<Path> squashedLogPath = Arg.aPath(Msg.arg_log_squashed);
        Arg<Path> revokedLogPath = Arg.aPath(Msg.arg_log_revoked);
        Arg<Path> anonLogPath = Arg.aPath(Msg.arg_log_anonymised);
        Arg<Path> bbErrorsPath = Arg.aPath(Msg.arg_bb_errors);
        Arg<Boolean> abortEarly = Arg.aFlag(Msg.arg_abort_early).setDefault(true);

        public IntegrityArgs() {
            super();
            args.add(bb);
            args.add(anonBoxPath);
            args.add(acceptedLogPath);
            args.add(squashedLogPath);
            args.add(revokedLogPath);
            args.add(anonLogPath);
            args.add(bbErrorsPath);
            args.add(abortEarly);
        }
    }
}
