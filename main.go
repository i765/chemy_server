package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"net"
	"net/http"
	"os"
	"strconv"
	"sync"

	"chemyserver/pkg/welcome"
)

//go:embed frontend/dist
var assets embed.FS

const (
	AnalogChannelsCount   = 112
	DiscreteChannelsCount = 320
)

type App struct {
	mu      sync.RWMutex
	Frame   Frame  `json:"frame"`
	Uptime  int    `json:"uptime"`
	Version string `json:"version"`
}

type Frame struct {
	AnalogValues   []int  `json:"analogValues"`
	DiscreteValues []bool `json:"discreteValues"`
}

func main() {
	app := &App{
		Frame: Frame{
			AnalogValues:   make([]int, AnalogChannelsCount),
			DiscreteValues: make([]bool, DiscreteChannelsCount),
		},
	}

	// Инициализация frontend.
	frontendFS, err := fs.Sub(assets, "frontend/dist")
	if err != nil {
		fmt.Println("Ошибка загрузки frontend:", err)
		return
	}

	// Получение конфигурации.
	config, err := welcome.Setup()
	if err != nil {
		fmt.Println("Ошибка настройки:", err)
		return
	}

	data, err := os.ReadFile("config.json")

	if err != nil {
		fmt.Println("Ошибка чтения конфигурации:", err)
		return
	}

	fmt.Println("Конфигурация успешно загружена.")

	// Выбор режима работы.
	if config.Emulation {

		// Emulator

	} else {

		fmt.Println("Открываю порт:", config.SerialPort)

		// Открытие порта.

		fmt.Println("Запуск чтения УСО...")

		// Запуск логики чтения УСО.

	}

	// Статические файлы frontend.
	fileServer := http.FileServer(http.FS(frontendFS))

	mux := http.NewServeMux()
	mux.Handle("/", fileServer)
	mux.HandleFunc("/config", configHandler(data))
	mux.HandleFunc("/frame", app.LiveHandler)

	// Формируем адрес сервера.
	address := net.JoinHostPort(config.ServerIP, strconv.Itoa(config.ServerPort))

	// Запускаем HTTP-сервер.
	fmt.Printf("Сервер запущен на http://%s\n", address)

	if err := http.ListenAndServe(address, mux); err != nil {
		fmt.Println("Ошибка сервера:", err)
	}

}

// Отдаем JSON конфигурации
func configHandler(data []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	}
}

func (a *App) LiveHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Доступ к данным идет через получателя (a)
	if err := json.NewEncoder(w).Encode(a.Frame); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

// Запись JSON.
func writeJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		return
	}

}
