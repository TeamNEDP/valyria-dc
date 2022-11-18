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
	Width  int       `json:"width" mapstructure:"width"`
	Height int       `json:"height" mapstructure:"height"`
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
	X        int    `json:"x" mapstructure:"x"`
	Y        int    `json:"y" mapstructure:"y"`
	Amount   int    `json:"amount" mapstructure:"amount"`
	Movement string `json:"movement" mapstructure:"movement"`
}

type GridChange struct {
	X    int     `json:"x" mapstructure:"x"`
	Y    int     `json:"y" mapstructure:"y"`
	Grid MapGrid `json:"grid" mapstructure:"grid"`
}

type GameTick struct {
	Operator    string       `json:"operator" mapstructure:"operator"`
	Changes     []GridChange `json:"changes" mapstructure:"changes"`
	Action      *MoveAction  `json:"action" mapstructure:"action"`
	ActionValid bool         `json:"action_valid" mapstructure:"action_valid"`
	ActionError string       `json:"action_error" mapstructure:"action_error"`
}

type UserGameStat struct {
	Rounds         int `json:"rounds" mapstructure:"rounds"`
	Moves          int `json:"moves" mapstructure:"moves"`
	SoldiersTotal  int `json:"soldiers_total" mapstructure:"soldiers_total"`
	SoldiersKilled int `json:"soldiers_killed" mapstructure:"soldiers_killed"`
	GridsTaken     int `json:"grids_taken" mapstructure:"grids_taken"`
}

type GameResult struct {
	Winner string       `json:"winner" mapstructure:"winner"`
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
