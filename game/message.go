package game

type Message struct {
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
}

type AuthData struct {
	Slots uint   `json:"slots" mapstructure:"slots"`
	Token string `json:"token" mapstructure:"token"`
}
