package archive

import (
	"archive/tar"
	"io"
	"os"
	"path/filepath"
)

func Dir(dir string, wr io.Writer) error {
	tr := tar.NewWriter(wr)

	walkFn := func(path string, info os.FileInfo, err error) error {
		if info.Mode().IsDir() {
			return nil
		}
		// Because of scoping we can reference the external root_directory variable
		rel_path := path[len(dir):]
		if len(rel_path) == 0 {
			return nil
		}

		fr, err := os.Open(path)
		if err != nil {
			return err
		}
		defer fr.Close()

		h, err := tar.FileInfoHeader(info, rel_path)
		if err != nil {
			return err
		}

		h.Name = rel_path
		if err = tr.WriteHeader(h); err != nil {
			return err
		}

		if _, err := io.Copy(tr, fr); err != nil {
			return err
		}
		return nil
	}

	if err := filepath.Walk(dir, walkFn); err != nil {
		return err
	}

	return tr.Close()
}
