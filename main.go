/**
 * This code was written by Dason Woodhouse and is licensed 
 * under the GNU General Public License, or GPLv3. Please
 * abide by these terms.
 */

package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"
)


//Separate array length from the actual number of members
type Pool struct {
	MaxMembers int `json:"-"`
	Members []int
}

/**
 * Returns false if can't add any more.
 * Returns true if addition was successful
 */
func(p * Pool) Push(ball int) bool {
	if len(p.Members) >= p.MaxMembers {
		return false
	}

	p.Members = append(p.Members, ball)

	return true
}

/**
 * Returns Ball{-1}, false, if can't pop.
 * Returns Ball, true if pop was successful.
 */
func(p * Pool) Pop() (int, bool) {
	num := len(p.Members)
	if num < 1 {
		return -1, false
	}

	last := p.Members[num-1]
	p.Members = p.Members[:num-1]
	return last, true
}

func(p * Pool) FlushPool(target * Pool ) {
	current, hasMore := p.Pop()

	for hasMore {
		if !target.Push(current) {
			panic("failed to add a ball back to the full pool. please report bug to dwood15 (github)")
		}

		current, hasMore = p.Pop()
	}
}

func populateFullPool(numBalls int) Pool {
	fullPool := make([]int, numBalls)
	for i :=0; i < numBalls; i++ {
		fullPool[i] = i
	}
	return Pool{
		MaxMembers:numBalls,
		Members:fullPool,
	}
}

func (p * Pool) PopFirst() (int, bool) {
	if len(p.Members) <= 0 {
		return -1, false
	}

	first := (p.Members)[0]
	newPool := p.Members
	p.Members = newPool[1:]

	return first, true
}

type Clock struct {
	Min Pool
	FiveMin Pool
	Hour Pool
	Main Pool
}

type ForJsonDisplay struct {
	Min []int `json:"Min"`
	FiveMin []int `json:"FiveMin"`
	Hour []int `json:"Hour"`
	Main []int `json:"Main"`
}

func (c * Clock) PrintJson() {
	display := ForJsonDisplay{
		c.Min.Members,
		c.FiveMin.Members,
		c.Hour.Members,
		c.Main.Members,
	}

  	fmt.Println("made by dwood15")
	dispText, _ := json.Marshal(display)
	fmt.Println(string(dispText))
}

/**
 * Adds a minute to the clock
 */
func (c * Clock)AddMinute() bool {
	newMin, success := c.Main.PopFirst()

	if !success {
		panic("failed to get the first ball from the pool.")
	}

	if c.Min.Push(newMin) {
		//We added a new minute! we're good.
		return false
	}

	//We made it this far, we have to:
	// 1) Flush minute pool in reverse order.

	c.Min.FlushPool(&c.Main)

	// 2) Add to the fifthPool now.
	if c.FiveMin.Push(newMin) {
		//We added to fifthPool successfully
		return false
	}

	c.FiveMin.FlushPool(&c.Main)

	if c.Hour.Push(newMin) {
		return false
	}

	c.Hour.FlushPool(&c.Main)

	c.Main.Push(newMin)
	//This means that the hour pool (hours 1-12) have been flushed
	return true
}


/**
 * Assumes the number of balls already been verified.
 */
func runSimulation(numBalls int, numMins int) {
	fullPool := populateFullPool(numBalls)

	//Up to 4 minutes
	minutePool := Pool{
		4,
		[]int{},
	}
	//Up to 48 minutes
	fiveMinPool := Pool{
		11,
		[]int{},
	}
	//Up to 12 hours
	hourPool := Pool{
		11,
		[]int{},
	}

	clock := Clock{
		Min:minutePool,
		FiveMin:fiveMinPool,
		Hour:hourPool,
		Main:fullPool,
	}

	if numMins > -1 {
		for remaining := numMins; remaining > 0; remaining-- {
			clock.AddMinute()
		}

		clock.PrintJson()
		fmt.Println("\nSimulation Complete!")
		return
	}

	//Keep Adding Minute until we get a repeat
	//I (Dason) could break this into a goroutine
	numMinutes := 0
	start := time.Now()
	for gotRepeat := false; gotRepeat != true;  {
		clock.AddMinute()
		numMinutes++

		if len(clock.Main.Members) == numBalls {
			gotRepeat = true

			for i := 0; i < numBalls; i++ {
				gotRepeat = i == clock.Main.Members[i]
			}
		}
	}

	elapsed := time.Since(start)

	balls := strconv.Itoa(numBalls)
	days := strconv.Itoa(numMinutes / (60 * 24))
	milli := elapsed / time.Millisecond
	seconds := float64(milli) / float64(1000)
	formattedSecs := strconv.FormatFloat(seconds, 'f', 3, 64)

	fmt.Println(balls + " cycle after " + days + " days.")
	fmt.Println("Completed in " + strconv.Itoa(int(milli)) + " milliseconds (" + formattedSecs + ") seconds" )

}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Please type 1 for Mode 1. 2 For Mode 2, then press Enter")
	scanner.Scan()
	text := scanner.Text()

	mode, err := strconv.Atoi(text)

	if err != nil || (mode != 1 && mode != 2) {
		fmt.Println("Invalid input entered. please restart and try again.")
		return
	}

	fmt.Println("Please type number of balls and press enter")
	scanner.Scan()
	text = scanner.Text()

  //I (Dason) won't use the := operator, because err has already been declared.
	var numBalls int

	numBalls, err = strconv.Atoi(text)

	if err!= nil || (numBalls < 27 || numBalls > 127) {
		fmt.Println("Invalid number for number of balls entered. Please restart and try again.")
		return
	}

	var numMins int
	if mode == 2 {
		fmt.Println("Please type number of minutes and press enter")
		scanner.Scan()
		text = scanner.Text()
		numMins, err = strconv.Atoi(text)

		if err != nil || numMins < 0 {
			fmt.Println("Invalid number for number of minutes entered. Please restart and try again.")
			return
		}
	} else {
		//This variable should be ignored if numMins is not provided.
		numMins = -1
	}
	runSimulation(numBalls, numMins)
}
