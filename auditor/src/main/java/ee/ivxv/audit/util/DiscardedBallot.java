package ee.ivxv.audit.util;

import com.fasterxml.jackson.annotation.JsonCreator;
import com.fasterxml.jackson.annotation.JsonIgnore;
import com.fasterxml.jackson.annotation.JsonProperty;

import java.util.List;

public class DiscardedBallot {
    private final String election;

    private final List<VoteJson> discarded;

    @JsonCreator
    public DiscardedBallot(
            @JsonProperty("election") String election,
            @JsonProperty("invalid") List<VoteJson> invalid) {
        this.election = election;
        this.discarded = invalid;
    }

    public List<VoteJson> getDiscarded() {
        return discarded;
    }

    @JsonIgnore
    public int getCount() {
        return discarded.size();
    }

    public static class VoteJson {
        private final String district;
        private final String station;
        private final String question;
        private final byte[] vote;

        @JsonCreator
        private VoteJson(
                @JsonProperty("district") String district,
                @JsonProperty("station") String station,
                @JsonProperty("question") String question,
                @JsonProperty("vote") byte[] vote) {
            this.district = district;
            this.station = station;
            this.question = question;
            this.vote = vote;
        }

        public byte[] getVote() {
            return vote;
        }
    }

}
