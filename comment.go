package main

import (
	bufio "bufio"
	"bytes"
	fmt "fmt"
	"io"

	os "os"
)

type StateType int

const (
	S0  StateType = 0 // initial state
	S1            = 1 // Escape symbol inside string
	S2            = 2 // first slash
	S3            = 3 // Star met inside multi line
	SS            = 4 // String
	SC            = 5 // Comment
	SMC           = 6 // Multiline comment
)

func (x StateType) String() string {
	switch x {
	case S0:
		return "S0"
	case S1:
		return "S1"
	case S2:
		return "S2"
	case S3:
		return "S3"
	case SS:
		return "SS"
	case SC:
		return "SC"
	case SMC:
		return "SMC"
	default:
		return "??"
	}

}

func StripComments(file io.Reader) (rv []byte, err error) {
	var state = S0
	var cache = " "
	var buffer bytes.Buffer

	reader := bufio.NewReader(file)
	cache = ""

	charcode, err := reader.ReadByte()
	for err == nil {

		character := string(charcode) // To one-symbol string

		if db0 {
			fmt.Fprintf(os.Stderr, "State = %s character >%s< cache >%s<\n", state, character, cache)
		}

		switch state {
		case S0: // Start
			switch character {
			case "\"":
				state = SS
				// fmt.Printf("%s", character)
				buffer.WriteString(character)
			case "/":
				state = S2
				cache = character
			default:
				// fmt.Printf("%s%s", cache, character)
				buffer.WriteString(cache)
				buffer.WriteString(character)
				cache = ""
			}
			cache = ""
		case S1: // escape symbol inside string
			state = SS
			//	fmt.Printf("%s", character)
			buffer.WriteString(character)
		case S2: // first slash met
			switch character {
			case "/":
				state = SC
				cache = ""
			case "*":
				state = SMC
				cache = ""
			default:
				state = S0
				// fmt.Printf("%s%s", cache, character)
				buffer.WriteString(cache)
				buffer.WriteString(character)
				cache = ""
			}
			//cache = character
		case S3: // star met inside multiline comment
			switch character {
			case "/":
				state = S0
			default:
				state = SMC
			}
		case SS: // String
			switch character {
			case "\"":
				state = S0
			case "\\":
				state = S1
			}
			// fmt.Printf("%s", character)
			buffer.WriteString(character)
		case SC: // Comment
			switch character {
			case "\n":
				state = S0
				// fmt.Printf("%s", character)
				buffer.WriteString(character)
			}
		case SMC: // Multiline comment
			switch character {
			case "*":
				state = S3
			}
		}

		charcode, err = reader.ReadByte()
	}

	if state == SS || state == S1 {
		err = fmt.Errorf("EOF inside quoted string")
		return
	}
	if state == SMC || state == S2 {
		err = fmt.Errorf("EOF inside multiline comment")
		return
	}
	if state != S0 {
		err = fmt.Errorf("EOF at unexpected point")
		return
	}

	if db0 {
		fmt.Fprintf(os.Stderr, "State = %d\n", state)
	}

	return buffer.Bytes(), nil
}

var db0 = false
