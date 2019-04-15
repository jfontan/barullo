package main

import (
	"log"

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

	sinBuf := make([]float64, bufferSize)
	envBuf := make([]float64, bufferSize)
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

	o := 2
	seq := barullo.NewSequence(sampleRate*4,
		[]barullo.Event{
			{sampleRate * 0, "C", o, barullo.NotePress},
			{sampleRate*0 + 22050, "C", o, barullo.NoteRelease},
			{sampleRate * 1, "D", o, barullo.NotePress},
			{sampleRate*1 + 22050, "D", o, barullo.NoteRelease},
			{sampleRate * 2, "E", o, barullo.NotePress},
			{sampleRate*2 + 22050, "E", o, barullo.NoteRelease},
			{sampleRate * 3, "F", o, barullo.NotePress},
			{sampleRate*3 + 22050, "F", o, barullo.NoteRelease},
		},
	)
	sig := barullo.NewSignal(barullo.Sin, sampleRate, seq)
	env := barullo.NewEnvelope(2000, 2000, 0.8, 10000, sinBuf, seq)
	var sampleOffset int64
	for {
		sig.Get(int(sampleOffset), sinBuf)
		env.Get(int(sampleOffset), envBuf)

		sampleOffset += int64(bufferSize)

		f64ToF32Copy(out, envBuf)

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
