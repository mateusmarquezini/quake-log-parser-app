package helper

const (
	INIT_GAME              = `InitGame`
	KILL_ACTION_LINE_REGEX = `Kill:\s+(\d+)\s+(\d+)\s+(\d+):\s+(.*?)\s+killed\s+(.*?)\s+by\s+(\S+)`
	END_MATCH_LINE_1       = "26  0:00 ------------------------------------------------------------"
	END_MATCH_LINE_2       = "ShutdownGame"
)
