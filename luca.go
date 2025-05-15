package main

import (
    "fmt"
    "math/rand"
    "os"
    "sort"
    "strings"
    "time"

    "github.com/charmbracelet/bubbles/textinput"
    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/lipgloss"
)

// ---------- Data Models ---------

type GameState int

const (
    StateIntro GameState = iota
    StateCustomize
    StateOneOnOne
    StateGroupDate
    StateDrama
    StateCeremony
    StateGameOver
    StateWin
)

// Contestant represents a contestant in the show (player or AI).

type Contestant struct {
    Name          string
    Charisma      int
    Attractiveness int
    Intelligence  int
    EyeColor      string
    HairColor     string
    Height        string
    Personality   string
    Score         float64
    IsPlayer      bool
}

type Bachelor struct {
    Name  string
    PrefC float64
    PrefA float64
    PrefI float64
    // simple preference: dislikes drama if true
    DislikesDrama bool
}

// Choice for one‑on‑one scenario

type Choice struct {
    Text string
    Stat string // "charisma", "attractiveness", "intelligence"
}

type Scenario struct {
    Description string
    Choices     []Choice
}

// Model for Bubble Tea

type Model struct {
    state GameState

    week int

    // text input for customization
    ti           textinput.Model
    questionStep int

    player        Contestant
    bachelor      Bachelor
    contestants   []Contestant // includes player & AI (current).

    // one‑on‑one UI
    currentScenario Scenario
    cursor          int // menu selection index
    outcomeText     string

    // drama flag
    dramaOccurred bool

    // group date info
    groupEventStat string
    groupWinnerIdx int
}

// ---------- Utility Functions ---------

func randName() string {
    names := []string{"Alex", "Jordan", "Taylor", "Sam", "Casey", "Jamie", "Morgan", "Riley", "Dana", "Skyler", "Cameron", "Quinn", "Sydney", "Avery", "Peyton", "Harper"}
    return names[rand.Intn(len(names))]
}

func randEye() string {
    eyes := []string{"blue", "green", "brown", "hazel", "gray"}
    return eyes[rand.Intn(len(eyes))]
}

func randHair() string {
    hair := []string{"blonde", "brunette", "black", "red", "auburn"}
    return hair[rand.Intn(len(hair))]
}

func randHeight() string {
    heights := []string{"5'2\"", "5'4\"", "5'6\"", "5'8\"", "5'10\"", "6'0\""}
    return heights[rand.Intn(len(heights))]
}

func randPersonality() string {
    pers := []string{"witty", "shy", "outgoing", "competitive", "thoughtful", "adventurous"}
    return pers[rand.Intn(len(pers))]
}

func generateContestant(isPlayer bool) Contestant {
    c := Contestant{
        Name:           randName(),
        Charisma:       rand.Intn(6) + 5,        // 5‑10
        Attractiveness: rand.Intn(6) + 5,
        Intelligence:   rand.Intn(6) + 5,
        EyeColor:       randEye(),
        HairColor:      randHair(),
        Height:         randHeight(),
        Personality:    randPersonality(),
        Score:          0,
        IsPlayer:       isPlayer,
    }
    return c
}

func generateBachelor() Bachelor {
    // random weights that sum to 1
    a, b, c := rand.Float64(), rand.Float64(), rand.Float64()
    sum := a + b + c
    return Bachelor{
        Name:          randName(),
        PrefC:         a / sum,
        PrefA:         b / sum,
        PrefI:         c / sum,
        DislikesDrama: rand.Intn(2) == 0,
    }
}

// base compatibility calculation
func baseScore(b Bachelor, c Contestant) float64 {
    return float64(c.Charisma)*b.PrefC + float64(c.Attractiveness)*b.PrefA + float64(c.Intelligence)*b.PrefI
}

// create initial model
func initialModel() Model {
    rand.Seed(time.Now().UnixNano())

    m := Model{
        state:        StateCustomize,
        week:         1,
        questionStep: 0,
    }

    // generate default player (will customize)
    m.player = generateContestant(true)

    // text input init
    ti := textinput.New()
    ti.Prompt      = "Name: "
    ti.Placeholder = "Enter your name"
    ti.Focus()
    ti.Width = 50  // ← make the box at least 20 characters wide
    m.ti = ti


    return m
}

