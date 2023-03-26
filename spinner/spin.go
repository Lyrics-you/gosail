package spinner

import (
	"fmt"
	"gosail/model"
	"time"

	"github.com/briandowns/spinner"
	"github.com/schollz/progressbar/v3"
)

type Spinner chan interface{}
type style uint8

var (
	tips                 = "gosail in progress ..."
	timeout              = 30
	spinnerStyle         = Bar
	Spin         Spinner = make(chan interface{}, 1)
	isSelection          = false
)

const (
	Tips style = 0
	Move style = 1
	Bar  style = 2
	// Progressor style = 2
)

func useTips(done Spinner) {
	fmt.Println(tips)
	<-done
}

func useBar(done Spinner) {
	bar := progressbar.NewOptions(timeout,
		// progressbar.OptionSetWriter(os.Stdout),
		// progressbar.OptionEnableColorCodes(true),
		// progressbar.OptionShowBytes(true),
		// progressbar.OptionSetWidth(15),
		// progressbar.OptionSetDescription("[cyan][1/3][reset] Writing moshable file..."),
		progressbar.OptionSetDescription("gosail in progress"),
		progressbar.OptionSetRenderBlankState(true),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "=",
			SaucerHead:    ">",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}))

	for i := 1; i <= timeout; i++ {
		select {
		case <-done:
			// bar.Set(i)
			// bar.Finish()
			bar.Clear()
			bar.Exit()
			fmt.Println()
			return
		default:
			time.Sleep(time.Second)
			bar.Add(1)
		}
	}
	// bar.Finish()
	bar.Clear()
	bar.Exit()
	fmt.Println()
}

func userMove(done Spinner) {
	// Build our new spinner
	// spinner.CharSets[12]
	newCharsSets := []string{">", "'>", ")'>", "))'>", ">))'>", " >))'>", "  >))'>", "   >))'>", "    >))'>", "   <'((<", "  <'((<", " <'((<", "<'((<", "'((<", "((<", "(<", "<"}
	s := spinner.New(newCharsSets, 200*time.Millisecond)
	// Start the spinner
	s.Start()
	<-done
	s.Stop()
}

func (ch *Spinner) spinnerRun() {
	switch spinnerStyle {
	case Move:
		userMove(*ch)
	case Bar:
		useBar(*ch)
	case Tips:
		useTips(*ch)
	default:
		useTips(*ch)
	}
}

func (ch *Spinner) SetStyle(s style) {
	spinnerStyle = s
}

func (ch *Spinner) SetTips(tip string) {
	if tip != "" {
		tips = tip
	}
}

func (ch *Spinner) SetTimeOut(t int) {
	timeout = t
}

func (ch *Spinner) SetIsSelection(selection bool) {
	isSelection = selection
}

func (ch *Spinner) Init(config *model.SpinConfig) {
	if config == nil {
		ch.SetStyle(Tips)
		ch.SetTips("")
	} else {
		ch.SetStyle(style(config.SpinType))
		ch.SetTips(config.SpinTips)
		ch.SetTimeOut(config.TimeOut)
		ch.SetIsSelection(config.IsSelect)
	}
}

func (ch *Spinner) Start() {
	go ch.spinnerRun()
}

func (ch *Spinner) Stop() {
	// done signal
	*ch <- struct{}{}

	if spinnerStyle == Bar && isSelection {
		time.Sleep(time.Second)
	} else {
		*ch <- struct{}{}
		<-*ch
	}
	// close(*ch)
}
