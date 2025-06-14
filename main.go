package main

import (
	"fmt"
	"os"
	"regexp"

	aphrodite "github.com/jonathon-chew/Aphrodite"
)

var RestrictedTokens = []string{
	"if ",
	"def ",
	":",
	"case",
	"match",
	" or ",
	" and ",
	" = ",
	" == ",
	"return ",
	"try",
	"while ",
	"with ",
	"in ",
	"is ",
	"else:",
	"elif ",
	"for ",
	" as ",
	"assert ",
	"True",
	"False",
	"{",
	"}",
	"[",
	"]",
	"(",
	")",
}

func PrintFile(fileName string) {
	// Print relevent sections with relevent colouring
	aphrodite.Colour("Colour", "Green", fmt.Sprintf("Read the file %s\n", fileName))

	// Read file
	fileBytes, err := os.ReadFile(fileName)
	if err != nil {
		aphrodite.Colour("Colour", "Red", fmt.Sprintf("Issue reading the file %s\n", err))
		return
	}

	var byteLength int
	var found bool

	// namedFunction := regexp.MustCompile(`\s[\w\.]*\(`)
	number, err := regexp.Compile("\\d")
	if err != nil {
		return
	}

	// Parse file
	for i := 0; i < len(fileBytes); i++ {
		fileByte := fileBytes[i]
		found = false

		for index, value := range RestrictedTokens {
			byteLength = len(value)

			/*       if i+byteLength > len(fileBytes[i:]) {
			  aphrodite.Colour("color", "red", "ByteLength is too long")
			  continue
			} */

			// TODO: Get this to look back and change what came before it

			if string(fileBytes[i:i+byteLength]) == value && !found {
				aphrodite.Colour("Colour", "Green", RestrictedTokens[index])
				i += byteLength - 1 // because the for loop will add one!
				found = true
			}

		}

		// This looks at the rest of the line and makes it a comment
		if fileBytes[i] == '#' {
			newLineIndex := i
			byteLength := 1
			for newLineIndex < len(fileBytes) && fileBytes[newLineIndex] != '\n' {
				newLineIndex += 1
				byteLength += 1
			}
			aphrodite.Colour("Colour", "Yellow", string(fileBytes[i:i+byteLength]))
			i += byteLength - 1
			found = true
		}

		// This looks at the rest of the line and makes it a '
		if fileBytes[i] == '\'' {
			newLineIndex := i + 1
			byteLength := 1
			for newLineIndex < len(fileBytes) && fileBytes[newLineIndex] != '\'' && fileBytes[newLineIndex] != '\n' {
				newLineIndex += 1
				byteLength += 1
			}
			if fileBytes[newLineIndex] != '\n' {
				aphrodite.Colour("Colour", "Yellow", string(fileBytes[i:i+byteLength+1]))
				i += byteLength
				found = true
			}
		}

		// This looks at the rest of the line and makes it a "
		if fileBytes[i] == '"' {
			newLineIndex := i + 1
			byteLength := 1
			for newLineIndex < len(fileBytes) && fileBytes[newLineIndex] != '"' && fileBytes[newLineIndex] != '\n' {
				newLineIndex += 1
				byteLength += 1
			}
			if fileBytes[newLineIndex] != '\n' {
				aphrodite.Colour("Colour", "Yellow", string(fileBytes[i:i+byteLength+1]))
				i += byteLength
				found = true
			}
		}

		matches := number.FindString(string(fileBytes[i]))
		if matches != "" && !found {
			aphrodite.Colour("Colour", "Red", string(fileBytes[i]))
			found = true
		}

		if !found {
			aphrodite.Colour("Colour", "White", string(fileByte))
		}

	}
}

func main() {

	if len(os.Args) == 1 {
		aphrodite.Colour("Colour", "Red", "[ERROR]: Not enough arguments provided\n")
		return
	}

	for _, fileName := range os.Args[1:] {
		PrintFile(fileName)
	}

}
