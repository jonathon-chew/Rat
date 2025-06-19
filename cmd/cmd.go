package cmd

import (
	"fmt"
	"os"
	"strings"

	aphrodite "github.com/jonathon-chew/Aphrodite"
)

func ParseArguments(Arguments []string) ([]string, []string) {
	var returnArray []string
	var fileName []string

	for index, argument := range Arguments {

		argument = strings.ToLower(argument)

		_, err := os.Stat(argument)
		if err == nil {
			fileName = append(fileName, argument)
		}

		if argument == "--help" || argument == "-help" {
			aphrodite.Colour("Colour", "Green", "Pass in at least one file\n")
			aphrodite.Colour("Colour", "Green", "You can use thte flag allow of force to force it for unknown / unsupported file types\n")
			aphrodite.Colour("Colour", "Green", "You can use the file type flag to choose the type of colour coding - eg comment / variable declaration \n")
			return []string{}, []string{}
		}

		if argument == "--allow" || argument == "--force" || argument == "-allow" || argument == "-force" {
			if !strings.Contains(argument, "--") {
				fmt.Println(argument)
			}
			returnArray = append(returnArray, "Allow")
		}

		if argument == "--file-type" || argument == "--filetype" || argument == "-file" {
			if index+1 < len(Arguments) {
				returnArray = append(returnArray, Arguments[index+1])
			}
		}
	}

	return fileName, returnArray
}
