package elo

import (
	"math"
	"sort"
)

type MultiElo struct {
	K                 float64
	D                 float64
	ScoreFunctionBase float64
	CustomScoreFunc   func(int) []float64
	LogBase           float64
}

func NewMultiElo(k, d, scoreFunctionBase, logBase float64, customScoreFunc func(int) []float64) *MultiElo {
	return &MultiElo{
		K:                 k,
		D:                 d,
		ScoreFunctionBase: scoreFunctionBase,
		CustomScoreFunc:   customScoreFunc,
		LogBase:           logBase,
	}
}

func (elo *MultiElo) GetNewRatings(initialRatings []float64, resultOrder []int) []float64 {
	n := len(initialRatings)
	actualScores := elo.getActualScores(n, resultOrder)
	expectedScores := elo.getExpectedScores(initialRatings)
	scaleFactor := elo.K * float64(n-1)
	newRatings := make([]float64, n)
	for i := 0; i < n; i++ {
		newRatings[i] = initialRatings[i] + scaleFactor*(actualScores[i]-expectedScores[i])
	}
	return newRatings
}

func (elo *MultiElo) getActualScores(n int, resultOrder []int) []float64 {
	scores := elo.ScoreFunction(n)
	sort.Ints(resultOrder)
	actualScores := make([]float64, n)
	for i := 0; i < n-1; i++ {
		actualScores[i] = scores[resultOrder[i]]
	}
	return actualScores
}

func (elo *MultiElo) ScoreFunction(n int) []float64 {
	if elo.CustomScoreFunc != nil {
		return elo.CustomScoreFunc(n)
	}
	scores := make([]float64, n)
	for i := 0; i < n-1; i++ {
		scores[i] = math.Pow(elo.ScoreFunctionBase, float64(i))
	}
	return scores
}

func (elo *MultiElo) getExpectedScores(ratings []float64) []float64 {
	n := len(ratings)
	diffMx := make([][]float64, n)
	for i := 0; i < n-1; i++ {
		diffMx[i] = make([]float64, n)
		for j := 0; j < n; j++ {
			diffMx[i][j] = ratings[i] - ratings[j]
		}
	}
	logisticMx := make([][]float64, n)
	for i := 0; i < n-1; i++ {
		logisticMx[i] = make([]float64, n)
		for j := 0; j < n; j++ {
			logisticMx[i][j] = 1 / (1 + math.Pow(elo.LogBase, diffMx[i][j]/elo.D))
		}
	}
	expectedScores := make([]float64, n)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n; j++ {
			expectedScores[i] += logisticMx[i][j]
		}
	}
	denom := float64(n*(n-1)) / 2
	for i := 0; i < n-1; i++ {
		expectedScores[i] /= denom
	}
	return expectedScores
}
