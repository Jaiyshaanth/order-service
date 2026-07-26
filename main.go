// Order service — returns placed orders.
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

var (
	serviceName = env("SERVICE_NAME", "order-service")
	version     = env("APP_VERSION", "dev")
	environment = env("ENVIRONMENT", "local")
)

// Order is a single placed order.
type Order struct {
	ID     int     `json:"id"`
	Item   string  `json:"item"`
	Amount float64 `json:"amount"`
	Status string  `json:"status"`
}

var orders = []Order{
	{ID: 1001, Item: "Filter Coffee Powder 500g", Amount: 349.0, Status: "SHIPPED"},
	{ID: 1002, Item: "Cotton Veshti", Amount: 799.0, Status: "PACKED"},
	{ID: 1003, Item: "Brass Kuthuvilakku", Amount: 1499.0, Status: "PLACED"},
}

func writeJSON(w http.ResponseWriter, code int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(payload)
}

// NewRouter wires up every route. Exported so tests can use it directly.
func NewRouter() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok", "service": serviceName})
	})

	mux.HandleFunc("/ready", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, map[string]any{"ready": true, "service": serviceName})
	})

	mux.HandleFunc("/version", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{
			"service":     serviceName,
			"version":     version,
			"environment": environment,
		})
	})

	mux.HandleFunc("/orders", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, map[string]any{"count": len(orders), "items": orders})
	})

	return mux
}

func main() {
	port := env("PORT", "3003")
	log.Printf("%s (%s/%s) listening on :%s", serviceName, version, environment, port)
	if err := http.ListenAndServe(":"+port, NewRouter()); err != nil {
		log.Fatal(err)
	}
}
