package service

import "bufio"

type Service interface {
	ProcessLine(scanner *bufio.Scanner)
	PrintMatchReport()
}
