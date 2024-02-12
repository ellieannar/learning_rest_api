/*
	Create a go program that handles the following
	/foo, handles GET and POST.
		For GET, return a json struct that looks like this "{'foo': 'Get'}",
		for POST, return this string "foo POST", no json
	/bar handles GET and DELETE.
		For GET, return the current timestamp,
		for DELETE, require that the user passes a json struct that looks like {"name": insert name here}, then returns "Hello" plus the name provided.
*/

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

func fooHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		response := map[string]string{"foo": "Get"}
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

func main() {
	http.HandleFunc("/foo", fooHandler)
	http.HandleFunc("/bar", barHandler)

	http.ListenAndServe(":3333", nil)
}
