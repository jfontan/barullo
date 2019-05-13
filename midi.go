package barullo

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"regexp"
	"strconv"
)

type MidiJson struct {
	Header struct {
		Name string `json:"name"`
		Ppq  int    `json:"ppq"`
		Meta []struct {
			Type  string `json:"type"`
			Ticks int    `json:"ticks"`
		} `json:"meta"`
		Tempos []struct {
			Ticks int     `json:"ticks"`
			Bpm   float64 `json:"bpm"`
		} `json:"tempos"`
		TimeSignatures []struct {
			Ticks         int   `json:"ticks"`
			TimeSignature []int `json:"timeSignature"`
			Measures      int   `json:"measures"`
		} `json:"timeSignatures"`
	} `json:"header"`
	Tracks []struct {
		Name       string `json:"name"`
		Channel    int    `json:"channel"`
		Instrument struct {
			Number int    `json:"number"`
			Name   string `json:"name"`
			Family string `json:"family"`
		} `json:"instrument"`
		Notes []struct {
			Time          float64 `json:"time"`
			Midi          int     `json:"midi"`
			Name          string  `json:"name"`
			Velocity      float64 `json:"velocity"`
			Duration      float64 `json:"duration"`
			Ticks         int     `json:"ticks"`
			DurationTicks int     `json:"durationTicks"`
		} `json:"notes"`
	} `json:"tracks"`
}

func GetEventsFromMidi(channel, sampleRate int, r io.Reader) []Event {
	reg := regexp.MustCompile(`([ABCDEFG]#?)(\d)`)
	var events []Event

	value := &MidiJson{}
	data, err := ioutil.ReadAll(r)
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(data, value); err != nil {
		panic(err)
	}

	track := value.Tracks[channel]
	fmt.Println("TRACK", track.Instrument.Name)
	for _, note := range track.Notes {
		noteAndOctave := reg.FindStringSubmatch(note.Name)
		stringNote := noteAndOctave[1]
		octave := 0
		if len(noteAndOctave) > 2 {
			octave, err = strconv.Atoi(noteAndOctave[2])
			if err != nil {
				panic(err)
			}
		}

		e := Event{
			Offset: int(note.Time * float64(sampleRate)),
			Note:   stringNote,
			Octave: octave,
			Key:    NotePress,
		}
		events = append(events, e)

		e = Event{
			Offset: int((note.Time + note.Duration) * float64(sampleRate)),
			Note:   stringNote,
			Octave: octave,
			Key:    NoteRelease,
		}
		events = append(events, e)
	}
	fmt.Println("EVTS", len(events))
	return events
}
