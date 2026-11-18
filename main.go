package main

import (
	"embed"
	"encoding/json"
	"net/http"
)

//go:embed frontend/dist
var assets embed.FS

// ANSI Codes
const (
	CheckMark = "\033[32m✔\033[0m"
	ErrorMark = "\033[31m✖\033[0m"
)

func main() {
	app := NewApp()
	app.run()
}

// Отдаем JSON конфигурации
func configHandler(data []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	}
}

func (a *App) frameHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	a.mu.RLock()
	err := json.NewEncoder(w).Encode(a.Frame)
	a.mu.RUnlock()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
