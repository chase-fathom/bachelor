package game

import (
		"github.com/charmbracelet/huh"
    "fmt"
		"math/rand"
		"os"
		"strconv"
		"sort"
)

func WaitForEnter(title string, desc string) {
	_ = huh.NewForm(
			huh.NewGroup(
					huh.NewNote().
            	Title(title).
            	Description(desc),

					huh.NewInput().
							Prompt("↩︎ Press enter to continue").
							Value(new(string)),
			),
	).Run()
}

func ClearScreen() {
	fmt.Print("\033[H\033[2J")
}








func RunIntroduction(state *GameState) {
	ClearScreen()
	WaitForEnter("🌹 The Bachelor Simulator 🌹", "Are you ready to compete against 24 other contestants for the heart of the Bachelor? In a game of personality, charm, and a little bit of luck, see if you can be the lucky contestant to find The One in beautiful Boston, Massachusetts.")
}







func IntroduceContestants(state *GameState) {
	ClearScreen()
	var names string
	for _, c := range state.Contestants {
		if c.IsPlayer {
			names += "\033[1;36m" + c.Name + "\033[0m: " + c.Personality + ", " + c.EyeColor + "-eyed, " + c.HairColor + "-haired, " + c.Height + " " + c.Noun + ".\n"
			} else {
			names += "\033[95m" + c.Name + "\033[0m: " + c.Personality + ", " + c.EyeColor + "-eyed, " + c.HairColor + "-haired, " + c.Height + " " + c.Noun + ".\n"
		}
		t := c.Attractiveness + c.Charisma + rand.Intn(3)
		state.Relationship[c.Name] += t
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewNote().
				Title("Meeting the Contestants").
				Description("Now introducing our wonderful contestants:\n\n" + names + "\n\nDo you have what it takes to win the Bachelor's love?"),
		),
	)

	form.Run()
}








func IntroduceBachelor(state *GameState) {
	b := state.Bachelor
	var reaction string
	switch b.Attractiveness {
	case 1:
		reaction = "The contestants seem pretty unimpressed. Do they really have to compete to win the hand of a guy like this?"
	case 2:
		reaction = "The contestants look around, hoping for someone else. He's not bad, but he's not great either. Guess he'll have to do."
	case 3:
		reaction = "Not bad. The contestants finally start to look serious now that they know there is something worth competing for."
	case 4:
		reaction = "The contestants start smiling and try to get his attention, realizing that this will be a tough fight. He is pretty special."
	case 5:
		reaction = "Most contestants giggle nervously, except for you, as you stare directly into the soul of the Bachelor. This man might be The One."

	}
	ClearScreen()
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewNote().
				Title("Meeting the Bachelor").
				Description("This season, our Bachelor is really something special. I introduce to you,\n\n\033[1;92m" + b.Name + " " + b.Personality + "!\033[0m\n\n" + reaction),
		),
	)

	form.Run()

	var response string
	form = huh.NewForm(
		huh.NewGroup(
			huh.NewNote().
				Description("After his initial arrival, " + b.Name + " is mingling with the contestants and getting to know them briefly. As he walks up to you, you have just a fleeting moment to ask him a question."),

			huh.NewInput().
				Title("What do you say to the Bachleor?").
				Placeholder("e.g. hey u up?").
				Value(&response),
		),
	)

	form.Run()
	c := state.PlayerCharacter
	t := c.Attractiveness + c.Charisma

	var br string
	rn := rand.Intn(2)
	if t < 4 {
		br = "It's like he didn't even see you. You hope that he just didn't hear you, but you spoke pretty loudly. Was it too loud? Or, maybe he'll come back to talk to you . . . as you wait, you come to accept that he's not coming back to meet you."
	} else if t < 7 {
		switch rn {
		case 0:
			br = "\"Ha, you're nervous,\" \033[1;92m" + b.Name + "\033[0m says. \"I like that.\""
		case 1:
			br = "\"You really know how to ask a question that stands out from the crowd, huh,\" \033[1;92m" + b.Name + "\033[0m says. \"I look forward to getting to know you better.\""
		}
	} else {
		switch rn {
		case 0:
			br = "\"Woah, I've never thought about it like that before,\" \033[1;92m" + b.Name + "\033[0m says. He blushes and walks away, but looks back over his shoulder at you afterwards."
		case 1:
			br = "\"I totally agree. I've never met someone who thinks so much like me,\" \033[1;92m" + b.Name + "\033[0m says. He goes on to meet the other contestants, but you can tell he's still thinking about you."
		}
	}
	form = huh.NewForm(
		huh.NewGroup(
			huh.NewNote().
				Description(br),
		),
	)

	form.Run()
}











