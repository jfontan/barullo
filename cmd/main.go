package main

import (
	"log"
	"math"

	"barullo"

	"github.com/gordonklaus/portaudio"
)

var (
	channelNum      = 1
	bitDepthInBytes = 2
	bufferSize      = 64 * 10
)

const (
	sampleRate = 44100
)

func main() {
	mPortaudio()
}

func mPortaudio() {
	portaudio.Initialize()
	defer portaudio.Terminate()

	buf := make([]float64, bufferSize)
	out := make([]float32, bufferSize)

	stream, err := portaudio.OpenDefaultStream(0, 1, sampleRate, bufferSize, &out)
	if err != nil {
		log.Fatal(err)
	}
	defer stream.Close()

	if err := stream.Start(); err != nil {
		log.Fatal(err)
	}
	defer stream.Stop()

	noteLength := sampleRate / 4

	o := 2
	seq := barullo.NewSequence(noteLength*4,
		[]barullo.Event{
			{noteLength * 0, "C", o, barullo.NotePress},
			{noteLength*0 + noteLength/2, "C", o, barullo.NoteRelease},
			{noteLength * 1, "D", o, barullo.NotePress},
			{noteLength*1 + noteLength/2, "D", o, barullo.NoteRelease},
			{noteLength * 2, "E", o, barullo.NotePress},
			{noteLength*2 + noteLength/2, "E", o, barullo.NoteRelease},
			{noteLength * 3, "F", o, barullo.NotePress},
			{noteLength*3 + noteLength/2, "F", o, barullo.NoteRelease},
		},
	)
	// sig := barullo.NewSignal(barullo.Sin, sampleRate, seq)
	// sig := barullo.NewPulse(0.5, sampleRate, seq)
	sig := barullo.NewTriangle(0.0, sampleRate, seq)
	env := barullo.NewEnvelope(2000/4, 2000/4, 0.8, 10000/4, buf, seq)
	lp := barullo.NewLPFilter(500.8, 0.8)
	var sampleOffset int64
	for {
		sig.Get(int(sampleOffset), buf)
		env.Get(int(sampleOffset), buf)

		freq := math.Sin(float64(sampleOffset)/44100.0)*500.0 + 500.0
		res := math.Sin(float64(sampleOffset)/44100.0*4.0)*math.Sqrt2/4.0 + math.Sqrt2/4.0
		lp.Set(freq, res)
		lp.Get(int(sampleOffset), buf)

		sampleOffset += int64(bufferSize)

		f64ToF32Copy(out, buf)

		if err := stream.Write(); err != nil {
			log.Printf("error writing to stream : %v\n", err)
		}
	}
}

func f64ToF32Copy(dst []float32, src []float64) {
	for i := range src {
		dst[i] = float32(src[i])
	}
}