// prepare game after customization complete
func (m *Model) startGame() {
    // generate bachelor and AI contestants
    m.bachelor = generateBachelor()

    totalContestants := 5 // player + 4 AI
    m.contestants = []Contestant{m.player}
    for len(m.contestants) < totalContestants {
        c := generateContestant(false)
        // ensure unique names to avoid confusion
        duplicate := false
        for _, exist := range m.contestants {
            if exist.Name == c.Name {
                duplicate = true
                break
            }
        }
        if !duplicate {
            c.Score = baseScore(m.bachelor, c)
            m.contestants = append(m.contestants, c)
        }
    }

    // compute player's initial score
    for i := range m.contestants {
        if m.contestants[i].IsPlayer {
            m.contestants[i].Score = baseScore(m.bachelor, m.contestants[i])
            m.player = m.contestants[i]
            break
        }
    }

    // pick first scenario
    m.currentScenario = randomScenario()
    m.cursor = 0
    m.state = StateOneOnOne
}

// random one‑on‑one scenario selection
func randomScenario() Scenario {
    scenarios := []Scenario{
        {
            Description: "You and the Bachelor enjoy a cozy dinner. He asks about your passions.",
            Choices: []Choice{
                {"Open up sincerely about your career ambitions", "intelligence"},
                {"Flirt playfully about your future together", "charisma"},
                {"Ask him about his own goals instead", "charisma"},
            },
        },
        {
            Description: "You go hiking with the Bachelor and stop at a scenic viewpoint.",
            Choices: []Choice{
                {"Share an adventurous travel story", "charisma"},
                {"Compliment the view and his company", "attractiveness"},
                {"Discuss environmental conservation efforts", "intelligence"},
            },
        },
        {
            Description: "You both attend a private cooking class together.",
            Choices: []Choice{
                {"Take charge and show your cooking skills", "intelligence"},
                {"Joke around and taste‑test ingredients playfully", "charisma"},
                {"Present the Bachelor with a beautifully plated dish", "attractiveness"},
            },
        },
    }
    return scenarios[rand.Intn(len(scenarios))]
}

// ---------- Bubble Tea Update & View ---------

func (m Model) Init() tea.Cmd {
  return textinput.Blink
}


func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch m.state {
    case StateCustomize:
        return updateCustomize(m, msg)
    case StateOneOnOne:
        return updateOneOnOne(m, msg)
    case StateGroupDate:
        return updateGroupDate(m, msg)
    case StateDrama:
        return updateDrama(m, msg)
    case StateCeremony:
        return updateCeremony(m, msg)
    case StateGameOver, StateWin:
        if key, ok := msg.(tea.KeyMsg); ok {
            if key.String() == "q" || key.Type == tea.KeyCtrlC || key.Type == tea.KeyEnter {
                return m, tea.Quit
            }
        }
    }
    return m, nil
}

// --- Update helpers ---

func updateCustomize(m Model, msg tea.Msg) (tea.Model, tea.Cmd) {
    var cmd tea.Cmd

    switch msg := msg.(type) {
    case tea.KeyMsg:
        if msg.Type == tea.KeyCtrlC {
            return m, tea.Quit
        }

        if msg.Type == tea.KeyEnter {
            input := strings.TrimSpace(m.ti.Value())
            if input == "" {
                // keep existing default
                input = m.ti.Placeholder // placeholder is field name; we won't use
            }
            // store input based on question step
            switch m.questionStep {
            case 0:
                if input != "" {
                    m.player.Name = input
                }
                m.questionStep++
                m.ti.SetValue("")
                m.ti.Placeholder = "Eye color (press Enter to keep random)"
            case 1:
                if input != "" {
                    m.player.EyeColor = input
                }
                m.questionStep++
                m.ti.SetValue("")
                m.ti.Placeholder = "Hair color (press Enter to keep random)"
            case 2:
                if input != "" {
                    m.player.HairColor = input
                }
                m.questionStep++
                m.ti.SetValue("")
                m.ti.Placeholder = "Height (press Enter to keep random)"
            case 3:
                if input != "" {
                    m.player.Height = input
                }
                // customization done
                m.startGame()
            }
        }
    }

    m.ti, cmd = m.ti.Update(msg)
    return m, cmd
}

func updateOneOnOne(m Model, msg tea.Msg) (tea.Model, tea.Cmd) {
    if key, ok := msg.(tea.KeyMsg); ok {
        switch key.String() {
        case "ctrl+c", "q":
            return m, tea.Quit
        case "up", "k":
            if m.cursor > 0 {
                m.cursor--
            }
        case "down", "j":
            if m.cursor < len(m.currentScenario.Choices)-1 {
                m.cursor++
            }
        case "enter":
            // resolve choice
            choice := m.currentScenario.Choices[m.cursor]
            delta := resolveChoice(&m, choice)
            m.outcomeText = fmt.Sprintf("You chose: %s ( %+0.1f points )", choice.Text, delta)
            m.cursor = 0
            // proceed to group date after showing outcome (require another Enter)
            m.state = StateGroupDate
        }
    }
    return m, nil
}

