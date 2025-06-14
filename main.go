package main

import (
	"fmt"
	"github.com/jonathon-chew/Aphrodite"
	"os"
	"regexp"
)

var RestrictedTokens = []string{
	"if ",
	"def ",
	":",
	"#",
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

	// Parse file
	for i := 0; i < len(fileBytes); i++ {
		fileByte := fileBytes[i]
		found = false

		for index, value := range RestrictedTokens {
			byteLength = len(value)

			// This looks at the rest of the line and makes it a comment
			if fileBytes[i] == '#' {
				newLineIndex := i
				for fileBytes[newLineIndex] != '\n' {
					newLineIndex += 1
					byteLength += 1
				}
				aphrodite.Colour("Colour", "Yellow", string(fileBytes[i:i+byteLength]))
				i += byteLength // because the for loop will add one!
				found = true
			}

			// TODO: Get this to look back and change what came before it
			// if fileBytes[i] == '(' {
			//   newLineIndex := i
			//   for fileBytes[newLineIndex] != ' ' {
			//     newLineIndex -= 1
			//     byteLength += 1
			//   }
			//   aphrodite.Colour("Colour", "Yellow", string(fileBytes[i-byteLength:i]))
			//   i += 1 // because the for loop will add one!
			//   found = true
			// }

			if string(fileBytes[i:i+byteLength]) == value && found != true {
				aphrodite.Colour("Colour", "Green", RestrictedTokens[index])
				i += byteLength - 1 // because the for loop will add one!
				found = true
			}

		}

		number, _ := regexp.Match("\\d", []byte(string(fileBytes[i])))
		if number {
			aphrodite.Colour("Colour", "Red", string(fileBytes[i]))
			found = true
		}

		if found == false {
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
