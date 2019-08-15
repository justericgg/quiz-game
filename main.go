package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {

	fileName := flag.String("csv", "problems.csv", "a csv file in the format of 'question,answer'")
	limitSec := flag.Int("limit", 30, "a time limit for the quiz game in seconds")
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

	timer := time.NewTimer(time.Duration(*limitSec) * time.Second)

	correct := 0

	for i, p := range problems {
		fmt.Printf("problem#%d: %s = ", i+1, p.question)
		userAnswerChan := answerQuiz()
		select {
		case <-timer.C:
			fmt.Printf("\nYou scored %d out of %d.\n", correct, len(problems))
			return
		case userAnswer := <-userAnswerChan:

			if userAnswer.Error != nil && userAnswer.Error.Error() != "unexpected newline" {
				fmt.Println("exception occur")
				os.Exit(1)
			}

			if userAnswer.Answer == p.answer {
				correct++
			}
		}
	}

	fmt.Printf("You scored %d out of %d.\n", correct, len(problems))
}

type UserAnswer struct {
	Error  error
	Answer string
}

func answerQuiz() <-chan UserAnswer {
	answerChan := make(chan UserAnswer)

	go func() {
		defer close(answerChan)

		var userAnswer string
		_, err := fmt.Scanln(&userAnswer)

		answer := UserAnswer{
			Error:  err,
			Answer: userAnswer,
		}
		answerChan <- answer
	}()

	return answerChan
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
