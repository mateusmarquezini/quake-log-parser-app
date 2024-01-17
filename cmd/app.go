package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/mateusmarquezini/quake-log-parser-app/domain"
	"github.com/mateusmarquezini/quake-log-parser-app/helper"
)

const (
	filePath      = "log/qgames.log"
	matchIDPrefix = "game_"
	worldID       = "<world>" // <world> isn't a player
)

func main() {

	file, err := helper.ReadFile(filePath)
	if err != nil {
		println(err.Error())
		return
	}

	scanner := bufio.NewScanner(file)

	var match = &domain.Match{}
	var matches []map[string]any
	var killEventsInMatch []domain.KillDetails
	var matchID int
	response := make(map[string]any)

	for scanner.Scan() {
		if strings.Contains(scanner.Text(), helper.INIT_GAME) {
			matchID++
			match = &domain.Match{
				MatchID:      matchIDPrefix + fmt.Sprint(matchID),
				Players:      []string{},
				Kills:        map[string]int{},
				KillsByMeans: map[string]int{},
			}
		}

		matchKillDetails := domain.ParseKillDetails(scanner)

		if matchKillDetails.Cause != "" {
			killEventsInMatch = append(killEventsInMatch, matchKillDetails)
		}

		if helper.HasMatchEnded(scanner) {
			if len(killEventsInMatch) == 0 {
				// if there is no kills, the current match is not considered valid
				matchID -= 1
			} else {
				for _, event := range killEventsInMatch {

					incrementKillsByCause(match, event.Cause)

					addPlayerIfNotExists(match, event.KillerName)
					addPlayerIfNotExists(match, event.VictimName)

					if !isTheSamePlayer(event) {
						if isAValidPlayer(event) {
							incrementKillCount(match, event.KillerName)
						} else {
							decrementKillCount(match, event.VictimName)
						}
					}
				}

				match.TotalKills = len(killEventsInMatch)
				response[match.MatchID] = match
				matches = append(matches, response)
				response = make(map[string]any)
				killEventsInMatch = nil
			}
		}
	}

	jsonData, err := json.MarshalIndent(matches, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling matches to JSON: %v", err)
		return
	}

	fmt.Println(string(jsonData))

	defer file.Close()
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
	if playerName != worldID && !helper.ContainsPlayer(match.Players, playerName) {
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

// func victimHasPositiveKills(match domain.Match, victimName string) bool {
// 	if kills, ok := match.Kills[victimName]; ok {
// 		return kills > 0
// 	}
// 	return false
// }
