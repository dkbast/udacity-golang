package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Customer struct (Model)
/*
ID
Name
Role
Email
Phone
Contacted (i.e., indication of whether or not the customer has been contacted)
*/
type Customer struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Role      string `json:"role"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Contacted bool   `json:"contacted"`
}

var database = map[string]Customer{
	"1": {ID: "1", Name: "John Doe", Role: "CEO", Email: "jodo@example.com", Phone: "1234567890", Contacted: false},
	"2": {ID: "2", Name: "Jane Doe", Role: "CTO", Email: "jado@example.com", Phone: "1234567890", Contacted: false},
	"3": {ID: "3", Name: "John Smith", Role: "CFO", Email: "josm@example.com", Phone: "1234567890", Contacted: false},
	"4": {ID: "4", Name: "Jane Smith", Role: "CMO", Email: "jasm@example.com", Phone: "1234567890", Contacted: false},
}

// getCustomers returns all customers
func getCustomers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// all the values from the database map:
	var customers []Customer
	for _, customer := range database {
		customers = append(customers, customer)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(customers)
}

// getCustomer returns a single customer
func getCustomer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	// find the customer in the database map and return it, if none is found return an empty json and a 404 status code
	customer, ok := database[params["id"]]
	if ok {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(customer)
		return
	}
	http.Error(w, "Customer not found", http.StatusNotFound)

}

// addCustomer adds a new customer
func addCustomer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var customer Customer
	_ = json.NewDecoder(r.Body).Decode(&customer)
	var tempCustomerId = strconv.Itoa(len(database) + 1)
	if _, ok := database[tempCustomerId]; ok {
		customer.ID = tempCustomerId
	} else {
		// if the id is already in the database, generate a new one
		for {
			tempCustomerId = strconv.Itoa(len(database) + 1)
			if _, ok := database[tempCustomerId]; ok {
				continue
			} else {
				customer.ID = tempCustomerId
				break
			}
		}
	}

	database[tempCustomerId] = customer
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(customer)
}

// updateCustomer updates an existing customer, if the customer does not exist it returns an empty json and a 404 status code
func updateCustomer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) // Get params
	var updatedCustomer Customer
	_ = json.NewDecoder(r.Body).Decode(&updatedCustomer)
	// find the customer in the database map and return it, if none is found return an empty json and a 404 status code
	_, ok := database[params["id"]]
	if ok {
		database[params["id"]] = updatedCustomer
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(updatedCustomer)
		return
	}
	http.Error(w, "Customer not found", http.StatusNotFound)
}

// deleteCustomer deletes a customer, if the customer does not exist it returns an empty json and a 404 status code
func deleteCustomer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) // Get params
	// find the customer in the database map and return it, if none is found return an empty json and a 404 status code
	_, ok := database[params["id"]]
	if ok {
		delete(database, params["id"])
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(database)
		return
	}
	http.Error(w, "Customer not found", http.StatusNotFound)

}

func main() {
	fmt.Println("Hello Udacity!")
	router := mux.NewRouter()
	router.HandleFunc("/customers", getCustomers).Methods("GET")
	router.HandleFunc("/customers/{id}", getCustomer).Methods("GET")
	router.HandleFunc("/customers", addCustomer).Methods("POST")
	router.HandleFunc("/customers/{id}", updateCustomer).Methods("PUT")
	router.HandleFunc("/customers/{id}", deleteCustomer).Methods("DELETE")
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))
	fmt.Println("Server is running on port 8000")
	log.Fatal(http.ListenAndServe(":8000", router))
}
