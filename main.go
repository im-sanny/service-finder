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

func ServiceGet(w http.ResponseWriter, r *http.Request) { // take req and write response
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

func ServicePost(w http.ResponseWriter, r *http.Request) {
	// client send new service data as JSON in the request body
	w.Header().Set("Content-Type", "application/json") // response will be set to json while taking data from output
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var s Service                                              // this will create an empty service variable
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil { // decode the JSON request body into the service variable body
		http.Error(w, "Failed to decode", http.StatusBadRequest)
		return
	}
	newId := len(services) + 1     // creating new id
	s.ID = newId                   // assigning new id to new service
	services = append(services, s) // appending new service to the in-memory services slice

	w.WriteHeader(http.StatusCreated) // set the success status 201 created
	json.NewEncoder(w).Encode(s)      // encode the created service as JSON and send it back
}

func ServiceUpdate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPut {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}

	var s Service
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		http.Error(w, "Failed to decode", http.StatusBadRequest)
		return
	}

	for i := range services {
		if services[i].ID == id { // if the id we got is == to services slice id then
			s.ID = id                    // empty service variable id = id that we got through request
			services[i] = s              // replace the existing service at this index with the new data stored in s.
			json.NewEncoder(w).Encode(s) // encode and return the data that service variable got
			return
		}
	}
	http.Error(w, "Service not found", http.StatusNotFound)
}

func ServicePatch(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPatch {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}

	var s Service
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		http.Error(w, "Failed to decode", http.StatusBadRequest)
		return
	}

	for i := range services { // loop over range of services for index of services
		if services[i].ID == id { // if services index id == id then
			if s.Name != "" { // if s.name isn't empty
				services[i].Name = s.Name // then update the existing service's Name with the new Name.
			}
			if s.Category != "" {
				services[i].Category = s.Category
			}
			json.NewEncoder(w).Encode(services[i]) // encode whatever service got by index then return it
			return
		}
	}
	http.Error(w, "Service not found", http.StatusNotFound)
}

func ServiceDelete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}

	// if the id matched with index id then exclude that from the slice
	for i := range services {
		if services[i].ID == id {
			services = append(services[:i], services[i+1:]...) // services[:i] means everything before index i, services[i+1:] means everything after index i
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	http.Error(w, "Service not found", http.StatusNotFound)
}

func main() {
	mux := http.NewServeMux()

	// service routes
	mux.HandleFunc("GET /service", ServiceGet)
	mux.HandleFunc("GET /service/{id}", ServiceId)
	mux.HandleFunc("POST /service", ServicePost)
	mux.HandleFunc("PUT /service/{id}", ServiceUpdate)
	mux.HandleFunc("PATCH /service/{id}", ServicePatch)
	mux.HandleFunc("DELETE /service/{id}", ServiceDelete)

	fmt.Printf("Server running on port :8080")
	http.ListenAndServe(":8080", mux)
}
