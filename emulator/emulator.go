package emulator

import (
	"math"
	"math/rand"
	"time"
)

// ---------------------------------------------------------------------
// Параметры эмуляции
// ---------------------------------------------------------------------
const (
	analogChannels   = 112
	discreteChannels = 320
	analogChangeStep = 0.1
)

type Stream struct {
	ACData []int  `json:"ACData"`
	DCData []bool `json:"DCData"`
}

type Emulator struct {
	stream    Stream
	FrameTime int
	step      float64
	t         float64
}

func NewEmulator() *Emulator {
	return &Emulator{
		step: analogChangeStep,
	}
}

func (e *Emulator) Run() error {
	for {
		e.stream = e.Next()
		// runtime.EventsEmit(ctx, "stream", e.stream)
		time.Sleep(time.Millisecond * time.Duration(e.FrameTime))
	}
}

func (e *Emulator) Next() Stream {
	stream := Stream{
		ACData: make([]int, analogChannels),
		DCData: make([]bool, discreteChannels),
	}

	for i := 0; i < analogChannels; i++ {
		s := math.Sin(e.t + float64(i)*0.1)
		stream.ACData[i] = int((s + 1) * 511.5)
	}

	// Обрабатываем дискретные каналы
	if e.stream.DCData == nil {
		// Инициализация при первом запуске - все каналы получают случайные значения
		for i := 0; i < discreteChannels; i++ {
			stream.DCData[i] = rand.Intn(2) == 1
		}
	} else {
		// Копируем текущее состояние всех дискретных каналов
		copy(stream.DCData, e.stream.DCData)

		// Выбираем один случайный дискретный канал для изменения
		if discreteChannels > 0 {
			randomIndex := rand.Intn(discreteChannels)
			// Меняем значение только этого канала на противоположное
			stream.DCData[randomIndex] = !stream.DCData[randomIndex]
		}
	}

	e.t += e.step

	return stream
}
