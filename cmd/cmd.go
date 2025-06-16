package cmd

import (
	"fmt"
	_ "os"
	"strings"
)

func ParseArguments(Arguments []string) string {
	for _, argument := range Arguments {
		if argument == strings.ToLower("--allow") || argument == strings.ToLower("--force") {
			if !strings.Contains(argument,"--") {
				fmt.Println(argument)
			}
			return "Allow"
		}
	}

	return "Don't Allow"
}
