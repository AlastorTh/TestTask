package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"sync"
)

type Server struct {
	http.Server
	mut    sync.Mutex
	queues map[string][]string
}

func (s *Server) handleRequest(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Starting handler\n")
	switch r.Method {
	case http.MethodGet:
		io.WriteString(w, "Is type GET")
		key := r.URL.Path[1:]
		s.mut.Lock()
		_, ok := s.queues[key]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if len(s.queues[key]) == 0 {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		io.WriteString(w, s.queues[key][0])
		if len(s.queues[key]) == 1 {
			s.queues[key] = []string{}
		} else {
			s.queues[key] = s.queues[key][1:]
		}
		fmt.Println(s.queues)
		s.mut.Unlock()
	case http.MethodPut:

		key := r.URL.Path[1:]

		val, ok := r.URL.Query()["v"]
		if !ok {
			io.WriteString(w, "nothing")
		}
		fmt.Println(val)

		s.mut.Lock()
		s.queues[key] = append(s.queues[key], val[0])
		fmt.Println(s.queues)
		s.mut.Unlock()
		io.WriteString(w, fmt.Sprintf("queue %s value %s input successful", key, val[0]))
	}
}

func main() {
	port := flag.Int("port", 8080, "port to start the server on")
	flag.Parse()
	if *port < 0 || *port > 65536 {
		fmt.Println("Port cmd flag is out of range!")
		fmt.Println("Designating a default port 8080")
		*port = 8080
	}
	mux := http.NewServeMux()
	s := Server{http.Server{Addr: fmt.Sprintf(":%d", *port), Handler: mux}, sync.Mutex{}, make(map[string][]string)}
	mux.HandleFunc("/", s.handleRequest)
	//err := http.ListenAndServe(*port, mux)

	fmt.Println(s.ListenAndServe())

}
