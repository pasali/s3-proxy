package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func main() {
	handler, err := ConfiguredProxyHandler()
	if err != nil {
		fmt.Printf("fatal: %v\n", err)
		return
	}
	http.Handle("/", handler)

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "healthy")
	})

	port := flag.Int("port", 8080, "Port to listen on")

	flag.Parse()

	portStr := strconv.FormatInt(int64(*port), 10)

	log.Println("s3-proxy is listening on port " + portStr)
	log.Fatal(http.ListenAndServe(":"+portStr, nil))
}
