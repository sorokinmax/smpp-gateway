package main

import (
	"bufio"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

// ReadFileLines reads file by line
func ReadFileLines(path string) (lines []string, err error) {
	file, err := os.Open(path)

	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Start reading from the file with a reader.
	reader := bufio.NewReader(file)

	var line string
	for {
		line, err = reader.ReadString('\n')

		log.Printf(" > Read %d characters\n", len(line))

		// Process the line here.
		lines = strings.Split(strings.Replace(line, "\r\n", "\n", -1), "\n")
		if err != nil {
			break
		}
	}

	if err != io.EOF {
		log.Fatalf(" > Failed!: %v\n", err)
	}
	err = nil

	return lines, err
}

func ReadFile(filename string) string {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	return string(data)
}
