package main

import (
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var rePart = regexp.MustCompile(`\d+|\D+`)

type file struct {
	os.FileInfo
	path string
	ext  string
}

func readdir(path string) ([]*file, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	names, err := f.Readdirnames(-1)
	f.Close()

	files := make([]*file, 0, len(names))
	for _, fname := range names {
		fpath := filepath.Join(path, fname)

		lstat, err := os.Lstat(fpath)
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return files, err
		}
		if lstat.IsDir() {
			continue
		}

		ext := filepath.Ext(fpath)
		switch ext[1:] {
		case "bmp", "jpg", "jpeg", "gif", "png":
			files = append(files, &file{
				FileInfo: lstat,
				path:     fpath,
				ext:      ext,
			})
		}

	}

	return files, err
}

type byName []*file

func (a byName) Len() int      { return len(a) }
func (a byName) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a byName) Less(i, j int) bool {
	return naturalSort(strings.ToLower(a[i].Name()), strings.ToLower(a[j].Name()))
}

func naturalSort(s1, s2 string) bool {
	parts1 := rePart.FindAllString(s1, -1)
	parts2 := rePart.FindAllString(s2, -1)

	for i := 0; i < len(parts1) && i < len(parts2); i++ {
		if parts1[i] == parts2[i] {
			continue
		}

		num1, err1 := strconv.Atoi(parts1[i])
		num2, err2 := strconv.Atoi(parts2[i])

		if err1 == nil && err2 == nil {
			return num1 < num2
		}

		return parts1[i] < parts2[i]
	}

	return len(parts1) < len(parts2)
}
