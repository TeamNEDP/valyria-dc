package game

import (
	"database/sql/driver"
	"encoding/json"
)

type GridType = string

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
	ID     string     `json:"ID" mapstructure:"ID"`
	Script UserScript `json:"script" mapstructure:"script"`
}

type GameSetting struct {
	Map   GameMap             `json:"map" mapstructure:"map"`
	Users map[string]GameUser `json:"users" mapstructure:"users"`
}

func (gameSetting *GameSetting) Scan(src interface{}) error {
	return json.Unmarshal([]byte(src.(string)), &gameSetting)
}

func (gameSetting GameSetting) Value() (driver.Value, error) {
	val, err := json.Marshal(gameSetting)
	return string(val), err
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

func (gameResult *GameResult) Scan(src interface{}) error {
	return json.Unmarshal([]byte(src.(string)), &gameResult)
}

func (gameResult GameResult) Value() (driver.Value, error) {
	val, err := json.Marshal(gameResult)
	return string(val), err
}

type GameTicks struct {
	Ticks []GameTick `json:"ticks"`
}

func (gameTicks *GameTicks) Scan(src interface{}) error {
	return json.Unmarshal([]byte(src.(string)), &gameTicks)
}

func (gameTicks GameTicks) Value() (driver.Value, error) {
	val, err := json.Marshal(gameTicks)
	return string(val), err
}
