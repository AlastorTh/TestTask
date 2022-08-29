package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
)

func getRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got / request\n")
	io.WriteString(w, "This is my website!\n")
}

func main() {
	port := flag.String("port", ":8080", "port to start the server on")
	flag.Parse()

	//s := &http.Server{Addr: *port, ReadTimeout: 10 * time.Second, WriteTimeout: 10 * time.Second}
	mux := http.NewServeMux()
	mux.HandleFunc("/", getRoot)
	err := http.ListenAndServe(*port, mux)
	fmt.Println(err)
}
