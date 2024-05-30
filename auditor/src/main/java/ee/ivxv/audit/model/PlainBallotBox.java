package ee.ivxv.audit.model;

import com.fasterxml.jackson.annotation.JsonCreator;
import com.fasterxml.jackson.annotation.JsonProperty;


import java.util.ArrayList;
import java.util.LinkedHashMap;
import java.util.Map;
import java.util.List;

/**
 * JSON serializable structure for holding the decrypted ballot box.
 */
public record PlainBallotBox(String election, Map<String, Map<String, List<String>>> byParish,
                             Map<String, List<String>> byDistrict) {
    @JsonCreator
    public PlainBallotBox(@JsonProperty("election") String election,
                          @JsonProperty("byparish") Map<String, Map<String, List<String>>> byParish,
                          @JsonProperty("bydistrict") Map<String, List<String>> byDistrict) {
        this.election = election;
        this.byParish = byParish;
        this.byDistrict = byDistrict;
    }

    public Map<String, List<String>> computeByDistrict() {
        Map<String, List<String>> res = new LinkedHashMap<>();
        byParish().forEach((d, sMap) -> {
            List<String> pList = res.computeIfAbsent(d, tmp -> new ArrayList<>());
            sMap.forEach((s, plains) -> pList.addAll(plains));
        });
        return res;
    }
}
