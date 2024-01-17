package main

import (
	"bufio"

	"github.com/mateusmarquezini/quake-log-parser-app/helper"
	service "github.com/mateusmarquezini/quake-log-parser-app/service/match"
)

const filePath = "log/qgames.log"

func main() {

	file, err := helper.ReadFile(filePath)
	if err != nil {
		println(err.Error())
		return
	}

	scanner := bufio.NewScanner(file)

	s := service.NewMatchService()

	for scanner.Scan() {
		s.ProcessLine(scanner)
	}

	s.PrintMatchReport()

	defer file.Close()
}
