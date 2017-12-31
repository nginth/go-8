package main

import (
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
)

// Sound - CHIP-8 sound device
type Sound struct {
	stream beep.StreamSeekCloser
}

func newSound() *Sound {
	f, err := os.Open("sound/beep.wav")
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