func RunFirstImpression(state *GameState) {
	ClearScreen()
	sort.Slice(state.Contestants, func(i, j int) bool {
		a := state.Contestants[i]
		b := state.Contestants[j]
		return state.Relationship[a.Name] > state.Relationship[b.Name]
	})

	form := huh.NewForm(
    huh.NewGroup(
        huh.NewNote().
            Title("0. First Impressions").
            Description("As the Bachlor \033[1;92m" + state.Bachelor.Name + "\033[0m leaves for the day, the contestants assemble at Boston City Hall to reflect on their first interactions and begin to plan their strategies. A large screen displays the order of events for the season:\n\t1. Cape Cod\n\t2. New England Aquarium\n\t3. The Berkshires\n\t4. Martha's Vineyard\n\nAfter a few minutes, however, the screen updates to show something different..."),
    ),
	)

	form.Run()



	var rankings string
	var playerPosition int
	var pos string
	for i, c := range state.Contestants {
		pos = strconv.Itoa(i+1)
		if i == 0 {
			rankings += "🌹 "
		}
		if c.IsPlayer {
			rankings += pos + ". \033[1;36m" + c.Name + "\033[0m \n"
			playerPosition = i
			} else {
			rankings += pos + ". \033[95m" + c.Name + "\033[0m \n"
		}
	}
	var response string
	if playerPosition < 10 {
		response = "You're already in the Top 10, \033[1;36m" + state.PlayerCharacter.Name + "\033[0m, and that's before he's really even got to know your incredible personality! You've got a great chance at this." 
	} else if playerPosition < 20 {
		response = "Maybe you didn't stand out as much as you'd hoped, but at least you're not in the Bottom 5. You'll have a few chances to shine at Cape Cod."
	} else {
		response = "Uh oh, \033[1;36m" + state.PlayerCharacter.Name + "\033[0m, you're already in the Bottom 5. You'll have to work some miracles at Cape Cod to have a chance of staying on the show."
	}
	form = huh.NewForm(
    huh.NewGroup(
        huh.NewNote().
            Title("0. First Impressions").
            Description("LEADERBOARD:\n" + rankings + "\n" + response + "\n\nRegardless, you head to bed for the night and prepare for the big day tomorrow."),
    ),
	)

	form.Run()
}

// 25 - Cape Cod
func RunSession1(state *GameState) {
	ClearScreen()
	var opt string
	form := huh.NewForm(
    huh.NewGroup(
        huh.NewNote().
            Title("1. Cape Cod").
            Description("Ah, to be back at the Cape. The May weather is perfect, just cool enough that you can enjoy the heat of the sun on your skin, and the ocean is a sparkling crystalline blue. As you arrive, you see \033[1;92m" + state.Bachelor.Name + "\033[0m waiting on the upper balcony of the beautiful beachside home where you all will be spending the day. As the contestants get settled for the day, everyone separates to participate in different activities."),

				huh.NewSelect[string]().
						Title("What will you spend the day doing?").
						Options(
							huh.NewOption("Hike in the hills nearby", "hike"),
							huh.NewOption("Play beach volleyball with the other contestants", "volleyball"),
							huh.NewOption("Relax on the beach", "relax"),
							).
						Value(&opt),
    ),
	)

	form.Run()
	ghike, gvolley, grelax := AssignToGroups(state)
	if opt == "hike" {
		rn := rand.Intn(20) + state.PlayerCharacter.Strength > 10
		if rn {
			form = huh.NewForm(
				huh.NewGroup(
						huh.NewNote().
								Title("1. Cape Cod").
								Description("On your hike, something."),
		
				),
			)
		
			form.Run()
		}
	} else if opt == "volleyball" {

	} else if opt == "relax" {

	}

	RunElimination(state, 10, 15, "1. First Rose Ceremony")
}






