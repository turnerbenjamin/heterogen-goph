package staticAssets

import (
	"compress/gzip"
	"embed"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"strings"
)

//go:embed *
var FileSystem embed.FS

func CompressFiles() {
	files, err := getAllFilenames(&FileSystem)
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

	gzipWriter := gzip.NewWriter(dotGz)
	defer gzipWriter.Close()

	// Copy the contents of the original file to the gzip writer
	_, err = io.Copy(gzipWriter, src)
	if err != nil {
		log.Fatal(err)
	}

	gzipWriter.Flush()
}

func getAllFilenames(efs *embed.FS) (files []string, err error) {
	if err := fs.WalkDir(efs, ".", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		files = append(files, path)
		return nil
	}); err != nil {
		return nil, err
	}

	return files, nil
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
