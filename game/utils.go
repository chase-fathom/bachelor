package game

import (
    "fmt"
    "github.com/charmbracelet/lipgloss"
)

var titleStyle = lipgloss.NewStyle().
    Bold(true).
    Foreground(lipgloss.Color("#FF5FA2")).
    Padding(1, 2)

func PrintTitle(title string) {
    fmt.Println(titleStyle.Render(title))
}
