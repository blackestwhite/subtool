package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

type Subtitle struct {
	Index int    `json:"index"`
	Start string `json:"start"`
	End   string `json:"end"`
	Text  string `json:"text"`
}

func main() {
	// Define flags for input and output filenames
	inputFile := flag.String("i", "", "Input filename (e.g., vid.srt)")
	outputFile := flag.String("o", "", "Output filename (e.g., output.json)")

	// Parse the command-line arguments
	flag.Parse()

	// Check if required flags are provided
	if *inputFile == "" || *outputFile == "" || !strings.HasSuffix(*outputFile, ".json") || !strings.HasSuffix(*inputFile, ".srt") {
		fmt.Println("Usage: subtool -i <input filename> -o <output filename>")
		return
	}

	// Read the SRT file content
	srtContent, err := os.ReadFile(*inputFile)
	if err != nil {
		log.Fatal(err)
	}

	// Normalize line endings to '\n'
	srtContentStr := strings.ReplaceAll(string(srtContent), "\r\n", "\n")

	// Split the SRT content into subtitles
	subtitles := strings.Split(srtContentStr, "\n\n")

	var subtitleList []Subtitle

	// Process each subtitle
	for _, subtitleBlock := range subtitles {
		// Use regex to extract index, start, end, and text
		re := regexp.MustCompile(`(\d+)\n(\d{2}:\d{2}:\d{2},\d{3}) --> (\d{2}:\d{2}:\d{2},\d{3})\n((.+\s)+)`)
		matches := re.FindStringSubmatch(subtitleBlock)

		if len(matches) == 6 {
			index := atoi(matches[1])
			start := matches[2]
			end := matches[3]
			text := matches[4]

			subtitle := Subtitle{
				Index: index,
				Start: start,
				End:   end,
				Text:  text,
			}

			subtitleList = append(subtitleList, subtitle)
		}
	}

	// Convert subtitles to JSON
	jsonData, err := json.Marshal(subtitleList)
	if err != nil {
		log.Fatal(err)
	}

	// Write JSON to a new file
	err = os.WriteFile(*outputFile, jsonData, 0644)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Conversion successful. JSON file created.")
}

func atoi(s string) int {
	result := 0
	for _, c := range s {
		result = result*10 + int(c-'0')
	}
	return result
}
