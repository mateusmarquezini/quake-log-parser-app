package domain

import (
	"bufio"
	"regexp"

	"github.com/mateusmarquezini/quake-log-parser-app/helper"
)

type KillDetails struct {
	KillerName string
	VictimName string
	Cause      string
}

func ParseKillDetails(scanner *bufio.Scanner) KillDetails {
	match := regexp.MustCompile(helper.KILL_ACTION_LINE_REGEX).FindStringSubmatch(scanner.Text())

	if len(match) > 0 {
		return KillDetails{
			KillerName: match[4],
			VictimName: match[5],
			Cause:      match[6],
		}
	}
	return KillDetails{}
}
