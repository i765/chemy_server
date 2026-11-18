package main

import (
	"chemyserver/pkg/emulator"
	"chemyserver/pkg/reader"
	"chemyserver/pkg/welcome"
	"fmt"
	"io/fs"
	"net"
	"net/http"
	"os"
	"strconv"
	"sync"
)

type App struct {
	Frame          Frame `json:"frame"`
	ChannelsConfig []byte
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

func (a *App) run() {

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

	// Инициализация frontend.
	frontendFS, err := fs.Sub(assets, "frontend/dist")

	if err != nil {
		fmt.Println(ErrorMark+" Ошибка загрузки frontend:", err)
		os.Exit(1)
		return
	}

	// Статические файлы frontend.
	fileServer := http.FileServer(http.FS(frontendFS))

	mux := http.NewServeMux()
	mux.Handle("/", fileServer)
	mux.HandleFunc("/config", configHandler(a.ChannelsConfig))
	mux.HandleFunc("/frame", a.frameHandler)

	// Формируем адрес сервера.
	address := net.JoinHostPort(appConfig.ServerIP, strconv.Itoa(appConfig.ServerPort))

	// Запускаем HTTP-сервер.
	fmt.Printf(CheckMark+" Сервер запущен на http://%s 🚀\n", address)

	if err := http.ListenAndServe(address, mux); err != nil {
		fmt.Println(ErrorMark+" Ошибка сервера:", err)
		os.Exit(1)
	}

}

func (a *App) startReadingProcess(seialPort string) {

	r := reader.NewReader(seialPort, 19200, "N", 8)

	ch := make(chan reader.Frame)

	go r.Run(ch)

	for values := range ch {
		a.mu.Lock()
		a.Frame.AnalogValues = values.AnalogValues
		a.Frame.DiscreteValues = values.DiscreteValues
		a.mu.Unlock()
	}
}

func (a *App) startEmulatorProcess() {

	e := emulator.NewEmulator()

	ch := make(chan emulator.Frame)

	go e.Run(ch)

	for values := range ch {
		a.mu.Lock()
		a.Frame.AnalogValues = values.AnalogValues
		a.Frame.DiscreteValues = values.DiscreteValues
		a.mu.Unlock()
	}

}
