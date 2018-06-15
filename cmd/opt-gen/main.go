package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"gopkg.in/mab-go/opt.v0/cmd/opt-gen/generator"
)

var (
	flagOut = flag.String("out", "", "The output `file` where results should be written. If omitted, results are\n    \twritten to 'opt_t.go', where 't' is the first type specified.")
)

func usage() {
	_, err := fmt.Fprintln(os.Stderr, `usage: opt-gen [flags] config_file

arguments:
  config_file
    	The file system path to the config file.

flags:`)
	if err != nil {
		panic(err)
	}

	flag.PrintDefaults()
	_, err = fmt.Fprintln(os.Stderr, `
examples:
  opt-gen config.json
  opt-gen -out=string.go config-string.json
  opt-gen -out=string.go -package=optstring config-string.json`)
	if err != nil {
		panic(err)
	}
}

func fail(msgFmt string, v ...interface{}) {
	_, err := fmt.Fprintf(os.Stderr, msgFmt+"\n", v...)
	if err != nil {
		panic(err)
	}
	os.Exit(1)
}

func getOutputNames(n string) (source string, test string) {
	source = n
	if source == "" {
		source = "opt_gen.go"
	}

	base := strings.TrimSuffix(source, ".go")
	test = base + "_test.go"

	return
}

func main() {
	flag.Usage = usage
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		usage()
		os.Exit(1)
	}

	configPath := args[0]
	config := generator.Config{}
	if err := config.LoadJSONFile(configPath); err != nil {
		fail("Could not read config file. Cause: \"%+v\".", err)
	}

	result, genErr := generator.Generate(config)
	if genErr != nil {
		log.Fatalf("ERROR: failed to generate source code: %v", genErr)
	}

	outSource, outTest := getOutputNames(*flagOut)
	if err := ioutil.WriteFile(outSource, result.Source, 0644); err != nil {
		fail("failed to write output: %s", err)
	}

	if err := ioutil.WriteFile(outTest, result.Test, 0644); err != nil {
		fail("failed to write output: %s", err)
	}
}
