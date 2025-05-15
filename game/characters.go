package game

import (
	"github.com/charmbracelet/huh"
	"fmt"
	"math/rand"
	"time"
	"strings"
)

type Character struct {
    Name          string
    Charisma      int // 1–5
    Attractiveness int // 1–5
    Intelligence  int // 1–5
    EyeColor      string
    HairColor     string
    Height        string
		Personality		string
    IsPlayer      bool
}

var eyeColors = []string{"Blue", "Green", "Brown", "Hazel"}
var hairColors = []string{"Blonde", "Brunette", "Black", "Red"}
var heights = []string{"5'2\"", "5'5\"", "5'7\"", "5'10\"", "6'0\"", "6'2\""}
var names = []string{"Ashley", "Jordan", "Kai", "Riley", "Skylar", "Taylor", "Morgan", "Cameron", "Jesse", "Avery"}
var personalities = []string{"witty", "shy", "outgoing", "competitive", "thoughtful", "adventurous"}

func GenerateRandomContestant(name string) Character {
	return Character{
			Name:          name,
			Charisma:      rand.Intn(5) + 1,
			Attractiveness: rand.Intn(5) + 1,
			Intelligence:  rand.Intn(5) + 1,
			EyeColor:      eyeColors[rand.Intn(len(eyeColors))],
			HairColor:     hairColors[rand.Intn(len(hairColors))],
			Height:        heights[rand.Intn(len(heights))],
			Personality:		personalities[rand.Intn(len(personalities))],
			IsPlayer:      false,
	}
}

func GenerateRandomContestants(n int) []Character {
	rand.Seed(time.Now().UnixNano()) // Ensure different results each run
	usedNames := map[string]bool{}
	var contestants []Character

	for len(contestants) < n {
			name := names[rand.Intn(len(names))]
			if usedNames[name] {
					continue // avoid duplicates
			}
			usedNames[name] = true
			contestants = append(contestants, GenerateRandomContestant(name))
	}

	return contestants
}

func ReactToPlayerChoice(c Character, state *GameState) {
	charisma := c.Charisma
	intelligence := c.Intelligence
	attractiveness := c.Attractiveness
	relationshipScore := state.Relationship[c.Name]

	fmt.Printf("\nYou approach %s...\n", c.Name)

	if relationshipScore > 3 {
			fmt.Printf("%s smiles brightly at you. They seem excited to chat.\n", c.Name)
	} else {
			fmt.Printf("%s seems a little distant. Their eyes wander as you talk.\n", c.Name)
	}

	// Personality-driven dialogue based on stats
	if charisma >= 4 {
			fmt.Println("You strike up a charming conversation, and they laugh at your jokes.")
	} else if attractiveness >= 4 {
			fmt.Println("You both share some compliments. It's a light and flirty conversation.")
	} else if intelligence >= 4 {
			fmt.Println("You dive into a deep, intellectual conversation. They seem impressed with your knowledge.")
	}

	// Adding some tension: If their relationship with you is low
	if relationshipScore <= 2 {
			fmt.Println("Despite the conversation, they seem uninterested... Maybe something's off.")
	}

	// Adjust relationship based on player stats and interaction outcome
	if charisma > 3 {
			state.Relationship[c.Name] += 1
	} else if attractiveness > 3 {
			state.Relationship[c.Name] += 1
	} else {
			state.Relationship[c.Name] -= 1
	}

	// Special outcome for high charisma (maybe some *romantic* bonus points)
	if charisma > 4 && relationshipScore >= 3 {
			fmt.Println("Your charisma sparks some chemistry. They're blushing a little... something is growing here.")
			state.Relationship[c.Name] += 2
	}
}



func CreatePlayerCharacter(state *GameState) {
	var c Character
	c.IsPlayer = true
	fmt.Print("\033[H\033[2J")

	err := huh.NewForm(
			huh.NewGroup(
					huh.NewNote().
            	Title("🌹 The Bachelor Simulator 🌹").
            	Description("Enter as a contestant in the bachelor."),

					huh.NewInput().
							Title("What's your name?").
							Placeholder("e.g. Ellory").
							Value(&c.Name),
			),
	).Run()
	if err != nil {
		fmt.Println("Cancelled.")
		return
	}

	for {
		err = huh.NewForm(
			huh.NewGroup(
				huh.NewNote().
            	Title("Stats").
            	Description("Assign a total of 9 stat points to different attributes."),

				huh.NewSelect[int]().
						Title("Charisma (1 = very low, 5 = very high)").
						Options(intOptions(1, 5)...).
						Value(&c.Charisma),
	
				huh.NewSelect[int]().
						Title("Attractiveness").
						Options(intOptions(1, 5)...).
						Value(&c.Attractiveness),
	
				huh.NewSelect[int]().
						Title("Intelligence").
						Options(intOptions(1, 5)...).
						Value(&c.Intelligence),
			),
		).Run()
		if err != nil {
			fmt.Println("Cancelled.")
			return
		}
		total := c.Charisma + c.Intelligence + c.Attractiveness
		if total <= 9 {
				fmt.Print("\033[H\033[2J")
				break // we're good!
		}
		fmt.Printf("\n❗ Your total was %d. Please distribute at most 9 points.\n\n", total)
	}
	err = huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
						Title("Personality").
						Options(func() []huh.Option[string] {
							opts := make([]huh.Option[string], 0, len(personalities))
							for _, p := range personalities {
									opts = append(opts, huh.NewOption(capitalize(p), p))
							}
							return opts
						}()...).
						Value(&c.Personality),

					huh.NewSelect[string]().
							Title("Eye Color").
							Options(
									huh.NewOption("Blue", "Blue"),
									huh.NewOption("Green", "Green"),
									huh.NewOption("Brown", "Brown"),
									huh.NewOption("Hazel", "Hazel"),
							).
							Value(&c.EyeColor),

					huh.NewSelect[string]().
							Title("Hair Color").
							Options(
									huh.NewOption("Blonde", "Blonde"),
									huh.NewOption("Brunette", "Brunette"),
									huh.NewOption("Black", "Black"),
									huh.NewOption("Red", "Red"),
							).
							Value(&c.HairColor),

					huh.NewInput().
							Title("Height").
							Placeholder("e.g. 5'11\"").
							Value(&c.Height),
		),
	).Run()
	if err != nil {
		fmt.Println("Cancelled.")
		return
	}

	state.PlayerCharacter = c
}

func intOptions(min, max int) []huh.Option[int] {
	var opts []huh.Option[int]
	for i := min; i <= max; i++ {
			opts = append(opts, huh.NewOption(fmt.Sprintf("%d", i), i))
	}
	return opts
}

func IntroduceContestants(state *GameState) {
	state.Contestants = GenerateRandomContestants(3)

		for _, c := range state.Contestants {
				state.Relationship[c.Name] = 0
		}
}

func capitalize(s string) string {
	if len(s) == 0 {
			return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

