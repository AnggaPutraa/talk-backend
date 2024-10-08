package ws

type Message struct {
	Content  string `json:"content"`
	RoomId   string `json:"room_id"`
	Username string `json:"username"`
}

type Room struct {
	Id      string             `json:"id"`
	Name    string             `json:"name"`
	Clients map[string]*Client `json:"clients"`
}
