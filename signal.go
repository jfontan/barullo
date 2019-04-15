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

func NewSignal(gen SignalKind, sampleRate float64, seq *Sequence) *Signal {
	return &Signal{
		gen:        gen,
		seq:        seq,
		sampleRate: sampleRate,
	}
}

type Signal struct {
	gen        SignalKind
	seq        *Sequence
	sampleRate float64
}

func (s *Signal) Get(offset int, buf []float64) {
	e := s.seq.Get(offset)
	s.gen(e.Frequency(), s.sampleRate, offset, buf)
}
