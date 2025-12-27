package main

import (
	"fmt"
	"os"
	"regexp"
	"slices"
	"strings"

	aphrodite "github.com/jonathon-chew/Aphrodite"
	"github.com/jonathon-chew/Rat/cmd"
	plainFile "github.com/jonathon-chew/Rat/plain_file"
)

var FileType = map[string]string{"py": "python", "go": "golang", "ps1": "powershell", "json": "json", "js": "javascript"}

// Token slices sorted by length (longest first) to avoid partial matches
var PythonRestrictedTokens = []string{"else:", "elif ", "match ", "assert ", "return ", "while ", "with ", "True", "False", "def ", "if ", "for ", "try", "in ", "is ", " or ", " and ", " == ", " as ", " = ", "case", ":", "{", "}", "[", "]", "(", ")", "import ", "from "}

var JavascriptRestrictedTokens = []string{"const ", "var ", "let ", "new ", "&&", "||", "null "}

var GolangRestrictedTokens = []string{"package ", "import ", "interface ", "default ", "return ", "struct ", "range ", "switch ", "defer ", "select ", "func ", "else ", "case ", "type ", "var ", "const ", "if ", "for ", "go ", "chan ", "map", "make", "len", "cap", "new", "nil", "true", "false", ":=", "==", "!=", "<=", ">=", "&&", "||", "{", "}", "[", "]", "(", ")"}

var supportedFileTypes = []string{"python", "powershell", "json", "javascript", "golang"}

// Pre-compiled regex for possessive comma detection (optimization #2)
var possessiveCommaRegex = regexp.MustCompile(`[a-zA-Z]'[a-zA-Z]`)

// Token byte slices for faster matching (optimization #9)
var pythonTokenBytes = makeTokenBytes(PythonRestrictedTokens)
var javascriptTokenBytes = makeTokenBytes(JavascriptRestrictedTokens)
var golangTokenBytes = makeTokenBytes(GolangRestrictedTokens)

// makeTokenBytes converts string tokens to byte slices for faster comparison
func makeTokenBytes(tokens []string) [][]byte {
	result := make([][]byte, len(tokens))
	for i, token := range tokens {
		result[i] = []byte(token)
	}
	return result
}

func notPossessiveComma(checkBytes []byte, fileExtension string) bool {
	// fmt.Printf("I've been passed: %s\n", checkBytes)
	if fileExtension != "python" {
		return possessiveCommaRegex.Match(checkBytes)
	}

	if fileExtension == "python" && checkBytes[0] == 'f' {
		// fmt.Println("I'm an f string")
		return false
	}

	return possessiveCommaRegex.Match(checkBytes)
}

// flushBuffer outputs buffered content with the specified color
func flushBuffer(buffer []byte, color string) {
	if len(buffer) > 0 {
		aphrodite.PrintColour(color, string(buffer))
	}
}

// findTokenMatch checks if a token matches at the current position
func findTokenMatch(fileBytes []byte, i int, tokens [][]byte) (int, []byte, bool) {
	for _, tokenBytes := range tokens {
		tokenLen := len(tokenBytes)
		if i+tokenLen <= len(fileBytes) {
			// Fast byte comparison
			match := true
			for j := 0; j < tokenLen; j++ {
				if fileBytes[i+j] != tokenBytes[j] {
					match = false
					break
				}
			}
			if match {
				return tokenLen, tokenBytes, true
			}
		}
	}
	return 0, nil, false
}

// findStringEnd finds the end of a string, handling escape sequences properly
func findStringEnd(fileBytes []byte, start int, quote byte) int {
	i := start + 1
	for i < len(fileBytes) {
		if fileBytes[i] == quote {
			// Check if it's escaped - count backslashes
			backslashCount := 0
			j := i - 1
			for j >= start && fileBytes[j] == '\\' {
				backslashCount++
				j--
			}
			// If even number of backslashes (or zero), the quote is not escaped
			if backslashCount%2 == 0 {
				return i + 1
			}
		}
		i++
	}
	return len(fileBytes)
}

// handleLineComment handles # and // style comments
func handleLineComment(fileBytes []byte, i int, commentChar byte) (int, bool) {
	if i >= len(fileBytes) || fileBytes[i] != commentChar {
		return i, false
	}

	// Check for // (two slashes)
	if commentChar == '/' {
		if i+1 >= len(fileBytes) || fileBytes[i+1] != '/' {
			return i, false
		}
	}

	newLineIndex := i
	for newLineIndex < len(fileBytes) && fileBytes[newLineIndex] != '\n' {
		newLineIndex++
	}
	if newLineIndex < len(fileBytes) && fileBytes[newLineIndex] == '\n' {
		newLineIndex++
	}
	return newLineIndex, true
}

