package main

import (
	"math"
)

// confidenceScore calculates both the confidence score (0.0 to 1.0) and perplexity (>= 1.0)
// from a slice of token logprobs.
//
// Perplexity (PPL) is calculated as exp(-avg(logprob)).
// Confidence is the geometric mean token probability: C = 1 / PPL = exp(avg(logprob)).
// When PPL = 1.0, confidence is 1.0 (100%).
// As PPL increases (> 1.0), confidence decreases smoothly towards 0.0.
func confidenceScore(logprobs []Logprobs) float32 {
	if len(logprobs) == 0 {
		return 1.0
	}

	var sumLogprobs float64
	for _, v := range logprobs {
		sumLogprobs += v.Logprob
	}

	avgLogprob := sumLogprobs / float64(len(logprobs))
	perplexity := math.Exp(-avgLogprob)

	// Clamp perplexity to minimum 1.0 in case of float rounding errors
	if perplexity < 1.0 || math.IsNaN(perplexity) {
		perplexity = 1.0
	}

	// Option A: Geometric mean token probability (C = 1 / PPL)
	confidence := 1.0 / perplexity
	if confidence > 1.0 {
		confidence = 1.0
	} else if confidence < 0.0 || math.IsNaN(confidence) {
		confidence = 0.0
	}

	return float32(confidence)
}
