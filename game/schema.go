package game

type GridType = rune

type MapGrid struct {
	Type     GridType `json:"type"`
	Soldiers *int     `json:"soldiers,omitempty"`
}

type GameMap struct {
	Width  uint      `json:"width"`
	Height uint      `json:"height"`
	Grids  []MapGrid `json:"grids"`
}

type UserScript struct {
	Type    string  `json:"type"`
	Content *string `json:"content,omitempty"`
}

type GameUser struct {
	ID     string     `json:"id"`
	Script UserScript `json:"script"`
}

type GameSetting struct {
	ID    string              `json:"id"`
	Map   GameMap             `json:"map"`
	Users map[string]GameUser `json:"users"`
}

type MoveAction struct {
	X        uint `json:"x"`
	Y        uint `json:"y"`
	Amount   uint `json:"amount"`
	Movement rune `json:"movement"`
}

type GridChange struct {
	X    uint    `json:"x"`
	Y    uint    `json:"y"`
	Grid MapGrid `json:"grid"`
}

type GameTick struct {
	Operator    rune         `json:"operator"`
	Changes     []GridChange `json:"changes"`
	Action      *MoveAction  `json:"action"`
	ActionValid bool         `json:"action_valid"`
}

type UserGameStat struct {
	Rounds         uint `json:"rounds"`
	Moves          uint `json:"moves"`
	SoldiersTotal  uint `json:"soldiers_total"`
	SoldiersKilled uint `json:"soldiers_killed"`
	GridsTaken     uint `json:"grids_taken"`
}

type GameResult struct {
	Winner rune         `json:"winner"`
	RStat  UserGameStat `json:"r_stat"`
	BStat  UserGameStat `json:"b_stat"`
}