// handleBlockComment handles /* */ style comments
func handleBlockComment(fileBytes []byte, i int) (int, bool) {
	if i+1 >= len(fileBytes) || fileBytes[i] != '/' || fileBytes[i+1] != '*' {
		return i, false
	}

	// Find closing */
	end := i + 2
	for end+1 < len(fileBytes) {
		if fileBytes[end] == '*' && fileBytes[end+1] == '/' {
			return end + 2, true
		}
		end++
	}
	// If not closed, treat rest of file as comment
	return len(fileBytes), true
}

// handleTripleQuotedString handles Python triple-quoted strings (""" or ”')
func handleTripleQuotedString(fileBytes []byte, i int, quote byte) (int, bool) {
	if i+2 >= len(fileBytes) {
		return i, false
	}

	// Check for triple quotes
	if fileBytes[i] == quote && fileBytes[i+1] == quote && fileBytes[i+2] == quote {
		// Find closing triple quotes
		end := i + 3
		for end+2 < len(fileBytes) {
			if fileBytes[end] == quote && fileBytes[end+1] == quote && fileBytes[end+2] == quote {
				return end + 3, true
			}
			end++
		}
		// If not closed, treat rest of file as string
		return len(fileBytes), true
	}
	return i, false
}

func PrintFile(fileName, fileExtension string) {
	// Print relevant sections with relevant colouring
	aphrodite.PrintColour("Green", fmt.Sprintf("\nRead the file %s\n", fileName))

	// Read file
	fileBytes, err := os.ReadFile(fileName)
	if err != nil {
		aphrodite.PrintColour("Red", fmt.Sprintf("Issue reading the file %s\n", err))
		return
	}

	if len(fileBytes) == 0 {
		return
	}

	// Buffering for performance: accumulate characters of the same color
	var currentBuffer []byte
	currentColor := "White"

	// Helper function to change color and flush buffer
	changeColor := func(newColor string, newContent []byte) {
		if newColor != currentColor {
			flushBuffer(currentBuffer, currentColor)
			currentBuffer = currentBuffer[:0]
			currentColor = newColor
		}
		currentBuffer = append(currentBuffer, newContent...)
	}

	// Determine which tokens to check based on file extension (optimization #4)
	var tokensToCheck [][]byte
	var tokenColor string
	switch fileExtension {
	case "python":
		tokensToCheck = pythonTokenBytes
		tokenColor = "Green"
	case "javascript":
		tokensToCheck = javascriptTokenBytes
		tokenColor = "Blue"
	case "golang":
		tokensToCheck = golangTokenBytes
		tokenColor = "Blue"
	default:
		tokensToCheck = nil
	}

	// Parse file
	for i := 0; i < len(fileBytes); i++ {
		found := false

		// Check language-specific tokens (optimization #4)
		if tokensToCheck != nil {
			if tokenLen, tokenBytes, matched := findTokenMatch(fileBytes, i, tokensToCheck); matched {
				changeColor(tokenColor, tokenBytes)
				i += tokenLen - 1
				found = true
			}
		}

		// Handle block comments /* */ for golang/javascript (improvement #1)
		if !found && (fileExtension == "golang" || fileExtension == "javascript") {
			if end, ok := handleBlockComment(fileBytes, i); ok {
				changeColor("Yellow", fileBytes[i:end])
				i = end - 1
				found = true
			}
		}

		// Handle line comments: # for python/powershell
		if !found && fileBytes[i] == '#' && (fileExtension == "python" || fileExtension == "powershell") {
			if end, ok := handleLineComment(fileBytes, i, '#'); ok {
				changeColor("Yellow", fileBytes[i:end])
				i = end - 1
				found = true
			}
		}

		// Handle line comments: // for golang/javascript
		if !found && (fileExtension == "golang" || fileExtension == "javascript") {
			if end, ok := handleLineComment(fileBytes, i, '/'); ok {
				changeColor("Yellow", fileBytes[i:end])
				i = end - 1
				found = true
			}
		}

		// Handle Python triple-quoted strings (improvement #10)
		if !found && fileExtension == "python" {
			if end, ok := handleTripleQuotedString(fileBytes, i, '"'); ok {
				changeColor("Yellow", fileBytes[i:end])
				i = end - 1
				found = true
			} else if end, ok := handleTripleQuotedString(fileBytes, i, '\''); ok {
				changeColor("Yellow", fileBytes[i:end])
				i = end - 1
				found = true
			}
		}

		// Handle Go raw strings (backticks) (improvement #6)
		if !found && fileExtension == "golang" && fileBytes[i] == '`' {
			end := i + 1
			for end < len(fileBytes) && fileBytes[end] != '`' {
				end++
			}
			if end < len(fileBytes) {
				end++
			}
			changeColor("Yellow", fileBytes[i:end])
			i = end - 1
			found = true
		}

		// Handle single quotes with improved escape handling (improvement #5)
		if !found && fileBytes[i] == '\'' {
			// Check if it's a possessive comma (like "can't")
			isPossessive := false
			end := findStringEnd(fileBytes, i-1, '\'')
			checkStr := fileBytes[i:end]
			if len(checkStr) >= 2 {
				isPossessive = notPossessiveComma(checkStr, fileExtension)
			}

			if !isPossessive {
				end := findStringEnd(fileBytes, i, '\'')
				changeColor("Yellow", fileBytes[i:end])
				i = end - 1
				found = true
			}
		}

		// Handle double quotes with improved escape handling (improvement #5)
		if !found && fileBytes[i] == '"' && (i == 0 || fileBytes[i-1] != '\\') {
			end := findStringEnd(fileBytes, i, '"')
			changeColor("Yellow", fileBytes[i:end])
			i = end - 1
			found = true
		}

		// Handle PowerShell/PHP variables
		if !found && (fileExtension == "powershell" || fileExtension == "php") && fileBytes[i] == '$' {
			nextSpaceChr := i + 1
			for nextSpaceChr < len(fileBytes) && fileBytes[nextSpaceChr] != ' ' && fileBytes[nextSpaceChr] != '=' && fileBytes[nextSpaceChr] != '.' && fileBytes[nextSpaceChr] != '(' {
				nextSpaceChr++
			}
			changeColor("Blue", fileBytes[i:nextSpaceChr])
			i = nextSpaceChr - 1
			found = true
		}

		// Handle digits (optimized: no regex, simple byte comparison)
		if !found && fileBytes[i] >= '0' && fileBytes[i] <= '9' {
			changeColor("Red", []byte{fileBytes[i]})
			found = true
		}

		// Default: white text
		if !found {
			changeColor("White", []byte{fileBytes[i]})
		}
	}

	// Flush any remaining buffered content
	flushBuffer(currentBuffer, currentColor)
}

