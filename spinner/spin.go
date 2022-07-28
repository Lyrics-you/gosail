package spinner

import (
	"time"

	"github.com/briandowns/spinner"
)

type Spinner chan interface{}

var Spin Spinner = make(chan interface{}, 1)

func (ch *Spinner) spinnerRun() {
	// Build our new spinner
	// spinner.CharSets[12]
	newCharsSets := []string{">", "'>", ")'>", "))'>", ">))'>", " >))'>", "  >))'>", "   >))'>", "    >))'>", "   <'((<", "  <'((<", " <'((<", "<'((<", "'((<", "((<", "(<", "<"}
	s := spinner.New(newCharsSets, 200*time.Millisecond)
	// Start the spinner
	s.Start()
	defer s.Stop()
	<-*ch
	*ch <- struct{}{}
}

func (ch *Spinner) Start() {
	go ch.spinnerRun()
}

func (ch *Spinner) Stop() {
	*ch <- struct{}{}
	<-*ch
}
