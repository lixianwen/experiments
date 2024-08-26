// Tutors:
//     https://gowebexamples.com/
//     https://www.soberkoder.com/go-rest-api-gorilla-mux/

package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Car struct {
	// Is it the best?
	Model    string `json:"model"`
	Color    string `json:"car_color"`
	MaxSpeed int    `json:"max_speed"`
}

type Laptop struct {
	Chip  string  `json:"chip"`
	Color string  `json:"color"`
	Inch  float64 `json:"inch"`
}

func BuyCar(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.Path, r.Method)

	dream := Car{
		Model:    "Tesla Model 3 Standard Edition",
		Color:    "white",
		MaxSpeed: 200,
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(dream); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func SoldCarAndBuyLaptop(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.Path, r.Method)

	var car Car
	if err := json.NewDecoder(r.Body).Decode(&car); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	mba := Laptop{
		Chip:  "Apple M4",
		Color: "silver",
		Inch:  13.6,
	}
	if err := json.NewEncoder(w).Encode(mba); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
