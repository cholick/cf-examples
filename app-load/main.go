package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"
)

type config struct {
	ListenPort  string
	Host        string
	RequestPort string
	Delay       time.Duration
}

func handler(w http.ResponseWriter, r *http.Request) {
	behavior := rand.Intn(100)
	if behavior < 90 {
		println("Returning 200")
		w.WriteHeader(200)
		fmt.Fprintf(w, "Success\n")
	} else if behavior < 95 {
		println("Returning 404")
		w.WriteHeader(404)
	} else if behavior < 97 {
		println("Returning 400")
		w.WriteHeader(400)
		fmt.Fprintf(w, "Bad request\n")
	} else {
		println("Returning 403")
		w.WriteHeader(403)
		fmt.Fprintf(w, "Not authorized\n")
	}
}

func periodic(config config) {
	uri := fmt.Sprintf("http://%v:%v", config.Host, config.RequestPort)
	println("Making request to " + uri)
	_, err := http.Get(uri)
	if err != nil {
		println(err.Error())
	}
	time.Sleep(config.Delay)
	periodic(config)
}

func parseConfig() config {
	vcap_raw := os.Getenv("VCAP_APPLICATION")
	var vcap map[string]interface{}
	json.Unmarshal([]byte(vcap_raw), &vcap)
	println(fmt.Sprintf("Unmarshalled %v", vcap))

	uris := vcap["application_uris"].([]interface{})
	requestPort := os.Getenv("REQUEST_PORT")
	if requestPort == "" {
		requestPort = "80"
	}

	delay, err := time.ParseDuration(os.Getenv("DELAY"))
	if err != nil {
		println(err.Error())
		delay = 250 * time.Millisecond
	}

	return config{
		ListenPort:  os.Getenv("PORT"),
		Host:        uris[0].(string),
		RequestPort: requestPort,
		Delay:       delay,
	}
}

func main() {
	config := parseConfig()
	go func() {
		time.Sleep(5 * time.Second)
		periodic(config)
	}()

	println(fmt.Sprintf("Listening on %v", config.ListenPort))
	http.HandleFunc("/", handler)
	panic(http.ListenAndServe(":"+config.ListenPort, nil))
}
