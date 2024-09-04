package staticAssets

import (
	"compress/gzip"
	"embed"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

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
		destinationPath := fmt.Sprintf("%s.gz", srcPath)

		if isCompressed(srcPath, destinationPath) {
			continue
		}

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
	fmt.Println("COMPRESSING", srcPath)

	src, err := os.Open(srcPath)
	if err != nil {
		log.Fatal(err)
	}
	defer src.Close()

	dotGz, err := os.Create(destinationPath)
	if err != nil {
		log.Fatal(err)
	}
	defer dotGz.Close()

	gzipWriter, err := gzip.NewWriterLevel(dotGz, gzip.BestCompression)
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

func isExcludedType(path string) bool {
	excludedTypes := []string{".go", ".webp", ".png", ".gz"}
	for _, excludedType := range excludedTypes {
		if strings.HasSuffix(path, excludedType) {
			return true
		}
	}
	return false
}
