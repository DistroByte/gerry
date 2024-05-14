package utils

import (
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"os"
)

type MarkovChain interface {
	AddWordTransition(from string, to string)
	AddWordTransitions(mappings map[string]string)
	ChooseNextWord(word string)
	Flush()
}

type MarkovChainImpl struct {
	transitions map[string]map[string]float64
}

func (mc *MarkovChainImpl) AddWordTransition(from, to string) {
	if mc.transitions[from] == nil {
		mc.transitions[from] = make(map[string]float64)
	} else {
		mc.transitions[from][to]++
	}

	fmt.Printf("[Before] %s: %v\n", from, mc.transitions[from])

	total := 0.0
	for _, count := range mc.transitions[from] {
		total += count
	}
	for word, count := range mc.transitions[from] {
		mc.transitions[from][word] = count / total
	}
	fmt.Printf("[After] %s: %v\n", from, mc.transitions[from])
}

func (mc *MarkovChainImpl) AddWordTransitions(mappings map[string]string) {
	for from, to := range mappings {
		mc.AddWordTransition(from, to)
	}
}

func absDifF(a float64, b float64) float64 {
	if a > b {
		return a - b
	}
	if b > a {
		return b - a
	}
	return 0
}

func findNearestKeyForWeight(weight float64, mapping map[string]float64) string {
	keys := make([]string, len(mapping))

	// 1 in 10 chance
	madnessChance := (rand.IntN(2)) == 0 && len(keys) > 1

	i := 0
	for k := range mapping {
		keys[i] = k
		i++
	}

	if len(keys) == 0 {
		return "EOF"
	}

	var bestKey string = keys[0]
	var lastDistance float64 = 100.0

	for key, keyWeight := range mapping {
		currentDistance := absDifF(keyWeight, weight)
		if currentDistance < lastDistance {
			lastDistance = currentDistance
			bestKey = key
		}
	}

	if madnessChance && len(keys) > 1 {
		bestKey = keys[rand.IntN(len(keys)-1)]
	}

	return bestKey
}

func (mc *MarkovChainImpl) ChooseNextWord(word string) string {
	transitions := mc.transitions[word]
	var weightTotal = 0.0

	for _, weight := range transitions {
		weightTotal += weight
	}
	segment := weightTotal * rand.Float64()

	return findNearestKeyForWeight(segment, transitions)
}

func (mc *MarkovChainImpl) ChooseFirstWord() string {
	keys := make([]string, 0, len(mc.transitions))
	for k := range mc.transitions {
		keys = append(keys, k)
	}

	return keys[rand.IntN(len(keys))]
}

func (mc *MarkovChainImpl) LoadModel(file string) {
	fmt.Printf("Training from file %s\n", file)
	transitions, err := loadModel(file)
	if err != nil {
		panic(err)
	}
	if mc.transitions == nil {
		mc.transitions = transitions
	} else {
		mc.Merge(transitions)
	}
}

func loadModel(file string) (map[string]map[string]float64, error) {
	var chain map[string]map[string]float64
	f, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	f.Close()

	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(f.Name())

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &chain)
	if err != nil {
		fmt.Println("Creating new database since none was found")
		return map[string]map[string]float64{}, nil
	}

	return chain, nil
}

// You are a liar
// You are, are a, a liar

func (mc *MarkovChainImpl) ImportMessage(message []string) {
	fmt.Printf("Training on the message %s\n", message)
	var transitions = make(map[string]string, len(message)-1)
	for i := range len(message) - 1 {
		transitions[message[i]] = message[i+1]
	}
	mc.AddWordTransitions(transitions)
}

func (mc MarkovChainImpl) Flush() {
	fmt.Println("Flushing!")
	data, err := json.Marshal(&mc.transitions)
	if err != nil {
		fmt.Printf("Error marshalling to JSON: %v\n", err)
	}

	err = os.WriteFile("markovChain.json", data, 0644)
	if err != nil {
		fmt.Printf("Error writing to file: %v\n", err)
	}
}

func (mc *MarkovChainImpl) Merge(other map[string]map[string]float64) {
	for word, transitions := range other {
		if mc.transitions[word] == nil {
			mc.transitions[word] = make(map[string]float64)
		}
		for otherWord, otherWeight := range transitions {
			weight, exists := mc.transitions[word][otherWord]
			if exists {
				mc.transitions[word][otherWord] = (weight + otherWeight) / 2
			} else {
				mc.transitions[word][otherWord] = otherWeight
			}
		}
	}
}

func NewMarkovChain() *MarkovChainImpl {
	return &MarkovChainImpl{} // Return a pointer to MarkovChainImpl
}
