package game

import (
    "fmt"
		"github.com/charmbracelet/huh"
)

func RunEpisode1(state *GameState) {
	fmt.Println("🌴 Episode 1: The Beach Party 🏖️")
	fmt.Println("The contestants are all gathered by the beach, ready for a group date! The sun is shining, and the vibes are immaculate.")
	
	// Dynamically generate the introduction based on relationship scores
	for _, contestant := range state.Contestants {
			relationshipScore := state.Relationship[contestant.Name]
			if contestant.IsPlayer {
					fmt.Println("\nIt's your turn to shine! The contestants gather around you.")
			} else {
					if relationshipScore > 2 {
							fmt.Printf("\n%s greets you with a warm smile, their relationship score is high!\n", contestant.Name)
					} else {
							fmt.Printf("\n%s barely acknowledges your presence. Their relationship score is low.\n", contestant.Name)
					}
			}
	}

	// Decision time: Player decides whom to talk to
	var choice string
	form := huh.NewForm(
			huh.NewGroup(
					huh.NewSelect[string]().
							Title("Who would you like to spend time with?").
							Options(
									huh.NewOption("Ashley", "Ashley"),
									huh.NewOption("Jordan", "Jordan"),
									huh.NewOption("Kai", "Kai"),
							).
							Value(&choice),
			),
	)
	form.Run()

	// Consequences of choice: Dialogue based on character stats
	for _, c := range state.Contestants {
			if c.Name == choice {
					ReactToPlayerChoice(c, state)
			}
	}

	// End of Episode 1 Wrap-up: Update relationships
	fmt.Println("\n🌹 The group date ends, and everyone heads back to the mansion.")
	for _, c := range state.Contestants {
			fmt.Printf("%s's current relationship score with you is: %d\n", c.Name, state.Relationship[c.Name])
	}
}

// Add this function
func RunRoseCeremony(state *GameState) {
    fmt.Println("\n🌹 The Final Rose Ceremony 🌹")
    var chosen string

    options := []string{}
    for _, c := range state.Contestants {
        if !isEliminated(state, c.Name) {
            options = append(options, c.Name)
        }
    }

    if len(options) == 0 {
        fmt.Println("Everyone has been eliminated. You are alone now. 😢")
        return
    }

    form := huh.NewForm(
        huh.NewGroup(
            huh.NewSelect[string]().
                Title("Who do you give your final rose to?").
                Options(buildOptions(options)...).
                Value(&chosen),
        ),
    )

    if err := form.Run(); err != nil {
        fmt.Println("Ceremony interrupted.")
        return
    }

    fmt.Printf("\nYou chose %s. They smile and say, \"I've been waiting for this moment 💕.\"\n", chosen)
}

// Helper to check if someone was eliminated
func isEliminated(state *GameState, name string) bool {
    for _, e := range state.Eliminated {
        if e == name {
            return true
        }
    }
    return false
}

// Helper to convert names to huh.Option
func buildOptions(names []string) []huh.Option[string] {
    opts := []huh.Option[string]{}
    for _, name := range names {
        opts = append(opts, huh.NewOption(name, name))
    }
    return opts
}

