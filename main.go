package main

import (
	"encoding/json"
	"fmt"
	"net/http"
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

func main() {
	mux := http.NewServeMux()

	// service routes
	mux.HandleFunc("/service", GetServices)

	fmt.Printf("Server running on port :8080")
	http.ListenAndServe(":8080", mux)
}
