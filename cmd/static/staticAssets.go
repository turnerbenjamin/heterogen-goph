package staticAssets

import (
	"compress/gzip"
	"embed"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/andybalholm/brotli"
	"github.com/turnerbenjamin/heterogen-go/internal/helpers"
)

//go:embed *
var FileSystem embed.FS

func CompressFiles() {
	files, err := helpers.GetFilesFromDir(&FileSystem)
	if err != nil {
		log.Fatal(err)
	}

	for _, path := range files {
		if isExcludedType(path) {
			continue
		}

		srcPath := fmt.Sprintf("cmd/static/%s", path)

		//Compress GZip
		destinationPath := fmt.Sprintf("%s.gz", srcPath)
		compressFile(srcPath, destinationPath)

		//Compress Brotli
		destinationPath = fmt.Sprintf("%s.br", srcPath)
		compressFile(srcPath, destinationPath)

	}
}

func isCompressed(srcPath string, destinationPath string) bool {

	srcInfo, err := os.Stat(srcPath)
	if err != nil {
		log.Fatal(err)
	}

	destinationInfo, err := os.Stat(destinationPath)
	if err != nil {
		return false
	}

	return srcInfo.ModTime().Before(destinationInfo.ModTime())
}

func compressFile(srcPath string, destinationPath string) {
	if isCompressed(srcPath, destinationPath) {
		return
	}

	src, err := os.Open(srcPath)
	if err != nil {
		log.Fatal(err)
	}
	defer src.Close()

	destination, err := os.Create(destinationPath)
	if err != nil {
		log.Fatal(err)
	}
	defer destination.Close()

	if strings.HasSuffix(destinationPath, ".gz") {
		compressGzip(src, destination)
	}

	if strings.HasSuffix(destinationPath, ".br") {
		compressBrotli(src, destination)
	}

}

func compressGzip(src *os.File, dWriter io.Writer) {

	gzipWriter, err := gzip.NewWriterLevel(dWriter, gzip.BestCompression)
	if err != nil {
		log.Fatal(err)
	}
	defer gzipWriter.Close()

	// Copy the contents of the original file to the gzip writer
	_, err = io.Copy(gzipWriter, src)
	if err != nil {
		log.Fatal(err)
	}

	gzipWriter.Flush()
}

func compressBrotli(src *os.File, dWriter io.Writer) {

	brotliWriter := brotli.NewWriterLevel(dWriter, brotli.BestCompression)
	defer brotliWriter.Close()

	// Copy the contents of the original file to the gzip writer
	_, err := io.Copy(brotliWriter, src)
	if err != nil {
		log.Fatal(err)
	}

	brotliWriter.Flush()
}

func isExcludedType(path string) bool {
	excludedTypes := []string{".go", ".webp", ".png", ".gz", ".br"}
	for _, excludedType := range excludedTypes {
		if strings.HasSuffix(path, excludedType) {
			return true
		}
	}
	return false
}
