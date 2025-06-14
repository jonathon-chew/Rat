package cmd

import (
	"fmt"
	_ "os"
	"strings"
)

func ParseArguments(Arguments []string) string {
	for _, argument := range Arguments {
		fmt.Println(argument)
		if argument == strings.ToLower("--allow") || argument == strings.ToLower("--force") {
			return "Allow"
		}
	}

	return "Don't Allow"
}
