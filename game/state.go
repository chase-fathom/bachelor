package game

type GameState struct {
    PlayerCharacter Character
    Bachelor        Character
    Contestants     []Character
    Episode         int
    Relationship    map[string]int
    Eliminated      []string
}

func NewGameState() GameState {
    return GameState{
        Episode:      1,
        Relationship: make(map[string]int),
    }
}
