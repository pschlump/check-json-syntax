package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/pschlump/godebug"
	"github.com/pschlump/json" //	"encoding/json"

	jsonSyntaxErroLib "github.com/pschlump/check-json-syntax/lib"
)

func printSyntaxError(js string, err error) {
	es := jsonSyntaxErroLib.GenerateSyntaxError(js, err)
	fmt.Printf("%s", es)
}

// Debug If true turns on debugging output
var Debug = flag.Bool("debug", false, "Debug flag") // 0

// GenListing shows line numbers in the listing
var GenListing = flag.Bool("list", false, "Add Line Numbers") // 1

// PrettyPrint JSON output - will print with tabs the JSON
var PrettyPrint = flag.Bool("pretty", false, "Add Line Numbers") // 2
func init() {
	flag.BoolVar(Debug, "D", false, "Debug flag")                                     // 0
	flag.BoolVar(GenListing, "l", false, "Add Line Numbers")                          // 1
	flag.BoolVar(PrettyPrint, "p", false, "Prtty print JSON if syntatically correct") // 2
}

func main() {

	flag.Parse()
	fns := flag.Args()

	jsonSyntaxErroLib.Debug = Debug

	if len(fns) == 0 {
		fmt.Fprintf(os.Stderr, "Usage: Must list files on command line to check\n")
		flag.Usage()
		os.Exit(1)
	}

	for _, fn := range fns {

		data, err := ioutil.ReadFile(fn)
		if *GenListing {
			GenerateListing(data)
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to open %s for input, Error:%s\n", fn, err)
		} else {
			hasTabs := jsonSyntaxErroLib.CheckForTabs(data)
			if hasTabs {
				fmt.Printf("Warning: File contains tab characters - Go allows this but some JSON parsers will not allow this\n%s", jsonSyntaxErroLib.TabListing(data))
			}
			isvv, isww, ismm := false, false, false
			var vv map[string]interface{}
			var ww []map[string]interface{}
			var mm []interface{}
			if *Debug {
				fmt.Printf("AT: %s\n", godebug.LF())
			}
			// Try a hash of name and values first
			err = json.Unmarshal([]byte(data), &vv)
			isvv = (err == nil)
			if err != nil {
				if *Debug {
					fmt.Printf("AT: %s\n", godebug.LF())
				}
				err = nil
				// Try an array of hash of name and values first
				err = json.Unmarshal([]byte(data), &ww)
				isww = (err == nil)
			}
			if err != nil {
				if *Debug {
					fmt.Printf("AT: %s\n", godebug.LF())
				}
				err = nil
				// Try an array of values
				err = json.Unmarshal([]byte(data), &mm)
				ismm = (err == nil)
			}
			if *Debug {
				fmt.Printf("AT: %s, isvv=%v isww=%v ismm=%v\n", godebug.LF(), isvv, isww, ismm)
			}
			if err != nil {
				printSyntaxError(string(data), err)
			} else if *PrettyPrint {
				var s []byte
				if isvv {
					s, err = json.MarshalIndent(vv, "", "\t")
				} else if isww {
					s, err = json.MarshalIndent(ww, "", "\t")
				} else if ismm {
					s, err = json.MarshalIndent(mm, "", "\t")
				} else {
					s = data
				}
				if err != nil {
					s = data
				}
				fmt.Printf("%s\n", s)
			} else {
				fmt.Printf("%s: Syntax OK\n", fn)
			}
		}

	}
}

// GenerateListing print the listing to standard output
func GenerateListing(data []byte) {
	lineNo := 1
	lines := strings.Split(string(data), "\n")
	for _, s := range lines {
		fmt.Printf("%3d: %s\n", lineNo, s)
		lineNo++
	}
}
