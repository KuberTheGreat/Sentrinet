package realtime

type WSMessage struct {
	Type string `json:"type"`
	Data interface{} `json:"data"`
}