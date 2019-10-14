package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

func main() {
	csvFilename := flag.String("csv", "problems.csv", "a csv file with questions and answers\n")
	timeLimit := flag.Int("limit", 30, "time limit of the quiz in seconds")
	shuffle := flag.Bool("rand", true, "specify if the questions should be shuffled or not")
	flag.Parse()

	file, err := os.Open(*csvFilename)

	if err != nil {
		exit(fmt.Sprintf("Failed to open the CSV file: %s\n", *csvFilename))
	}

	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	if err != nil {
		exit("Failed to parse the provided csv file")
	}

	problems := parseLines(lines)

	if *shuffle {
		rand.Seed(time.Now().UnixNano())
		problems = shuffleProblems(problems)
	}

	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

	correct := 0
	for i, p := range problems {
		fmt.Printf("Problem #%d: %s = \n", i+1, p.q)
		answerChan := make(chan string)

		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)
			answerChan <- answer
		}()

		select {
		case <-timer.C:
			fmt.Printf("\nYou scored %d out of %d.\n", correct, len(problems))
			return
		case answer := <-answerChan:
			if answer == p.a {
				correct++
			}
		}

	}

}

func shuffleProblems(problems []problem) []problem {

	shuffled := make([]problem, 0)
	count := len(problems)
	var k int
	for i := 0; i < count; i++ {
		k = rand.Intn(len(problems))
		shuffled = append(shuffled, problems[k])
		problems = append(problems[:k], problems[k+1:]...)
	}

	return shuffled
}

func parseLines(lines [][]string) []problem {
	ret := make([]problem, len(lines))
	for i, line := range lines {
		ret[i] = problem{q: line[0], a: strings.TrimSpace(line[1])}
	}
	return ret
}

type problem struct {
	q string
	a string
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
