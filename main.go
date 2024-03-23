/*
	Create a go program that handles the following
	/foo, handles GET and POST.
		For GET, return a json struct that looks like this "{'foo': 'Get'}",
		for POST, return this string "foo POST", no json
	/bar handles GET and DELETE.
		For GET, return the current timestamp,
		for DELETE, require that the user passes a json struct that looks like {"name": insert name here}, then returns "Hello" plus the name provided.

	Add 2 metrics to your rest api.
		A metric named "Up", while the program is running, it randomly either pushes a 0 or a 1 every 60 seconds
		A metric that counts every time a GET is called on the /foo endpoint
*/

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var mu sync.Mutex

var up = prometheus.NewGauge(
	prometheus.GaugeOpts{
		Name: "up_time",
		Help: "Pushes 1 or 0 every 60 seconds while running",
	},
)

var getCount = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "get_request_count",
		Help: "Tracks the total number of times GET is called on the /foo endpoint",
	},
)

func fooHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		response := map[string]string{"foo": "Get"}
		getCount.Inc()
		json.NewEncoder(w).Encode(response)

	case http.MethodPost:
		io.WriteString(w, "foo POST")

	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}

}
func barHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		io.WriteString(w, time.Now().String())

	case http.MethodDelete:
		var requestData map[string]string
		err := json.NewDecoder(r.Body).Decode(&requestData)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		name, exists := requestData["name"]
		if !exists {
			http.Error(w, "Name not provided", http.StatusBadRequest)
			return
		}
		fmt.Fprint(w, "Hello, ", name)

	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func UpFunction() {
	go func() {
		for {
			mu.Lock()
			up.Set(float64(rand.Intn(2))) // Set up to be either 0 or 1 randomly
			mu.Unlock()
			time.Sleep(60 * time.Second)
		}
	}()

}

func main() {
	prometheus.Register(getCount)
	prometheus.Register(up)

	UpFunction()

	http.HandleFunc("/foo", fooHandler)
	http.HandleFunc("/bar", barHandler)
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":3333", nil)
}
