# Multiplayer Elo Rating System

This is a simple implementation of the Elo rating system for multiplayer games. The system is based on the [Elo rating system](https://en.wikipedia.org/wiki/Elo_rating_system) which is a method for calculating the relative skill levels of players in two-player games such as chess. The system is used in many games and sports to rank players based on their performance.

[![Go Report Card](https://goreportcard.com/badge/github.com/distrobyte/multielo)](https://goreportcard.com/report/github.com/distrobyte/multielo)

## Installation

```bash
go get github.com/distrobyte/multielo
```

## Usage

```go
package main

import (
    "fmt"
    elo "github.com/distrobyte/multielo"
)

func main() {
    // Create a new league
    league := elo.NewLeague()

    // Add players to the league
    league.AddPlayer("player1")
    league.AddPlayer("player2")

    // Get the ELO of a player
    player := league.GetPlayer("player1")
    fmt.Println(player.ELO)
}
```
