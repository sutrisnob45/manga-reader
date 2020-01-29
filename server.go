package main

import (
	"net/http"
	"os"
	"sort"
)

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
		max-width: 80%;
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
