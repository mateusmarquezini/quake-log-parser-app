package service

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/mateusmarquezini/quake-log-parser-app/domain"
)

type MatchService struct {
	matchID           int
	matches           []map[string]any
	response          map[string]any
	killEventsInMatch []domain.KillDetails
	match             *domain.Match
}

func NewMatchService() *MatchService {
	return &MatchService{
		matches:  []map[string]any{},
		response: make(map[string]any),
	}
}

func (s *MatchService) ProcessLine(scanner *bufio.Scanner) {
	if hasMatchStarted(scanner) {
		s.startNewMatch()
	}

	matchKillDetails := parseKillDetails(scanner)

	if matchKillDetails.Cause != "" {
		s.killEventsInMatch = append(s.killEventsInMatch, matchKillDetails)
	}

	if hasMatchEnded(scanner) {
		s.endMatch()
	}
}

func (s *MatchService) startNewMatch() {
	s.matchID++
	s.match = &domain.Match{
		MatchID:      matchIDPrefix + fmt.Sprint(s.matchID),
		Players:      []string{},
		Kills:        map[string]int{},
		KillsByMeans: map[string]int{},
	}
}

func (s *MatchService) endMatch() {
	// If there are no kill events in the match, we consider the match as invalid and decrement the matchID
	if len(s.killEventsInMatch) == 0 {
		s.matchID -= 1
	} else {
		for _, event := range s.killEventsInMatch {
			incrementKillsByCause(s.match, event.Cause)

			addPlayerIfNotExists(s.match, event.KillerName)
			addPlayerIfNotExists(s.match, event.VictimName)

			if !isTheSamePlayer(event) {
				if isAValidPlayer(event) {
					incrementKillCount(s.match, event.KillerName)
				} else {
					decrementKillCount(s.match, event.VictimName)
				}
			}
		}

		s.match.TotalKills = len(s.killEventsInMatch)
		s.response[s.match.MatchID] = s.match
		s.matches = append(s.matches, s.response)
		s.response = make(map[string]any)
		s.killEventsInMatch = nil
	}
}

func (s *MatchService) PrintMatchReport() {
	matchesJson, err := json.MarshalIndent(s.matches, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling matches to JSON: %v", err)
		return
	}

	fmt.Println(string(matchesJson))
}

func incrementKillCount(match *domain.Match, playerName string) {
	match.Kills[playerName]++
}

func decrementKillCount(match *domain.Match, playerName string) {
	if match.Kills[playerName] > 0 {
		match.Kills[playerName]--
	}
}

func incrementKillsByCause(match *domain.Match, cause string) {
	match.KillsByMeans[cause]++
}

func addPlayerIfNotExists(match *domain.Match, playerName string) {
	if playerName != worldID && !containsPlayer(match.Players, playerName) {
		match.Players = append(match.Players, playerName)
		addPlayerToKillsRanking(match, playerName)
	}
}

func isAValidPlayer(kill domain.KillDetails) bool {
	return kill.KillerName != worldID
}

func addPlayerToKillsRanking(match *domain.Match, playerName string) {
	if _, exists := match.Kills[playerName]; !exists {
		match.Kills[playerName] = 0
	}
}

func isTheSamePlayer(kill domain.KillDetails) bool {
	return kill.KillerName == kill.VictimName
}

func parseKillDetails(scanner *bufio.Scanner) domain.KillDetails {
	match := regexp.MustCompile(killActionLineRegex).FindStringSubmatch(scanner.Text())

	if len(match) > 0 {
		return domain.KillDetails{
			KillerName: match[4],
			VictimName: match[5],
			Cause:      match[6],
		}
	}
	return domain.KillDetails{}
}

func containsPlayer(players []string, name string) bool {
	for _, player := range players {
		if player == name {
			return true
		}
	}
	return false
}

func hasMatchEnded(scanner *bufio.Scanner) bool {
	return strings.Contains(scanner.Text(), endMatchLine1) ||
		strings.Contains(scanner.Text(), endMatchLine2)
}

func hasMatchStarted(scanner *bufio.Scanner) bool {
	return strings.Contains(scanner.Text(), initGame)
}
