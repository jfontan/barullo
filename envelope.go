package barullo

func NewEnvelope(a, d int, s float64, r int, in []float64) *Envelope {
	return &Envelope{
		attack:  a,
		decay:   d,
		sustain: s,
		release: r,
		in:      in,
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

	input *Node
	in    []float64
}

type KeyStatus int

const (
	NoteOff KeyStatus = iota
	NotePress
	NotePressed
	NoteRelease
)

func seqKey(offset int) KeyStatus {
	pos := offset % 44100

	if pos == 0 {
		return NotePress
	} else if pos < 20000 {
		return NotePressed
	} else if pos == 20000 {
		return NoteRelease
	}

	return NoteOff
}

func (e *Envelope) recalculate() {
	e.attackEnd = e.start + e.attack
	e.attackSpeed = 1.0 / float64(e.attack)
	e.decayEnd = e.attackEnd + e.decay
	e.decaySpeed = (1.0 - e.sustain) / float64(e.decay)
	e.releaseEnd = e.noteOff + e.release
	e.releaseSpeed = e.sustain / float64(e.release)
}

func (e *Envelope) Get(offset int, buf []float64) error {
	for i := range buf {
		pos := offset + i

		k := seqKey(pos)
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
			// println("error")
			e.vol = 0.0
		}

		buf[i] = e.in[i] * e.vol
	}

	return nil
}
