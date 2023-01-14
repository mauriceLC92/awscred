package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type credentials struct {
	profileName string
	accessKeyId string
	accesskey   string
}

func main() {
	file, err := os.Open(os.ExpandEnv("$HOME/.aws/credentials"))
	if err != nil {
		log.Fatalln("error opening credentials file:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "[") {
			fmt.Println(strings.Trim(line, "[]"))
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "error reading from credentials file", err)
	}
}
