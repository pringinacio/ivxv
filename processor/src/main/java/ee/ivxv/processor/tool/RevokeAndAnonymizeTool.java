package ee.ivxv.processor.tool;

import ee.ivxv.common.M;
import ee.ivxv.common.cli.Arg;
import ee.ivxv.common.cli.Args;
import ee.ivxv.common.cli.Tool;
import ee.ivxv.common.model.*;
import ee.ivxv.common.service.container.Container;
import ee.ivxv.common.service.container.DataFile;
import ee.ivxv.common.service.i18n.Message;
import ee.ivxv.common.service.i18n.MessageException;
import ee.ivxv.common.service.report.Reporter;
import ee.ivxv.common.service.report.Reporter.RevokeAction;
import ee.ivxv.common.util.*;
import ee.ivxv.processor.Msg;
import ee.ivxv.processor.ProcessorContext;
import ee.ivxv.processor.util.ReportHelper;

import java.nio.file.Path;
import java.util.*;
import java.util.concurrent.ConcurrentHashMap;
import java.util.stream.Collectors;

public class RevokeAndAnonymizeTool implements Tool.Runner<RevokeAndAnonymizeTool.RevokeAndAnonymizeArgs> {

    private static final String OUT_BB_TMPL = "bb-4.json";
    private static final String OUT_RR_TMPL = "revocation-report.csv";
    private static final String OUT_IVLJSON_TMPL = "ivoterlist.json";
    private static final String OUT_LOG_DISCRIMINATOR_REVOKE = "revoke";
    private static final String OUT_LOG_DISCRIMINATOR_ANONYMIZE = "anonymize";
    private static final Map<String, Object> EMPTY = new HashMap<>();
    final ProcessorContext ctx;
    final I18nConsole console;
    final ReportHelper reporter;
    private final ToolHelper tool;

    final Map<String, Map<String, Object>> excluded = new ConcurrentHashMap<>();

    public RevokeAndAnonymizeTool(ProcessorContext ctx) {
        this.ctx = ctx;
        console = new I18nConsole(ctx.i.console, ctx.i.i18n);
        reporter = new ReportHelper(ctx, console);
        tool = new ToolHelper(console, ctx.container, ctx.bbox);
    }

    @Override
    public boolean run(RevokeAndAnonymizeArgs args) throws Exception {
        tool.checkBbChecksum(args.bb.value(), args.bbChecksum.value());

        BallotBox bb = tool.readJsonBb(args.bb.value(), BallotBox.Type.INVALID_CIPHERTEXTS_REMOVED);
        DistrictList dl = tool.readJsonDistricts(args.districts.value());
        RlLoader loader = new RlLoader(bb);
        Path out = args.out.value();

        applyRevocationLists(bb, args.revLists.value(), loader);

        Path OUT_IVLJSON = Util.prefixedPath(bb.getElection(), OUT_IVLJSON_TMPL);
        Path OUT_RR = Util.prefixedPath(bb.getElection(), OUT_RR_TMPL);
        Path OUT_RR_ANONYMOUS = Util.prefixedPath(bb.getElection(), OUT_RR_TMPL + ".anonymous");
        Path OUT_BB = Util.prefixedPath(bb.getElection(), OUT_BB_TMPL);

        reporter.writeIVoterList(out.resolve(OUT_IVLJSON), null, bb, dl);
        reporter.writeRevocationReport(out.resolve(OUT_RR), bb.getElection(), loader.revRecords,
                Reporter.AnonymousFormatter.NOT_ANONYMOUS);
        reporter.writeRevocationReport(out.resolve(OUT_RR_ANONYMOUS), bb.getElection(), loader.revRecords,
                Reporter.AnonymousFormatter.REVOCATION_REPORT_CSV);

        // There should still be an empty .log2 file even if there are no revoked records
        reporter.writeEmptyLogFiles(args.out.value(), OUT_LOG_DISCRIMINATOR_REVOKE, Reporter.LogType.LOG2, bb);
        reporter.writeLog2(out, bb.getElection(), OUT_LOG_DISCRIMINATOR_REVOKE, loader.getLog2Records());

        // There should still be an empty .log3 file even if there are no anonymised records
        reporter.writeEmptyLogFiles(args.out.value(), OUT_LOG_DISCRIMINATOR_ANONYMIZE, Reporter.LogType.LOG3, bb);
        reporter.writeLog3(args.out.value(), bb, OUT_LOG_DISCRIMINATOR_ANONYMIZE,
                (voterId, qid) -> !excluded.getOrDefault(voterId, EMPTY).containsKey(qid));

        AnonymousBallotBox abb = anonymize(bb);

        tool.writeJsonBb(abb, out.resolve(OUT_BB));

        return true;
    }

    private AnonymousBallotBox anonymize(BallotBox bb) {
        console.println();
        console.println(Msg.m_anonymizing_ballot_box);

        return bb.anonymize();
    }

