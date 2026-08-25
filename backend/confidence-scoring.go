package main

import (
	"math"
	"sync"
)

func confidenceScore(logprobs []float32) (float32, float32) {
	var avg_cf, prplx float32
	var wg sync.WaitGroup
	wg.Add(2)

	go func(avg_cf *float32, wg *sync.WaitGroup) {
		defer wg.Done()

		total := 0.00
		count := 0
		for _, p := range logprobs {
			total += math.Exp(float64(p))
			count += 1
		}

		*avg_cf = float32(total) / float32(count)
	}(&avg_cf, &wg)

	go func(prplx *float32, wg *sync.WaitGroup) {
		defer wg.Done()

		total := float32(0.00)
		for _, p := range logprobs {
			total += p
		}
		*prplx = float32(math.Exp(float64(total)))
	}(&prplx, &wg)

	return avg_cf, prplx
}
