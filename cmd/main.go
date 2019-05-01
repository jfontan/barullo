package main

import (
	"log"
	"os"

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

	f, err := os.Open("./assets/mario.json")
	if err != nil {
		log.Fatal(err)
	}

	events := barullo.GetEventsFromMidi(2, sampleRate, f)

	seq := barullo.NewSequence(events[len(events)-1].Offset+1, events)
	sig := barullo.NewPulse(0.5, sampleRate, seq)
	env := barullo.NewEnvelope(2000/4, 2000/4, 0.8, 10000/4, buf, seq)
	lp := barullo.NewLPFilter(500.8, 0.8)
	var sampleOffset int64
	for {
		sig.Get(int(sampleOffset), buf)
		env.Get(int(sampleOffset), buf)
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