// 15 - Aqaurium
func RunSession2(state *GameState) {
	ClearScreen()

	RunElimination(state, 7, 8, "2. Second Rose Ceremony")
}






// 8 - Berkshires
func RunSession3(state *GameState) {
	ClearScreen()
	RunElimination(state, 5, 3, "3. Third Rose Ceremony")
}






// 3 - Martha's Vineyard
func RunFantasySuites(state *GameState) {
	ClearScreen()
	RunElimination(state, 2, 1, "4. Final Rose Ceremony")
}







// 1
func RunProposal(state *GameState) {

}




func AssignToGroups(state *GameState) (group1, group2, group3 []Character) {
	// Seed the random generator (once in your program)
	rand.Seed(time.Now().UnixNano())

	// Make a copy so you don't shuffle the original order
	contestants := append([]Character(nil), state.Contestants...)
	rand.Shuffle(len(contestants), func(i, j int) {
		contestants[i], contestants[j] = contestants[j], contestants[i]
	})

	// Split into 3 roughly equal groups
	for i, c := range contestants {
		if c.IsPlayer {
			continue
		}
		switch i % 3 {
		case 0:
			group1 = append(group1, c)
		case 1:
			group2 = append(group2, c)
		case 2:
			group3 = append(group3, c)
		}
	}

	return
}





func RunElimination(state *GameState, num int, numIn int, title string) {
	sort.Slice(state.Contestants, func(i, j int) bool {
		a := state.Contestants[i]
		b := state.Contestants[j]
		return state.Relationship[a.Name] > state.Relationship[b.Name]
	})
	top := state.Contestants[:len(state.Contestants)-num]
	bottom := state.Contestants[len(state.Contestants)-num:]
	for _, c := range bottom {
		state.Eliminated = append(state.Eliminated, c.Name)
	}

	state.Contestants = top
	var rankings string
	var pos string
	eliminated := false
	for i, c := range top {
		pos = strconv.Itoa(i+1)
		if c.IsPlayer {
			rankings += "🌹 " + pos + ". \033[1;36m" + c.Name + "\033[0m \n"
			} else {
			rankings += "🌹 " + pos + ". \033[95m" + c.Name + "\033[0m \n"
		}
	}
	for i, c := range bottom {
		pos = strconv.Itoa(i+len(top)+1)
		if c.IsPlayer {
			rankings += "❌ " + pos + ". \033[1;36m" + c.Name + "\033[0m \n"
			eliminated = true
			} else {
			rankings += "❌ " + pos + ". \033[91m" + c.Name + "\033[0m \n"
		}
	}

	form := huh.NewForm(
    huh.NewGroup(
        huh.NewNote().
            Title(title).
            Description("LEADERBOARD:\n" + rankings + "\n"),
    ),
	)

	form.Run()

	if eliminated {
		form := huh.NewForm(
			huh.NewGroup(
					huh.NewNote().
							Title("The End").
							Description("Unfortunately, you've been eliminated. Maybe it was the outfit you wore, the comment you made, or just something fundamental about you as a person. Better luck next time!"),
			),
		)
	
		form.Run()
		os.Exit(0)
	}
}

// func RunEpisode1(state *GameState) {
// 	fmt.Println("🌴 Episode 1: The Beach Party 🏖️")
// 	fmt.Println("The contestants are all gathered by the beach, ready for a group date! The sun is shining, and the vibes are immaculate.")
	
