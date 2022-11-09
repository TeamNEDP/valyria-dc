package game

type GridType = rune

type MapGrid struct {
	Type     GridType `json:"type" mapstructure:"type"`
	Soldiers *int     `json:"soldiers,omitempty" mapstructure:"soldiers"`
}

type GameMap struct {
	Width  uint      `json:"width" mapstructure:"width"`
	Height uint      `json:"height" mapstructure:"height"`
	Grids  []MapGrid `json:"grids" mapstructure:"grids"`
}

type UserScript struct {
	Type    string  `json:"type" mapstructure:"type"`
	Content *string `json:"content,omitempty" mapstructure:"content"`
}

type GameUser struct {
	ID     string     `json:"id" mapstructure:"id"`
	Script UserScript `json:"script" mapstructure:"script"`
}

type GameSetting struct {
	ID    string              `json:"id" mapstructure:"id"`
	Map   GameMap             `json:"map" mapstructure:"map"`
	Users map[string]GameUser `json:"users" mapstructure:"users"`
}

type MoveAction struct {
	X        uint `json:"x" mapstructure:"x"`
	Y        uint `json:"y" mapstructure:"y"`
	Amount   uint `json:"amount" mapstructure:"amount"`
	Movement rune `json:"movement" mapstructure:"movement"`
}

type GridChange struct {
	X    uint    `json:"x" mapstructure:"x"`
	Y    uint    `json:"y" mapstructure:"y"`
	Grid MapGrid `json:"grid" mapstructure:"grid"`
}

type GameTick struct {
	Operator    rune         `json:"operator" mapstructure:"operator"`
	Changes     []GridChange `json:"changes" mapstructure:"changes"`
	Action      *MoveAction  `json:"action" mapstructure:"action"`
	ActionValid bool         `json:"action_valid" mapstructure:"action_valid"`
}

type UserGameStat struct {
	Rounds         uint `json:"rounds" mapstructure:"rounds"`
	Moves          uint `json:"moves" mapstructure:"moves"`
	SoldiersTotal  uint `json:"soldiers_total" mapstructure:"soldiers_total"`
	SoldiersKilled uint `json:"soldiers_killed" mapstructure:"soldiers_killed"`
	GridsTaken     uint `json:"grids_taken" mapstructure:"grids_taken"`
}

type GameResult struct {
	Winner rune         `json:"winner" mapstructure:"winner"`
	RStat  UserGameStat `json:"r_stat" mapstructure:"r_stat"`
	BStat  UserGameStat `json:"b_stat" mapstructure:"b_stat"`
}
