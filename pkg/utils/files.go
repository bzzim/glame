package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
)

// CopyFile copies a file from src to dst. If src and dst files exist, and are
// the same, then return success. Copy the file contents from src to dst.
func CopyFile(src, dst string) (err error) {
	sfi, err := os.Stat(src)
	if err != nil {
		return
	}
	if !sfi.Mode().IsRegular() {
		// cannot copy non-regular files (e.g., directories,
		// symlinks, devices, etc.)
		return fmt.Errorf("CopyFile: non-regular source file %s (%q)", sfi.Name(), sfi.Mode().String())
	}
	dfi, err := os.Stat(dst)
	if err != nil {
		if !os.IsNotExist(err) {
			return
		}
	} else {
		if !(dfi.Mode().IsRegular()) {
			return fmt.Errorf("CopyFile: non-regular destination file %s (%q)", dfi.Name(), dfi.Mode().String())
		}
		if os.SameFile(sfi, dfi) {
			return
		}
	}

	err = copyFileContents(src, dst)
	return
}

func MakeDir(path string) error {
	err := os.Mkdir(path, os.ModeDir|os.ModePerm)
	if err == nil {
		return nil
	}
	if os.IsExist(err) {
		// check that the existing path is a directory
		info, err := os.Stat(path)
		if err != nil {
			return err
		}
		if !info.IsDir() {
			return errors.New("path exists but is not a directory")
		}
		return nil
	}
	return err
}

// copyFileContents copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file.
func copyFileContents(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}

func WriteToFileIfNotExists(content, dst string) (err error) {
	if _, err := os.Stat(dst); !errors.Is(err, os.ErrNotExist) {
		return nil
	}
	f, err := os.Create(dst)

	if err != nil {
		return err
	}

	defer f.Close()

	_, err = f.WriteString(content)
	return err
}

func ReadFile(dst string) (string, error) {
	content, err := os.ReadFile(dst)
	return string(content), err
}

func WriteToFile(content, dst string) error {
	f, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(content)
	return err
}

func ReadJson[K any](filename string) (*K, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var queries K
	err = json.Unmarshal(content, &queries)
	if err != nil {
		return nil, err
	}

	return &queries, nil
}

func SaveJson[K any](filename string, j *K) error {
	data, err := json.Marshal(j)
	if err != nil {
		return err
	}
	err = WriteToFile(string(data), filename)
	return err
}
