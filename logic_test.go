package main

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	aphrodite "github.com/jonathon-chew/Aphrodite"
)

func TestPossessiveComma(t *testing.T) {
	t.Logf("Testing the possesive comma func")

	tests := map[string]bool{
		// Valid English contractions / possessives (apostrophe between letters)
		"can't":       true, // contraction: apostrophe between letters
		"he's":        true, // contraction: apostrophe between letters
		"rock'n'roll": true, // multiple apostrophes within a single word
		"won't":       true, // contraction with omitted letters
		"it's":        true, // possessive / contraction form

		// Invalid: apostrophe at word boundary (not a contraction)
		"'start": false, // leading apostrophe, not between letters
		"end'":   false, // trailing apostrophe, not between letters

		// No apostrophe present
		"no quote": false, // no apostrophe to evaluate

		// Python single-quoted string literals (apostrophes are string delimiters)
		"'hello'":       false, // simple single-quoted string literal
		"'can't'":       false, // contraction inside a string literal, not punctuation
		"'rock'n'roll'": false, // apostrophes interpreted as string delimiters, not word punctuation

		// Python f-strings using single quotes
		"f'hello'":  false, // f-string with single-quoted literal
		"f'can't'":  false, // contraction inside f-string literal
		"f'{name}'": false, // f-string expression, apostrophes are delimiters

		// Python f-strings using double quotes (apostrophe is literal text)
		"f\"he's\"": false, // apostrophe inside double-quoted f-string

		// Escaped apostrophes inside Python strings
		"'can\\'t'":  false, // escaped apostrophe within single-quoted string literal
		"f'can\\'t'": false, // escaped apostrophe within f-string literal
	}

	supportedFileTypes = []string{"python"}

	type ReportCard struct {
		Good int
		Bad  int
		Seen int
	}

	var report ReportCard

	var counter int
	for _, supportedType := range supportedFileTypes {
		aphrodite.PrintInfo(fmt.Sprintf("Looking at: %s\n______________\n", supportedType))
		for test, result := range tests {
			// checkStr := string(test[0:3]) // impliment when back
			number := strconv.Itoa(counter)
			if notPossessiveComma([]byte(test), supportedType) == result {
				// aphrodite.PrintColour("Green", number+" "+test+"\n")
				report.Good += 1
				report.Seen += 1
			} else {
				aphrodite.PrintColour("Red", number+" "+test+"\n")
				report.Bad += 1
				report.Seen += 1
				t.Error("")
			}

			counter++
		}
		fmt.Print("\n")
	}

	t.Log("Report: ", report.Good, report.Bad, report.Seen, "\n")
}

func Benchmark_Full_Process(b *testing.B) {
	os.Args = append(os.Args, "./tests/main.go")

	main()
}

func Benchmark_ParseFile(b *testing.B) {

	files := []struct {
		fileName string
		fileType string
	}{
		{
			fileName: "./tests/main.go",
			fileType: "golang",
		},
	}

	for _, test := range files {
		PrintFile(test.fileName, test.fileType)
	}
}
