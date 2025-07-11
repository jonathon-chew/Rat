package main

import (
	"fmt"
	"os"
	"regexp"
	"slices"
	"strings"

	aphrodite "github.com/jonathon-chew/Aphrodite"
	"github.com/jonathon-chew/Rat/cmd"
)

var FileType = map[string]string{"py": "python", "go": "golang", "ps1": "powershell", "json": "json", "js": "javascript"}

var PythonRestrictedTokens = []string{"if ", "def ", ":", "case", "match ", " or ", " and ", " = ", " == ", "return ", "try", "while ", "with ", "in ", "is ", "else:", "elif ", "for ", " as ", "assert ", "True", "False", "{", "}", "[", "]", "(", ")"}

var JavascriptRestrictedTokens = []string{"const ", "var ", "let ", "new ", "&&", "||"}

func notPossessiveComma(checkBytes string) bool {

	// namedFunction := regexp.MustCompile(`\s[\w\.]*\(`)
	possessive_comma, err := regexp.Compile(`[a-zA-Z]'[a-zA-Z]`)
	if err != nil {
		return false
	}

	matches := possessive_comma.FindString(checkBytes)
	return matches != ""
}

func PrintFile(fileName, fileExtension string) {
	// Print relevent sections with relevent colouring
	aphrodite.PrintColour("Green", fmt.Sprintf("\nRead the file %s\n", fileName))

	// Read file
	fileBytes, err := os.ReadFile(fileName)
	if err != nil {
		aphrodite.PrintColour("Red", fmt.Sprintf("Issue reading the file %s\n", err))
		return
	}

	var byteLength int
	var found bool

	// namedFunction := regexp.MustCompile(`\s[\w\.]*\(`)
	number, err := regexp.Compile(`\d`)
	if err != nil {
		return
	}

	// Parse file
	for i := 0; i < len(fileBytes); i++ {
		fileByte := fileBytes[i]
		found = false

		for index, value := range PythonRestrictedTokens {
			byteLength = len(value)

			if byteLength+1 < len(fileBytes[i:]) {
				// fmt.Printf("Looking at: %d, and the rest of the file is %d, %t", i+byteLength+1, len(fileBytes[i:]), (i+byteLength > len(fileBytes[i:])))
				// fmt.Print("\n")
				if string(fileBytes[i:i+byteLength]) == value && !found {
					aphrodite.PrintColour("Green", PythonRestrictedTokens[index])
					i += byteLength - 1 // because the for loop will add one!
					found = true
				}
			}
		}

		if fileExtension == "javascript" {
			for index, value := range JavascriptRestrictedTokens {
				byteLength = len(value)

				if byteLength+1 < len(fileBytes[i:]) {
					if string(fileBytes[i:i+byteLength]) == value && !found {
						aphrodite.PrintColour("Blue", JavascriptRestrictedTokens[index])
						i += byteLength - 1
						found = true
					}
				}
			}
		}

		// This looks at the rest of the line and makes it a comment
		if fileBytes[i] == '#' && (fileExtension == "python" || fileExtension == "powershell") {
			newLineIndex := i
			for newLineIndex < len(fileBytes) && fileBytes[newLineIndex] != '\n' {
				newLineIndex++
			}

			// We want to include the newline if it exists
			if newLineIndex < len(fileBytes) && fileBytes[newLineIndex] == '\n' {
				newLineIndex++ // include newline in slice
			}

			comment := string(fileBytes[i:newLineIndex])
			aphrodite.PrintColour("Yellow", comment)

			i = newLineIndex - 1 // -1 because the for loop will increment `i`
			found = true
		}

		if fileBytes[i] == '/' && fileBytes[i+1] == '/' && (fileExtension == "golang" || fileExtension == "javascript") {
			newLineIndex := i
			for newLineIndex < len(fileBytes) && fileBytes[newLineIndex] != '\n' {
				newLineIndex++
			}

			// We want to include the newline if it exists
			if newLineIndex < len(fileBytes) && fileBytes[newLineIndex] == '\n' {
				newLineIndex++ // include newline in slice
			}

			comment := string(fileBytes[i:newLineIndex])
			aphrodite.PrintColour("Yellow", comment)

			i = newLineIndex - 1 // -1 because the for loop will increment `i`
			found = true
		}

		// (#8) TODO: Work out how to do escape quotes

		// This looks at the rest of the line and makes it a '
		if fileBytes[i] == '\'' && string(fileBytes[i-1]) != string('\\') {
			if len(fileBytes[i:]) > 1 {
				if notPossessiveComma(string(fileBytes[i-1 : i+1])) {
					nextSingleQuote := i + 1 // Plus 1 because the current byte is a ' we're looking for the next one
					for nextSingleQuote < len(fileBytes) && fileBytes[nextSingleQuote] != '\'' && string(fileBytes[i-1]) != string('\\') {
						nextSingleQuote++
					}

					// We want to include the newline if it exists
					if nextSingleQuote < len(fileBytes) && fileBytes[nextSingleQuote] == '\'' {
						nextSingleQuote++ // include newline in slice
					}

					comment := string(fileBytes[i:nextSingleQuote])
					aphrodite.PrintColour("Yellow", comment)

					i = nextSingleQuote - 1 // -1 because the for loop will increment `i`
					found = true
				}
			}
		}

		// This looks at the rest of the line and makes it a "
		if fileBytes[i] == '"' && fileBytes[i-1] != '\\' {
			nextDoubleQuote := i + 1 // Plus 1 because the current byte is a ' we're looking for the next one
			for nextDoubleQuote < len(fileBytes) && fileBytes[nextDoubleQuote] != '"' && string(fileBytes[i-1]) != string('\\') {
				nextDoubleQuote++
			}

			// We want to include the newline if it exists
			if nextDoubleQuote < len(fileBytes) && fileBytes[nextDoubleQuote] == '"' {
				nextDoubleQuote++ // include newline in slice
			}

			comment := string(fileBytes[i:nextDoubleQuote])
			aphrodite.PrintColour("Yellow", comment)

			i = nextDoubleQuote - 1 // -1 because the for loop will increment `i`
			found = true
		}

		// (#1) TODO: Colour code functions

		// (#2) TODO: Colour code variables \s[a-zA-Z]+[\s|=]
		if fileExtension == "powershell" || fileExtension == "php" {
			if fileBytes[i] == '$' {
				nextSpaceChr := i + 1 // Plus 1 because the current byte is a ' we're looking for the next one
				for nextSpaceChr < len(fileBytes) && (fileBytes[nextSpaceChr] != ' ' && fileBytes[nextSpaceChr] != '=' && fileBytes[nextSpaceChr] != '.' && fileBytes[nextSpaceChr] != '(') {
					nextSpaceChr++
				}

				// We want to include the newline if it exists
				if nextSpaceChr < len(fileBytes) && (fileBytes[nextSpaceChr] != ' ' && fileBytes[nextSpaceChr] != '=') {
					nextSpaceChr++ // include newline in slice
				}

				comment := string(fileBytes[i:nextSpaceChr])
				aphrodite.PrintColour("Blue", comment)

				i = nextSpaceChr - 1 // -1 because the for loop will increment `i`
				found = true
			}
		}

		matches := number.FindString(string(fileBytes[i]))
		if matches != "" && !found {
			aphrodite.PrintColour("Red", string(fileBytes[i]))
			found = true
		}

		if !found {
			aphrodite.PrintColour("White", string(fileByte))
		}

	}
}

