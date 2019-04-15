package barullo

import (
	"math"
	"strings"
)

type KeyStatus int

const (
	NoteOff KeyStatus = iota
	NotePress
	NotePressed
	NoteRelease
)

var notes = map[string]float64{
	"C":  16.35,
	"C#": 17.32,
	"D":  18.35,
	"D#": 19.45,
	"E":  20.6,
	"F":  21.83,
	"F#": 23.12,
	"G":  24.50,
	"G#": 25.96,
	"A":  27.50,
	"A#": 29.14,
	"B":  30.87,
}

type Event struct {
	Offset int
	Note   string
	Octave int
	Key    KeyStatus
}

func (e *Event) Frequency() float64 {
	freq, ok := notes[strings.ToUpper(e.Note)]
	if !ok {
		return 0.0
	}

	if e.Octave < 1 {
		return freq
	}

	return math.Pow(freq, float64(e.Octave))
}

type Sequence struct {
	note   string
	octave int
	key    KeyStatus

	events []Event
	length int
	pos    int
}

func NewSequence(length int, events []Event) *Sequence {
	return &Sequence{
		key:    NoteOff,
		events: events,
		length: length,
	}
}

func (s *Sequence) Get(offset int) Event {
	offset = offset % s.length
	now := s.events[s.pos]
	prev := s.events[(len(s.events)+s.pos-1)%len(s.events)] // warp around

	if offset < now.Offset || (s.pos == 0 && offset > prev.Offset) {
		return Event{
			Offset: offset,
			Note:   s.note,
			Octave: s.octave,
			Key:    s.key,
		}
	}

	event := Event{
		Offset: offset,
		Note:   now.Note,
		Octave: now.Octave,
		Key:    now.Key,
	}

	s.note = now.Note
	s.octave = now.Octave
	s.key = NoteOff

	if now.Key == NotePress || now.Key == NotePressed {
		s.key = NotePressed
	}

	// do not increment possition if we are in the event offest
	if offset == now.Offset {
		return event
	}
	s.pos = (s.pos + 1) % len(s.events)

	switch event.Key {
	case NoteOff:
		println("NoteOff")
	case NotePress:
		println("NotePress")
	case NotePressed:
		println("NotePressed")
	case NoteRelease:
		println("NoteRelease")
	}

	return event
}
