package main

import (
	"log"
	"math"
	"os"

	"barullo"

	"github.com/gordonklaus/portaudio"
)

var (
	channelNum      = 1
	bitDepthInBytes = 2
	bufferSize      = 64
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

	f, err := os.Open("./assets/mario.json")
	if err != nil {
		log.Fatal(err)
	}

	events := barullo.GetEventsFromMidi(2, sampleRate, f)

	seq := barullo.NewSequence(events[len(events)-1].Offset+1, events)
	sig := barullo.NewSignal(barullo.Sin, sampleRate, seq)
	env := barullo.NewEnvelope(2000, 2000, 0.8, 10000, buf, seq)
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
