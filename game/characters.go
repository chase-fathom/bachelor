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
    Strength  		int // 1–5
    EyeColor      string
    HairColor     string
    Height        string
		Personality		string
		Noun					string
    IsPlayer      bool
		IsBachelor		bool
}

// var eyeColors = []string{"Blue", "Green", "Brown", "Hazel"}
// var hairColors = []string{"Blonde", "Brown", "Black", "Red"}
// var heights = []string{"4'9\"", "4'10\"", "4'11\"", "5'0\"", "5'1\"", "5'2\"", "5'3\"", "5'4\"", "5'5\"", "5'6\"", "5'7\"", "5'8\"", "5'9\"", "5'10\"", "5'11\"", "6'0\"", "6'1\"", "6'2\"", "6'3\""}
// var personalities = []string{"witty", "shy", "outgoing", "competitive", "thoughtful", "adventurous"}

func GenerateContestants(state *GameState) {
	state.Contestants = GenerateRandomContestants(24, state)
	state.Contestants = append(state.Contestants, state.PlayerCharacter)
	ShuffleCharacters(state.Contestants)

	for _, c := range state.Contestants {
			state.Relationship[c.Name] = 0
	}

	state.Bachelor = GenerateBachelor()

}

func GenerateRandomContestant(name string) Character {
	return Character{
			Name:          name,
			Charisma:      rand.Intn(4) + 1,
			Attractiveness: rand.Intn(4) + 1,
			Strength:  rand.Intn(4) + 1,
			EyeColor:      eyeColors[rand.Intn(len(eyeColors))],
			HairColor:     hairColors[rand.Intn(len(hairColors))],
			Height:        heights[rand.Intn(len(heights))],
			Personality:		personalities[rand.Intn(len(personalities))],
			Noun:					beautyTerms[rand.Intn(len(beautyTerms))],
			IsPlayer:      false,
			IsBachelor:			false,
	}
}

func GenerateRandomContestants(n int, state *GameState) []Character {
	rand.Seed(time.Now().UnixNano()) // Ensure different results each run
	usedNames := map[string]bool{}
	usedNames[state.PlayerCharacter.Name] = true
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

func GenerateBachelor() Character {
	return Character{
		Name:          bachelorNames[rand.Intn(len(bachelorNames))],
		Charisma:      rand.Intn(5) + 1,
		Attractiveness: rand.Intn(5) + 1,
		Strength:  rand.Intn(5) + 1,
		EyeColor:      eyeColors[rand.Intn(len(eyeColors))],
		HairColor:     hairColors[rand.Intn(len(hairColors))],
		Height:        heights[rand.Intn(len(heights))],
		Personality:		bachelorPersonalities[rand.Intn(len(bachelorPersonalities))],
		Noun:					"bachelor",
		IsPlayer:      false,
		IsBachelor:			true,
}
}

func ShuffleCharacters(chars []Character) {
	rand.Shuffle(len(chars), func(i, j int) {
		chars[i], chars[j] = chars[j], chars[i]
	})
}

func ReactToPlayerChoice(c Character, state *GameState) {
	charisma := c.Charisma
	strength := c.Strength
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
	} else if strength >= 4 {
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
	c.Noun = "player"
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
						Title("Strength").
						Options(intOptions(1, 5)...).
						Value(&c.Strength),
			),
		).Run()
		if err != nil {
			fmt.Println("Cancelled.")
			return
		}
		total := c.Charisma + c.Strength + c.Attractiveness
		if total <= 9 {
				fmt.Print("\033[H\033[2J")
				break
		}
		fmt.Printf("\n❗ Your total was %d. Please distribute at most 9 points.\n\n", total)
	}
	err = huh.NewForm(
		huh.NewGroup(
			huh.NewNote().
            	Title("Attributes"),

			// huh.NewSelect[string]().
			// 			Title("Personality").
			// 			Options(func() []huh.Option[string] {
			// 				opts := make([]huh.Option[string], 0, len(personalities))
			// 				for _, p := range personalities {
			// 						opts = append(opts, huh.NewOption(capitalize(p), p))
			// 				}
			// 				return opts
			// 			}()...).
			// 			Value(&c.Personality),

				huh.NewInput().
							Title("Personality").
							Placeholder("e.g. contemplative").
							Value(&c.Personality),

					huh.NewInput().
						Title("Eye Color").
						Placeholder("e.g. brown").
						Value(&c.EyeColor),

					huh.NewInput().
							Title("Hair Color").
							Placeholder("e.g. black").
							Value(&c.HairColor),

					huh.NewInput().
							Title("Height").
							Placeholder("e.g. 5'11").
							Value(&c.Height),
		),
	).Run()
	if err != nil {
		fmt.Println("Cancelled.")
		return
	}
	state.Relationship[c.Name] = 0

	state.PlayerCharacter = c
}

func intOptions(min, max int) []huh.Option[int] {
	var opts []huh.Option[int]
	for i := min; i <= max; i++ {
			opts = append(opts, huh.NewOption(fmt.Sprintf("%d", i), i))
	}
	return opts
}

