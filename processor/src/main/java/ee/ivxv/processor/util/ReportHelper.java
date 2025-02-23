package ee.ivxv.processor.util;

import ee.ivxv.common.crypto.CorrectnessUtil.CiphertextCorrectness;
import ee.ivxv.common.model.Ballot;
import ee.ivxv.common.model.BallotBox;
import ee.ivxv.common.model.DistrictList;
import ee.ivxv.common.service.bbox.Ref;
import ee.ivxv.common.service.bbox.Result;
import ee.ivxv.common.service.i18n.Message;
import ee.ivxv.common.service.i18n.MessageException;
import ee.ivxv.common.service.report.Reporter;
import ee.ivxv.common.service.report.Reporter.LogType;
import ee.ivxv.common.service.report.Reporter.Record;
import ee.ivxv.common.util.I18nConsole;
import ee.ivxv.common.util.Util;
import ee.ivxv.processor.Msg;
import ee.ivxv.processor.ProcessorContext;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.*;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.ConcurrentLinkedQueue;
import java.util.concurrent.atomic.LongAdder;
import java.util.function.BiPredicate;
import java.util.function.Function;
import java.util.stream.Stream;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class ReportHelper {

    private static final Logger log = LoggerFactory.getLogger(ReportHelper.class);

    public static final String OUT_VL_ERR = "voterlist_errors.txt";
    public static final String OUT_BB_ERR = "ballotbox_errors.txt";

    private final ProcessorContext ctx;
    private final I18nConsole console;
    private final Map<Result, LongAdder> bbErrors = new ConcurrentHashMap<>();
    private final Map<String, Queue<String>> errors = new ConcurrentHashMap<>();

    public ReportHelper(ProcessorContext ctx, I18nConsole console) {
        this.ctx = ctx;
        this.console = console;
    }

    public void writeLog1(Path dir, BallotBox bb, String disc) {
        writeLogN(dir, bb, disc, LogType.LOG1, (x, y) -> true);
    }

    public void writeLog2(Path dir, BallotBox bb, String disc, BiPredicate<String, String> filter) {
        writeLogN(dir, bb, disc, LogType.LOG2, filter);
    }

    public void writeLog3(Path dir, BallotBox bb, String disc, BiPredicate<String, String> filter) {
        writeLogN(dir, bb, disc, LogType.LOG3, filter);
    }

    public void writeLog2(Path dir, String eid, String disc, List<Reporter.LogNRecord> records) {
        writeLogN(dir, eid, disc, LogType.LOG2, records.stream());
    }

    public void writeEmptyLogFiles(Path outputDir, String logDiscriminator, LogType type, BallotBox ballotBox) {
        Set<String> questionIDs = new HashSet<>();

        ballotBox.getBallots().forEach((voterId, voterBallots) -> voterBallots.getBallots()
                .forEach(voterBallot -> questionIDs.addAll(voterBallot.getVotes().keySet())));

        // Does not print the output files to the console to avoid double printing.
        // Potential improvement: get rid of the double printing issue to consistently print the
        // output filename.
        questionIDs.forEach(questionID -> ctx.reporter.writeEmptyLog(outputDir, ballotBox.getElection(),
                logDiscriminator, type, questionID));
        }

    private void writeLogN(Path dir, String eid, String disc, LogType type,
            Stream<Reporter.LogNRecord> records) {
        try {
            console.println();
            console.println(Msg.m_writing_log_n, type.value);
            Map<String, Path> paths = ctx.reporter.writeLogN(dir, eid, disc, type, records);
            paths.values().forEach(p -> console.println(Msg.m_output_file, p));
        } catch (Exception e) {
            throw new MessageException(e, Msg.e_writing_log_n, type.value, dir, e);
        }
    }

    private void writeLogN(Path dir, BallotBox bb, String disc, LogType type,
            BiPredicate<String, String> filter) {
        Map<String, List<Record>> rmap = new LinkedHashMap<>();

        bb.getBallots()
                .forEach((voterId, vb) -> vb.getBallots() //
                        .forEach(b -> b.getVotes().keySet().stream()
                                .filter(qid -> filter.test(voterId, qid)) //
                                .forEach(qid -> rmap.computeIfAbsent(qid, x -> new ArrayList<>())
                                        .add(ctx.reporter.newLog123Record(voterId, b, qid)))));

        try {
            console.println();
            console.println(Msg.m_writing_log_n, type.value);
            Map<String, Path> paths = ctx.reporter.writeRecords(dir, bb.getElection(), disc, type, rmap);
            paths.values().forEach(p -> console.println(Msg.m_output_file, p));
        } catch (Exception e) {
            throw new MessageException(e, Msg.e_writing_log_n, type.value, dir, e);
        }
    }

    public void writeIVoterList(Path jsonOut, Path pdfOut, BallotBox bb, DistrictList dl) {
        try {
            console.println();
            console.println(Msg.m_writing_ivoter_list);
            ctx.reporter.writeIVoterList(jsonOut, pdfOut, bb, dl);
            console.println(Msg.m_output_file, jsonOut);
            if (pdfOut != null) {
                console.println(Msg.m_output_file, pdfOut);
            }
        } catch (Exception e) {
            throw new MessageException(e, Msg.e_writing_ivoter_list, e);
        }
    }

    public void writeRevocationReport(Path out, String electionId, List<Reporter.Record> records, Reporter.AnonymousFormatter formatter) {
        try {
            console.println();
            console.println(Msg.m_writing_revocation_report);
            ctx.reporter.write(out, electionId, records, formatter);
            console.println(Msg.m_output_file, out);
        } catch (Exception e) {
            throw new MessageException(e, Msg.e_writing_revocation_report, out, e);
        }
    }

    public void reportVlErrors(Enum<?> key, Object... args) {
        reportErrors(OUT_VL_ERR, key, args);
    }

    public void reportBbError(Ref.BbRef ref, Result res, Object... args) {
        String voter = ref == null ? "?" : ref.voter;
        String ballot = ref == null ? "?" : ref.ballot;
        String r = String.format("%s/%s", voter, ballot);
        log.error("Error while reading ballot box: {}, {}, {}", voter, ballot, res);
        reportBbError(m -> new Message(Msg.e_bb_ballot_processing, voter, ballot, m), r, res, args);
    }

    public void reportRegError(Ref.RegRef ref, Result res, Object... args) {
        String r = ref == null ? "?" : ref.ref;
        log.error("Error while reading registration data: {}, {}", r, res);
        reportBbError(m -> new Message(Msg.e_reg_record_processing, r, m), r, res, args);
    }

    private void reportBbError(Function<Message, Message> provider, String ref, Result res,
            Object... args) {
        bbErrors.computeIfAbsent(res, r -> new LongAdder()).increment();
        Msg key = translate(res);
        if (key == null) {
            return;
        }
        Message innerMsg = new Message(key, args);
        Message msg = provider.apply(innerMsg);
        reportErrors(OUT_BB_ERR, s -> String.format("%s\t%s\t%s", ref, res, s), msg.key, msg.args);
    }

    public void reportAbbError(String voterId, Ballot b, String qid, CiphertextCorrectness res) {
        Msg key = translate(res);
        if (key == null) {
            return;
        }
        Message innerMsg = new Message(key);
        Message msg = new Message(Msg.e_bb_ciphertext_checking, voterId, b.getId(), qid, innerMsg);
        String ref = String.format("%s/%s", voterId, b.getId());

        reportErrors(OUT_BB_ERR, s -> String.format("%s\t%s\t%s", ref, res, s), msg.key, msg.args);
    }

    private void reportErrors(String type, Enum<?> key, Object... args) {
        reportErrors(type, s -> s, key, args);
    }

    private void reportErrors(String type, Function<String, String> fmt, Enum<?> key,
            Object... args) {
        if (!ctx.args.quiet.value()) {
            console.println(key, args);
        }
        errors.computeIfAbsent(type, x -> new ConcurrentLinkedQueue<>())
                .add(fmt.apply(console.i18n.get(key, args)));
    }

    private Msg translate(Result res) {
        switch (res) {
            case INVALID_FILE_NAME:
                return Msg.e_bb_invalid_file_name;
            case INVALID_FILE_SIZE:
                return Msg.e_bb_invalid_file_size;
            case MISSING_FILE:
                return Msg.e_bb_missing_file;
            case REPEATED_FILE:
                return Msg.e_bb_repeated_file;
            case UNKNOWN_FILE_TYPE:
                return Msg.e_bb_unknown_file_type;
            case INVALID_BALLOT_SIGNATURE:
                return Msg.e_ballot_signature_invalid;
            case MISSING_VOTER_SIGNATURE:
                return Msg.e_ballot_missing_voter_signature;
            case VOTER_NOT_FOUND:
                return Msg.e_active_voter_not_found;
            case VOTERLIST_NOT_FOUND:
                return Msg.e_active_voterlist_not_found;
            case TIME_BEFORE_START:
                return Msg.e_time_before_start;
            case REG_RESP_INVALID:
                return Msg.e_reg_resp_invalid;
            case REG_REQ_INVALID:
                return Msg.e_reg_req_invalid;
            case REG_RESP_NOT_UNIQUE:
                return Msg.e_reg_resp_not_unique;
            case REG_REQ_NOT_UNIQUE:
                return Msg.e_reg_req_not_unique;
            case REG_NO_NONCE:
                return Msg.e_reg_resp_no_nonce;
            case REG_NONCE_NOT_SIG:
                return Msg.e_reg_resp_nonce_not_sig;
            case REG_NONCE_ALG_MISMATCH:
                return Msg.e_reg_resp_nonce_alg_mismatch;
            case REG_NONCE_SIG_INVALID:
                return Msg.e_reg_resp_nonce_sig_invalid;
            case UNKNOWN_FILE_IN_VOTE_CONTAINER:
                return Msg.e_unknown_file_in_vote_container;
            case TECHNICAL_ERROR:
                return Msg.e_tehcnical_error;
            case REG_RESP_REQ_UNMATCH:
                return Msg.e_reg_resp_req_unmatch;
            case REG_REQ_WITHOUT_BALLOT:
                return Msg.e_reg_req_without_ballot;
            case BALLOT_WITHOUT_REG_REQ:
                return Msg.e_ballot_without_reg_req;
            case SAME_TIME_AS_LATEST:
                return Msg.e_same_time_as_latest;
            case INVALID_SIGNATURE_PROFILE:
                return Msg.e_invalid_signature_profile;
            case OK:
                return null;
            default:
                throw new RuntimeException("Unhandled ballot box processing result: " + res);
        }
    }

    private Msg translate(CiphertextCorrectness res) {
        switch (res) {
            case INVALID:
                return Msg.e_ciphertext_invalid;
            case INVALID_BYTES:
                return Msg.e_ciphertext_invalid_bytes;
            case INVALID_GROUP:
                return Msg.e_ciphertext_invalid_group;
            case INVALID_POINT:
                return Msg.e_ciphertext_invalid_point;
            case INVALID_QR:
                return Msg.e_ciphertext_invalid_qr;
            case INVALID_RANGE:
                return Msg.e_ciphertext_invalid_range;
            case VALID:
                return null;
            default:
                throw new RuntimeException("Unhandled ciphertext correctness result: " + res);
        }
    }

    public long countBbErrors() {
        return bbErrors.values().stream().mapToLong(LongAdder::sum).sum();
    }

    public long countBbErrors(Result type) {
        return bbErrors.getOrDefault(type, new LongAdder()).sum();
    }

    public void writeVlErrors(Path out) {
        writeErrors(out, OUT_VL_ERR, Msg.e_vl_error_report);
    }

    public void writeBbErrors(Path out) {
        writeErrors(out, OUT_BB_ERR, Msg.e_bb_error_report);
    }

    private void writeErrors(Path out, String file, Enum<?> key) {
        Path path = out.resolve(file);
        try {
            // The file should always be created for consistency...
            Util.createFile(path);

            // ...however the file may be empty if there is nothing to report.
            Queue<String> errs = errors.get(file);
            if (errs == null || errs.isEmpty()) {
                return;
            }

            // Only mention that there were errors if there indeed were any.
            console.println(key, path);
            Files.write(path, errs);
        } catch (Exception e) {
            log.error("Error occurred while writing error report {}: {}", file, e.getMessage(), e);
            throw new MessageException(Msg.e_writing_error_report, path, e);
        }
    }

}
