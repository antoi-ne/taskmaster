package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	flagName     string
	flagLifetime int
	flagExitCode int
	flagLogOut   bool
	flagLogErr   bool
	flagSignal   int
)

func init() {
	flag.StringVar(&flagName, "name", "default", "define the program name. used for logging.")
	flag.IntVar(&flagLifetime, "lifetime", -1, "seconds until the program exits.")
	flag.IntVar(&flagExitCode, "exitcode", 0, "program return value.")
	flag.BoolVar(&flagLogOut, "stdout", false, "prints output to stdout every second.")
	flag.BoolVar(&flagLogErr, "stderr", false, "prints output to stderr every second.")
	flag.IntVar(&flagSignal, "signal", -1, "will exit n seconds after receiving SIGUSR1. set to -1 to disable.")
}

func main() {
	flag.Parse()

	if flagSignal >= 0 {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGUSR1)

		go func() {
			<-sigs
			time.Sleep(time.Duration(flagSignal) * time.Second)
			os.Exit(flagExitCode)
		}()
	}

	if flagLifetime < 0 {
		infiniteLifetime()
	} else {
		finiteLifetime()
	}
}

func finiteLifetime() {
	ticker := time.NewTicker(time.Second)
	timer := time.NewTimer(time.Second * time.Duration(flagLifetime))

	for {
		select {
		case t := <-ticker.C:
			periodicLog(t)

		case <-timer.C:
			os.Exit(flagExitCode)
		}
	}
}

func infiniteLifetime() {
	for t := range time.Tick(time.Second) {
		periodicLog(t)
	}
}

func periodicLog(t time.Time) {
	if flagLogOut {
		fmt.Printf("[%s](%s) logging to stdout\n", flagName, t.Format(time.RFC822Z))
	}

	if flagLogErr {
		fmt.Fprintf(os.Stderr, "[%s](%s) logging to stderr\n", flagName, t.Format(time.RFC822Z))
	}
}
