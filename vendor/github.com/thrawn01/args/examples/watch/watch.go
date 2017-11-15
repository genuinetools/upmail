package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"time"

	"github.com/thrawn01/args"
)

func main() {

	parser := args.NewParser(args.Name("watch"))
	parser.AddOption("--bind").Alias("-b").Default("localhost:8080").
		Help("Interface to bind the server too")
	parser.AddOption("--complex-example").Alias("-ce").IsBool().
		Help("Run the more complex example")
	parser.AddOption("--config-file").Alias("-c").
		Help("The Config file to load and watch our config from")

	// Add a connection string to the database group
	parser.AddOption("--connection-string").InGroup("database").Alias("-cS").
		Default("mysql://username@hostname:MyDB").
		Help("Connection string used to connect to the database")

	// Store the password in the config and not passed via the command line
	parser.AddConfig("password").InGroup("database").Help("database password")
	// Specify a config file version, when this version number is updated, the user is signaling to the application
	// that all edits are complete and the application can reload the config
	parser.AddConfig("version").IsInt().Default("0").Help("config file version")

	appConf, err := parser.Parse(nil)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}

	// Simple handler that prints out our config information
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		conf := appConf.GetOpts()

		db := conf.Group("database")
		payload, err := json.Marshal(map[string]string{
			"bind":     conf.String("bind"),
			"mysql":    db.String("connection-string"),
			"password": conf.Group("database").String("password"),
		})
		if err != nil {
			fmt.Println("error:", err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(payload)
	})

	var cancelWatch args.WatchCancelFunc
	if appConf.Bool("complex-example") {
		cancelWatch = complex(parser)
	} else {
		cancelWatch = simple(parser)
	}

	// Shut down the watcher when done
	defer cancelWatch()

	// Listen and serve requests
	log.Printf("Listening for requests on %s", appConf.String("bind"))
	err = http.ListenAndServe(appConf.String("bind"), nil)
	if err != nil {
		log.Println(err)
		os.Exit(-1)
	}
}

// Simple example always updates the config when file changes are detected.
func simple(parser *args.ArgParser) args.WatchCancelFunc {
	// Get our current config
	appConf := parser.GetOpts()
	configFile := appConf.String("config-file")

	// Watch the file every time.Second and call func(err error){} when the file is modified
	cancelWatch, err := args.WatchFile(configFile, time.Second, func(err error) {
		if err != nil {
			fmt.Printf("Watch Error %s\n", err.Error())
			return
		}

		// You can safely ignore the returned Options{} object here.
		// the next call to GetOpts() from within the handler will
		// pick up the newly parsed config
		appConf, err = parser.FromINIFile(configFile)
		if err != nil {
			fmt.Printf("Failed to load updated config - %s\n", err.Error())
			return
		}
	})

	if err != nil {
		fmt.Printf("Unable to start watch '%s' -  %s", configFile, err.Error())
	}
	return cancelWatch
}

// The complex example allows a user to write the config file multiple times, possibly applying edits incrementally.
// When the user is ready for the application to apply the config changes, modify the 'version' value and the
// new config is applied.
func complex(parser *args.ArgParser) args.WatchCancelFunc {
	// Get our current config
	appConf := parser.GetOpts()
	configFile := appConf.String("config-file")

	// Watch the file every time.Second and call func(err error){} when the file is modified
	cancelWatch, err := args.WatchFile(configFile, time.Second, func(err error) {
		if err != nil {
			fmt.Printf("Watch Error %s\n", err.Error())
			return
		}

		// load the file from disk
		content, err := args.LoadFile(configFile)
		if err != nil {
			fmt.Printf("Failed to load config - %s\n", err.Error())
		}

		// Parse the file contents
		newConfig, err := parser.ParseINI(content)
		if err != nil {
			fmt.Printf("Failed to update config - %s\n", err.Error())
			return
		}

		// Only "Apply" the newConfig when the version changes
		if appConf.Int("version") != newConfig.Int("version") {
			// Apply the newConfig values to the parser rules
			appConf, err = parser.Apply(newConfig)
			if err != nil {
				fmt.Printf("Probably a type cast error - %s\n", err.Error())
				return
			}
		}
	})
	if err != nil {
		log.Fatal(fmt.Sprintf("Failed to watch '%s' -  %s", configFile, err.Error()))
	}
	return cancelWatch
}
