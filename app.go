package main

import (
	"chemyserver/pkg/emulator"
	"chemyserver/pkg/reader"
	"chemyserver/pkg/welcome"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

type App struct {
	Frame          Frame `json:"frame"`
	ChannelsConfig []byte
	ctx            context.Context
	mu             sync.RWMutex
}

type Frame struct {
	AnalogValues   []int  `json:"analogValues"`
	DiscreteValues []bool `json:"discreteValues"`
}

func NewApp() *App {
	// Читаем файл конфигурации каналов
	configData, err := os.ReadFile("config.json")

	if err != nil {
		fmt.Println(ErrorMark+" Ошибка чтения конфигурации:", err)
		os.Exit(1)
	}

	return &App{
		ChannelsConfig: configData,
	}

}

func (a *App) run(ctx context.Context) {

	a.ctx = ctx

	// Получение конфигурации.
	appConfig, err := welcome.Setup()

	if err != nil {
		fmt.Println(ErrorMark+" Ошибка настройки:", err)
		os.Exit(1)
	}

	// Выбор режима работы.
	if appConfig.Emulation {

		// Запуск эмуляции
		go a.startEmulatorProcess()

	} else {

		// Запуск чтения УСО…
		go a.startReadingProcess(appConfig.SerialPort)

	}

	go a.startServerProcess(appConfig.ServerIP, appConfig.ServerPort)

	<-ctx.Done()

	fmt.Println("Завершение работы приложения…")

	time.Sleep(time.Second * 1)

}

func (a *App) startReadingProcess(seialPort string) {

	r := reader.NewReader(seialPort, 19200, "N", 8)

	go r.Run(a.ctx)

	for values := range r.Stream {
		a.mu.Lock()
		a.Frame.AnalogValues = values.AnalogValues
		a.Frame.DiscreteValues = values.DiscreteValues
		a.mu.Unlock()
	}

	fmt.Println("Синхронизация потоков остановлена.")

}

func (a *App) startEmulatorProcess() {

	e := emulator.NewEmulator()

	go e.Run(a.ctx)

	for values := range e.Stream {
		a.mu.Lock()
		a.Frame.AnalogValues = values.AnalogValues
		a.Frame.DiscreteValues = values.DiscreteValues
		a.mu.Unlock()
	}

	fmt.Println("Синхронизация потоков остановлена.")

}

func (a *App) startServerProcess(ip string, port int) {
	subFS, _ := fs.Sub(assets, "frontend/dist")

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServerFS(subFS))
	mux.HandleFunc("/config", configHandler(a.ChannelsConfig))
	mux.HandleFunc("/frame", a.frameHandler)

	// Формируем адрес сервера.
	address := net.JoinHostPort(ip, strconv.Itoa(port))

	// 1. Создаем сервер
	srv := &http.Server{
		Addr:    address,
		Handler: mux,
	}

	fmt.Printf(CheckMark+" Сервер запущен на http://%s 🚀\n", address)

	err := srv.ListenAndServe()

	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Printf("HTTP server error: %v", err)
	}

}

// Отдаем JSON конфигурации
func configHandler(data []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	}
}

// Отдаем JSON поток данных
func (a *App) frameHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	a.mu.RLock()
	err := json.NewEncoder(w).Encode(a.Frame)
	a.mu.RUnlock()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
