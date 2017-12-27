package main

import (
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
)

func playSound() {
	f, err := os.Open("sound/beep.wav")
	check(err)

	s, format, err := wav.Decode(f)
	check(err)

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	done := make(chan struct{})

	speaker.Play(beep.Seq(s, beep.Callback(func() {
		//time.Sleep(time.Second)
		close(done)
	})))

	<-done
}
