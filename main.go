package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type lemgram_data struct {
	compound_analysis  bool
	lemgram            string
	part_of_speech     []string
	relative_frequency float64
}

func unique_slice(a []string) []string {
	unique_map := make(map[string]bool)

	for _, key := range a {
		unique_map[key] = true
	}

	unique_slice := make([]string, 0)

	for key := range unique_map {
		unique_slice = append(unique_slice, key)
	}

	return unique_slice
}

func parse_lemgram(s string) ([]string, []string) {
	lemgram_slice := make([]string, 0)
	aaaaaaa_slice := make([]string, 0)

	for _, lemgram := range strings.FieldsFunc(s, func(r rune) bool { return r == '|' }) {
		lemgram_slice = append(lemgram_slice, strings.Split(lemgram, "..")[0])
		aaaaaaa_slice = append(aaaaaaa_slice, strings.Split(lemgram, "..")[1])
	}

	return unique_slice(lemgram_slice), unique_slice(aaaaaaa_slice)
}

func Round(x, unit float64) float64 {
	return math.Round(x/unit) * unit
}

func generate_lemgram_map() map[string]lemgram_data {
	lemgram_map := make(map[string]lemgram_data)

	matches, _ := filepath.Glob(".corpus/*.txt")
	for _, text := range matches {
		file, _ := os.Open(text)
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			row := scanner.Text()
			column := strings.Split(row, "\t")

			lemgram, a := parse_lemgram(column[2])

			for _, key := range lemgram {
				t := new(lemgram_data)
				if !lemgram_map[key].compound_analysis && strings.Contains(column[3], "+") {
					t.compound_analysis = true
				}
				t.lemgram = key
				t.part_of_speech = append(lemgram_map[key].part_of_speech, a...)
				t.relative_frequency, _ = strconv.ParseFloat(column[5], 64)
				t.relative_frequency += lemgram_map[key].relative_frequency

				lemgram_map[key] = *t
			}
		}
	}

	for key, value := range lemgram_map {
		value.part_of_speech = unique_slice(value.part_of_speech)
		value.relative_frequency /= float64(len(matches))
		value.relative_frequency = Round(value.relative_frequency, 0.0001)
		lemgram_map[key] = value
	}

	return lemgram_map
}

func alphabetic(s string) bool {
	for _, rune := range s {
		if !strings.ContainsRune("abcdefghijklmnopqrstuvwxyzåäö_", rune) {
			return false
		}
	}

	return true
}

func selective(input []string) bool {
	allowed := []string{"ab.", "abm.", "al.", "av.", "avm.", "nn.", "nnm.", "pn.", "pp.", "ppm.", "vb.", "vbm."}

	for _, s := range input {
		if containsAny(s, allowed) {
			return true
		}
	}
	return false
}

func containsAny(s string, checks []string) bool {
	for _, check := range checks {
		if strings.Contains(s, check) {
			return true
		}
	}
	return false
}

func new_map(m map[string]lemgram_data) map[string]lemgram_data {
	new_map := make(map[string]lemgram_data)

	for key, value := range m {
		if value.isValid() {
			continue
		}

		new_map[key] = value
	}

	return new_map
}

func (data *lemgram_data) isValid() bool {
	return len([]rune(data.lemgram)) < 4 ||
		data.relative_frequency > 0061.0658 ||
		data.relative_frequency < 0000.0118 ||
		!alphabetic(data.lemgram) ||
		!selective(data.part_of_speech)
}

func print_lemgram_slice(m map[string]lemgram_data, n uint) {
	for _, value := range m {
		if n == 0 {
			break
		}

		fmt.Printf("%09.4f %v\n", value.relative_frequency, strings.ReplaceAll(value.lemgram, "_", " "))

		n--
	}
}

func main() {
	print_lemgram_slice(new_map(generate_lemgram_map()), 10)
}
