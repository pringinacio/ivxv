package ee.ivxv.audit;

import ch.qos.cal10n.BaseName;
import ch.qos.cal10n.LocaleData;
import ee.ivxv.common.util.NameHolder;

@BaseName("i18n.audit-msg")
@LocaleData(defaultCharset = "UTF-8", value = {})
public enum Msg implements NameHolder {
    /*-
     * The part of the enum name until the first '_' (including) is excluded from the getName().
     * This is a means to provide multiple translations for the same tool or argument name.
     */

    // App
    app_audit,

    // Error messages
    e_proof_verif_false, e_proof_verif_exception,
    e_file_missing,

    // Tools
    tool_decrypt, tool_mixer, tool_convert, tool_integrity,

    // Tool arguments
    arg_abort_early, //
    arg_proofs, arg_invalidity_proofs, arg_discarded, //
    arg_hash, arg_links("l"), arg_out("o"), arg_pbb("p"), arg_pub("p"), //
    arg_revoke("r"), arg_seed("s"), arg_storage("s"), arg_signaturepub, arg_threads("t"), //
    arg_input_bb, arg_output_bb, arg_protinfo, arg_proofdir, arg_threaded, //
    arg_ballotbox, arg_ballotbox_checksum, arg_anon_bb, arg_plain_bb, //
    arg_log_accepted, arg_log_squashed, arg_log_revoked, arg_log_anonymised, arg_bb_errors, //
    arg_tally, arg_candidates, arg_districts, arg_questioncount, //

    // Messages
    m_yes, m_no, //
    m_pub_loading, m_pub_loaded, m_failurecount, m_verify_start, m_verify_finish, //
    m_decrypt_bb_has_proofs, m_decrypt_bb_has_invalids, m_decrypt_one_per_file, //
    m_decrypt_verified_valid, m_decrypt_verified_invalid, //
    m_decrypt_consistent_valids_begin, //
    m_decrypt_consistent_valids, m_decrypt_consistent_invalids, m_decrypt_dec_consistent, //
    m_decrypt_success, m_decrypt_failure, //
    m_shuffle_proof_loading, m_shuffle_proof_failed_reason, m_shuffle_proof_succeeded, //
    m_shuffle_proof_failed, m_convert_publickey_failed, m_convert_publickey_succ, //
    m_convert_bb_to_bt_failed, m_convert_bb_to_bt_succ, m_convert_bt_to_bb_failed, //
    m_convert_bt_to_bb_succ, m_shuffle_step, m_shuffle_read, m_shuffle_read_prot_info, //
    m_shuffle_read_pubkey, m_shuffle_read_pc, m_shuffle_read_posc, m_shuffle_read_posr, //
    m_shuffle_read_ciphs, m_shuffle_read_shuffled, //
    m_shuffle_verify, m_shuffle_verify_params, m_shuffle_verify_ni, m_shuffle_verify_permutation, //
    m_shuffle_verify_rerandomisation, //
    m_plain_loading, m_plain_loaded, m_plain_count, //
    m_tally_loading, m_tally_loaded, m_tally_match, //
    m_discarded_loading, m_discarded_loaded, m_discarded_count, //
    m_bb_loading, m_bb_loaded, m_anon_loading, m_anon_loaded, //
    m_log_accepted, m_log_squashed, m_log_revoked, m_log_anonymised, m_bb_errors, //
    m_integrity_match_anon, m_integrity_match_anon_logs, //
    m_integrity_match_valid, m_integrity_match_logs, //
    m_integrity_bb_consistent, m_integrity_bb_inconsistent, m_integrity_recurring_ct, //
    m_tally_start, m_tally_done, m_out_tally, m_tally_read, //

    e_abb_invalid_question_count, //

    ;


    private final String shortName;

    Msg() {
        this(null);
    }

    Msg(String shortName) {
        this.shortName = shortName;
    }

    @Override
    public String getShortName() {
        return shortName;
    }

    @Override
    public String getName() {
        return extractName(name());
    }

    @Override
    public Enum<?> getKey() {
        return this;
    }
}
