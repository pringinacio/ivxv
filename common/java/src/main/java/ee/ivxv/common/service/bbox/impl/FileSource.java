package ee.ivxv.common.service.bbox.impl;

import java.io.InputStream;
import java.util.List;
import java.util.function.BiConsumer;
import java.util.function.Consumer;

/**
 * Simple file system abstraction.
 */
public interface FileSource {

    /**
     * Iterates over all files and calls the processor on them.
     *
     * @param processor
     */
    void processFiles(BiConsumer<String, InputStream> processor);

    /**
     * Lists all file names and calls the processor on them.
     *
     * @param processor
     */
    default void list(Consumer<String> processor) {
        processFiles((name, in) -> {
            processor.accept(name);
        });
    }

    /**
     * Lists all file names and their sizes and then calls the processor on them.
     * @param processor
     */
    void listFileNamesAndSizes(BiConsumer<String, Long> processor);

    /**
     * Count all found files.
     *
     * @return amount of found files.
     */
    default int countFiles() {
        return 0;
    }

    /**
     * Count files matching some specific suffix(es).
     *
     * @param suffixes the suffixes to match for
     * @return the amount of found files
     */
    default int countFilesWithSuffix(List<String> suffixes) { return 0; }
}
