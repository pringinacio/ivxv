package ee.ivxv.audit.model;


import com.fasterxml.jackson.annotation.JsonCreator;
import com.fasterxml.jackson.annotation.JsonProperty;
import ee.ivxv.common.model.CandidateList;
import ee.ivxv.common.model.DistrictList;

import java.util.HashMap;
import java.util.LinkedHashMap;
import java.util.Map;

/**
 * JSON serializable structure for holding the tally of the votes.
 */
public class Tally {
    public static final String INVALID_VOTE_ID = "invalid";
    private final String election;
    private final Map<String, Map<String, Map<String, Integer>>> byParish;
    private final Map<String, Map<String, Integer>> byDistrict;


    @JsonCreator Tally(@JsonProperty("election") String election,
                       @JsonProperty("byparish") Map<String, Map<String, Map<String, Integer>>> byParish,
                       @JsonProperty("bydistrict") Map<String, Map<String, Integer>> byDistrict) {
        this.election = election;
        this.byParish = byParish;
        this.byDistrict = byDistrict;
    }

    /**
     * Get the election identifier.
     *
     * @return
     */
    public String getElection() {
        return election;
    }

    /**
     * @return Returns a map from district id to a map from station id to a map from candidate id to
     *         number of received votes.
     */
    public Map<String, Map<String, Map<String, Integer>>> getByParish() {
        return byParish;
    }

    /**
     * @return Returns a map from district id to a map from candidate id to number of received
     *         votes.
     */
    public Map<String, Map<String, Integer>> getByDistrict() {
        return byDistrict;
    }

    public Map<String, Map<String, Integer>> computeByDistrict() {
        Map<String, Map<String, Integer>> res = new LinkedHashMap<>();
        getByParish().forEach((d, sMap) -> {
            Map<String, Integer> ccMap = res.computeIfAbsent(d, tmp -> new LinkedHashMap<>());
            sMap.forEach((s, cMap) -> cMap.forEach((c, count) -> {
                ccMap.compute(c, (cc, ccount) -> ccount == null ? count : ccount + count);
            }));

        });
        return res;
    }
}

