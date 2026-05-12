package internal

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	aphrodite "github.com/jonathon-chew/Aphrodite"
)

func Ping(rest_of_command []string) {

	ipRegex := "(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)"
	ipAddress := regexp.MustCompile(ipRegex)

	timeRegex := `time=([0-9]+(?:\.[0-9]+)?)\s*ms`
	timeAddress := regexp.MustCompile(timeRegex)

	cmd := exec.Command("ping", rest_of_command...)

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	ErrRunningCommand := cmd.Run()
	if ErrRunningCommand != nil {
		aphrodite.PrintError("error running command: ping: " + ErrRunningCommand.Error() + "\n")
		aphrodite.PrintError(stderr.String())
		return
	}

	// cmd.Start()

	if stderr.Len() > 0 {
		fmt.Print("error: from the command output: ", stderr.String())
	}

	scanner := bufio.NewScanner(&stdout)

	for scanner.Scan() {
		line := scanner.Text()

		// Replace every IP match with colored version
		coloredLine := ipAddress.ReplaceAllStringFunc(line, func(ip string) string {
			return aphrodite.ReturnInfo(ip)
		})

		match := timeAddress.FindStringSubmatch(coloredLine)
		if len(match) > 1 {
			timeValue := match[1]

			number_of_seconds, _ := strconv.Atoi(strings.Split(timeValue, ".")[0])

			var colored string
			switch {
			case number_of_seconds > 25:
				colored, _ = aphrodite.ReturnColour("red", timeValue)
			case number_of_seconds > 10:
				colored, _ = aphrodite.ReturnColour("yellow", timeValue)
			case number_of_seconds > 5:
				colored, _ = aphrodite.ReturnColour("green", timeValue)
			default:
				colored, _ = aphrodite.ReturnColour("blue", timeValue)
			}

			// replace ONLY the number inside the original string
			coloredLine = strings.Replace(coloredLine, timeValue, colored, 1)
		}

		fmt.Println(coloredLine)
	}

	if err := scanner.Err(); err != nil {
		aphrodite.PrintError(err.Error())
	}
}
