package ee.ivxv.audit.util;

import ee.ivxv.common.crypto.hash.HashType;


public class RawBallotWithDigest {
    private final String id;

    private final byte[] ballot;

    private final byte[] rawDigest;

    public RawBallotWithDigest(byte[] ballot, String id) {
        this.id = id;
        this.ballot = ballot;
        this.rawDigest = HashType.SHA256.getFunction().digest(ballot);
    }

    public String getId() {
        return this.id;
    }

    public byte[] getBallot() {
        return this.ballot;
    }

    public byte[] getRawDigest() {
        return this.rawDigest;
    }
}
