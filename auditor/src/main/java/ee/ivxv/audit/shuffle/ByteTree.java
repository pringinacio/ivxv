package ee.ivxv.audit.shuffle;

import ee.ivxv.common.math.ECGroupElement;
import ee.ivxv.common.math.GroupElement;
import ee.ivxv.common.math.ModPGroup;
import ee.ivxv.common.math.ModPGroupElement;
import ee.ivxv.common.math.ProductGroupElement;
import ee.ivxv.common.util.Util;
import java.io.DataInputStream;
import java.io.File;
import java.io.FileInputStream;
import java.io.IOException;
import java.io.OutputStream;
import java.io.UnsupportedEncodingException;
import java.math.BigInteger;
import java.nio.ByteBuffer;
import java.nio.ByteOrder;
import java.nio.file.Path;
import java.util.HexFormat;
import org.bouncycastle.util.Arrays;

/**
 * ByteTree decodes the ByteTree (BT) format as defined in Verificatum user manual.
 */
public class ByteTree {
    public int getLength() {
        return 0;
    }

    public int getEncodedLength() {
        return 0;
    }

    public byte[] getEncoded() {
        return null;
    }

    public void writeEncoded(OutputStream out) throws IOException {
        // Overridden
    }

    public boolean isLeaf() {
        return false;
    }

    public String getPrefix() {
        return "";
    }

    public String toString(int indent) {
        return "";
    }

    @Override
    public String toString() {
        return toString(0);
    }

    /**
     * Node is a recursive object in a ByteTree which consists of an array of nodes and leafs.
     */
    public static class Node extends ByteTree {
        private final ByteTree[] nodes;

        public static final byte PREFIX = 0;

        /**
         * Initialize Node from an array of ByteTree objects.
         *
         * @param nodes Array of ByteTree objects.
         */
        public Node(ByteTree[] nodes) {
            this.nodes = nodes;
        }

        /**
         * Initialize Node from an array of BigInteger objects.
         *
         * @param ints Array of BigInteger objects
         */
        public Node(BigInteger[] ints) {
            nodes = new Leaf[ints.length];
            for (int i = 0; i < nodes.length; i++) {
                nodes[i] = new Leaf(ints[i]);
            }
        }

        /**
         * Initialize Node from an array of GroupElements.
         *
         * @param elements Array of GroupElements
         */
        public Node(GroupElement[] elements) {
            nodes = new ByteTree[elements.length];
            for (int i = 0; i < nodes.length; i++) {
                if (elements[i] instanceof ProductGroupElement) {
                    nodes[i] = new Node((ProductGroupElement) elements[i]);
                } else {
                    nodes[i] = new Leaf(elements[i]);
                }
            }
        }

        /**
         * Initialize Node from ProductGroupElement.
         *
         * @param element A ProductGroupElement element
         */
        public Node(ProductGroupElement element) {
            nodes = new ByteTree[element.getElements().length];
            for (int i = 0; i < nodes.length; i++) {
                GroupElement ge = element.getElements()[i];
                if (ge instanceof ProductGroupElement) {
                    nodes[i] = new Node((ProductGroupElement) ge);
                } else {
                    nodes[i] = new Leaf(ge);
                }
            }
        }

        /**
         * Return the array of ByteTree objects.
         *
         * @return Array of ByteTree objects.
         */
        public ByteTree[] getNodes() {
            return nodes;
        }

        /**
         * Get the number of elements in the ByteTree array.
         *
         * @return Length of array.
         */
        @Override
        public int getLength() {
            return nodes.length;
        }

        /**
         * Get the byte-length of the Node.
         *
         * @return Byte-length of the Node.
         */
        @Override
        public int getEncodedLength() {
            int sum = 5;
            for (ByteTree n : nodes) {
                sum += n.getEncodedLength();
            }
            return sum;
        }

