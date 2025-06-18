package cmd

import (
	"fmt"
	"os"
	"strings"
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

		if argument == "--allow" || argument == "--force" {
			if !strings.Contains(argument, "--") {
				fmt.Println(argument)
			}
			returnArray = append(returnArray, "Allow")
		}

		if argument == "--file-type" || argument == "--fileType" || argument == "-file" {
			if index+1 < len(Arguments) {
				returnArray = append(returnArray, Arguments[index+1])
			}
		}
	}

	return fileName, returnArray
}
