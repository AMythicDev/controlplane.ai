package main

import (
	"math"
	"sync"
)

func confidenceScore(logprobs []Logprobs) float64 {
	px := 0.0
	for _, v := range logprobs {
		px += v.Logprob
	}
	if len(logprobs) > 0 {
		return math.Exp(-px / len(logprobs))
	} else {
		return 1.0
	}
}
