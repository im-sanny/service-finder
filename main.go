package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type Service struct {
	ID       int
	Name     string
	Category string
}

type Provider struct {
	ID          int
	Name        string
	Phone       int
	Location    string
	Description string
}

var services = []Service{
	{ID: 1, Name: "Electrician", Category: "Electrical"},
	{ID: 2, Name: "Plumber", Category: "Plumbing"},
	{ID: 3, Name: "Painter", Category: "Painting"},
	{ID: 4, Name: "Mason", Category: "Construction"},
	{ID: 5, Name: "Carpenter", Category: "Woodworking"},
}

var providers = []Provider{
	{ID: 1, Name: "John Smith", Phone: 1234567890, Location: "New York, NY", Description: "Expert electrician with 10 years experience"},
	{ID: 2, Name: "Sarah Johnson", Phone: 2345678901, Location: "Los Angeles, CA", Description: "Licensed plumber, emergency services available"},
	{ID: 3, Name: "Mike Davis", Phone: 3456789012, Location: "Chicago, IL", Description: "Professional painter, interior and exterior"},
	{ID: 4, Name: "Emily Brown", Phone: 4567890123, Location: "Houston, TX", Description: "Experienced mason specializing in brick and stone work"},
	{ID: 5, Name: "David Wilson", Phone: 5678901234, Location: "Phoenix, AZ", Description: "Carpentry and custom woodworking"},
}

func GetServices(w http.ResponseWriter, r *http.Request) { // take req and write response
	w.Header().Set("Content-Type", "application/json") // this will set header and content type as json
	if r.Method != http.MethodGet {                    // if the req method does't match then send error
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// encode services into json and give that as response, data isn't coming in json that's why it needs to be encoded
	if err := json.NewEncoder(w).Encode(services); err != nil { // if encoding/writing fails this will send a error
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func ServiceId(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// first get id, it'll come as string so need to convert it into int, then run loop in-memory db to get matching id, set header, content type encode the response and show the result
	idStr := r.PathValue("id")     // collect string id from request
	id, err := strconv.Atoi(idStr) // convert string id to int
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	for _, v := range services { // run loop to check each service
		if id == v.ID { // if requested ID matches this service ID then
			w.Header().Set("Content-Type", "application/json") // set response content type to json
			json.NewEncoder(w).Encode(v)                       // then encode the response and return the response
			return
		}
	}
}

func main() {
	mux := http.NewServeMux()

	// service routes
	mux.HandleFunc("/service", GetServices)
	mux.HandleFunc("/service/{id}", ServiceId)

	fmt.Printf("Server running on port :8080")
	http.ListenAndServe(":8080", mux)
}
