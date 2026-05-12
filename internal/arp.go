package internal

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"regexp"

	aphrodite "github.com/jonathon-chew/Aphrodite"
)

func Arp(rest_of_command []string) {

	ipRegex := "(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)"
	ipAddress := regexp.MustCompile(ipRegex)

	macRegex := "(?:[0-9A-Fa-f]{1,2}[:-]){5}[0-9A-Fa-f]{1,2}"
	macAddress := regexp.MustCompile(macRegex)

	nameRegex := "^(.+?){1}\\("
	nameValue := regexp.MustCompile(nameRegex)

	cmd := exec.Command("arp", rest_of_command...)

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	ErrRunningCommand := cmd.Run()
	if ErrRunningCommand != nil {
		aphrodite.PrintError("error running command: arp: " + ErrRunningCommand.Error() + "\n")
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

		// Replace every IP match with colored version
		coloredLine = macAddress.ReplaceAllStringFunc(coloredLine, func(mac string) string {
			returnMe, _ := aphrodite.ReturnColour("blue", mac)
			return returnMe
		})

		coloredLine = nameValue.ReplaceAllStringFunc(coloredLine, func(name string) string {
			returnMe, _ := aphrodite.ReturnColour("yellow", name)
			return returnMe
		})

		fmt.Println(coloredLine)
	}

	if err := scanner.Err(); err != nil {
		aphrodite.PrintError(err.Error())
	}
}
