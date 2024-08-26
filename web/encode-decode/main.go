package main

import "net/http"

func main() {
	http.HandleFunc("/encode", BuyCar)
	http.HandleFunc("/decode", SoldCarAndBuyLaptop)
	http.ListenAndServe(":80", nil)
}
