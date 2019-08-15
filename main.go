package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {

	fileName := flag.String("csv", "problems.csv", "a csv file in the format of 'question,answer'")
	flag.Parse()

	file, err := os.Open(*fileName)
	if err != nil {
		fmt.Printf("can't open csv file %s\n", *fileName)
		os.Exit(1)
	}

	csvReader := csv.NewReader(file)
	lines, err := csvReader.ReadAll()
	if err != nil {
		fmt.Printf("read csv file error")
		os.Exit(1)
	}

	problems := parseLine(lines)

	correct := 0

	for i, p := range problems {
		fmt.Printf("problem#%d: %s\n", i+1, p.question)
		var userAnswer string
		_, err := fmt.Scanln(&userAnswer)
		if err != nil && err.Error() != "unexpected newline" {
			fmt.Println(err)
			os.Exit(1)
		}

		if userAnswer == p.answer {
			correct++
		}
	}

	fmt.Printf("You scored %d out of %d.\n", correct, len(problems))
}

func parseLine(lines [][]string) []problem {

	problems := make([]problem, len(lines))

	for i, v := range lines {
		problems[i] = problem{
			question: v[0],
			answer:   strings.TrimSpace(v[1]),
		}
	}

	return problems
}

type problem struct {
	question string
	answer   string
}
