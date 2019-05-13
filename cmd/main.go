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

	e1 := barullo.GetEventsFromMidi(1, sampleRate, f)
	f.Seek(0, os.SEEK_SET)
	e2 := barullo.GetEventsFromMidi(2, sampleRate, f)

	eLength1 := e1[len(e1)-1].Offset + 6000
	eLength2 := e2[len(e2)-1].Offset + 6000

	eLength := eLength1
	if eLength2 > eLength1 {
		eLength = eLength2
	}

	seq1 := barullo.NewSequence(eLength, e1)
	seq2 := barullo.NewSequence(eLength, e2)

	sig1 := barullo.NewPulse(0.5, sampleRate, seq1)
	env1 := barullo.NewEnvelope(2000/4, 2000/4, 0.8, 10000/4, sig1, seq1)
	lp1 := barullo.NewLPFilter(env1, 8000.3, 1.1)

	sig2 := barullo.NewPulse(0.5, sampleRate, seq2)
	env2 := barullo.NewEnvelope(2000/4, 2000/4, 0.8, 10000/4, sig2, seq1)
	lp2 := barullo.NewLPFilter(env2, 8000.3, 1.1)

	mixer := barullo.NewMixer([]barullo.Node{lp1, lp2}, []float64{1.0, 1.0})

	var sampleOffset int64
	for {
		mixer.Get(int(sampleOffset), buf)

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