        /**
         * Encode the Node as bytes.
         *
         * @return Node encoded as bytes
         */
        @Override
        public byte[] getEncoded() {
            byte[] ret = new byte[getEncodedLength()];
            ret[0] = PREFIX;
            System.arraycopy(ByteTree.parse_from_int(getLength()), 0, ret, 1, 4);
            int pt = 5;
            for (int i = 0; i < getLength(); i++) {
                byte[] enc = getNodes()[i].getEncoded();
                System.arraycopy(enc, 0, ret, pt, enc.length);
                pt += enc.length;
            }
            return ret;
        }

        /**
         * Encode the Node and write it to out.
         *
         * @param out Stream to write the ByteTree description of the Node.
         * @throws IOException When writing to out fails
         */
        @Override
        public void writeEncoded(OutputStream out) throws IOException {
            out.write(PREFIX);
            out.write(ByteTree.parse_from_int(getLength()));
            for (int i = 0; i < getLength(); i++) {
                getNodes()[i].writeEncoded(out);
            }
        }

        /**
         * @return false
         */
        @Override
        public boolean isLeaf() {
            return false;
        }

        /**
         * Parse an array of bytes into a node, starting at offset.
         *
         * @param b Array of bytes to be parsed.
         * @param offset The offset in the array to start parsing
         * @return A decoded Node instance.
         */
        public static Node parse(byte[] b, int offset) {
            int len = parse_int_fullbytes(b, offset);
            ByteTree[] bts = new ByteTree[len];
            int pt = 0;
            for (int i = 0; i < len; i++) {
                bts[i] = ByteTree.parse(b, offset + pt + 5);
                pt += bts[i].getEncodedLength();
            }
            Node n = new Node(bts);
            return n;
        }

        /**
         * Parse a DataInputStream into a node.
         *
         * The given DataInputStream must be pointed to the beginning of a valid Node element.
         *
         * @param is DataInputStream to be used for reading a node, seeked to the beginning of Node
         *        definition.
         * @return A decoded Node instance.
         * @throws IOException When reading from the stream fails.
         */
        public static Node parse(DataInputStream is) throws IOException {
            int len = read_length(is);
            ByteTree[] bts = new ByteTree[len];
            for (int i = 0; i < len; i++) {
                bts[i] = ByteTree.parse(is);
            }
            Node n = new Node(bts);
            return n;
        }

        /**
         * @return "NODE".
         */
        @Override
        public String getPrefix() {
            return "NODE";
        }

        /**
         * Human-friendly representation of the Node object, with starting indentation.
         *
         * @param indent The indentation level
         * @return String representation of instance
         */
        @Override
        public String toString(int indent) {
            String spaces = new String(new char[indent]).replace("\0", " ");
            String header = String.format("%s%s %d", spaces, getPrefix(), getLength());
            String[] subs = new String[getLength() + 1];
            subs[0] = header;
            for (int i = 0; i < getLength(); i++) {
                subs[i + 1] = nodes[i].toString(indent + 1);
            }
            return String.join("\n", subs);
        }
    }

    /**
     * Leaf represents an abstract object.
     */
    public static class Leaf extends ByteTree {
        private final byte[] value;
        public static final byte PREFIX = 1;

        /**
         * Initialize a leaf from an abstract byte array.
         *
         * @param value A byte array
         */
        public Leaf(byte[] value) {
            this.value = value;
        }

        /**
         * Initialize a leaf from a String.
         *
         * @param value A string to initialize Leaf.
         */
        public Leaf(String value) {
            byte[] encoded = null;
            try {
                encoded = value.getBytes("US-ASCII");
            } catch (UnsupportedEncodingException e) {
                // this encoding is supported
            }
            this.value = encoded;
        }

        /**
         * Initialize a leaf from a BigInteger.
         *
         * @param value A BigInteger to initialize Leaf.
         */
        public Leaf(BigInteger value) {
            this.value = value.toByteArray();
        }

