package ee.ivxv.key.model;

import com.fasterxml.jackson.annotation.JsonProperty;
import ee.ivxv.common.model.CandidateList;
import ee.ivxv.common.model.DistrictList;

import java.util.LinkedHashMap;
import java.util.HashMap;
import java.util.ArrayList;
import java.util.Map;
import java.util.List;

/**
 * JSON serializable structure for holding the decrypted ballot box.
 */
public class PlainBallotBox {
    private final String election;
    private final Map<String, Map<String, List<String>>> byParish = new HashMap<>();

    /**
     * Initialize using values.
     *
     * @param election   Election identifier.
     * @param candidates List of candidates.
     * @param districts  List of districts.
     */
    public PlainBallotBox(String election, CandidateList candidates, DistrictList districts) {
        this.election = election;
        init(candidates, districts);
    }

    private void init(CandidateList candidates, DistrictList districts) {
        // No candidate validity check: we assume the check was done during ballot validation.
        districts.getDistricts().forEach((dId, d) -> {
            Map<String, List<String>> plaintexts = new LinkedHashMap<>();
            byParish.put(dId, plaintexts);
            d.getParish().forEach(p -> {
                List<String> pts = new ArrayList<>();
                plaintexts.put(p, pts);
            });
        });
    }

    public String getElection() {
        return election;
    }

    @JsonProperty("byparish")
    public Map<String, Map<String, List<String>>> getByParish() {
        return byParish;
    }

    @JsonProperty("bydistrict")
    public Map<String, List<String>> getByDistrict() {
        Map<String, List<String>> res = new LinkedHashMap<>();
        getByParish().forEach((d, sMap) -> {
            List<String> pList = res.computeIfAbsent(d, tmp -> new ArrayList<>());
            sMap.forEach((s, plains) -> pList.addAll(plains));
        });
        return res;
    }
}
