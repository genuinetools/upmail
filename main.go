package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/smtp"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/jfrazelle/upmail/email"
	"github.com/sourcegraph/checkup"
)

const (
	// BANNER is what is printed for help/info output
	BANNER = "upmail - %s\n"
	// VERSION is the binary version.
	VERSION = "v0.1.0"
)

var (
	configFile string
	recipient  string
	interval   string

	smtpServer   string
	smtpSender   string
	smtpUsername string
	smtpPassword string

	debug   bool
	version bool
)

func init() {
	// parse flags
	flag.StringVar(&configFile, "config", "checkup.json", "config file location")
	flag.StringVar(&recipient, "recipient", "", "recipient for email notifications")
	flag.StringVar(&interval, "interval", "10m", "check interval (ex. 5ms, 10s, 1m, 3h)")

	flag.StringVar(&smtpServer, "server", "", "SMTP server for email notifications")
	flag.StringVar(&smtpSender, "sender", "", "SMTP default sender email address for email notifications")
	flag.StringVar(&smtpUsername, "username", "", "SMTP server username")
	flag.StringVar(&smtpPassword, "password", "", "SMTP server password")

	flag.BoolVar(&version, "version", false, "print version and exit")
	flag.BoolVar(&version, "v", false, "print version and exit (shorthand)")
	flag.BoolVar(&debug, "d", false, "run in debug mode")

	flag.Usage = func() {
		fmt.Fprint(os.Stderr, fmt.Sprintf(BANNER, VERSION))
		flag.PrintDefaults()
	}

	flag.Parse()

	if version {
		fmt.Printf("%s", VERSION)
		os.Exit(0)
	}

	// set log level
	if debug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	if configFile == "" {
		usageAndExit("Config file cannot be empty.", 1)
	}
	if recipient == "" {
		usageAndExit("Recipient cannot be empty.", 1)
	}
	if smtpServer == "" {
		usageAndExit("SMTP server cannot be empty.", 1)
	}
}

func main() {
	var ticker *time.Ticker
	// On ^C, or SIGTERM handle exit.
	s := make(chan os.Signal, 1)
	signal.Notify(s, os.Interrupt)
	signal.Notify(s, syscall.SIGTERM)
	go func() {
		for sig := range s {
			ticker.Stop()
			logrus.Infof("Received %s, exiting.", sig.String())
			os.Exit(0)
		}
	}()

	configBytes, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatal(err)
	}

	var c checkup.Checkup
	err = json.Unmarshal(configBytes, &c)
	if err != nil {
		log.Fatal(err)
	}

	n := email.Notifier{
		Recipient: recipient,
		Server:    smtpServer,
		Sender:    smtpSender,
		Auth: smtp.PlainAuth(
			"",
			smtpUsername,
			smtpPassword,
			strings.SplitN(smtpServer, ":", 2)[0],
		),
	}
	c.Notifier = n

	// parse the duration
	dur, err := time.ParseDuration(interval)
	if err != nil {
		logrus.Fatalf("parsing %s as duration failed: %v", interval, err)
	}

	logrus.Infof("Starting checks that will send emails to: %s", recipient)
	ticker = time.NewTicker(dur)

	for range ticker.C {
		if err := c.CheckAndStore(); err != nil {
			logrus.Warnf("check failed: %v", err)
		}
	}
}

func usageAndExit(message string, exitCode int) {
	if message != "" {
		fmt.Fprintf(os.Stderr, message)
		fmt.Fprintf(os.Stderr, "\n\n")
	}
	flag.Usage()
	fmt.Fprintf(os.Stderr, "\n")
	os.Exit(exitCode)
}