        /**
         * Initialize a leaf from a GroupElement.
         *
         * @param value A GroupElement to initialize Leaf.
         */
        public Leaf(GroupElement value) {
            if (value instanceof ModPGroupElement) {
                this.value = getEncoded((ModPGroupElement) value);
            } else if (value instanceof ECGroupElement) {
                this.value = getEncoded((ECGroupElement) value);
            } else if (value instanceof ProductGroupElement) {
                throw new IllegalArgumentException("Use Node for ProductGroupElement");
            } else {
                throw new IllegalArgumentException("Invalid group");
            }
        }

        /**
         * Get the byte array used to initialize the Leaf.
         *
         * @return A byte array.
         */
        public byte[] getValue() {
            return value;
        }

        /**
         * Return a String representation of the underlying byte array.
         *
         * @return String representing byte array.
         */
        public String getString() {
            return Util.toString(getValue());
        }

        /**
         * Return a BigInteger representation of the underlying byte array.
         *
         * @return BigInteger representing byte array.
         */
        public BigInteger getBigInteger() {
            return new BigInteger(value);
        }

        /**
         * Get the length of the underlying byte array.
         *
         * @return Length of byte array.
         */
        @Override
        public int getLength() {
            return value.length;
        }

        /**
         * Get the length of the whole Leaf object represented as byte array.
         *
         * @return Length of Leaf instance representation as byte array.
         */
        @Override
        public int getEncodedLength() {
            return getLength() + 5;
        }

        /**
         * Get the value with corresponding headers.
         *
         * @return A byte array.
         */
        @Override
        public byte[] getEncoded() {
            byte[] ret = new byte[getEncodedLength()];
            ret[0] = PREFIX;
            System.arraycopy(ByteTree.parse_from_int(getLength()), 0, ret, 1, 4);
            System.arraycopy(getValue(), 0, ret, 5, getLength());
            return ret;
        }

        /**
         * Encode the Leaf and write it to out.
         *
         * @param out Stream to write the ByteTree description of the Leaf.
         * @throws IOException When writing to out fails
         */
        @Override
        public void writeEncoded(OutputStream out) throws IOException {
            out.write(PREFIX);
            out.write(ByteTree.parse_from_int(getLength()));
            out.write(getValue());
        }

        private byte[] getEncoded(ModPGroupElement value) {
            // Verificatum verifier manual says that enough enough bytes are needed such that the
            // element fits in. In implementation, it uses BigInteger.toByteArray().length, which is
            // not always the same (it has a bit for magnitude).
            byte[] ret = new byte[((ModPGroup) value.getGroup()).getOrder().toByteArray().length];
            byte[] bvalue = value.getValue().toByteArray();
            System.arraycopy(bvalue, 0, ret, ret.length - bvalue.length, bvalue.length);
            return ret;
        }

        private byte[] getEncoded(ECGroupElement value) {
            throw new UnsupportedOperationException("ECGroupElement not supported currently");
        }

        /**
         * @return true
         */
        @Override
        public boolean isLeaf() {
            return true;
        }

        /**
         * Parse an array of bytes into a leaf.
         *
         * @param b The array of bytes to parse.
         * @param offset Starting offset to start parsing from.
         * @return The Leaf constructed from bytes.
         */
        public static Leaf parse(byte[] b, int offset) {
            int len = parse_int_fullbytes(b, offset);
            byte[] leafbytes = Arrays.copyOfRange(b, offset + 5, offset + 5 + len);
            return new Leaf(leafbytes);
        }

        /**
         * Parse a data stream into a leaf.
         *
         * Parses the given DataInputStream into the leaf. The given stream must be pointed to the
         * beginning of the definition (with length).
         *
         * @param is The given input stream, seeked to the beginning of the Leaf definition.
         * @return Leaf constructed from the input stream
         * @throws IOException When reading from the stream fails.
         */
        public static Leaf parse(DataInputStream is) throws IOException {
            int len = read_length(is);
            byte[] leafbytes = new byte[len];
            int read = is.read(leafbytes);
            if (read == -1) {
                throw new IOException("Unexpected end of stream");
            } else if (read != len) {
                throw new IOException("Short read");
            }
            return new Leaf(leafbytes);
        }

