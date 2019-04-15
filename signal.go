package barullo

import "math"

const twoPi = math.Pi * 2.0

type SignalKind func(freq, sampleRate float64, offset int, buf []float64)

var _ SignalKind = Sin

func Sin(freq, sampleRate float64, offset int, buf []float64) {
	duration := len(buf)
	for i := 0; i < duration; i++ {
		pos := float64(offset + i)
		buf[i] = math.Sin((twoPi / sampleRate) * freq * pos)
	}
}

func NewSignal(gen SignalKind, freq, sampleRate float64) *Signal {
	return &Signal{
		gen:        gen,
		freq:       freq,
		sampleRate: sampleRate,
	}
}

type Signal struct {
	gen              SignalKind
	freq, sampleRate float64
}

func (s *Signal) Get(offset int, buf []float64) {
	s.gen(s.freq, s.sampleRate, offset, buf)
}
