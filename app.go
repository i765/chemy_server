package main

import (
	"chemyserver/pkg/emulator"
	"chemyserver/pkg/reader"
	"chemyserver/pkg/welcome"
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"os"
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
	appConfig := welcome.Setup()

	// Выбор режима работы.
	if appConfig.Emulation {

		// Запуск эмуляции
		go a.startEmulatorProcess()

	} else {

		// Запуск чтения УСО…
		go a.startReadingProcess(appConfig)

	}
	// Запуск чтения сервера
	go a.startServerProcess(appConfig.ServerAddress, ctx)

	<-ctx.Done()

	fmt.Println("\n🔥 Завершение работы…")
	time.Sleep(time.Second * 3)

}

func (a *App) startReadingProcess(appConfig *welcome.AppConfig) {

	baudRate := os.Getenv("BAUD")
	if baudRate == "" {
		baudRate = "19200" // значение по умолчанию
	}

	serialParity := os.Getenv("PARITY")
	if serialParity == "" {
		serialParity = "NONE" // значение по умолчанию
	}

	bufferSize := os.Getenv("BUFFERSIZE")
	if bufferSize == "" {
		bufferSize = "16" // значение по умолчанию
	}

	r := reader.NewReader(appConfig.SerialPort, baudRate, serialParity, bufferSize)

	go r.Run(a.ctx)

	for values := range r.Stream {
		a.mu.Lock()
		a.Frame.AnalogValues = values.AnalogValues
		a.Frame.DiscreteValues = values.DiscreteValues
		a.mu.Unlock()
	}

	fmt.Printf("\033[32m✔\033[0m %s освобожден.", appConfig.SerialPort)

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

	fmt.Println("\033[32m✔\033[0m Эмуляция остановлена.")

}

func (a *App) startServerProcess(address string, ctx context.Context) {
	// Подключаем фронтенд.
	subFS, _ := fs.Sub(assets, "frontend/dist")

	// Настраеваем маршрутизацию.
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServerFS(subFS))
	mux.HandleFunc("/config", configHandler(a.ChannelsConfig))
	mux.HandleFunc("/frame", a.frameHandler)

	// Создаем сервер.
	srv := &http.Server{
		Addr:    address,
		Handler: enableCORS(mux),
	}

	// Ожидаение заполнения буфера.
	for {
		time.Sleep(50 * time.Millisecond)
		if len(a.Frame.AnalogValues) == reader.AnalogChannels && len(a.Frame.DiscreteValues) == reader.DiscreteChannels {
			fmt.Println("\r\033[32m✔\033[0m УСО готово.")
			break
		}
	}

	// Запускаем его в отдельном процессе.
	go func() {
		fmt.Printf("🚀 Сервер запущен на http://%s, Ctrl+C для остановки.", address)
		// Здесь происходит блокирующая операция ListenAndServe.
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Println("Ошибка сервера:", err)
		}
	}()

	//Graceful Shutdown.

	<-ctx.Done()

	if err := srv.Shutdown(ctx); err != nil {
		fmt.Printf("Сервер завершился с ошибкой при Shutdown: %v", err)
	}

	fmt.Println("\033[32m✔\033[0m Сервер успешно остановлен.")

}

// Отдаем JSON конфигурацию каналов
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

// Разрешаем CORS
func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
