package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"

	"bytes"
	"encoding/json"
	"io/ioutil"

	"path"

	"github.com/pkg/errors"
	"github.com/thrawn01/args"
)

type SharedStruct struct {
	Metadata string
	Url      *url.URL
}

func main() {
	parser := args.NewParser(args.Name("http-client"),
		args.Desc("Example http client client"))

	parser.AddOption("--verbose").Alias("-v").Count().
		Help("Be verbose")
	parser.AddOption("--endpoint").Default("http://localhost:1234").Env("API_ENDPOINT").
		Help("The HTTP endpoint our client will talk too")

	parser.AddCommand("super-chickens", func(subParser *args.ArgParser, data interface{}) (int, error) {
		subParser.AddCommand("create", createChickens)
		subParser.AddCommand("list", listChickens)
		subParser.AddCommand("delete", deleteChickens)

		// Apply our super metadata =)
		shared := data.(*SharedStruct)
		shared.Metadata = "super"

		// Run the sub-commands
		return subParser.ParseAndRun(nil, data)
	})

	// Add our non super chicken actions
	parser.AddCommand("create", createChickens)
	parser.AddCommand("list", listChickens)
	parser.AddCommand("delete", deleteChickens)

	// Parse the command line
	opts, err := parser.Parse(nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	// Build our url from what the user passed in
	parts, err := url.Parse(opts.String("endpoint"))
	if err != nil {
		fmt.Fprint(os.Stderr, "url endpoint '%s' is invalid", opts.String("endpoint"), err.Error())
		os.Exit(1)
	}
	// This is a shows how you can pass in arbitrary struct
	// data to your sub-commands allows you to share data and resources
	// with sub-commands with out resorting to global variables
	super := &SharedStruct{Metadata: "", Url: parts}

	// Run the command chosen by the user
	retCode, err := parser.RunCommand(super)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	os.Exit(retCode)
}

func createChickens(subParser *args.ArgParser, data interface{}) (int, error) {
	subParser.AddArgument("name").Required().Help("The name of the chicken to create")
	opts := subParser.ParseSimple(nil)
	if opts == nil {
		return 1, nil
	}

	shared := data.(*SharedStruct)

	// Create the payload
	payload, err := json.Marshal(map[string]string{
		"name":     opts.String("name"),
		"metadata": shared.Metadata,
	})
	if err != nil {
		return 1, errors.Wrap(err, "while marshalling JSON")
	}

	// Create the new Request
	req, err := http.NewRequest("POST", joinUrl(shared, "chickens"), bytes.NewBuffer(payload))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1, errors.Wrap(err, "while creating new http request")
	}
	resp, err := sendRequest(opts, req, &payload)
	if err != nil {
		return 1, err
	}
	fmt.Println(resp)
	return 0, nil
}

func listChickens(subParser *args.ArgParser, data interface{}) (int, error) {
	opts := subParser.GetOpts()

	shared := data.(*SharedStruct)
	req, err := http.NewRequest("GET", joinUrl(shared, "chickens"), nil)
	if err != nil {
		return 1, errors.Wrap(err, "while creating new http request")
	}
	resp, err := sendRequest(opts, req, nil)
	if err != nil {
		return 1, err
	}
	fmt.Println(resp)
	return 0, nil
}

func deleteChickens(subParser *args.ArgParser, data interface{}) (int, error) {
	subParser.AddArgument("name").Required().Help("The name of the chicken to delete")
	opts := subParser.ParseSimple(nil)
	if opts == nil {
		return 1, nil
	}

	shared := data.(*SharedStruct)
	req, err := http.NewRequest("DELETE", joinUrl(shared, "chickens", opts.String("name")), nil)
	if err != nil {
		return 1, errors.Wrap(err, "while creating new http request")
	}
	resp, err := sendRequest(opts, req, nil)
	if err != nil {
		return 1, err
	}
	fmt.Println(resp)
	return 0, nil
}

func joinUrl(shared *SharedStruct, slugs ...string) string {
	uri := shared.Url
	results := []string{shared.Url.Path}
	for _, slug := range slugs {
		results = append(results, slug)
	}
	uri.Path = path.Join(slugs...)
	return uri.String()
}

func sendRequest(opts *args.Options, req *http.Request, payload *[]byte) (string, error) {
	req.Header.Set("Content-Type", "application/json")

	// Print a curl representation of the request if verbose
	curl := args.CurlString(req, payload)
	if opts.Bool("verbose") {
		fmt.Printf("-- %s\n", curl)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "Client Error")
	}
	defer resp.Body.Close()

	// Read in the entire response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Wrap(err, "ReadAll Error")
	}
	return string(body), nil
}
