package commands

import (
	"log/slog"
	"strconv"

	"github.com/DistroByte/gerry/elo"
	"github.com/DistroByte/gerry/internal/models"
)

var K = 32.0
var D = 400.0
var ScoringFunctionBase = 1.0
var LogBase = 10.0

func KartingCommand(args []string, message models.Message) string {
	// calculate the new elo for each racer
	elo := elo.NewMultiElo(K, D, ScoringFunctionBase, LogBase, nil)

	// get the initial ratings
	initialRatings := make([]float64, len(args))
	for i, arg := range args {
		rating, err := strconv.ParseFloat(arg, 64)
		if err != nil {
			slog.Error(err.Error())
			return "Invalid rating"
		}
		initialRatings[i] = rating
	}

	// get the result order
	resultOrder := make([]int, len(args))
	for i := range resultOrder {
		resultOrder[i] = i
	}

	// get the new ratings
	newRatings := elo.GetNewRatings(initialRatings, resultOrder)

	// format the response
	response := "New ratings: "
	for i, rating := range newRatings {
		response += strconv.FormatFloat(rating, 'f', 2, 64)
		if i < len(newRatings)-1 {
			response += ", "
		}

	}

	// return the response

	return response
}
