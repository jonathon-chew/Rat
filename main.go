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

func PrintFile(fileName, fileExtension string) {
	// Print relevent sections with relevent colouring
	aphrodite.Colour("Colour", "Green", fmt.Sprintf("\nRead the file %s\n", fileName))

	// Read file
	fileBytes, err := os.ReadFile(fileName)
	if err != nil {
		aphrodite.Colour("Colour", "Red", fmt.Sprintf("Issue reading the file %s\n", err))
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
					aphrodite.Colour("Colour", "Green", PythonRestrictedTokens[index])
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
						aphrodite.Colour("Colour", "Blue", JavascriptRestrictedTokens[index])
						i += byteLength - 1
						found = true
					}
				}
			}
		}

		// This looks at the rest of the line and makes it a comment
		if fileBytes[i] == '#' && fileExtension == "python" {
			newLineIndex := i
			for newLineIndex < len(fileBytes) && fileBytes[newLineIndex] != '\n' {
				newLineIndex++
			}

			// We want to include the newline if it exists
			if newLineIndex < len(fileBytes) && fileBytes[newLineIndex] == '\n' {
				newLineIndex++ // include newline in slice
			}

			comment := string(fileBytes[i:newLineIndex])
			aphrodite.Colour("Colour", "Yellow", comment)

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
			aphrodite.Colour("Colour", "Yellow", comment)

			i = newLineIndex - 1 // -1 because the for loop will increment `i`
			found = true
		}

		// (#3) TODO: Look to combine single and double quotes in the case of python this is treated as the same thing!

		// (#6) TODO: The ' and " doesn't work in json as it doesn't know when to stop while in each other 's in "" makes everything a comment!

		// This looks at the rest of the line and makes it a '
		if fileBytes[i] == '\'' && string(fileBytes[i-1]) != string('\\') {
			nextSingleQuote := i + 1 // Plus 1 because the current byte is a ' we're looking for the next one
			for nextSingleQuote < len(fileBytes) && fileBytes[nextSingleQuote] != '\'' && string(fileBytes[i-1]) != string('\\') {
				nextSingleQuote++
			}

			// We want to include the newline if it exists
			if nextSingleQuote < len(fileBytes) && fileBytes[nextSingleQuote] == '\'' {
				nextSingleQuote++ // include newline in slice
			}

			comment := string(fileBytes[i:nextSingleQuote])
			aphrodite.Colour("Colour", "Yellow", comment)

			i = nextSingleQuote - 1 // -1 because the for loop will increment `i`
			found = true
		}

		// This looks at the rest of the line and makes it a "
		if fileBytes[i] == '"' && string(fileBytes[i-1]) != string('\\') {
			nextDoubleQuote := i + 1 // Plus 1 because the current byte is a ' we're looking for the next one
			for nextDoubleQuote < len(fileBytes) && fileBytes[nextDoubleQuote] != '"' && string(fileBytes[i-1]) != string('\\') {
				nextDoubleQuote++
			}

			// We want to include the newline if it exists
			if nextDoubleQuote < len(fileBytes) && fileBytes[nextDoubleQuote] == '"' {
				nextDoubleQuote++ // include newline in slice
			}

			comment := string(fileBytes[i:nextDoubleQuote])
			aphrodite.Colour("Colour", "Yellow", comment)

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
				aphrodite.Colour("Colour", "Blue", comment)

				i = nextSpaceChr - 1 // -1 because the for loop will increment `i`
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

	fileNames, flags := cmd.ParseArguments(os.Args[1:])

	// (#4) TODO: Add more language suuport

	// (#5) TODO: JSON can use a lot of what python can BUT the test example

	// (#7) TODO: Look at whether to return instead of print, and print at the end for speed? Look at is printing including escape for every chr and slowing it down as well!

	if len(fileNames) == 0 {
		aphrodite.Colour("Colour", "Red", "[ERROR]: No file could be detected\n")
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
				aphrodite.Colour("Colour", "Red", fmt.Sprintf("\n[ERROR]: File extension %s is not yet supported\n", fileExtension[len(fileExtension)-1]))
			}
		}
	}
}
