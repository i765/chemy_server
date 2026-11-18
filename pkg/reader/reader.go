package reader

import (
	"context"
	"fmt"
	"time"

	"go.bug.st/serial"
)

// ---------------------------------------------------------------------
// Параметры протокола
// ---------------------------------------------------------------------
const (
	bufferSize         = 32
	packetSize         = 32
	analogChannels     = 112
	totalAnalogSize    = 224 // 112 * 2 байта
	discretePacketSize = 20
	totalDiscreteSize  = 40  // 2 пакета по 20 байт
	discreteChannels   = 320 // 40 байт по 8 бит
	batchSize          = 10  // количество пакетов в сглаживающем буфере
)

// ---------------------------------------------------------------------
// Состояния парсера
// ---------------------------------------------------------------------
type rxState int

const (
	stateIdle rxState = iota
	stateRecvAnalog
	stateRecvDiscrete
)

// ---------------------------------------------------------------------
// Reader – инкапсулирует открытый порт и логику разбора пакетов
// ---------------------------------------------------------------------
type Reader struct {
	port  serial.Port
	state rxState

	// аналоговые данные
	analogData   []byte // промежуточный буфер (224 байта)
	analogMarker byte
	analogIndex  int
	analogOffset int
	analogReady  bool

	// дискретные данные
	discreteData   []byte // 40 байт
	discreteMarker byte
	discreteIndex  int
	discreteOffset int
	discreteReady  bool

	batchCount int
	bufferSize int
	acSums     []int
	dcSums     []int

	Frame  Frame
	Stream chan Frame
}

type Frame struct {
	AnalogValues   []int
	DiscreteValues []bool
}

// NewReader открывает порт (блокирует, пока порт не появится)
// Возвращает *Reader, готовый к использованию
func NewReader(portName string, baud int, parity string, bufferSize int) *Reader {

	// mode - настройки для открытия порта
	mode := &serial.Mode{
		BaudRate: baud,
		Parity:   parseParity(parity),
	}

	var p serial.Port

	// открываем порт
	for {

		var err error

		// ждем 100 мс
		time.Sleep(100 * time.Millisecond)

		// пытаемся открыть порт
		p, err = serial.Open(portName, mode)

		if err == nil {
			// если порт открыт, то
			// устанавливаем таймаут чтения для порта
			err = p.SetReadTimeout(500 * time.Millisecond)
		}

		// если порт открылся, то
		if err == nil {
			// выводим информацию о том, что порт готов к использованию
			fmt.Printf("\033[32m✔\033[0m Чтение данных через %s... ", portName)
			break
		}

		// если порт не открылся, то
		// выводим информацию об ошибке

		fmt.Printf("\033[31m✖\033[0m Не удалось открыть: %s: %v\n", portName, err)

		// ждем 5 сек
		time.Sleep(time.Second * 5)

	}

	return &Reader{
		port:         p,
		state:        stateIdle,
		analogData:   make([]byte, totalAnalogSize),
		discreteData: make([]byte, totalDiscreteSize),
		acSums:       make([]int, analogChannels),
		dcSums:       make([]int, discreteChannels),
		batchCount:   0,
		bufferSize:   bufferSize,
		Stream:       make(chan Frame),
	}

}

// Парсинг string в serial.Parity
func parseParity(parity string) serial.Parity {
	switch parity {
	case "N":
		return serial.NoParity
	case "E":
		return serial.EvenParity
	case "O":
		return serial.OddParity
	case "M":
		return serial.MarkParity
	case "S":
		return serial.SpaceParity
	default:
		return serial.NoParity
	}
}

// close – закрывает порт и канал
func (r *Reader) close() {
	if r.port != nil {
		r.port.Close()        // закрываем порт
		defer close(r.Stream) // закрываем канал
	}
}

// reset – возвращает парсер в исходное состояние
func (r *Reader) reset() {
	r.state = stateIdle
	r.analogMarker = 0
	r.analogIndex = 0
	r.analogOffset = 0
	r.discreteMarker = 0
	r.discreteIndex = 0
	r.discreteOffset = 0
	r.analogReady = false
	r.discreteReady = false
	fmt.Println("Выполнен сброс парсера")
}

