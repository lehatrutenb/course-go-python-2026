package service

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/lehatrutenb/course-go-python-2026/homeworks/hw1/dice/internal/config"
	"github.com/lehatrutenb/course-go-python-2026/homeworks/hw1/dice/internal/tools"
	"go.uber.org/zap"
)

type DiceRoller struct {
	cfg config.Config
	lg  *zap.SugaredLogger
}

func NewDiceRoller(cfg config.Config, lg *zap.SugaredLogger) DiceRoller {
	return DiceRoller{cfg: cfg, lg: lg}
}

func (dr DiceRoller) RollDice() {
	rolls := make(chan chan int, dr.cfg.Count)
	go func() {
		closeCh := sync.Once{}
		defer closeCh.Do(func() {
			close(rolls)
		})
		for i := 0; i < dr.cfg.Count; i++ {
			curRolls := make(chan int, dr.cfg.DiceAmt)
			closeCh := sync.Once{}
			defer closeCh.Do(func() {
				close(curRolls)
			})
			rolls <- curRolls
			for j := 0; j < dr.cfg.DiceAmt; j++ {
				curRolls <- tools.GenerateRandomInt(1, dr.cfg.DiceSides)
			}
			closeCh.Do(func() {
				close(curRolls)
			})
		}
		closeCh.Do(func() {
			close(rolls)
		})
	}()
	dr.prettyPrintIterations(rolls)
}

// I suppose it is not expected that we will write own tests for printing, so not write complex code.

func (dr DiceRoller) needToPrintIteration() bool {
	return dr.cfg.Count != 1
}

func (dr DiceRoller) prettyPrintIterations(rolls <-chan chan int) {
	if dr.cfg.Verbose {
		if dr.needToPrintIteration() {
			fmt.Printf("Rolling %dd%d (%d iterations):\n", dr.cfg.DiceAmt, dr.cfg.DiceSides, dr.cfg.Count)
		} else {
			fmt.Printf("Rolling %dd%d:\n", dr.cfg.DiceAmt, dr.cfg.DiceSides)
		}
	}
	i := 0
	for iteration := range rolls {
		if dr.needToPrintIteration() {
			fmt.Printf("Iteration %d:\n", i+1)
		}
		out := dr.prettyPrintRolls(iteration)
		for line := range out {
			if dr.needToPrintIteration() {
				fmt.Printf("  %s\n", line)
			} else {
				fmt.Println(line)
			}
		}
		i++
	}
}

func (dr DiceRoller) prettyPrintRolls(rolls <-chan int) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		sum := 0
		i := 0
		if dr.cfg.Verbose {
			for roll := range rolls {
				out <- fmt.Sprintf("Dice %d: %d", i+1, roll)
				sum += roll
				i++
			}
			out <- fmt.Sprintf("Sum: %d", sum)
			out <- fmt.Sprintf("Average: %f", float64(sum)/float64(i)) // TODO: do we need rounding?
			return
		}
		if dr.cfg.Sum {
			for roll := range rolls {
				sum += roll
			}
			out <- fmt.Sprintf("Sum: %d", sum)
			return
		}

		sb := strings.Builder{}
		sb.WriteString("Rolls: ")
		for roll := range rolls {
			if i > 0 {
				sb.WriteString("+")
			}
			sb.WriteString(strconv.Itoa(roll))
			sum += roll
			i++
		}
		if i > 1 {
			sb.WriteString("=")
			sb.WriteString(strconv.Itoa(sum))
		}
		out <- sb.String()
	}()
	return out
}