        /**
         * @return "LEAF"
         */
        @Override
        public String getPrefix() {
            return "LEAF";
        }

        /**
         * Return a human-friendly String representation of the Leaf with indentation.
         *
         * @param indent Indentation of the String.
         */
        @Override
        public String toString(int indent) {
            String spaces = new String(new char[indent]).replace("\0", " ");
            return String.format("%s%s %d %s", spaces, getPrefix(), getLength(),
                    HexFormat.of().formatHex(getValue()).toUpperCase());
        }
    }

    /**
     * Parse an array of bytes into ByteTree instance. Internally, either Node or Leaf is
     * constructed depending on the prefix.
     *
     * @param b Byte array to be parsed.
     * @param offset Starting offset of the byte array to start parsing from.
     * @return A ByteTree representing the byte array.
     */
    public static ByteTree parse(byte[] b, int offset) {
        if (offset < 0) {
            throw new IllegalArgumentException("Offset must be non-negatove");
        }
        if (b == null || b.length <= offset + 5) {
            throw new IllegalArgumentException("Non-existing bytetree");
        }

        if (b[offset] == Node.PREFIX) {
            return Node.parse(b, offset);
        } else if (b[offset] == Leaf.PREFIX) {
            return Leaf.parse(b, offset);
        }
        throw new IllegalArgumentException("Invalid bytetree value");
    }

    /**
     * Parse an input stream into ByteTree instance.
     *
     * Depending on the prefix, either Node or Leaf instance is constructed.
     *
     * @param is Given input stream to construct the ByteTree instance from, seeked to the
     *        beginning.
     * @return A ByteTree representing the byte array.
     * @throws IOException When reading from the stream fails or if the prefix is invalid.
     */
    public static ByteTree parse(DataInputStream is) throws IOException {
        switch (is.read()) {
            case -1:
                throw new IOException("Unexpected end of stream");
            case Node.PREFIX:
                return Node.parse(is);
            case Leaf.PREFIX:
                return Leaf.parse(is);
            default:
                throw new IOException("Unexpected token in input stream");
        }
    }

    /**
     * Short-hand method for {@link #parse(byte[], int)} with {@literal offset} 0.
     *
     * @param b Byte array to be parsed
     * @return A ByteTree representing the byte array.
     */
    public static ByteTree parse(byte[] b) {
        return parse(b, 0);
    }

    /**
     * Parse a file at a path into a ByteTree instance.
     *
     * Read a file from the given location into a ByteTree instance. In practice, depending on the
     * prefix, either Node or Leaf is constructed. This method is useful when the files are large
     * and do not fit into byte arrays.
     *
     * @param path Location of file
     * @return A ByteTree instance representing the byte array.
     * @throws IOException When the path is invalid, or the corresponding file is not valid ByteTree
     *         representation.
     */
    public static ByteTree parse(Path path) throws IOException {
        File file = path.toFile();
        FileInputStream fis = new FileInputStream(file);
        DataInputStream dis = new DataInputStream(fis);
        return parse(dis);
    }

    /**
     * Parse value from the byte array into integer.
     *
     * @param b Byte array to be parsed.
     * @param offset Starting offset to start parsing at.
     * @return Integer representation of the byte array.
     */
    private static int parse_int(byte[] b, int offset) {
        if (b.length < offset + 4) {
            throw new IllegalArgumentException("Bytetree int must be in four bytes");
        }
        return ByteBuffer.wrap(b, offset, 4).order(ByteOrder.BIG_ENDIAN).getInt();
    }

    private static int parse_int_fullbytes(byte[] b, int offset) {
        int len = parse_int(b, offset + 1);
        return len;
    }

    private static byte[] parse_from_int(int v) {
        return ByteBuffer.allocate(4).order(ByteOrder.BIG_ENDIAN).putInt(v).array();
    }

    private static int read_length(DataInputStream is) throws IOException {
        return is.readInt();
    }
}
