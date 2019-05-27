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
	bufferSize      = 64 * 1000
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

	f, err := os.Open("./assets/take_on_me.json")
	if err != nil {
		log.Fatal(err)
	}

	e1 := barullo.GetEventsFromMidi(0, sampleRate, f)
	f.Seek(0, os.SEEK_SET)
	e2 := barullo.GetEventsFromMidi(1, sampleRate, f)
	f.Seek(0, os.SEEK_SET)
	e3 := barullo.GetEventsFromMidi(4, sampleRate, f)

	eLength1 := e1[len(e1)-1].Offset + 6000
	eLength2 := e2[len(e2)-1].Offset + 6000
	eLength3 := e3[len(e3)-1].Offset + 6000

	eLength := eLength1
	if eLength2 > eLength1 {
		eLength = eLength2
	}
	if eLength2 > eLength {
		eLength = eLength3
	}

	seq1 := barullo.NewSequence(eLength, e1)
	seq2 := barullo.NewSequence(eLength, e2)
	seq3 := barullo.NewSequence(eLength, e3)

	sig1 := barullo.NewPulse(0.25, sampleRate, seq1)
	env1 := barullo.NewEnvelope(2000/4, 2000/4, 0.8, 10000/4, sig1, seq1)
	lp1 := barullo.NewLPFilter(env1, 300.3, 1.1)

	sig2 := barullo.NewTriangle(0.8, sampleRate, seq2)
	env2 := barullo.NewEnvelope(100, 0, 1.0, 500, sig2, seq2)
	lp2 := barullo.NewLPFilter(env2, 8000.3, 6.2)

	sig3 := barullo.NewPulse(0.1, sampleRate, seq3)
	env3 := barullo.NewEnvelope(2000/4, 2000/4, 0.8, 10000/4, sig3, seq3)
	lp3 := barullo.NewLPFilter(env3, 8000.3, 1.1)

	mixer := barullo.NewMixer([]barullo.Node{lp1, lp2, lp3}, []float64{0.8, 0.8, 0.8})
	// mixer := barullo.NewMixer([]barullo.Node{lp1, lp2, lp3}, []float64{0.0, 0.0, 1.0})

	var sampleOffset int64
	sampleOffset = 44100 * 8
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
