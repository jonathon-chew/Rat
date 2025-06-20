package cmd

import (
	"fmt"
	"os"
	"strings"

	aphrodite "github.com/jonathon-chew/Aphrodite"
)

func ParseArguments(Arguments []string) ([]string, []string) {
	var returnArray []string
	var fileNames []string

	for index, argument := range Arguments {

		argument = strings.ToLower(argument)

		if argument == "*" {
			files, err := os.ReadDir(".")
			if err == nil {
				for _, eachFile := range files {
					fileName := eachFile.Name()
					if !eachFile.IsDir() {
						fileNames = append(fileNames, fileName)
					}
				}
			}
		} else if strings.Contains(argument, "*") {
			if len(argument) >= 1 {
				passedInFileType := argument[1:]
				files, err := os.ReadDir(".")
				if err == nil {
					for _, eachFile := range files {
						fileName := eachFile.Name()
						if !eachFile.IsDir() && strings.Contains(fileName, passedInFileType) {
							fileNames = append(fileNames, fileName)
						}
					}
				}
			}
		}

		_, err := os.Stat(argument)
		if err == nil {
			fileNames = append(fileNames, argument)
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

	return fileNames, returnArray
}
