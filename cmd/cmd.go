package cmd

import (
	"fmt"
	"os"
	"strings"

	aphrodite "github.com/jonathon-chew/Aphrodite"
)

type Settings struct {
	Caseinsensative bool
}

func ParseArguments(Arguments []string) ([]string, []string, Settings) {
	var fileNames, flagsArray []string
	var file, findWord string
	var new_settings Settings

	for index := 0; index < len(Arguments); index++ {
		argument := strings.ToLower(Arguments[index])

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

		// -----------
		// Flag section
		//------------
		if argument == "--help" || argument == "-help" {
			aphrodite.PrintBold("Purple", "Arguments:\n")
			aphrodite.PrintBold("Cyan", "*\n")
			aphrodite.PrintInfo("This can be used to do this on all files in the current directory, or all files of a certain type ie *.py would be all python files, it will only print on supported files unless the force flag is also used\n\n")
			aphrodite.PrintBold("Cyan", "--help\n")
			aphrodite.PrintInfo("Show all the flags\n\n")
			aphrodite.PrintBold("Cyan", "--force --allow\n")
			aphrodite.PrintInfo("You can use the flag allow of force to force it for unknown / unsupported file types\n\n")
			aphrodite.PrintBold("Cyan", "--file-type\n")
			aphrodite.PrintInfo("You can use --file-type to dictate the file type!\n\n")
			aphrodite.PrintBold("Cyan", "--file\n")
			aphrodite.PrintInfo("You can use this to pass in a file of any type if you just want to highlight a specific word in the file and make it stand out\n\n")
			aphrodite.PrintBold("Cyan", "--word\n")
			aphrodite.PrintInfo("You use this in conjunction with the file flag to denote which word in the file you are looking for\n\n")

			return []string{}, []string{"help_menu"}, new_settings
		}

		if argument == "--allow" || argument == "--force" || argument == "-allow" || argument == "-force" {
			if !strings.Contains(argument, "--") {
				fmt.Println(argument)
			}
			flagsArray = append(flagsArray, "Allow")
		}

		if argument == "--file-type" || argument == "--filetype" || argument == "-filetype" || argument == "-file-type" {
			if index+1 < len(Arguments) {
				flagsArray = append(flagsArray, Arguments[index+1])
			}
		}

		if argument == "--file" || argument == "-f" {
			if index+1 >= len(Arguments) {
				aphrodite.PrintError("Missing filename after --file\n")
				return []string{}, []string{}, new_settings
			}
			next := Arguments[index+1]
			if _, err := os.Stat(next); err != nil {
				aphrodite.PrintError("File not found or inaccessible:" + next + "\n")
				return []string{}, []string{}, new_settings
			}
			file = next
			index++
		}

		if argument == "--word" || argument == "-w" {
			if index+1 >= len(Arguments) {
				aphrodite.PrintError("Missing word after --word\n")
				return []string{}, []string{}, new_settings
			}
			findWord = Arguments[index+1]
			index++
		}

		if argument == "--case-insensative" || argument == "-caseinsensative" || argument == "-cs" {
			new_settings.Caseinsensative = true
		}

		if argument == "--version" || argument == "-v" {
			versionNumber := "v0.1.4"
			aphrodite.PrintInfo("Version Number is: " + versionNumber + "\n")
		}
	}

	if findWord != "" && file != "" {
		return []string{file}, []string{"Plain_File", findWord}, new_settings
	}

	return fileNames, flagsArray, new_settings
}