    private void applyRevocationLists(BallotBox bb, List<Path> paths, RlLoader loader) {
        console.println();
        console.println(Msg.m_applying_revocation_lists);

        bb.revokeDoubleVotes(paths.stream().map(p -> () -> loader.load(p)), loader::collect);
        loader.reportAfterLoading();

        console.println();
        console.println(M.m_bb_type, bb.getType());
        console.println(M.m_bb_numof_ballots, bb.getNumberOfBallots());
    }

    private class RlLoader {

        private final BallotBox bb;
        final List<Reporter.Record> revRecords = new ArrayList<>();
        /**
         * The LOG2 records. At most one entry per ballot. Using ordered map to retain the original
         * iteration order in the output.
         */
        private final Map<String, Reporter.LogNRecord> log2Records = new LinkedHashMap<>();

        int ballotCount;
        RevocationList current;
        String operator;
        int currentCount;

        RlLoader(BallotBox bb) {
            this.bb = bb;
            ballotCount = bb.getNumberOfBallots();
            bb.getBallots().keySet().forEach(vid -> log2Records.put(vid, null));
        }

        RevocationList load(Path path) {
            reportAfterLoading();

            loadRevocationList(path);

            reportBeforeLoading();

            currentCount = 0;

            return current;
        }

        private void reportBeforeLoading() {
            if (current.isRevoke()) {
                console.println(Msg.m_rl_revoke_start);
                console.println(Msg.m_rl_revoke_ballots_before, ballotCount);
            } else {
                console.println(Msg.m_rl_restore_start);
                console.println(Msg.m_rl_restore_ballots_before, ballotCount);
            }
        }

        void reportAfterLoading() {
            if (current != null) {
                if (current.isRevoke()) {
                    ballotCount -= currentCount;
                    console.println(Msg.m_rl_revoke_count, currentCount);
                    console.println(Msg.m_rl_revoke_ballots_after, ballotCount);
                    console.println(Msg.m_rl_revoke_done);
                } else {
                    ballotCount += currentCount;
                    console.println(Msg.m_rl_restore_count, currentCount);
                    console.println(Msg.m_rl_restore_ballots_after, ballotCount);
                    console.println(Msg.m_rl_restore_done);
                }
            }
        }

        private void loadRevocationList(Path path) {
            try {
                console.println();
                console.println(Msg.m_rl_loading, path);
                ctx.container.requireContainer(path);
                Container c = ctx.container.read(path.toString());
                console.println(Msg.m_rl_loaded);

                ContainerHelper ch = new ContainerHelper(console, c);
                DataFile file = ch.getSingleFileAndReport(Msg.m_rl_arg_for_cont);

                console.println(Msg.m_rl_checking_integrity);
                RevocationList rl = Json.read(file.getStream(), RevocationList.class);
                if (rl.getElection() != null && !rl.getElection().equals(bb.getElection())) {
                    throw new MessageException(Msg.e_rl_election_id, rl.getElection(),
                            bb.getElection());
                }
                console.println(Msg.m_rl_data_is_integrous);

                Msg totalMsg = rl.isRevoke() ? Msg.m_rl_revoke_total : Msg.m_rl_restore_total;
                console.println(totalMsg, rl.getPersons().size());

                current = rl;
                operator = ch.getSignerNames();
            } catch (Exception e) {
                throw new MessageException(e, Msg.e_rl_read_error, path, e);
            }
        }

        void collect(String voterId, Ballot b, boolean revoke, boolean success) {
            if (success) {
                RevokeAction action = revoke ? RevokeAction.REVOKED : RevokeAction.RESTORED;

                revRecords.add(ctx.reporter.newRevocationRecord(action, voterId, b, operator));
                log2Records.put(voterId, revoke ? ctx.reporter.newLog123Record(voterId, b) : null);
                currentCount++;
            } else {
                Msg key = b == null ? Msg.e_rl_voter_not_found_in_bb
                        : revoke ? Msg.e_rl_ballot_already_revoked
                        : Msg.e_rl_ballot_already_restored;
                Message innerMsg = new Message(key, voterId);
                Message msg = new Message(Msg.e_rl_processing_error, innerMsg);
                console.println(msg.key, msg.args);
            }
        }

        List<Reporter.LogNRecord> getLog2Records() {
            return log2Records.values().stream().filter(r -> r != null)
                    .collect(Collectors.toList());
        }
    }

    public static class RevokeAndAnonymizeArgs extends Args {

        Arg<Path> bb = Arg.aPath(Msg.arg_ballotbox, true, false);
        Arg<Path> bbChecksum = Arg.aPath(Msg.arg_ballotbox_checksum, true, false);
        Arg<Path> districts = Arg.aPath(Msg.arg_districts, true, false);
        Arg<List<Path>> revLists =
                Arg.listOfPaths(Msg.arg_revocationlists, true, false).setOptional();

        Arg<Path> out = Arg.aPath(Msg.arg_out, false, null);

        public RevokeAndAnonymizeArgs() {
            args.add(bb);
            args.add(bbChecksum);
            args.add(districts);
            args.add(revLists);
            args.add(out);
        }

    }

}
