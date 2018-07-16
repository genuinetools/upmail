package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/smtp"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"google.golang.org/appengine"

	"github.com/genuinetools/pkg/cli"
	"github.com/genuinetools/upmail/email"
	"github.com/genuinetools/upmail/version"
	"github.com/sirupsen/logrus"
	"github.com/sourcegraph/checkup"
)

var (
	configFile string
	recipient  string
	interval   time.Duration

	ae bool

	mailgunAPIKey string
	mailgunDomain string

	smtpServer   string
	smtpSender   string
	smtpUsername string
	smtpPassword string

	debug bool
)

func main() {
	// Create a new cli program.
	p := cli.NewProgram()
	p.Name = "upmail"
	p.Description = "Email notification hook for https://github.com/sourcegraph/checkup"

	// Set the GitCommit and Version.
	p.GitCommit = version.GITCOMMIT
	p.Version = version.VERSION

	// Setup the global flags.
	p.FlagSet = flag.NewFlagSet("global", flag.ExitOnError)
	p.FlagSet.StringVar(&configFile, "config", "checkup.json", "config file location")
	p.FlagSet.StringVar(&recipient, "recipient", "", "recipient for email notifications")
	p.FlagSet.DurationVar(&interval, "interval", 10*time.Minute, "check interval (ex. 5ms, 10s, 1m, 3h)")

	p.FlagSet.BoolVar(&ae, "appengine", false, "enable the server for running in Google App Engine")

	p.FlagSet.StringVar(&mailgunAPIKey, "mailgun", "", "Mailgun API Key to use for sending email (optional)")
	p.FlagSet.StringVar(&mailgunDomain, "mailgun-domain", "", "Mailgun Domain to use for sending email (optional)")

	p.FlagSet.StringVar(&smtpServer, "server", "", "SMTP server for email notifications")
	p.FlagSet.StringVar(&smtpSender, "sender", "", "SMTP default sender email address for email notifications")
	p.FlagSet.StringVar(&smtpUsername, "username", "", "SMTP server username")
	p.FlagSet.StringVar(&smtpPassword, "password", "", "SMTP server password")

	p.FlagSet.BoolVar(&debug, "d", false, "enable debug logging")

	// Set the before function.
	p.Before = func(ctx context.Context) error {
		// Set the log level.
		if debug {
			logrus.SetLevel(logrus.DebugLevel)
		}

		if len(configFile) < 1 {
			return fmt.Errorf("Config file cannot be empty")
		}
		if len(recipient) < 1 {
			return fmt.Errorf("Recipient cannot be empty")
		}
		if len(smtpServer) < 1 && len(mailgunAPIKey) < 1 && len(mailgunDomain) < 1 {
			return fmt.Errorf("SMTP server OR Mailgun API Key cannot be empty")
		}

		return nil
	}

	// Set the main program action.
	p.Action = func(ctx context.Context, args []string) error {
		ticker := time.NewTicker(interval)

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
			logrus.Fatal(err)
		}

		var c checkup.Checkup
		err = json.Unmarshal(configBytes, &c)
		if err != nil {
			logrus.Fatal(err)
		}

		n := email.Notifier{
			MailgunAPIKey: mailgunAPIKey,
			MailgunDomain: mailgunDomain,
			Recipient:     recipient,
			Server:        smtpServer,
			Sender:        smtpSender,
			Auth: smtp.PlainAuth(
				"",
				smtpUsername,
				smtpPassword,
				strings.SplitN(smtpServer, ":", 2)[0],
			),
		}
		c.Notifier = n

		logrus.Infof("Starting checks that will send emails to: %s", recipient)

		if ae {
			// setup necessary app engine health checks and listener
			go appengine.Main()
		}

		logrus.Info("Performing initial check")
		if err := c.CheckAndStore(); err != nil {
			logrus.Fatalf("CheckAndStore failed: %v", err)
		}

		for range ticker.C {
			if err := c.CheckAndStore(); err != nil {
				logrus.Warnf("CheckAndStore failed: %v", err)
			}
		}

		return nil
	}

	// Run our program.
	p.Run()
}
