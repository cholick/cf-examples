package main

import (
	"net/http"
	"os"
	"log"
	"fmt"
)

func main() {
	http.HandleFunc("/", index)

	port := getPort()

	fmt.Println("Listening on port [" + port + "]")
	err := http.ListenAndServe(":" + port, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func index(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(res, "Success: Go")
}

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	return port
}
