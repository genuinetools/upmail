package main

import (
	"fmt"
	"os"

	"github.com/thrawn01/args"
)

type Config struct {
	PowerLevel  int
	Message     string
	StringSlice []string
	Verbose     int
	DbHost      string
	TheQuestion string
	TheAnswer   int
}

func main() {
	var conf Config

	// Create the parser with program name 'example'
	// and environment variables prefixed with APP_
	parser := args.NewParser(args.Name("demo"), args.EnvPrefix("APP_"),
		args.Desc("This is a demo app to showcase some features of args"))

	// Store Integers directly into a struct with a default value
	parser.AddOption("--power-level").Alias("-p").StoreInt(&conf.PowerLevel).
		Env("POWER_LEVEL").Default("10000").Help("set our power level")

	// Command line options can begin with -name, --name or even ++name
	// Most non word characters are supported
	parser.AddOption("++config-file").Alias("+c").IsString().
		Default("/path/to/config").Help("path to config file")

	// Use the args.Env() function to define an environment variable
	// NOTE: Since the parser was passed args.EnvPrefix("APP_") the actual
	// environment variable name is 'APP_MESSAGE'
	parser.AddOption("--message").Alias("-m").StoreStr(&conf.Message).
		Env("MESSAGE").Default("over-ten-thousand").Help("send a message")

	// Pass a comma separated list of strings and get a []string slice
	parser.AddOption("--slice").Alias("-s").StoreStringSlice(&conf.StringSlice).Env("LIST").
		Default("one,two,three").Help("list of messages")

	// Count the number of times an option is seen
	parser.AddOption("--verbose").Alias("-v").Count().StoreInt(&conf.Verbose).Help("be verbose")

	// Set bool to true if the option is present on the command line
	parser.AddOption("--debug").Alias("-d").IsTrue().Help("turn on Debug")

	// Specify the type of the arg with IsInt(), IsString(), IsBool() or IsTrue()
	parser.AddOption("--help").Alias("-h").IsTrue().Help("show this help message")

	// Add Required arguments
	parser.AddArgument("the-question").Required().
		StoreStr(&conf.TheQuestion).Help("Before you have an answer")

	// Add Optional arguments
	parser.AddArgument("the-answer").IsInt().Default("42").
		StoreInt(&conf.TheAnswer).Help("It must be 42")

	// 'Conf' options are not set via the command line but can be set
	// via a config file or an environment variable
	parser.AddConfig("twelve-factor").Env("TWELVE_FACTOR").Help("Demo of config options")

	// Define a 'database' subgroup
	db := parser.InGroup("database")

	// Add command line options to the subgroup
	db.AddOption("--host").Alias("-dH").StoreStr(&conf.DbHost).
		Default("localhost").Help("database hostname")

	// Add subgroup specific config. 'Conf' options are not set via the
	// command line but can be set via a config file or anything that calls parser.Apply()
	db.AddConfig("debug").IsTrue().Help("enable database debug")

	// 'Conf' option names are not allowed to start with a non word character like
	// '--' or '++' so they can not be confused with command line options
	db.AddConfig("database").IsString().Default("myDatabase").Help("name of database to use")

	// If no type is specified, defaults to 'IsString'
	db.AddConfig("user").Help("database user")
	db.AddConfig("pass").Help("database password")

	// Pass our own argument list, or nil to parse os.Args[]
	opt := parser.ParseOrExit(nil)

	// NOTE: ParseOrExit() is just a convenience, you can call
	// parser.Parse(nil) directly and handle the errors
	// yourself if you have more complicated use case

	// Demo default variables in a struct
	fmt.Printf("Power        '%d'\n", conf.PowerLevel)
	fmt.Printf("Message      '%s'\n", conf.Message)
	fmt.Printf("String Slice '%s'\n", conf.StringSlice)
	fmt.Printf("DbHost       '%s'\n", conf.DbHost)
	fmt.Printf("TheAnswer    '%d'\n", conf.TheAnswer)
	fmt.Println("")

	fmt.Println("")
	fmt.Println("==================")
	fmt.Println(" Direct Cast")
	fmt.Println("==================")

	// Fetch values by using the Cast functions
	fmt.Printf("Power               '%d'\n", opt.Int("power-level"))
	fmt.Printf("Message             '%s'\n", opt.String("message"))
	fmt.Printf("String Slice        '%s'\n", opt.StringSlice("slice"))
	fmt.Printf("Verbose             '%d'\n", opt.Int("verbose"))
	fmt.Printf("Debug               '%t'\n", opt.Bool("debug"))
	fmt.Printf("TheAnswer           '%d'\n", opt.Int("the-answer"))
	fmt.Printf("TheAnswer as String '%s'\n", opt.String("the-answer"))

	fmt.Println("")
	fmt.Println("==================")
	fmt.Println(" Database Group")
	fmt.Println("==================")

	// Fetch Group values
	dbAddOption := opt.Group("database")
	fmt.Printf("CAST DB Host   '%s'\n", dbAddOption.String("host"))
	fmt.Printf("CAST DB Debug  '%t'\n", dbAddOption.Bool("debug"))
	fmt.Printf("CAST DB User   '%s'\n", dbAddOption.String("user"))
	fmt.Printf("CAST DB Pass   '%s'\n", dbAddOption.String("pass"))

	fmt.Println("")

	iniFile := []byte(`
		power-level=20000
		message=OVER-TEN-THOUSAND!
		slice=three,four,five,six
		verbose=5
		debug=true

		[database]
		debug=false
		host=mysql.thrawn01.org
		user=my-username
		pass=my-password
	`)

	// Make configuration simple by reading arguments from an INI file
	opts, err := parser.FromINI(iniFile)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}

	fmt.Println("")
	fmt.Println("==================")
	fmt.Println("From INI file")
	fmt.Println("==================")

	// Values from the config file are used only if the argument is not present
	// on the commandline
	fmt.Printf("INI Power      '%d'\n", conf.PowerLevel)
	fmt.Printf("INI Message    '%s'\n", conf.Message)
	fmt.Printf("INI Slice      '%s'\n", conf.StringSlice)
	fmt.Printf("INI Verbose    '%d'\n", conf.Verbose)
	fmt.Printf("INI Debug      '%t'\n", opts.Bool("debug"))
	fmt.Println("")

	// Create an Options object from a map
	opts = parser.NewOptionsFromMap(
		map[string]interface{}{
			"int":    1,
			"bool":   true,
			"string": "one",
			"endpoints": map[string]interface{}{
				"endpoint1": "host1",
				"endpoint2": "host2",
				"endpoint3": "host3",
			},
			"deeply": map[string]interface{}{
				"nested": map[string]interface{}{
					"thing": "foo",
				},
				"foo": "bar",
			},
		},
	)

	// Small demo of how Group() works with options
	opts.String("string")                       // == "one"
	opts.Group("endpoints")                     // == *args.Options
	opts.Group("endpoints").String("endpoint1") // == "host1"
	opts.Group("endpoints").ToMap()             // map[string]interface{} {"endpoint1": "host1", ...}
	opts.StringMap("endpoints")                 // map[string]string {"endpoint1": "host1", ...}
	opts.KeySlice("endpoints")                  // [ "endpoint1", "endpoint2", ]
	opts.StringSlice("endpoints")               // [ "host1", "host2", "host3" ]

	// Leaves the door open for IntSlice(), IntMap(), etc....
}
