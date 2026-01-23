package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type Inventory struct {
	ID          int    `json:"id"`
	Nama        string `json:"nama"`
	Description string `json:"description"`
}

var inventory = []Inventory{
	{ID: 1, Nama: "monitor", Description: "Asus"},
	{ID: 2, Nama: "mouse", Description: "Rexus"},
	{ID: 3, Nama: "komputer", Description: "i7 Gen 12"},
}

func getInventoryByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/inventory/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Inventory ID", http.StatusBadRequest)
		return
	}

	for _, p := range inventory {
		if p.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(p)
			return
		}
	}

	http.Error(w, "Inventory belum ada", http.StatusNotFound)
}

// PUT localhost:8080/api/inventory/{id}
func updateInventory(w http.ResponseWriter, r *http.Request) {

	// get id dari request
	idStr := strings.TrimPrefix(r.URL.Path, "/api/inventory/")

	// ganti int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Inventory ID", http.StatusBadRequest)
		return
	}

	// get data dari request
	var updateInventory Inventory
	err = json.NewDecoder(r.Body).Decode(&updateInventory)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// loop inventory, cari id, ganti sesuai data dari request
	for i := range inventory {
		if inventory[i].ID == id {
			updateInventory.ID = id
			inventory[i] = updateInventory

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(updateInventory)
			return
		}
	}
	http.Error(w, "Inventory belum ada", http.StatusNotFound)
}

func deleteInventory(w http.ResponseWriter, r *http.Request) {
	// get id
	idStr := strings.TrimPrefix(r.URL.Path, "/api/inventory/")
	// ganti id int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Inventory ID", http.StatusBadRequest)
		return
	}
	// loop inventory cari ID, dapet index yang mau dihapus
	for i, p := range inventory {
		if p.ID == id {
			// bikin slice baru dengan data sebelum dan sesudah index
			inventory = append(inventory[:i], inventory[i+1:]...)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"message": "sukses delete",
			})

			return
		}
	}

	http.Error(w, "Inventory belum ada", http.StatusNotFound)
}

func main() {

	// GET localhost:8080/api/inventory/{id}
	// PUT localhost:8080/api/inventory/{id}
	// DELETE localhost:8080/api/inventory/{id}
	http.HandleFunc("/api/inventory/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			getInventoryByID(w, r)
		case "PUT":
			updateInventory(w, r)
		case "DELETE":
			deleteInventory(w, r)
		}

	})

	// GET localhost:8080/api/inventory
	// POST localhost:8080/api/inventory
	http.HandleFunc("/api/inventory", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(inventory)
		case "POST":
			// baca data dari request
			var Inventorybaru Inventory
			err := json.NewDecoder(r.Body).Decode(&Inventorybaru)
			if err != nil {
				http.Error(w, "Invalid request", http.StatusBadRequest)
				return
			}

			// masukkin data ke dalam variable inventory
			Inventorybaru.ID = len(inventory) + 1
			inventory = append(inventory, Inventorybaru)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated) // 201
			json.NewEncoder(w).Encode(Inventorybaru)
		}

	})

	// localhost:8080/health
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "OK",
			"message": "API Running",
		})
	})
	fmt.Println("Server running di localhost:8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("gagal running server")
	}
}
