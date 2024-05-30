package ee.ivxv.common.service.bbox.impl;

import static ee.ivxv.common.util.Util.CHARSET;

import ee.ivxv.common.service.bbox.InvalidBboxException;
import java.io.InputStream;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.List;
import java.util.function.BiConsumer;
import java.util.function.Consumer;
import java.util.zip.ZipEntry;
import java.util.zip.ZipFile;
import java.util.zip.ZipInputStream;

public class ZipSource implements FileSource {

    private final Path path;

    public ZipSource(Path path) {
        this.path = path;
    }

    private static void processZippedStream(Path path,
            BiConsumer<ZipEntry, InputStream> processor) {
        try (ZipInputStream zis = new ZipInputStream(Files.newInputStream(path), CHARSET)) {
            for (ZipEntry ze; (ze = zis.getNextEntry()) != null;) {
                if (ze.isDirectory()) {
                    continue;
                }
                processor.accept(ze, zis);
            }
        } catch (Exception e) {
            throw new RuntimeException(new InvalidBboxException(path, e));
        }
    }

    private static void processZipFile(Path path, Consumer<ZipEntry> processor) {
        try (ZipFile zip = openZipFile(path)) {
            zip.stream().forEach(ze -> {
                if (ze.isDirectory()) {
                    return;
                }
                processor.accept(ze);
            });
        } catch (Exception e) {
            throw new InvalidBboxException(path, e);
        }
    }

    private static ZipFile openZipFile(Path path) {
        try {
            return new ZipFile(path.toFile());
        } catch (Exception e) {
            throw new InvalidBboxException(path, e);
        }
    }

    @Override
    public void processFiles(BiConsumer<String, InputStream> processor) {
        processZippedStream(path, (ze, in) -> {
            processor.accept(ze.getName(), in);
        });
    }

    @Override
    public void list(Consumer<String> processor) {
        processZipFile(path, ze -> {
            processor.accept(ze.getName());
        });
    }

    @Override
    public void listFileNamesAndSizes(BiConsumer<String, Long> processor) {
        processZipFile(path, ze -> processor.accept(ze.getName(), ze.getSize()));
    }

    @Override
    public int countFiles() {
        try (ZipFile zipFile = new ZipFile(this.path.toFile())){
            return zipFile.size();
        } catch (Exception e) {
            throw new InvalidBboxException(this.path, e);
        }
    }

    @Override
    public int countFilesWithSuffix(List<String> suffixes) {
        try (ZipFile zipFile = new ZipFile(this.path.toFile())){
            List<String> fileContent = zipFile.stream().filter(ze -> {
                String fileName = ze.getName();
                for (String suffix: suffixes) if (fileName.endsWith(suffix)) return true;
                return false;
            }).map(ZipEntry::getName).toList();
            return fileContent.size();
        } catch (Exception e) {
            throw new InvalidBboxException(this.path, e);
        }
    }
}
