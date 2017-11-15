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

	"google.golang.org/appengine"

	"github.com/Sirupsen/logrus"
	"github.com/jessfraz/upmail/email"
	"github.com/jessfraz/upmail/version"
	"github.com/sourcegraph/checkup"
)

const (
	// BANNER is what is printed for help/info output
	BANNER = "upmail - %s\n"
)

var (
	configFile string
	recipient  string
	interval   string

	ae bool

	mailgunAPIKey string
	mailgunDomain string

	mandrillAPIKey string

	smtpServer   string
	smtpSender   string
	smtpUsername string
	smtpPassword string

	debug bool
	vrsn  bool
)

func init() {
	// parse flags
	flag.StringVar(&configFile, "config", "checkup.json", "config file location")
	flag.StringVar(&recipient, "recipient", "", "recipient for email notifications")
	flag.StringVar(&interval, "interval", "10m", "check interval (ex. 5ms, 10s, 1m, 3h)")

	flag.BoolVar(&ae, "appengine", false, "enable the server for running in Google App Engine")

	flag.StringVar(&mailgunAPIKey, "mailgun", "", "Mailgun API Key to use for sending email (optional)")
	flag.StringVar(&mailgunDomain, "mailgun-domain", "", "Mailgun Domain to use for sending email (optional)")

	flag.StringVar(&mandrillAPIKey, "mandrill", "", "Mandrill API Key to use for sending email (optional)")

	flag.StringVar(&smtpServer, "server", "", "SMTP server for email notifications")
	flag.StringVar(&smtpSender, "sender", "", "SMTP default sender email address for email notifications")
	flag.StringVar(&smtpUsername, "username", "", "SMTP server username")
	flag.StringVar(&smtpPassword, "password", "", "SMTP server password")

	flag.BoolVar(&vrsn, "version", false, "print version and exit")
	flag.BoolVar(&vrsn, "v", false, "print version and exit (shorthand)")
	flag.BoolVar(&debug, "d", false, "run in debug mode")

	flag.Usage = func() {
		fmt.Fprint(os.Stderr, fmt.Sprintf(BANNER, version.VERSION))
		flag.PrintDefaults()
	}

	flag.Parse()

	if vrsn {
		fmt.Printf("upmail version %s, build %s", version.VERSION, version.GITCOMMIT)
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
	if smtpServer == "" && mandrillAPIKey == "" && mailgunAPIKey == "" && mailgunDomain == "" {
		usageAndExit("SMTP server OR Mailgun API Key OR  Mandrill API Key cannot be empty.", 1)
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
		MailgunAPIKey:  mailgunAPIKey,
		MailgunDomain:  mailgunDomain,
		MandrillAPIKey: mandrillAPIKey,
		Recipient:      recipient,
		Server:         smtpServer,
		Sender:         smtpSender,
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

	if ae {
		// setup necessary app engine health checks and listener
		go appengine.Main()
	}

	logrus.Info("Performing initial check")
	if err := c.CheckAndStore(); err != nil {
		logrus.Fatalf("check failed: %v", err)
	}

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
