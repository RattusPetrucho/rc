package main

import (
	"log"

	"github.com/RattusPetrucho/rc/ui"
)

func main() {
	win, err := ui.New()
	if err != nil {
		log.Fatal(err)
	}
	defer win.Close()

	if err := win.MainLoop(); err != nil && err != ui.ErrQuit {
		log.Panicln(err)
	}
}