// 	// Dynamically generate the introduction based on relationship scores
// 	for _, contestant := range state.Contestants {
// 			relationshipScore := state.Relationship[contestant.Name]
// 			if contestant.IsPlayer {
// 					fmt.Println("\nIt's your turn to shine! The contestants gather around you.")
// 			} else {
// 					if relationshipScore > 2 {
// 							fmt.Printf("\n%s greets you with a warm smile, their relationship score is high!\n", contestant.Name)
// 					} else {
// 							fmt.Printf("\n%s barely acknowledges your presence. Their relationship score is low.\n", contestant.Name)
// 					}
// 			}
// 	}

// 	// Decision time: Player decides whom to talk to
// 	var choice string
// 	form := huh.NewForm(
// 			huh.NewGroup(
// 					huh.NewSelect[string]().
// 							Title("Who would you like to spend time with?").
// 							Options(
// 									huh.NewOption("Ashley", "Ashley"),
// 									huh.NewOption("Jordan", "Jordan"),
// 									huh.NewOption("Kai", "Kai"),
// 							).
// 							Value(&choice),
// 			),
// 	)
// 	form.Run()

// 	// Consequences of choice: Dialogue based on character stats
// 	for _, c := range state.Contestants {
// 			if c.Name == choice {
// 					ReactToPlayerChoice(c, state)
// 			}
// 	}

// 	// End of Episode 1 Wrap-up: Update relationships
// 	fmt.Println("\n🌹 The group date ends, and everyone heads back to the mansion.")
// 	for _, c := range state.Contestants {
// 			fmt.Printf("%s's current relationship score with you is: %d\n", c.Name, state.Relationship[c.Name])
// 	}
// }

// // Add this function
// func RunRoseCeremony(state *GameState) {
//     fmt.Println("\n🌹 The Final Rose Ceremony 🌹")
//     var chosen string

//     options := []string{}
//     for _, c := range state.Contestants {
//         if !isEliminated(state, c.Name) {
//             options = append(options, c.Name)
//         }
//     }

//     if len(options) == 0 {
//         fmt.Println("Everyone has been eliminated. You are alone now. 😢")
//         return
//     }

//     form := huh.NewForm(
//         huh.NewGroup(
//             huh.NewSelect[string]().
//                 Title("Who do you give your final rose to?").
//                 Options(buildOptions(options)...).
//                 Value(&chosen),
//         ),
//     )

//     if err := form.Run(); err != nil {
//         fmt.Println("Ceremony interrupted.")
//         return
//     }

//     fmt.Printf("\nYou chose %s. They smile and say, \"I've been waiting for this moment 💕.\"\n", chosen)
// }

// Helper to check if someone was eliminated
func isEliminated(state *GameState, name string) bool {
    for _, e := range state.Eliminated {
        if e == name {
            return true
        }
    }
    return false
}

func EliminateWeighted(state *GameState, count int) {
	// Build a list of weighted entries
	type weightedEntry struct {
		Index int
		Weight int
	}

	var entries []weightedEntry
	for i, c := range state.Contestants {
		score := state.Relationship[c.Name] + 1
		entries = append(entries, weightedEntry{i, score})
	}

	// Normalize weights
	totalWeight := 0
	for _, e := range entries {
		totalWeight += e.Weight
	}

	// Sample without replacement
	eliminatedIndices := map[int]bool{}
	for len(eliminatedIndices) < count && len(eliminatedIndices) < len(entries) {
		r := rand.Intn(totalWeight)
		acc := 0
		for _, e := range entries {
			if eliminatedIndices[e.Index] {
				continue
			}
			acc += e.Weight
			if r < acc {
				eliminatedIndices[e.Index] = true
				totalWeight -= e.Weight
				break
			}
		}
	}

	// Build the new list of contestants
	var remaining []Character
	for i, c := range state.Contestants {
		if !eliminatedIndices[i] {
			remaining = append(remaining, c)
		} else {
			state.Eliminated = append(state.Eliminated, c.Name)
		}
	}
	state.Contestants = remaining
}

// Helper to convert names to huh.Option
func buildOptions(names []string) []huh.Option[string] {
    opts := []huh.Option[string]{}
    for _, name := range names {
        opts = append(opts, huh.NewOption(name, name))
    }
    return opts
}

