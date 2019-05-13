package barullo

type Mixer struct {
	inputs  []Node
	volumes []float64
	bufIn   []float64
}

func NewMixer(inputs []Node, volumes []float64) *Mixer {
	v := volumes
	if v == nil {
		v = make([]float64, len(inputs))
		for i := range v {
			v[i] = 1.0
		}
	}

	return &Mixer{
		inputs:  inputs,
		volumes: v,
	}
}

func (m *Mixer) Get(offset int, buf []float64) {
	if len(m.bufIn) != len(buf) {
		m.bufIn = make([]float64, len(buf))
	}

	for i := range buf {
		buf[i] = 0.0
	}

	for i, input := range m.inputs {
		input.Get(offset, m.bufIn)
		vol := m.volumes[i]

		for j, v := range m.bufIn {
			buf[j] += v * vol
		}
	}
}
