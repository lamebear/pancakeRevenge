package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

// Order represents a stack of pancakes, and the number of flips to show all pancakes happy side up.
type Order struct {
	NumFlips int
	Pancakes []string
}

func main() {
	currCase := 1
	fOutput, err := os.Create("output.txt")

	if err != nil {
		log.Panicln("error creating file output.txt")
	}

	defer func() {
		if err := fOutput.Close(); err != nil {
			log.Panicln("error attempting to close file because: " + err.Error())
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)

	// skip first line: using scanner, so dont need to know how many lines to iterate over
	scanner.Scan()

	for scanner.Scan() {
		currPancakes := scanner.Text()

		if err := displayNumberOfFlips(currCase, currPancakes, fOutput); err != nil {
			log.Println(err.Error())
		}

		currCase += 1
	}
}

// areValidPancakes checks whether the current pancakes are all valid (- or +).
func areValidPancakes(pancakes string) bool {
	if len(pancakes) == 0 {
		return false
	}

	re := regexp.MustCompile(`[^-+]`)

	return !re.Match([]byte(pancakes))
}

// displayNumberOfFlips prints the number of minimum flips to get the pancake stack facing happy-side up.
func displayNumberOfFlips(currCase int, pancakes string, file *os.File) error {
	if areValidPancakes(pancakes) {
		o := processOrder(pancakes)

		output := fmt.Sprintf("Case #%d: %d\n", currCase, o.NumFlips)

		if _, err := file.WriteString(output); err != nil {
			return err
		}

		if err := file.Sync(); err != nil {
			return err
		}
	} else {
		msg := fmt.Sprintf("skipping case %d due to character other than - or + detected", currCase)

		return errors.New(msg)
	}

	return nil
}

// processOrder flips 1 or more pancakes in the stack, before returning a completed Order.
func processOrder(pancakeOrder string) Order {
	var prevFace string
	o := Order{
		NumFlips: 0,
		Pancakes: strings.Split(pancakeOrder, ""),
	}

	for i, p := range o.Pancakes {
		currFace := string(p)

		// the first pancake face will always be the current and previous face
		if prevFace == "" {
			prevFace = currFace

			continue
		}

		if currFace != prevFace {
			o.flipPancakes(o.Pancakes[0:i])

			prevFace = currFace
		}
	}

	// final flip from all blank side, to happy side up
	if o.Pancakes[0] == "-" {
		o.flipPancakes(o.Pancakes)
	}

	return o
}

// flipPancakes flips over the provided pancakes, and places them back down on top of any non-flipped pancakes.
func (f *Order) flipPancakes(stack []string) {
	preFlip := strings.Join(stack, "")

	flippedPancakes := strings.Map(func(r rune) rune {
		switch {
		case r == '-':
			return '+'
		case r == '+':
			return '-'
		}

		return r
	}, preFlip)

	flippedStack := flipStack(flippedPancakes)

	for i, p := range flippedStack {
		f.Pancakes[i] = string(p)
	}

	f.NumFlips += 1
}

// flipStack takes the current pancake stack, and flips it over so the top pancake is now at the bottom.
func flipStack(p string) []string {
	s := strings.Split(p, "")
	t := len(s)
	left := t/2 - 1

	for left >= 0 {
		right := t - 1 - left
		s[left], s[right] = s[right], s[left]
		left--
	}

	return s
}
