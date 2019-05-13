package barullo

func NewEnvelope(a, d int, s float64, r int, in Node, seq *Sequence) *Envelope {
	return &Envelope{
		attack:  a,
		decay:   d,
		sustain: s,
		release: r,
		input:   in,
		seq:     seq,
	}
}

type Envelope struct {
	start   int
	noteOff int

	attack  int
	decay   int
	sustain float64
	release int

	attackEnd    int
	attackSpeed  float64
	decayEnd     int
	decaySpeed   float64
	releaseEnd   int
	releaseSpeed float64
	vol          float64

	input Node
	buf   []float64
	seq   *Sequence
}

func (e *Envelope) recalculate() {
	e.attackEnd = e.start + e.attack
	e.attackSpeed = 1.0 / float64(e.attack)
	e.decayEnd = e.attackEnd + e.decay
	e.decaySpeed = (1.0 - e.sustain) / float64(e.decay)
	e.releaseEnd = e.noteOff + e.release
	e.releaseSpeed = e.sustain / float64(e.release)
}

func (e *Envelope) Get(offset int, buf []float64) {
	if len(e.buf) != len(buf) {
		e.buf = make([]float64, len(buf))
	}

	e.input.Get(offset, e.buf)

	for i := range buf {
		pos := offset + i

		event := e.seq.Get(pos)
		k := event.Key
		switch k {
		case NotePress:
			e.start = pos
			e.vol = 0
			e.recalculate()

		case NoteRelease:
			e.noteOff = pos
			e.vol = e.sustain
			e.recalculate()
		}

		if k == NotePress || k == NotePressed {
			if pos >= e.start && pos < e.attackEnd {
				e.vol += e.attackSpeed
			} else if pos >= e.attackEnd && pos < e.decayEnd {
				e.vol -= e.decaySpeed
			} else {
				e.vol = e.sustain
			}
		} else if k == NoteOff || k == NoteRelease {
			if pos < e.releaseEnd {
				e.vol -= e.releaseSpeed
			} else {
				e.vol = 0.0
			}
		} else {
			e.vol = 0.0
		}

		buf[i] = e.buf[i] * e.vol
	}
}
