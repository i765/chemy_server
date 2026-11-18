package welcome

import (
	"errors"
	"fmt"
	"net"
	"sort"
	"strconv"
	"strings"

	"github.com/cqroot/prompt"
	"go.bug.st/serial"
)

const (
	emulationOption = "Emul"
	defaultServerIP = "127.0.0.1"
	defaultPort     = "9765"
)

// Config содержит параметры подключения.
type Config struct {
	Emulation  bool
	SerialPort string
	ServerIP   string
	ServerPort int
}

func Setup() (*Config, error) {
	fmt.Println("АО \"ЦС \"Звездочка\", Объект #160-161.")
	fmt.Println("Сервер обработки данных УСО, версия 0.3.0")
	fmt.Println("Разработчик Верещагин Павел Артемович, эл. почта: ipashu@gmail.com")
	fmt.Println("")

	// Получаем доступные порты.
	ports, err := serial.GetPortsList()
	if err != nil {
		return nil, fmt.Errorf("get serial ports: %w", err)
	}

	// Сортируем порты.
	sort.Strings(ports)

	// Добавляем режим эмуляции.
	ports = append(ports, emulationOption)

	p := prompt.New()

	// Выбираем порт УСО или эмуляцию.
	selectedPort, err := p.
		Ask("Порт УСО:").
		Choose(ports)

	if err != nil {
		return nil, fmt.Errorf("select serial port: %w", err)
	}

	emulation := selectedPort == emulationOption

	// Настройка IP-адреса сервера.
	serverIP, err := p.
		Ask("IP-адрес сервера:").
		Input(defaultServerIP)

	if err != nil {
		return nil, fmt.Errorf("input server IP: %w", err)
	}

	serverIP = strings.TrimSpace(serverIP)

	if net.ParseIP(serverIP) == nil {
		return nil, errors.New("неправильный IP-адрес")
	}

	// Настройка сетевого порта.
	serverPortStr, err := p.
		Ask("Порт сервера:").
		Input(defaultPort)

	if err != nil {
		return nil, fmt.Errorf("input server port: %w", err)
	}

	serverPort, err := strconv.Atoi(strings.TrimSpace(serverPortStr))
	if err != nil || serverPort < 1 || serverPort > 65535 {
		return nil, errors.New("порт должен быть от 1 до 65535")
	}

	// Формируем конфигурацию.
	config := &Config{
		Emulation:  emulation,
		ServerIP:   serverIP,
		ServerPort: serverPort,
	}

	// COM-порт нужен только в режиме чтения.
	if !emulation {
		config.SerialPort = selectedPort
	}

	return config, nil
}
