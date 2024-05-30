package ee.ivxv.common.service.bbox.impl;

import java.io.ByteArrayInputStream;
import java.io.InputStream;
import java.util.function.BiConsumer;
import java.util.function.Consumer;
import java.util.zip.ZipEntry;
import java.util.zip.ZipInputStream;

import static ee.ivxv.common.util.Util.CHARSET;

public class ZipSourceRaw implements FileSource {

    private final byte[] zipFileAsBytes;

    public ZipSourceRaw(byte[] zipFileAsBytes) {
        this.zipFileAsBytes = zipFileAsBytes;
    }

    private void processZipStream(BiConsumer<ZipEntry, InputStream> processor) {
        try (ByteArrayInputStream byteArrayInputStream = new ByteArrayInputStream(this.zipFileAsBytes);
             ZipInputStream zis = new ZipInputStream(byteArrayInputStream, CHARSET)) {
            for (ZipEntry ze; (ze = zis.getNextEntry()) != null; ) {
                if (ze.isDirectory()) {
                    continue;
                }
                processor.accept(ze, zis);
            }
        } catch (Exception e) {
            throw new RuntimeException(e);
        }
    }

    private void processZipFilesFromZipStream(Consumer<ZipEntry> processor) {
        try (ByteArrayInputStream byteArrayInputStream = new ByteArrayInputStream(this.zipFileAsBytes);
             ZipInputStream zis = new ZipInputStream(byteArrayInputStream, CHARSET)) {

            for (ZipEntry ze; (ze = zis.getNextEntry()) != null; ) {
                if (ze.isDirectory()) {
                    continue;
                }
                processor.accept(ze);
            }
        } catch (Exception e) {
            throw new RuntimeException(e);
        }
    }

    @Override
    public void processFiles(BiConsumer<String, InputStream> processor) {
        processZipStream((ze, in) -> processor.accept(ze.getName(), in));
    }

    @Override
    public void list(Consumer<String> processor) {
        processZipFilesFromZipStream(ze -> processor.accept(ze.getName()));
    }

    @Override
    public void listFileNamesAndSizes(BiConsumer<String, Long> processor) {
        // TODO: unimplemented
    }
}

