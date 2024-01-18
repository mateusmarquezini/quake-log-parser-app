package service

const (
	matchIDPrefix       = "game_"
	worldID             = "<world>" // <world> isn't a player
	initGame            = `InitGame`
	killActionLineRegex = `Kill:\s+(\d+)\s+(\d+)\s+(\d+):\s+(.*?)\s+killed\s+(.*?)\s+by\s+(\S+)`
	endMatchLine1       = "26  0:00 ------------------------------------------------------------"
	endMatchLine2       = "ShutdownGame"
)
