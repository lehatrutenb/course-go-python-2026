package tools

import (
	"math/rand"
	"time"
)

// Unused value, but i suppose globally good practice
var randomGenerator = rand.New(rand.NewSource(time.Now().UnixNano()))

func GenerateRandomInt(min, max int) int {
	return randomGenerator.Intn(max-min+1) + min
}

func SetRandomGeneratorSeed(seed int64) {
	randomGenerator = rand.New(rand.NewSource(seed))
}
