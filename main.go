package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

var rePart = regexp.MustCompile(`\d+|\D+`)

func handlerIndex(w http.ResponseWriter, r *http.Request) {
	var message, css string
	path, _ := os.Getwd()
	files, _ := readdir(path)
	if len(files) == 0 {
		os.Exit(1)
	}
	css = `<style>
	body{
		background-color: black;
	}
	img{
		width: 80%;
		display: block;
		margin: 30px auto;
	}
	</style>`
	message += "<html><head>"
	message += "<meta name='viewport' content='width=device-width, minimum-scale=0.1'>"
	message += css
	message += "<title>" + path + "</title></head><body>"

	sort.Sort(byName(files))
	for _, file := range files {
		message += "<img loading=lazy src='img/" + file.Name() + "'>"
	}

	message += "</body></html>"
	w.Write([]byte(message))
}

func main() {
	path, _ := os.Getwd()
	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir(path))))
	http.HandleFunc("/", handlerIndex)
	var address = "localhost:9000"
	fmt.Printf("server started at %s\n", address)
	err := http.ListenAndServe(address, nil)
	if err != nil {
		fmt.Println(err.Error())
	}
}

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
		fmt.Println(ext)
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
