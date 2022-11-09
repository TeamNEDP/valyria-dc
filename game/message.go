package game

type Message struct {
	Event string                 `json:"event"`
	Data  map[string]interface{} `json:"data"`
}