// compute delta and update player score
func resolveChoice(m *Model, ch Choice) float64 {
    statVal := 0
    switch ch.Stat {
    case "charisma":
        statVal = m.player.Charisma
    case "attractiveness":
        statVal = m.player.Attractiveness
    case "intelligence":
        statVal = m.player.Intelligence
    }
    var weight float64
    switch ch.Stat {
    case "charisma":
        weight = m.bachelor.PrefC
    case "attractiveness":
        weight = m.bachelor.PrefA
    case "intelligence":
        weight = m.bachelor.PrefI
    }
    // simple: delta = weight * statVal / 2 + random noise [‑1,1]
    delta := weight*float64(statVal)/2 + (rand.Float64()*2 - 1)
    // update player's score in contestants slice
    for i := range m.contestants {
        if m.contestants[i].IsPlayer {
            m.contestants[i].Score += delta
            m.player = m.contestants[i]
            break
        }
    }
    return delta
}

func updateGroupDate(m Model, msg tea.Msg) (tea.Model, tea.Cmd) {
    // On first entry to this state, simulate group date once; afterwards wait for Enter to continue
    if m.groupEventStat == "" {
        // pick random event stat
        stats := []string{"charisma", "attractiveness", "intelligence"}
        m.groupEventStat = stats[rand.Intn(len(stats))]
        simulateGroupDate(&m)
        return m, nil
    }

    if key, ok := msg.(tea.KeyMsg); ok {
        if key.String() == "enter" {
            m.state = StateDrama
            m.groupEventStat = ""
            m.dramaOccurred = false
        } else if key.String() == "q" || key.Type == tea.KeyCtrlC {
            return m, tea.Quit
        }
    }
    return m, nil
}

func simulateGroupDate(m *Model) {
    // compute performance = relevant stat + random(0‑3)
    bestIdx := 0
    bestPerf := -1
    for i := range m.contestants {
        var stat int
        switch m.groupEventStat {
        case "charisma":
            stat = m.contestants[i].Charisma
        case "attractiveness":
            stat = m.contestants[i].Attractiveness
        case "intelligence":
            stat = m.contestants[i].Intelligence
        }
        perf := stat + rand.Intn(4)
        if perf > bestPerf {
            bestPerf = perf
            bestIdx = i
        }
    }
    m.groupWinnerIdx = bestIdx

    // apply points: winner +6, others +2, last +0
    for i := range m.contestants {
        gain := 2.0
        if i == bestIdx {
            gain = 6.0
        }
        m.contestants[i].Score += gain * rand.Float64() // small randomization
        if m.contestants[i].IsPlayer {
            m.player = m.contestants[i]
        }
    }
    m.outcomeText = fmt.Sprintf("Group date focused on %s. %s impressed the Bachelor!", m.groupEventStat, m.contestants[bestIdx].Name)
}

func updateDrama(m Model, msg tea.Msg) (tea.Model, tea.Cmd) {
    // on entry, maybe create drama once
    if !m.dramaOccurred {
        handleDrama(&m)
        return m, nil
    }

    if key, ok := msg.(tea.KeyMsg); ok {
        if key.String() == "enter" {
            m.state = StateCeremony
        } else if key.String() == "q" || key.Type == tea.KeyCtrlC {
            return m, tea.Quit
        }
    }
    return m, nil
}

func handleDrama(m *Model) {
    m.dramaOccurred = true
    // 50% chance event involves player
    playerInvolved := rand.Intn(2) == 0
    if playerInvolved {
        // simple: player argued; if bachelor dislikes drama -> penalty else minor gain
        if m.bachelor.DislikesDrama {
            delta := -3.0
            for i := range m.contestants {
                if m.contestants[i].IsPlayer {
                    m.contestants[i].Score += delta
                    m.player = m.contestants[i]
                }
            }
            m.outcomeText = "You got into a brief argument during the cocktail party. The Bachelor dislikes drama and seems disappointed."
        } else {
            delta := 2.0
            for i := range m.contestants {
                if m.contestants[i].IsPlayer {
                    m.contestants[i].Score += delta
                    m.player = m.contestants[i]
                }
            }
            m.outcomeText = "Your fiery spirit caught the Bachelor’s eye! He seems intrigued by the drama."
        }
    } else {
        // two AI contestants drama; they each lose small
        idx1 := rand.Intn(len(m.contestants))
        idx2 := rand.Intn(len(m.contestants))
        for idx2 == idx1 || m.contestants[idx2].IsPlayer {
            idx2 = rand.Intn(len(m.contestants))
        }
        m.contestants[idx1].Score -= 2
        m.contestants[idx2].Score -= 2
        m.outcomeText = fmt.Sprintf("%s and %s had a heated argument, turning the Bachelor off.", m.contestants[idx1].Name, m.contestants[idx2].Name)
    }
}

