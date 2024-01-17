package helper

import (
	"bufio"
	"strings"
)

// ContainsPlayer checks if a player is in a list of players
func ContainsPlayer(players []string, name string) bool {
	for _, player := range players {
		if player == name {
			return true
		}
	}
	return false
}

// HasMatchEnded checks if a match has ended
func HasMatchEnded(scanner *bufio.Scanner) bool {
	return strings.Contains(scanner.Text(), END_MATCH_LINE_1) ||
		strings.Contains(scanner.Text(), END_MATCH_LINE_2)
}
