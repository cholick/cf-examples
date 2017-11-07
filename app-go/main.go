package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
)

type panicHandler struct {
	http.Handler
}

func (h panicHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer handleFatalError()
	h.Handler.ServeHTTP(w, r)
}

func handleFatalError() {
	e := recover()
	if e == nil {
		return
	}

	stack := make([]byte, 1<<16)
	stackSize := runtime.Stack(stack, true)
	prettyStack := string(stack[:stackSize])

	fmt.Fprintf(os.Stderr, "panic: %v", e)
	fmt.Fprintf(os.Stderr, prettyStack)
	os.Exit(1)
}

func index(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(res, "Success: Go.")
}

func notFound(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(404)
}

func crashHandler(res http.ResponseWriter, req *http.Request) {
	panic("oh no!")
}

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	return port
}

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/favicon.ico", notFound)
	http.Handle("/crash", panicHandler{http.HandlerFunc(crashHandler)})

	port := getPort()

	fmt.Println("Listening on port [" + port + "]")
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
