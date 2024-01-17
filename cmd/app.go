package main

import (
	"bufio"
	"encoding/json"
	"fmt"
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
	var matches []domain.Match
	var killEventsInMatch []domain.KillDetails
	var matchID int

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

					match.KillsByMeans[event.Cause]++

					addPlayerIfNotExists(match, event.KillerName)
					addPlayerIfNotExists(match, event.VictimName)

					if isAValidPlayer(event) && !isTheSamePlayer(event) {
						match.Kills[event.KillerName]++
					} else if victimHasPositiveKills(*match, event.VictimName) && !isTheSamePlayer(event) {
						match.Kills[event.VictimName]--
					}
				}

				match.TotalKills = len(killEventsInMatch)
				matches = append(matches, *match)
				killEventsInMatch = nil
			}
		}
	}

	jsonData, err := json.MarshalIndent(matches, "", "  ")
	if err != nil {
		println(err.Error())
		return
	}

	fmt.Println(string(jsonData))

	defer file.Close()
}

func addPlayerIfNotExists(match *domain.Match, playerName string) {
	if playerName != worldID && !helper.ContainsPlayer(match.Players, playerName) {
		match.Players = append(match.Players, playerName)
	}
}

func isAValidPlayer(kill domain.KillDetails) bool {
	return kill.KillerName != worldID
}

func isTheSamePlayer(kill domain.KillDetails) bool {
	return kill.KillerName == kill.VictimName
}

func victimHasPositiveKills(match domain.Match, victimName string) bool {
	if kills, ok := match.Kills[victimName]; ok {
		return kills > 0
	}
	return false
}
