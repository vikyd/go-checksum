package checksum

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// GoDirHash is the result of hash
type GoDirHash struct {
	HashSynthesized       string
	HashSynthesizedBase64 string
	GoCheckSum            string
}

// HashDir hash the specify dir
func HashDir(dir, prefix string) (GoDirHash, error) {
	files, err := DirFiles(dir, prefix)
	if err != nil {
		return GoDirHash{}, err
	}
	osOpen := func(name string) (io.ReadCloser, error) {
		return os.Open(filepath.Join(dir, strings.TrimPrefix(name, prefix)))
	}
	return Hash1(files, osOpen)
}

// Hash1 `1` means SHA256
func Hash1(files []string, open func(string) (io.ReadCloser, error)) (GoDirHash, error) {
	h := sha256.New()
	files = append([]string(nil), files...)
	sort.Strings(files)
	for _, file := range files {
		if strings.Contains(file, "\n") {
			return GoDirHash{}, errors.New("filenames with newlines are not supported")
		}
		r, err := open(file)
		if err != nil {
			return GoDirHash{}, err
		}
		hf := sha256.New()
		_, err = io.Copy(hf, r)
		r.Close()
		if err != nil {
			return GoDirHash{}, err
		}
		ss := hf.Sum(nil)
		// fmt.Println(fmt.Sprintf("%x\n", ss))
		fmt.Fprintf(h, "%x  %s\n", ss, file)
	}
	bArr := h.Sum(nil)
	hashSynthesized := fmt.Sprintf("%x", bArr)
	hashSynthesizedBase64 := Base64(bArr)

	allHash := GoDirHash{
		HashSynthesized:       hashSynthesized,
		HashSynthesizedBase64: hashSynthesizedBase64,
		GoCheckSum:            "h1:" + hashSynthesizedBase64,
	}

	return allHash, nil
}

// DirFiles loop all files, and join prefix
func DirFiles(dir, prefix string) ([]string, error) {
	var files []string
	dir = filepath.Clean(dir)
	err := filepath.Walk(dir, func(file string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			// skip .git dir
			if info.Name() == ".git" {
				return filepath.SkipDir
			}
			return nil
		}
		rel := file
		if dir != "." {
			rel = file[len(dir)+1:]
		}
		f := filepath.Join(prefix, rel)
		files = append(files, filepath.ToSlash(f))
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}
