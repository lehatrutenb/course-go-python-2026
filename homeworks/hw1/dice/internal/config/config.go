package config

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/lehatrutenb/course-go-python-2026/homeworks/hw1/dice/internal/tools"
	"go.uber.org/zap"
)

type Config struct {
	Count     int
	Seed      int64
	Sum       bool
	Verbose   bool
	DiceAmt   int
	DiceSides int
}

func NewConfig(lg *zap.SugaredLogger, count int, seed int64, sum bool, verbose bool, diceAmtSides string) (Config, error) {
	diceAmt, diceSides, err := parseDiceAmtSides(lg, diceAmtSides)
	if err != nil {
		return Config{}, err
	}
	tools.SetRandomGeneratorSeed(seed)
	return Config{Count: count, Seed: seed, Sum: sum, Verbose: verbose, DiceAmt: diceAmt, DiceSides: diceSides}, nil
}

func parseDiceAmtSides(lg *zap.SugaredLogger, diceAmtSides string) (int, int, error) {
	parts := strings.Split(diceAmtSides, "d")
	if len(parts) != 2 {
		lg.Errorw("failed to parse dice amount or sides", "diceAmtSides", diceAmtSides)
		return 0, 0, fmt.Errorf("invalid dice amount or sides: %s", diceAmtSides)
	}
	diceAmtInt, err := strconv.Atoi(parts[0])
	if err != nil {
		lg.Errorw("invalid dice amount", "diceAmtSides", diceAmtSides, "error", err)
		return 0, 0, err
	}
	diceSidesInt, err := strconv.Atoi(parts[1])
	if err != nil {
		lg.Errorw("invalid dice sides", "diceAmtSides", diceAmtSides, "error", err)
		return 0, 0, err
	}
	return diceAmtInt, diceSidesInt, nil
}
