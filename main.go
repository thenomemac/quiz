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
	csvFilename := flag.String("csv", "problems.csv", "a csv file in fomat 'Q&A'")
	timeLimit := flag.Int("limit", 30, "the time limit for the quiz in seconds.")
	flag.Parse()

	file, err := os.Open(*csvFilename)
	if err != nil {
		exit(fmt.Sprintf("Failed to open CSV: %s", *csvFilename))
	}

	r := csv.NewReader(file)

	lines, err := r.ReadAll()
	if err != nil {
		exit(fmt.Sprintf("Failed to parse CSV: %s", *csvFilename))
	}

	problems := parseLinesCSV(lines)

	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

	correct := 0
	for i, p := range problems {

		fmt.Printf("Problem #%d: %s = \n", i+1, p.q)

		ansChannel := make(chan string)
		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)
			ansChannel <- answer
		}()

		select {
		case <-timer.C:
			fmt.Println("Fudge! out of time.")
			fmt.Printf("\nYou scored %d out of %d.\nThat's %.01f%% correct.\n",
				correct, len(problems),
				float64(correct)/float64(len(problems))*100)
			return
		case answer := <-ansChannel:
			if answer == p.a {
				correct++
			}
		}
	}

	fmt.Printf("You scored %d out of %d.\nThat's %.01f%% correct.\n",
		correct, len(problems),
		float64(correct)/float64(len(problems))*100)
}

func parseLinesCSV(lines [][]string) []problem {
	ret := make([]problem, len(lines))
	for i, line := range lines {
		ret[i] = problem{
			q: line[0],
			a: strings.TrimSpace(line[1]),
		}
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
