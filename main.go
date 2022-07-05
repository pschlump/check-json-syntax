package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	jsonSyntaxErroLib "github.com/pschlump/check-json-syntax/lib"
	"github.com/pschlump/dbgo"
	"github.com/pschlump/json"
	"github.com/pschlump/jsondiff"
)

func printSyntaxError(js string, err error) {
	es := jsonSyntaxErroLib.GenerateSyntaxError(js, err)
	fmt.Printf("%s", es)
}

// Debug If true turns on debugging output
var Debug = flag.Bool("debug", false, "Debug flag") // 0

// GenListing shows line numbers in the listing
var GenListing = flag.Bool("list", false, "Add Line Numbers")                // 1
var IgnoreTabWarning = flag.Bool("ignore-tab-warning", false, "Ignore Tabs") // 1

// PrettyPrint JSON output - will print with tabs the JSON
var PrettyPrint = flag.Bool("pretty", false, "Add Line Numbers") // 2
// Add flag to check differences between two files.
var Diff = flag.Bool("diff", false, "Compare JSON files for differences") //
func init() {
	flag.BoolVar(Debug, "D", false, "Debug flag")                                     // 0
	flag.BoolVar(GenListing, "l", false, "Add Line Numbers")                          // 1
	flag.BoolVar(PrettyPrint, "p", false, "Prtty print JSON if syntatically correct") // 2
	flag.BoolVar(Diff, "d", false, "Compare JSON files for differences")              //
}

func main() {

	flag.Parse()
	fns := flag.Args()

	jsonSyntaxErroLib.Debug = Debug

	readInput := func(fn string) (rv []byte, err error) {
		if fn == "" {
			buf := bytes.NewBuffer(nil)
			io.Copy(buf, os.Stdin) // Error handling elided for brevity.
			rv = buf.Bytes()
		} else {
			rv, err = ioutil.ReadFile(fn)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Unable to open %s for input, Error:%s\n", fn, err)
			}
		}
		if false {
			bytes.Replace(rv, []byte("\t"), []byte(" "), -1)
		}
		return
	}

	exitVal := 0

	processData := func(fn string, data []byte) {
		if *GenListing {
			GenerateListing(data)
		}
		hasTabs := jsonSyntaxErroLib.CheckForTabs(data)
		if hasTabs {
			if !*IgnoreTabWarning {
				fmt.Printf("Warning: File contains tab characters - Go allows this but some JSON parsers will not allow this\n%s", jsonSyntaxErroLib.TabListing(data))
			}
		}
		isvv, isww, ismm := false, false, false
		var vv map[string]interface{}
		var ww []map[string]interface{}
		var mm []interface{}
		if *Debug {
			fmt.Printf("AT: %s\n", dbgo.LF())
		}
		// Try a hash of name and values first
		err := json.Unmarshal([]byte(data), &vv)
		isvv = (err == nil)
		if err != nil {
			if *Debug {
				fmt.Printf("AT: %s\n", dbgo.LF())
			}
			err = nil
			// Try an array of hash of name and values first
			err = json.Unmarshal([]byte(data), &ww)
			isww = (err == nil)
		}
		if err != nil {
			if *Debug {
				fmt.Printf("AT: %s\n", dbgo.LF())
			}
			err = nil
			// Try an array of values
			err = json.Unmarshal([]byte(data), &mm)
			ismm = (err == nil)
		}
		if *Debug {
			dbgo.Printf("AT: %(LF), isvv=%v isww=%v ismm=%v\n", isvv, isww, ismm)
		}
		if err != nil {
			printSyntaxError(string(data), err)
			exitVal = 1
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

	// if len(fns) == 0 {
	// 	fmt.Fprintf(os.Stderr, "Usage: Must list files on command line to check\n")
	// 	flag.Usage()
	// 	os.Exit(1)
	// }

	if *Diff {
		if len(fns) != 2 {
			fmt.Fprintf(os.Stderr, "Usage: Must have 2 files to compare\n")
			flag.Usage()
			os.Exit(1)
		}

		d := jsondiff.CompareFiles(fns[0], fns[1])
		fmt.Printf("Differences\n%s\n", jsondiff.Format(d))
		if d.HasDiff {
			os.Exit(1)
		}
		os.Exit(0)
	}

	if len(fns) == 0 {
		data, _ := readInput("")
		processData("--stdin--", data)
	} else {
		for _, fn := range fns {
			data, err := readInput(fn)
			if err == nil {
				processData(fn, data)
			}
		}
	}

	os.Exit(exitVal)
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
