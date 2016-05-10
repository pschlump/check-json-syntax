package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/pschlump/godebug"
)

type HintType struct {
	Pattern string
	Note    string
	re      *regexp.Regexp
}

var HintList []HintType
var hasTabs *regexp.Regexp

func init() {
	HintList = make([]HintType, 0, 25)
	HintList = append(HintList, HintType{Pattern: "invalid character '.' after object key:value", Note: "Check for missing comma(,) or colon(:) immediately preceding this"})
	HintList = append(HintList, HintType{Pattern: "invalid character '.' after object key$", Note: "Check for missing colon(:) betewen key and value"})
	HintList = append(HintList, HintType{Pattern: "unexpected end of JSON input", Note: "Check for missing end brace(}) or end array(])"})
	HintList = append(HintList, HintType{Pattern: "invalid character '\\\\'' looking for beginning of object key string", Note: "JSON only allows doublequotes(\")"})
	for ii, vv := range HintList {
		HintList[ii].re = regexp.MustCompile(vv.Pattern)
	}
	hasTabs = regexp.MustCompile("\t")
}

func printSyntaxError(js string, err error) {

	max := func(a, b int) int {
		if a < b {
			return b
		}
		return a
	}
	min := func(a, b int) int {
		if a < b {
			return a
		}
		return b
	}

	syntax, ok := err.(*json.SyntaxError)
	if !ok {
		fmt.Printf("%s\n", err)
		return
	}

	var sHint []string
	sErr := fmt.Sprintf("%s", err)
	for _, vv := range HintList {
		if vv.re.MatchString(sErr) {
			sHint = append(sHint, vv.Note)
		}
	}

	start, end := strings.LastIndex(js[:syntax.Offset], "\n")+1, len(js)
	start = min(max(start, 0), len(js)-1)
	end = min(max(end, 1), len(js))
	if idx := strings.Index(js[start:], "\n"); idx >= 0 {
		end = start + idx
	}

	if *Debug {
		fmt.Printf("AT: %s - start=%d end=%d len(js)=%d\n", godebug.LF(), start, end, len(js))
	}

	line, pos := strings.Count(js[:start], "\n"), int(syntax.Offset)-start-1
	pos = max(pos, 0)

	fmt.Printf("Error in line %d: %s \n", line, err)
	fmt.Printf("%s\n%s^\n", js[start:end], strings.Repeat(" ", pos))
	for _, vv := range sHint {
		fmt.Printf("%s\n", vv)
	}
}

var Debug = flag.Bool("debug", false, "Debug flag")           // 0
var GenListing = flag.Bool("list", false, "Add Line Numbers") // 2
func init() {
	flag.BoolVar(Debug, "D", false, "Debug flag")            // 0
	flag.BoolVar(GenListing, "l", false, "Add Line Numbers") // 2
}

func main() {

	flag.Parse()
	fns := flag.Args()

	if len(fns) == 0 {
		fmt.Fprintf(os.Stderr, "Usage: Must list files on command line to check\n")
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
			hasTabs := CheckForTabs(data)
			if hasTabs {
				fmt.Printf("Warning: File contains tab characters - Go allows this but some JSON parsers will not allow this\n")
				TabListing(data)
			}
			var vv map[string]interface{}
			var ww []map[string]interface{}
			var mm []interface{}
			if *Debug {
				fmt.Printf("AT: %s\n", godebug.LF())
			}
			// Try a hash of name and values first
			err = json.Unmarshal([]byte(data), &vv)
			if err != nil {
				if *Debug {
					fmt.Printf("AT: %s\n", godebug.LF())
				}
				err = nil
				// Try an array of hash of name and values first
				err = json.Unmarshal([]byte(data), &ww)
			}
			if err != nil {
				if *Debug {
					fmt.Printf("AT: %s\n", godebug.LF())
				}
				err = nil
				// Try an array of values
				err = json.Unmarshal([]byte(data), &mm)
			}
			if *Debug {
				fmt.Printf("AT: %s\n", godebug.LF())
			}
			if err != nil {
				printSyntaxError(string(data), err)
			} else {
				fmt.Printf("%s: Syntax OK\n", fn)
			}
		}

	}
}

func GenerateListing(data []byte) {
	line_no := 1
	lines := strings.Split(string(data), "\n")
	for _, s := range lines {
		fmt.Printf("%3d: %s\n", line_no, s)
		line_no++
	}
}

func CheckForTabs(data []byte) bool {
	if hasTabs.MatchString(string(data)) {
		return true
	}
	return false
}

func TabListing(data []byte) {
	line_no := 1
	lines := strings.Split(string(data), "\n")
	for _, s := range lines {
		if hasTabs.MatchString(s) {
			s = strings.Replace(s, "\t", "\\t", -1)
			fmt.Printf("%3d: %s\n", line_no, s)
		}
		line_no++
	}
}
