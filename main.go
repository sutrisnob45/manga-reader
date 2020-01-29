package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	path, _ := os.Getwd()
	http.HandleFunc("/", handlerIndex)
	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir(path))))
	var address = "localhost:9000"
	fmt.Printf("Server started at %s\n", address)
	fmt.Printf("Open your browser with address above to read your manga")
	err := http.ListenAndServe(address, nil)
	if err != nil {
		fmt.Println(err.Error())
	}
}