func capitalize(s string) string {
	if len(s) == 0 {
			return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

var heights = []string{"4'9\"", "4'10\"", "4'11\"", "5'0\"", "5'1\"", "5'2\"", "5'3\"", "5'4\"", "5'5\"", "5'6\"", "5'7\"", "5'8\"", "5'9\"", "5'10\"", "5'11\"", "6'0\"", "6'1\"", "6'2\"", "6'3\""}

var names = []string{
	"Aaliyah",
	"Adelina",
	"Alexis",
	"Alice",
	"Alyssa",
	"Amara",
	"Angelina",
	"Anastasia",
	"Ariana",
	"Autumn",
	"Bailey",
	"Barbara",
	"Beatrice",
	"Bianca",
	"Blake",
	"Brenda",
	"Brooklyn",
	"Camila",
	"Carmen",
	"Carolina",
	"Catherine",
	"Chloe",
	"Claire",
	"Clara",
	"Danica",
	"Delilah",
	"Diana",
	"Donna",
	"Elena",
	"Elizabeth",
	"Ellie",
	"Ellery",
	"Ellory",
	"Emily",
	"Emma",
	"Erica",
	"Evelyn",
	"Fiona",
	"Gabriella",
	"Grace",
	"Hannah",
	"Harper",
	"Heather",
	"Isabella",
	"Jade",
	"Jasmine",
	"Jessica",
	"Jordan",
	"Julia",
	"Kaitlyn",
	"Kate",
	"Katherine",
	"Kendall",
	"Kimberly",
	"Kylie",
	"Lauren",
	"Leah",
	"Lexi",
	"Libby",
	"Lily",
	"Lola",
	"Madeline",
	"Maya",
	"Megan",
	"Melanie",
	"Monica",
	"Naomi",
	"Natalie",
	"Olivia",
	"Paige",
	"Penelope",
	"Rachel",
	"Reagan",
	"Riley",
	"Rose",
	"Sadie",
	"Samantha",
	"Sara",
	"Sasha",
	"Scarlett",
	"Selena",
	"Sophie",
	"Stella",
	"Summer",
	"Taylor",
	"Victoria",
	"Violet",
	"Vivian",
	"Zoe",
	"Ashley",
	"Jordan",
	"Karlie",
	"Kai",
	"Riley",
	"Skylar",
	"Taylor",
	"Morgan",
	"Cameron",
	"Jesse",
	"Avery",
	"Adriana",
	"Bailey",
	"Brooke",
	"Charlotte",
	"Caitlyn",
	"Destiny",
	"Faith",
	"Gracie",
	"Kaitlin",
	"Lindsey",
	"Madison",
	"Maria",
	"Paige",
	"Riley",
	"Skylar",
	"Taylor",
	"Valeria",
	"Viviana",
}
var personalities = []string{
	"romantic",
	"bubbly",
	"shy",
	"confident",
	"jealous",
	"loyal",
	"ambitious",
	"outgoing",
	"awkward",
	"cynical",
	"dramatic",
	"sweet",
	"quirky",
	"intense",
	"flirty",
	"bold",
	"funny",
	"serious",
	"competitive",
	"chill",
	"mysterious",
	"emotional",
	"honest",
	"playful",
	"reserved",
	"thoughtful",
	"adventurous",
	"stylish",
	"sensitive",
	"charismatic",
}
var hairColors = []string{
	"blonde",
	"brunette",
	"black",
	"red",
	"auburn",
	"platinum",
	"chestnut",
	"dirty blonde",
	"silver",
	"pink",
}
var eyeColors = []string{
	"blue",
	"green",
	"brown",
	"hazel",
	"gray",
	"amber",
	"dark brown",
	"black",
}
var beautyTerms = []string{
	"stunner",
	"angel",
	"babe",
	"bombshell",
	"fox",
	"knockout",
	"vision",
	"goddess",
	"doll",
	"queen",
	"lady",
	"muse",
	"gem",
	"enchantress",
	"charm",
	"diva",
	"starlet",
	"cutie",
	"heartbreaker",
	"showstopper",
	"darling",
	"sweetheart",
	"vixen",
	"icon",
	"flame",
	"flower",
	"princess",
	"pearl",
	"treasure",
	"dream",
}

var bachelorNames = []string{
	"Aaron",
	"Alex",
	"Andrew",
	"Austin",
	"Blake",
	"Brady",
	"Brandon",
	"Caleb",
	"Cameron",
	"Carter",
	"Chase",
	"Chris",
	"Clayton",
	"Cole",
	"Connor",
	"Cory",
	"Daniel",
	"David",
	"Devin",
	"Drew",
	"Dylan",
	"Eli",
	"Eric",
	"Ethan",
	"Evan",
	"Gabe",
	"Grayson",
	"Greg",
	"Ian",
	"Jack",
	"Jacob",
	"Jake",
	"James",
	"Jason",
	"Jayden",
	"Jesse",
	"John",
	"Jordan",
	"Josh",
	"Julian",
	"Kyle",
	"Leo",
	"Liam",
	"Logan",
	"Lucas",
	"Mark",
	"Mason",
	"Matt",
	"Nate",
	"Paul",
	"Ryan",
	"Zach",
}

var bachelorPersonalities = []string{
	"the Gym Bro",
	"the Sensitive Cowboy",
	"the Spreadsheet Guy",
	"the Golden Retriever Man",
	"the Adrenaline Junkie",
	"the Overly Chill Surfer",
	"the Finance Bro",
	"the Guy Who Brings A Guitar",
	"the Emotional Himbo",
	"the One-Upper",
	"the Jealous Protector",
	"the Accidental Poet",
	"the Quiet Heartthrob",
	"the Drama Magnet",
	"the Guy With A Podcast",
	"the Wholesome Flirt",
	"the Conspiracy Theorist",
	"the Crypto Enthusiast",
	"the Wannabe Chef",
	"the Outdoorsy Loner",
	"the Motivational Speaker",
	"the Guy Who Peaked In High School",
	"the Armchair Philosopher",
	"the Loud Hugger",
	"the Deep-Feelings Dude",
	"the Guy Who Says 'No Worries' Too Much",
	"the Competitive Cuddler",
	"the Guy Who Makes Everything A Metaphor",
	"the Guy With Too Many Rings",
	"the One Who Can't Stop Talking About His Mom",
}
