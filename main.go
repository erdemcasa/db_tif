package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type TifEntry struct {
	ID        int       `json:"id"`
	TifLabel  string    `json:"tif_label"`
	Author    string    `json:"author"`
	IPAddress string    `json:"ip_address"`
	UserAgent string    `json:"user_agent"`
	CreatedAt time.Time `json:"created_at"`
}

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("sqlite3", "./db_tif.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	createTable()

	http.HandleFunc("/api/tifs", handleTifs)
	http.Handle("/", http.FileServer(http.Dir("./static")))

	log.Println("Server listening on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func createTable() {
	query := `
	CREATE TABLE IF NOT EXISTS tifs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		tif_label TEXT NOT NULL,
		author TEXT NOT NULL,
		ip_address TEXT DEFAULT '',
		user_agent TEXT DEFAULT '',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal("Error creating table:", err)
	}
}

func handleTifs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "GET" {
		rows, err := db.Query("SELECT id, tif_label, author, ip_address, user_agent, created_at FROM tifs ORDER BY created_at DESC")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var tifs []TifEntry
		for rows.Next() {
			var t TifEntry
			if err := rows.Scan(&t.ID, &t.TifLabel, &t.Author, &t.IPAddress, &t.UserAgent, &t.CreatedAt); err != nil {
				log.Println("Scan error:", err)
				continue
			}
			tifs = append(tifs, t)
		}

		if tifs == nil {
			tifs = []TifEntry{} // Return empty array instead of null
		}
		json.NewEncoder(w).Encode(tifs)

	} else if r.Method == "POST" {
		var t TifEntry
		if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		// Validations
		t.TifLabel = strings.TrimSpace(t.TifLabel)
		t.Author = strings.TrimSpace(t.Author)

		if t.Author == "" {
			http.Error(w, "Author is required", http.StatusBadRequest)
			return
		}
		if t.TifLabel == "" || !strings.HasSuffix(t.TifLabel, "tif") {
			http.Error(w, "tif_label must end with 'tif'", http.StatusBadRequest)
			return
		}

		// Extract IP Address and User Agent
		ip := r.Header.Get("X-Forwarded-For")
		if ip == "" {
			ip = r.RemoteAddr
		}
		// S'il y a un port dans RemoteAddr, on l'enlève pour plus de propreté (ex: 127.0.0.1:54321 -> 127.0.0.1)
		if strings.Contains(ip, ":") && !strings.Contains(ip, "[") {
			ip = strings.Split(ip, ":")[0]
		}
		
		t.IPAddress = ip
		t.UserAgent = r.UserAgent()

		stmt, err := db.Prepare("INSERT INTO tifs (tif_label, author, ip_address, user_agent) VALUES (?, ?, ?, ?)")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer stmt.Close()

		res, err := stmt.Exec(t.TifLabel, t.Author, t.IPAddress, t.UserAgent)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		id, _ := res.LastInsertId()
		t.ID = int(id)
		t.CreatedAt = time.Now()

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(t)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
