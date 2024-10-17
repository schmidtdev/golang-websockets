package types

type Message struct {
	Content string `json:"content"`
	Type    string `json:"type"`
	Channel string `json:"channel"`
	IP      string `json:"ip"`
}
