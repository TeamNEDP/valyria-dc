package game

type Message struct {
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
}

type AuthData struct {
	Slots uint   `json:"slots" mapstructure:"slots"`
	Token string `json:"token" mapstructure:"token"`
}

type GameStartData struct {
	ID      string      `json:"id"`
	Setting GameSetting `json:"setting"`
}

type GameUpdateData struct {
	ID   string   `json:"ID" mapstructure:"ID"`
	Tick GameTick `json:"tick" mapstructure:"tick"`
}

type GameEndData struct {
	ID     string     `json:"ID" mapstructure:"ID"`
	Result GameResult `json:"result" mapstructure:"result"`
}
