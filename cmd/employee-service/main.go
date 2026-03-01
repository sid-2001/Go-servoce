package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"erp-microservices/internal/common"
)

func main() {
	kafkaREST := envOr("KAFKA_REST_URL", "http://kafka-rest:8082")
	databaseURL := envOr("DATABASE_URL", "postgres://postgres:postgres@postgres:5432/erp?sslmode=disable")
	port := envOr("PORT", "8082")

	client := common.NewKafkaClient(kafkaREST)
	store := common.NewPostgresStore(databaseURL)
	mustInit(store, `CREATE TABLE IF NOT EXISTS employees (id TEXT PRIMARY KEY, data JSONB NOT NULL);`)

	// Consume department events to illustrate service-to-service event flow.
	client.StartConsumer("employee-group", "employee-consumer", "department-events", func(msg []byte) {
		log.Printf("employee-service observed department event: %s", string(msg))
	})

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) { w.Write([]byte("ok")) })
	mux.HandleFunc("/employees", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			rows, err := store.QueryJSON(`SELECT COALESCE(json_agg(data), '[]'::json) FROM employees;`)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			writeRawJSON(w, rows, http.StatusOK)
		case http.MethodPost:
			var payload common.Employee
			if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			raw, _ := json.Marshal(payload)
			sql := fmt.Sprintf(`INSERT INTO employees(id, data) VALUES (%s, %s::jsonb) ON CONFLICT (id) DO UPDATE SET data = EXCLUDED.data;`, common.QuoteLiteral(payload.ID), common.QuoteLiteral(string(raw)))
			if err := store.Exec(sql); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			client.Publish("employee-events", payload)
			writeJSON(w, payload, http.StatusCreated)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	log.Printf("employee-service running on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}

func mustInit(store *common.PostgresStore, sql string) {
	if err := store.Exec(sql); err != nil {
		log.Fatal(err)
	}
}

func writeRawJSON(w http.ResponseWriter, data []byte, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(data)
}

func writeJSON(w http.ResponseWriter, data any, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
