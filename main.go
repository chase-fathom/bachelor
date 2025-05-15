package main

import (
    "github.com/yourusername/bachelor-sim/game"
)

func main() {
    state := game.NewGameState()

    game.CreatePlayerCharacter(&state)
    game.IntroduceContestants(&state)
    game.RunEpisode1(&state)
    game.RunRoseCeremony(&state)
}