// Run – основной цикл чтения
func (r *Reader) Run(ctx context.Context) {

	// создаем буфер для чтения порта
	buf := make([]byte, bufferSize)

Loop: // Метка для цикла (для break внутри select)
	// основной цикл чтения
	for {
		select {

		case <-ctx.Done():
			break Loop // Выход из цикла

		default:
			n, err := r.port.Read(buf)

			// ошибка чтения
			if err != nil {
				r.reset()
				fmt.Printf("Ошибка чтения: %e", err)
				break Loop
			}

			// нет данных для чтения
			if n == 0 {
				continue
			}

			// разбор полученных байтов
			for i := 0; i < n; i++ {
				b := buf[i]
				switch r.state {
				case stateIdle:
					switch {
					case b >= 0xF0 && b <= 0xF6: // аналоговый пакет
						r.state = stateRecvAnalog
						r.analogMarker = b
						r.analogIndex = 0
						r.analogOffset = int(b-0xF0) * packetSize
					case b >= 0xE0 && b <= 0xE1: // дискретный пакет
						r.state = stateRecvDiscrete
						r.discreteMarker = b
						r.discreteIndex = 0
						r.discreteOffset = int(b-0xE0) * discretePacketSize
					}
				case stateRecvAnalog:
					if r.analogOffset+r.analogIndex >= totalAnalogSize {
						fmt.Println("Переполненние аналогового блока")
						r.reset()
						break Loop
					}
					r.analogData[r.analogOffset+r.analogIndex] = b
					r.analogIndex++
					if r.analogIndex == packetSize {
						if r.analogMarker == 0xF6 {
							vals, err := r.parseAnalogPack(r.analogData)
							if err != nil {
								fmt.Printf("Ошибка парсинга аналогового блока: %v", err)
								r.reset()
								break Loop
							}
							r.Frame.AnalogValues = vals
							r.analogReady = true
						}
						r.state = stateIdle
					}
				case stateRecvDiscrete:
					if r.discreteOffset+r.discreteIndex >= totalDiscreteSize {
						fmt.Println("Переполненние дискретного блока – СБРОС")
						r.reset()
						break Loop
					}
					r.discreteData[r.discreteOffset+r.discreteIndex] = b
					r.discreteIndex++
					if r.discreteIndex == discretePacketSize {
						if r.discreteMarker == 0xE1 {
							bools, err := r.bytesToBools(r.discreteData)
							if err != nil {
								fmt.Printf("Ошибка парсинга дискретного блока: %v", err)
								r.reset()
								break Loop
							}
							r.Frame.DiscreteValues = bools
							r.discreteReady = true
						}
						r.state = stateIdle
					}
				}
			}

			// Полный кадр получен
			if r.analogReady && r.discreteReady {

				// --- пакетное сглаживание ---
				for i := 0; i < analogChannels; i++ {
					r.acSums[i] += r.Frame.AnalogValues[i]
				}
				for i := 0; i < discreteChannels; i++ {
					if r.Frame.DiscreteValues[i] {
						r.dcSums[i]++
					}
				}
				r.batchCount++

				if r.batchCount == r.bufferSize {
					avg := Frame{
						AnalogValues:   make([]int, analogChannels),
						DiscreteValues: make([]bool, discreteChannels),
					}

					for i := 0; i < analogChannels; i++ {
						avg.AnalogValues[i] = r.acSums[i] / r.bufferSize
					}

					for i := 0; i < discreteChannels; i++ {
						avg.DiscreteValues[i] = r.dcSums[i] > r.bufferSize/2 // строгое большинство
					}

					r.Stream <- avg // отправляем данные в канал

					// полный сброс batch
					r.batchCount = 0
					for i := range r.acSums {
						r.acSums[i] = 0
					}
					for i := range r.dcSums {
						r.dcSums[i] = 0
					}
				}

				r.analogReady = false
				r.discreteReady = false
			}
		}
	}

	fmt.Println("\nОстановка чтения.")
	r.close()

}

// ---------------------------------------------------------------------
// Вспомогательные функции парсинга
// ---------------------------------------------------------------------

// parseAnalogPack – из 224‑байтового буфера получает 112 значений (10‑бит)
func (r *Reader) parseAnalogPack(buffer []byte) ([]int, error) {
	if len(buffer) != totalAnalogSize {
		return nil, fmt.Errorf("analog pack size mismatch")
	}

	values := make([]int, analogChannels)

	for i := 0; i < analogChannels; i++ {
		hi := buffer[i*2]
		lo := buffer[i*2+1]
		values[i] = int(lo) | (int(hi&0x03) << 8)
	}
	return values, nil
}

// bytesToBools – 40 байт → 320 логических каналов
func (r *Reader) bytesToBools(data []byte) ([]bool, error) {
	if len(data) != totalDiscreteSize {
		return nil, fmt.Errorf("discrete pack size mismatch")
	}
	result := make([]bool, len(data)*8)
	idx := 0
	for _, b := range data {
		for bit := 0; bit < 8; bit++ {
			result[idx] = (b & (1 << bit)) != 0
			idx++
		}
	}
	return result, nil
}
