package service

import (
	"fmt"
	"strconv"
	"strings"

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
	rolls := make([][]int, dr.cfg.Count)
	for i := 0; i < dr.cfg.Count; i++ {
		rolls[i] = make([]int, dr.cfg.DiceAmt)
		for j := 0; j < dr.cfg.DiceAmt; j++ {
			rolls[i][j] = tools.GenerateRandomInt(1, dr.cfg.DiceSides)
		}
	}
	dr.prettyPrintIterations(rolls)
}

// I suppose it is not expected that we will write own tests for printing, so not write complex code.

func (dr DiceRoller) needToPrintIteration() bool {
	return dr.cfg.Count != 1
}

func (dr DiceRoller) prettyPrintIterations(rolls [][]int) {
	for i, iteration := range rolls {
		if dr.needToPrintIteration() {
			fmt.Printf("Iteration %d:\n", i+1)
		}
		out := dr.prettyPrintRolls(iteration)
		for _, line := range out {
			if dr.needToPrintIteration() {
				fmt.Printf("  %s\n", line)
			} else {
				fmt.Println(line)
			}
		}
	}
}

func (dr DiceRoller) prettyPrintRolls(rolls []int) []string { // With channel it looks more elegant, but should run in separate goroutine
	out := make([]string, 0)
	sum := 0
	for _, roll := range rolls {
		sum += roll
	}
	if dr.cfg.Verbose {
		out = append(out, fmt.Sprintf("Rolling %dd%d:", dr.cfg.DiceAmt, dr.cfg.DiceSides))
		for i, roll := range rolls {
			out = append(out, fmt.Sprintf("Dice %d: %d", i, roll))
		}
		out = append(out, fmt.Sprintf("Sum: %d", sum))
		out = append(out, fmt.Sprintf("Average: %f", float64(sum)/float64(len(rolls)))) // TODO: do we need rounding?
		return out
	}
	if dr.cfg.Sum {
		out = append(out, fmt.Sprintf("Sum: %d", sum))
		return out
	}

	sb := strings.Builder{}
	sb.WriteString("Rolls: ")
	for i, roll := range rolls {
		if i > 0 {
			sb.WriteString("+")
		}
		sb.WriteString(strconv.Itoa(roll))
	}
	if len(rolls) > 1 {
		sb.WriteString("=")
		sb.WriteString(strconv.Itoa(sum))
	}
	out = append(out, sb.String())
	return out
}