func main() {

	if len(os.Args) == 1 {
		aphrodite.PrintColour("Red", "[ERROR]: Not enough arguments provided\n")
		return
	}

	fileNames, flags := cmd.ParseArguments(os.Args[1:])

	// (#4) TODO: Add more language suuport

	// (#5) TODO: JSON can use a lot of what python can BUT the test example

	// (#7) TODO: Look at whether to return instead of print, and print at the end for speed? Look at is printing including escape for every chr and slowing it down as well!

	if len(fileNames) == 0 {
		aphrodite.PrintColour("Red", "[ERROR]: No file could be detected\n")
		return
	}

	for _, fileName := range fileNames {
		fileExtension := strings.Split(fileName, ".")
		convertedFileType := FileType[fileExtension[len(fileExtension)-1]]

		for _, extension := range FileType {
			if slices.Contains(flags, extension) {
				convertedFileType = extension
			}
		}

		var supportedFileTypes = []string{"python", "powershell", "json", "javascript"}

		if slices.Contains(supportedFileTypes, convertedFileType) || slices.Contains(flags, "Allow") {
			PrintFile(fileName, convertedFileType)
		} else {
			// Check the file name isn't in the flags - things like --allow or --file-type
			if !slices.Contains(flags, fileName) {
				aphrodite.PrintColour("Red", fmt.Sprintf("\n[ERROR]: File extension %s is not yet supported\n", fileExtension[len(fileExtension)-1]))
			}
		}
	}
}
