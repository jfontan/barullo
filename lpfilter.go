package barullo

import "math"

/*

FROM: http://www.musicdsp.org/en/latest/Filters/38-lp-and-hp-filter.html

r  = rez amount, from sqrt(2) to ~ 0.1
f  = cutoff frequency
(from ~0 Hz to SampleRate/2 - though many
synths seem to filter only  up to SampleRate/4)

The filter algo:
out(n) = a1 * in + a2 * in(n-1) + a3 * in(n-2) - b1*out(n-1) - b2*out(n-2)

Lowpass:
      c = 1.0 / tan(pi * f / sample_rate);

      a1 = 1.0 / ( 1.0 + r * c + c * c);
      a2 = 2* a1;
      a3 = a1;
      b1 = 2.0 * ( 1.0 - c*c) * a1;
      b2 = ( 1.0 - r * c + c * c) * a1;

Hipass:
      c = tan(pi * f / sample_rate);

      a1 = 1.0 / ( 1.0 + r * c + c * c);
      a2 = -2*a1;
      a3 = a1;
      b1 = 2.0 * ( c*c - 1.0) * a1;
      b2 = ( 1.0 - r * c + c * c) * a1;

*/

type LPFilter struct {
	freq float64
	res  float64

	a1, a2, a3 float64
	b1, b2     float64

	in0, in1, in2 float64
	out1, out2    float64

	input Node
}

func NewLPFilter(input Node, freq, res float64) *LPFilter {
	f := &LPFilter{
		freq:  freq,
		res:   res,
		input: input,
	}

	f.recalculate()

	return f
}

func (f *LPFilter) Set(freq, res float64) {
	f.freq = freq
	f.res = res
	f.recalculate()
}

func (f *LPFilter) recalculate() {
	c := 1.0 / math.Tan(math.Pi*f.freq/44100.0)

	f.a1 = 1.0 / (1.0 + f.res*c + c*c)
	f.a2 = 2 * f.a1
	f.a3 = f.a1
	f.b1 = 2.0 * (1.0 - c*c) * f.a1
	f.b2 = (1.0 - f.res*c + c*c) * f.a1
}

func (f *LPFilter) Get(offset int, buf []float64) {
	f.input.Get(offset, buf)

	for i := 0; i < len(buf); i++ {
		f.in2 = f.in1
		f.in1 = f.in0
		f.in0 = buf[i]

		o := f.a1*f.in0 + f.a2*f.in1 + f.a3*f.in2 - f.b1*f.out1 - f.b2*f.out2

		f.out2 = f.out1
		f.out1 = o

		buf[i] = o
	}
}