func updateCeremony(m Model, msg tea.Msg) (tea.Model, tea.Cmd) {
    if key, ok := msg.(tea.KeyMsg); ok {
        if key.String() == "enter" {
            // eliminate lowest and proceed
            sort.Slice(m.contestants, func(i, j int) bool {
                return m.contestants[i].Score > m.contestants[j].Score
            })
            eliminated := m.contestants[len(m.contestants)-1]
            if eliminated.IsPlayer {
                m.state = StateGameOver
            } else if len(m.contestants) == 2 {
                // next elimination will decide winner, but if player top now, they win
                if m.contestants[0].IsPlayer {
                    m.state = StateWin
                } else {
                    m.state = StateGameOver
                }
            } else {
                // remove eliminated
                m.contestants = m.contestants[:len(m.contestants)-1]
                m.week++
                // prepare next week
                m.currentScenario = randomScenario()
                m.state = StateOneOnOne
            }
        } else if key.String() == "q" || key.Type == tea.KeyCtrlC {
            return m, tea.Quit
        }
    }
    return m, nil
}

// ---------- View ---------

var (
    titleStyle      = lipgloss.NewStyle().Bold(true).Underline(true)
    playerStyle     = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
    eliminatedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
    highlightStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true)
    wrapStyle       = lipgloss.NewStyle().MaxWidth(80)
)

func (m Model) View() string {
    var b strings.Builder

    switch m.state {
    case StateCustomize:
        b.WriteString(titleStyle.Render("Contestant Customization") + "\n\n")
        b.WriteString(m.ti.View() + "\n")
        b.WriteString("(Press Enter when done)\n")

    case StateOneOnOne:
        b.WriteString(titleStyle.Render(fmt.Sprintf("Week %d – One‑on‑One Date\n\n", m.week)))
        b.WriteString(wrapStyle.Render(m.currentScenario.Description) + "\n\n")
        for i, ch := range m.currentScenario.Choices {
            cursor := "  "
            text := ch.Text
            if i == m.cursor {
                cursor = "> "
                text = highlightStyle.Render(text)
            }
            b.WriteString(fmt.Sprintf("%s%s\n", cursor, text))
        }
        b.WriteString("\nUse ↑/↓ to select, Enter to confirm. q to quit.\n")

    case StateGroupDate:
        b.WriteString(titleStyle.Render(fmt.Sprintf("Week %d – Group Date\n\n", m.week)))
        b.WriteString(m.outcomeText + "\n\n")
        b.WriteString("Press Enter to continue.\n")

    case StateDrama:
        b.WriteString(titleStyle.Render(fmt.Sprintf("Week %d – Mansion Dynamics\n\n", m.week)))
        b.WriteString(m.outcomeText + "\n\n")
        b.WriteString("Press Enter to proceed to the Rose Ceremony.\n")

    case StateCeremony:
        b.WriteString(titleStyle.Render(fmt.Sprintf("Week %d – Rose Ceremony\n\n", m.week)))
        // sort a copy for display
        copyList := make([]Contestant, len(m.contestants))
        copy(copyList, m.contestants)
        sort.Slice(copyList, func(i, j int) bool { return copyList[i].Score > copyList[j].Score })
        for i, c := range copyList {
            line := fmt.Sprintf("%d. %s – %.1f", i+1, c.Name, c.Score)
            if c.IsPlayer {
                line = playerStyle.Render(line)
            }
            if i == len(copyList)-1 {
                line = eliminatedStyle.Render(line + "  (Eliminated)")
            } else {
                line += "  🌹"
            }
            b.WriteString(line + "\n")
        }
        b.WriteString("\nPress Enter to continue. q to quit.\n")

    case StateGameOver:
        b.WriteString(eliminatedStyle.Render("You have been eliminated. 😢\n"))
        b.WriteString("Thank you for playing! Press q to quit.\n")

    case StateWin:
        b.WriteString(playerStyle.Render("🎉 You received the final rose! You win! 🎉\n"))
        b.WriteString("Congratulations on finding love in the terminal. Press q to quit.\n")
    default:
        b.WriteString("Unknown state")
    }

    return b.String()
}

// ---------- main ----------

func main() {
    p := tea.NewProgram(initialModel())
    if err := p.Start(); err != nil {
        fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
        os.Exit(1)
    }
}
