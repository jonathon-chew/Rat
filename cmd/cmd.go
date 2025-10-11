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

	var file, findWord string

	for index := 0; index < len(Arguments); index++ {

		argument := Arguments[index]
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
			aphrodite.PrintColour("Green", "Pass in at least one file\n")
			aphrodite.PrintColour("Green", "You can use the flag allow of force to force it for unknown / unsupported file types\n")
			aphrodite.PrintColour("Green", "You can use --file-type to dictate the file type!")
			aphrodite.PrintColour("Green", "You can use the file type flag to choose the type of colour coding - eg comment / variable declaration \n")
			return []string{}, []string{"help_menu"}
		}

		if argument == "--allow" || argument == "--force" || argument == "-allow" || argument == "-force" {
			if !strings.Contains(argument, "--") {
				fmt.Println(argument)
			}
			returnArray = append(returnArray, "Allow")
		}

		if argument == "--file-type" || argument == "--filetype" || argument == "-filetype" || argument == "-file-type" {
			if index+1 < len(Arguments) {
				returnArray = append(returnArray, Arguments[index+1])
			}
		}

		if argument == "--file" || argument == "-f" {
			if index+1 >= len(Arguments) {
				aphrodite.PrintError("Missing filename after --file")
				return []string{}, []string{}
			}
			next := Arguments[index+1]
			if _, err := os.Stat(next); err != nil {
				aphrodite.PrintError(fmt.Sprintf("File not found or inaccessible: %s", next))
				return []string{}, []string{}
			}
			file = next
			index++
		}

		if argument == "--word" || argument == "-w" {
			if index+1 >= len(Arguments) {
				aphrodite.PrintError("Missing word after --word")
				return []string{}, []string{}
			}
			findWord = Arguments[index+1]
			index++
		}
	}

	if findWord != "" && file != "" {
		return []string{file}, []string{"Plain_File", findWord}
	}

	return fileNames, returnArray
}
