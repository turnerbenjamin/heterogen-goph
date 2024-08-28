package helpers

import (
	"io/fs"
	"os"
	"path/filepath"
)

type dirFs struct {
	root string
}

func GetDirFs(root string) dirFs {
	return dirFs{
		root: root,
	}
}

func GetFilesFromDir(dir fs.ReadDirFS) (files []string, err error) {
	if err := fs.WalkDir(dir, ".", func(path string, d fs.DirEntry, err error) error {
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

func (fs dirFs) ReadDir(name string) ([]fs.DirEntry, error) {
	return os.ReadDir(filepath.Join(fs.root, name))
}

func (m dirFs) Open(name string) (fs.File, error) {
	f, err := os.Open(filepath.Join(m.root, name))
	if err != nil {
		return nil, err
	}
	return fs.File(f), nil
}
