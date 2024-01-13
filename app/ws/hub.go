package ws

type Hub struct {
	Rooms      map[string]*Room
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan *Message
}

func NewHub() (*Hub, error) {
	hub := &Hub{
		Rooms:      make(map[string]*Room),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan *Message, 5),
	}
	return hub, nil
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			if _, ok := h.Rooms[client.RoomId]; ok {
				r := h.Rooms[client.RoomId]
				if _, ok := r.Clients[client.Id]; !ok {
					r.Clients[client.Id] = client
				}
			}
		case client := <-h.Unregister:
			if _, ok := h.Rooms[client.RoomId]; ok {
				if _, ok := h.Rooms[client.RoomId].Clients[client.Id]; ok {
					if len(h.Rooms[client.RoomId].Clients) != 0 {
						h.Broadcast <- &Message{
							Content:  "User left the chat",
							RoomId:   client.RoomId,
							Username: client.Username,
						}
					}
					delete(h.Rooms[client.RoomId].Clients, client.Id)
					close(client.Message)
				}
			}
		case message := <-h.Broadcast:
			if _, ok := h.Rooms[message.RoomId]; ok {
				for _, client := range h.Rooms[message.RoomId].Clients {
					client.Message <- message
				}
			}
		}
	}
}
