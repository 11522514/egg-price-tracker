package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type EggPrice struct {
	ID            int     `json:"id"`
	Date          string  `json:"date"`
	Location      string  `json:"location"`
	PricePerDozen float64 `json:"price_per_dozen"`
	Source        string  `json:"source"`
	CreatedAt     string  `json:"created_at"`
}

type Location struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type PriceComparison struct {
	Location      string  `json:"location"`
	CurrentPrice  float64 `json:"current_price"`
	NationalPrice float64 `json:"national_price"`
	Difference    float64 `json:"difference"`
	Percentage    float64 `json:"percentage"`
}

var db *sql.DB

func initDB() {
	var err error
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "password")
	dbName := getEnv("DB_NAME", "egg_tracker")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	log.Println("Database connection established")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getPricesHandler(w http.ResponseWriter, r *http.Request) {
	location := r.URL.Query().Get("location")
	limit := r.URL.Query().Get("limit")

	if limit == "" {
		limit = "30"
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		http.Error(w, "Invalid limit parameter", http.StatusBadRequest)
		return
	}

	var query string
	var args []interface{}

	if location != "" {
		query = `SELECT id, date, location, price_per_dozen, source, created_at 
				FROM egg_prices WHERE location = $1 
				ORDER BY date DESC LIMIT $2`
		args = []interface{}{location, limitInt}
	} else {
		query = `SELECT id, date, location, price_per_dozen, source, created_at 
				FROM egg_prices ORDER BY date DESC LIMIT $1`
		args = []interface{}{limitInt}
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var prices []EggPrice
	for rows.Next() {
		var p EggPrice
		err := rows.Scan(&p.ID, &p.Date, &p.Location, &p.PricePerDozen, &p.Source, &p.CreatedAt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		prices = append(prices, p)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(prices)
}

func getLocationsHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, name, type FROM locations ORDER BY name")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var locations []Location
	for rows.Next() {
		var l Location
		err := rows.Scan(&l.ID, &l.Name, &l.Type)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		locations = append(locations, l)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(locations)
}

func addPriceHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var price EggPrice
	if err := json.NewDecoder(r.Body).Decode(&price); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	query := `INSERT INTO egg_prices (date, location, price_per_dozen, source) 
			 VALUES ($1, $2, $3, $4) RETURNING id`

	err := db.QueryRow(query, price.Date, price.Location, price.PricePerDozen, price.Source).Scan(&price.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(price)
}

func getComparisonHandler(w http.ResponseWriter, r *http.Request) {
	// Get the most recent prices for each location
	query := `
		WITH latest_prices AS (
			SELECT DISTINCT ON (location) location, price_per_dozen, date
			FROM egg_prices 
			ORDER BY location, date DESC
		)
		SELECT 
			lp.location,
			lp.price_per_dozen,
			np.price_per_dozen as national_price
		FROM latest_prices lp
		CROSS JOIN (
			SELECT price_per_dozen 
			FROM latest_prices 
			WHERE location = 'NATIONAL'
		) np
		WHERE lp.location != 'NATIONAL'
		ORDER BY lp.location`

	rows, err := db.Query(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var comparisons []PriceComparison
	for rows.Next() {
		var c PriceComparison
		err := rows.Scan(&c.Location, &c.CurrentPrice, &c.NationalPrice)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		c.Difference = c.CurrentPrice - c.NationalPrice
		if c.NationalPrice > 0 {
			c.Percentage = (c.Difference / c.NationalPrice) * 100
		}

		comparisons = append(comparisons, c)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comparisons)
}

func main() {
	initDB()
	defer db.Close()

	r := mux.NewRouter()

	// API routes
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/prices", getPricesHandler).Methods("GET")
	api.HandleFunc("/prices", addPriceHandler).Methods("POST")
	api.HandleFunc("/locations", getLocationsHandler).Methods("GET")
	api.HandleFunc("/comparison", getComparisonHandler).Methods("GET")

	// Serve static files
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	// Enable CORS
	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"*"}),
	)(r)

	port := getEnv("PORT", "8080")
	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, corsHandler))
}
