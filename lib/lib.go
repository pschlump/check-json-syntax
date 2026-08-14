package jsonSyntaxErroLib

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/pschlump/dbgo"
	"github.com/pschlump/json"
)

// HintType is a single hint on what might have cause the error
type HintType struct {
	Pattern string
	Note    string
	re      *regexp.Regexp
}

// HintList a set of possible, sometimes naive, hints on how to find/fix error
var HintList []HintType

// Debug if true will turn on extra debugging output
var Debug *bool

func init() {
	HintList = make([]HintType, 0, 25)
	HintList = append(HintList, HintType{Pattern: "invalid character '.' after object key:value", Note: "Check for missing comma(,) or colon(:) immediately preceding this"})
	HintList = append(HintList, HintType{Pattern: "invalid character '.' after object key$", Note: "Check for missing colon(:) betewen key and value"})
	HintList = append(HintList, HintType{Pattern: "unexpected end of JSON input", Note: "Check for missing end brace(}) or end array(])"})
	HintList = append(HintList, HintType{Pattern: "invalid character '\\\\'' looking for beginning of object key string", Note: "JSON only allows doublequotes(\")"})
	HintList = append(HintList, HintType{Pattern: "invalid character ',' looking for beginning of value", Note: "Check for two commas in a row, may be a missing value"})
	HintList = append(HintList, HintType{Pattern: "invalid character ']' looking for beginning of value", Note: "Check for an extra comma at end of array"})
	for ii, vv := range HintList {
		HintList[ii].re = regexp.MustCompile(vv.Pattern)
	}
	v := false
	Debug = &v
}

var hasTabs *regexp.Regexp

func init() {
	hasTabs = regexp.MustCompile("\t")
}

// GenerateSyntaxError converts from the offset error message into a human
// readable syntax error.  If err is not a *json.SyntaxError the error is
// returned formatted as a plain message.
func GenerateSyntaxError(js string, err error) (rv string) {

	var syntax *json.SyntaxError
	if !errors.As(err, &syntax) {
		return fmt.Sprintf("%s\n", err)
	}

	var sHint []string
	sErr := err.Error()
	for _, vv := range HintList {
		if vv.re.MatchString(sErr) {
			sHint = append(sHint, vv.Note)
		}
	}

	// Offset is the byte offset just past the offending byte.  Clamp it
	// defensively to the length of the input (which may be empty).
	offset := min(int(syntax.Offset), len(js))

	// Line that contains the offending byte, and its end (next newline).
	start := strings.LastIndex(js[:offset], "\n") + 1
	end := len(js)
	if idx := strings.Index(js[start:], "\n"); idx >= 0 {
		end = start + idx
	}

	line := strings.Count(js[:start], "\n") + 1
	pos := max(offset-start-1, 0)

	if *Debug {
		fmt.Printf("AT: %s - start=%d end=%d len(js)=%d\n", dbgo.LF(), start, end, len(js))
	}

	rv += fmt.Sprintf("Error in line %d: %s\n", line, err)
	rv += fmt.Sprintf("%s\n%s^\n", js[start:end], strings.Repeat(" ", pos))
	for _, vv := range sHint {
		rv += vv + "\n"
	}
	return
}

// CheckForTabs returns true if data has tabs in it
func CheckForTabs(data []byte) bool {
	return hasTabs.MatchString(string(data))
}

// TabListing shows tabs as '\\t' instead of a whitepace
func TabListing(data []byte) (rv string) {
	lineNo := 1
	lines := strings.Split(string(data), "\n")
	for _, s := range lines {
		if hasTabs.MatchString(s) {
			s = strings.ReplaceAll(s, "\t", "\\t")
			rv += fmt.Sprintf("%3d: %s\n", lineNo, s)
		}
		lineNo++
	}
	return
}
