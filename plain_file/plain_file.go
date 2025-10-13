package plainFile

import (
	"fmt"
	"os"
	"strings"

	aphrodite "github.com/jonathon-chew/Aphrodite"
	"github.com/jonathon-chew/Rat/cmd"
)

func Parse_plain_file(file, findWord string, settings cmd.Settings) {

	fileBytes, err := os.ReadFile(file)
	if err != nil {
		aphrodite.PrintError(fmt.Sprintf("Error reading file: %v\n", err))
		return
	}

	var currentWord []byte
	for i := 0; i < len(fileBytes); i++ {
		b := fileBytes[i]

		if b != ' ' && b != '\n' && b != '\t' {
			currentWord = append(currentWord, b)
			continue
		}

		origionalWord := string(currentWord)
		var word string
		if settings.Caseinsensative {
			word = strings.ToLower(origionalWord)
		}

		if word == findWord {
			if i+1 < len(fileBytes) {
				aphrodite.PrintBold("Green", origionalWord+string(fileBytes[i]))
			} else {
				aphrodite.PrintBold("Green", origionalWord)
			}
		} else if len(currentWord) > 0 {
			if i+1 < len(fileBytes) {
				fmt.Print(word, string(fileBytes[i]))
			} else {
				fmt.Print(word)
			}
		}
		currentWord = currentWord[:0] // reset
	}

	// Handle case where file doesn’t end with whitespace
	if len(currentWord) > 0 {
		if string(currentWord) == findWord {
			aphrodite.PrintBold("Green", string(currentWord))
		} else {
			fmt.Print(string(currentWord))
		}
	}
}
