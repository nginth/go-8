package main

import (
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
)

// SoundDevice - a generic sound device interfaces
type SoundDevice interface {
	playSound()
}

// Sound - SoundDevice implementation with the github.com/faiface/beep library
type Sound struct {
	stream beep.StreamSeekCloser
}

func newSound(filename string) *Sound {
	f, err := os.Open(filename)
	check(err)

	s, format, err := wav.Decode(f)
	check(err)

	speaker.Init(
		format.SampleRate,
		format.SampleRate.N(time.Second/10),
	)

	return &Sound{stream: s}
}

func (sound *Sound) playSound() {
	speaker.Play(beep.Seq(sound.stream))
	sound.stream.Seek(0)
}
