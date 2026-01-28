package main

import (
	"encoding/json"
	"fmt"
	"kasir-api/database"
	"kasir-api/handlers"
	"kasir-api/repositories"
	"kasir-api/services"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

type Inventory struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

var inventory = []Inventory{
	{ID: 1, Name: "Handphone xiaomi", Description: "Redmi note 10 pro"},
	{ID: 2, Name: "Handphone Iphone", Description: "Iphone 13 pro max"},
	{ID: 3, Name: "Handphone Samsung", Description: "Galaxy S21"},
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

// ubah Config
type Config struct {
	Port   string `mapstructure:"PORT"`
	DBConn string `mapstructure:"DB_CONN"`
}

func main() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		_ = viper.ReadInConfig()
	}

	config := Config{
		Port:   viper.GetString("PORT"),
		DBConn: viper.GetString("DB_CONN"),
	}

	// Setup database
	db, err := database.InitDB(config.DBConn)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	// GET localhost:8080/api/inventory/{id}
	// PUT localhost:8080/api/inventory/{id}
	// DELETE localhost:8080/api/inventory/{id}
	productRepo := repositories.NewProductRepository(db)
	productService := services.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productService)

	http.HandleFunc("/api/produk", productHandler.HandleProducts)
	http.HandleFunc("/api/produk/", productHandler.HandleProductByID)

	http.HandleFunc("/api/inventory/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getInventoryByID(w, r)
		case http.MethodPut:
			updateInventory(w, r)
		case http.MethodDelete:
			deleteInventory(w, r)
		}

	})

	// GET localhost:8080/api/inventory
	// POST localhost:8080/api/inventory
	http.HandleFunc("/api/inventory", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(inventory)
		case http.MethodPost:
			var Inventorybaru Inventory
			err := json.NewDecoder(r.Body).Decode(&Inventorybaru)
			if err != nil {
				http.Error(w, "Invalid request", http.StatusBadRequest)
				return
			}
			Inventorybaru.ID = len(inventory) + 1
			inventory = append(inventory, Inventorybaru)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
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

	fmt.Println("Server running di localhost:" + config.Port)

	err = http.ListenAndServe(":"+config.Port, nil)
	if err != nil {
		fmt.Println("gagal running server")
	}
}
