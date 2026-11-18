package emulator

import (
	"context"
	"math"
	"math/rand"
	"time"
)

const (
	analogChannels   = 112
	discreteChannels = 320
	frameTime        = 1000 // мс
)

type Frame struct {
	AnalogValues   []int
	DiscreteValues []bool
}

// Emulator генерирует независимые правдоподобные значения каналов для демонстрации UI.
// Каналы не имеют семантической привязки друг к другу или к оборудованию.
type Emulator struct {
	rng *rand.Rand

	analog       [analogChannels]float64
	analogTarget [analogChannels]float64
	analogSpeed  [analogChannels]float64
	discrete     [discreteChannels]bool

	Stream chan Frame
}

func NewEmulator() *Emulator {
	e := &Emulator{
		rng:    rand.New(rand.NewSource(time.Now().UnixNano())),
		Stream: make(chan Frame),
	}

	for i := range e.analog {
		// Начальные значения распределены по всему диапазону, но не у границ.
		e.analog[i] = e.rnd(100, 923)
		e.analogTarget[i] = e.rnd(0, 1023)
		// Разная скорость изменения исключает синхронное движение каналов.
		e.analogSpeed[i] = e.rnd(0.004, 0.025)
	}

	for i := range e.discrete {
		e.discrete[i] = e.rng.Intn(2) == 1
	}

	return e
}

func (e *Emulator) rnd(min, max float64) float64 {
	return min + e.rng.Float64()*(max-min)
}

func clamp(value, min, max float64) float64 {
	return math.Max(min, math.Min(max, value))
}

func (e *Emulator) Run(ctx context.Context) {
	defer close(e.Stream)

	for {
		select {
		case <-ctx.Done():
			return
		case e.Stream <- e.Next():
			time.Sleep(frameTime * time.Millisecond)
		}
	}
}

// Next возвращает очередной кадр с неизменным количеством каналов.
func (e *Emulator) Next() Frame {
	frame := Frame{
		AnalogValues:   make([]int, analogChannels),
		DiscreteValues: make([]bool, discreteChannels),
	}

	for i := range e.analog {
		// Периодически выбираем новую цель. Само значение подходит к ней плавно.
		if e.rng.Float64() < 0.008 {
			e.analogTarget[i] = e.rnd(0, 1023)
		}

		value := e.analog[i]
		value += (e.analogTarget[i] - value) * e.analogSpeed[i]
		// Небольшой шум имитирует колебания измерений.
		value += e.rnd(-1.2, 1.2)
		e.analog[i] = clamp(value, 0, 1023)
		frame.AnalogValues[i] = int(math.Round(e.analog[i]))
	}

	for i := range e.discrete {
		if e.rng.Float64() < 0.00015 {
			e.discrete[i] = !e.discrete[i]
		}
		frame.DiscreteValues[i] = e.discrete[i]
	}

	return frame
}
