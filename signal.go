package barullo

import (
	"math"
	"math/rand"
)

const twoPi = math.Pi * 2.0

type SignalKind func(freq, sampleRate float64, offset int, buf []float64)

var _ SignalKind = Sin
var _ SignalKind = Noise

func Sin(freq, sampleRate float64, offset int, buf []float64) {
	duration := len(buf)
	for i := 0; i < duration; i++ {
		pos := float64(offset + i)
		buf[i] = math.Sin((twoPi / sampleRate) * freq * pos)
	}
}

var r *rand.Rand = rand.New(rand.NewSource(99))

func Noise(freq, sampleRate float64, offset int, buf []float64) {

	duration := len(buf)
	for i := 0; i < duration; i++ {
		//pos := float64(offset + i)
		buf[i] = r.Float64()
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

type Pulse struct {
	seq        *Sequence
	sampleRate float64
	duty       float64
}

func NewPulse(duty float64, sampleRate float64, seq *Sequence) *Pulse {
	return &Pulse{
		duty:       duty,
		sampleRate: sampleRate,
		seq:        seq,
	}
}

func (p *Pulse) Get(offset int, buf []float64) {
	for i := 0; i < len(buf); i++ {
		e := p.seq.Get(offset + i)
		cycle := p.sampleRate / e.Frequency()
		pos := math.Mod(float64(offset+i), cycle)

		if pos < p.duty*cycle {
			buf[i] = -1.0
		} else {
			buf[i] = 1.0
		}
	}
}

type Triangle struct {
	seq        *Sequence
	sampleRate float64
	duty       float64
}

func NewTriangle(duty float64, sampleRate float64, seq *Sequence) *Triangle {
	return &Triangle{
		duty:       duty,
		sampleRate: sampleRate,
		seq:        seq,
	}
}

func (p *Triangle) Get(offset int, buf []float64) {
	for i := 0; i < len(buf); i++ {
		e := p.seq.Get(offset + i)
		cycle := p.sampleRate / e.Frequency()
		pos := math.Mod(float64(offset+i), cycle)

		d := p.duty * cycle

		if pos < d {
			buf[i] = (pos/d)*2.0 - 1.0
		} else {
			buf[i] = (1.0-(pos-d)/(cycle-d))*2.0 - 1.0
		}
	}
}