func main() {

	if len(os.Args) == 1 {
		aphrodite.PrintColour("Red", "[ERROR]: Not enough arguments provided\n")
		return
	}

	fileNames, flags, settings := cmd.ParseArguments(os.Args[1:])

	// (#5) TODO: JSON can use a lot of what python can BUT the test example

	if len(flags) == 1 && flags[0] == "help_menu" {
		return // as to not return a error when the help menu has been selected
	}

	if len(fileNames) == 0 {
		aphrodite.PrintColour("Red", "[ERROR]: No file could be detected\n")
		return
	}

	for _, fileName := range fileNames {
		fileExtension := strings.Split(fileName, ".")
		// Handle files without extensions
		var convertedFileType string
		if len(fileExtension) > 1 {
			convertedFileType = FileType[fileExtension[len(fileExtension)-1]]
		} else {
			convertedFileType = ""
		}

		for _, extension := range FileType {
			if slices.Contains(flags, extension) {
				convertedFileType = extension
			}
		}

		if slices.Contains(supportedFileTypes, convertedFileType) || slices.Contains(flags, "Allow") {
			PrintFile(fileName, convertedFileType)
		} else if slices.Contains(flags, "Plain_File") {
			plainFile.Parse_plain_file(fileName, flags[1], settings)
		} else {
			// Check the file name isn't in the flags - things like --allow or --file-type
			if !slices.Contains(flags, fileName) {
				ext := "unknown"
				if len(fileExtension) > 1 {
					ext = fileExtension[len(fileExtension)-1]
				}
				aphrodite.PrintColour("Red", fmt.Sprintf("\n[ERROR]: File extension %s is not yet supported\n", ext))
			}
		}
	}
}
