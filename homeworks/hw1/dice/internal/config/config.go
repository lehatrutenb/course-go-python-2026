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
	return checkConfig(lg, Config{Count: count, Seed: seed, Sum: sum, Verbose: verbose, DiceAmt: diceAmt, DiceSides: diceSides})
}

func checkConfig(lg *zap.SugaredLogger, cfg Config) (Config, error) {
	if cfg.DiceAmt <= 0 || cfg.DiceAmt > 100 {
		lg.Errorw("Dice amount must be greater than 0 and less than 100", "diceAmt", cfg.DiceAmt)
		return Config{}, fmt.Errorf("Dice amount must be greater than 0 and less than 100")
	}
	if cfg.DiceSides <= 0 || cfg.DiceSides > 100 {
		lg.Errorw("Dice sides must be greater than 0 and less than 100", "diceSides", cfg.DiceSides)
		return Config{}, fmt.Errorf("Dice sides must be greater than 0 and less than 100")
	}
	if cfg.Sum && cfg.Verbose {
		lg.Error("Sum and verbose cannot be true at the same time")
		return Config{}, fmt.Errorf("Sum and verbose cannot be true at the same time")
	}
	return cfg, nil
}

func parseDiceAmtSides(lg *zap.SugaredLogger, diceAmtSides string) (int, int, error) {
	parts := strings.Split(diceAmtSides, "d")
	if len(parts) != 2 {
		lg.Errorw("Failed to parse dice amount or sides", "diceAmtSides", diceAmtSides)
		return 0, 0, fmt.Errorf("Invalid dice amount or sides: %s", diceAmtSides)
	}
	diceAmtInt, err := strconv.Atoi(parts[0])
	if err != nil {
		lg.Errorw("Invalid dice amount", "diceAmtSides", diceAmtSides, "error", err)
		return 0, 0, err
	}
	diceSidesInt, err := strconv.Atoi(parts[1])
	if err != nil {
		lg.Errorw("Invalid dice sides", "diceAmtSides", diceAmtSides, "error", err)
		return 0, 0, err
	}
	return diceAmtInt, diceSidesInt, nil
}
