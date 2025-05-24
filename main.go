package main

import (
    "github.com/yourusername/bachelor-sim/game"
)

func main() {
    state := game.NewGameState()

    game.RunIntroduction(&state)
    game.CreatePlayerCharacter(&state)

    game.GenerateContestants(&state)
    game.IntroduceContestants(&state)
    game.IntroduceBachelor(&state)

    game.RunFirstImpression(&state)
    game.RunSession1(&state)
    game.RunSession2(&state)
    game.RunSession3(&state)
    game.RunFantasySuites(&state)
    game.RunProposal(&state)
}
