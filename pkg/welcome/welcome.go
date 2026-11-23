package welcome

import (
	"errors"
	"fmt"
	"net"
	"os"
	"sort"
	"strconv"

	"github.com/cqroot/prompt"
	"go.bug.st/serial"
)

const (
	defaultServerIP = "127.0.0.1"
	defaultPort     = "9765"
)

// Config содержит параметры подключения.
type AppConfig struct {
	Emulation     bool
	SerialPort    string
	ServerAddress string
}

func Setup() *AppConfig {

	config := &AppConfig{}

	fmt.Println("АО \"ЦС \"Звездочка\", Объект #160-161.")
	fmt.Println("Сервер обработки данных УСО, версия 0.3.0")
	fmt.Println("Павел Верещагин, эл. почта: ipashu@gmail.com")
	fmt.Println()

	// Получаем доступные порты.
	ports, err := serial.GetPortsList()
	if err != nil {
		fmt.Println("\033[31m✖ Ошибка при попытке получения списка портов.\033[0m")
		os.Exit(1)
	}

	// Сортируем порты.
	sort.Strings(ports)

	p := prompt.New()

	// Выбираем источник.
	source, err := p.Ask("Источник данных:").Choose([]string{"УСО", "Эмулятор УСО"})
	CheckErr(err)

	config.Emulation = source != "УСО"

	if !config.Emulation {
		// Выбираем порт УСО.
		serialPort, err := p.Ask("Порт УСО:").Choose(ports)
		CheckErr(err)
		config.SerialPort = serialPort
	}

	// Настройка IP-адреса сервера.
	serverIP, err := p.Ask("IP-адрес сервера:").Input(defaultServerIP)
	CheckErr(err)

	if net.ParseIP(serverIP) == nil {
		fmt.Println("\033[31m✖ Указан неправильный IP-адрес! \033[0m")
		os.Exit(1)
	}

	// Настройка сетевого порта.
	serverPort, err := p.Ask("Порт сервера:").Input(defaultPort)
	CheckErr(err)

	portInt, _ := strconv.Atoi(serverPort) // Проверка на корректность номера порта

	if portInt > 1 && portInt < 65535 {
		config.ServerAddress = net.JoinHostPort(serverIP, serverPort) // Формируем адрес сервера.
	} else {
		fmt.Println("\033[31m✖ Указан неправильный порт сервера! \033[0m")
		os.Exit(1)
	}

	return config

}

// Обработка ошибок prompt
func CheckErr(err error) {
	if err != nil {
		if errors.Is(err, prompt.ErrUserQuit) {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		} else {
			panic(err)
		}
	}
}
