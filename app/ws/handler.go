package ws

import (
	"fmt"
	"net/http"

	"github.com/AnggaPutraa/talk-backend/exceptions"
	"github.com/gin-gonic/gin"
)

type WebSocketHandler struct {
	hub Hub
}

func NewWebSocketHandler(h Hub) *WebSocketHandler {
	return &WebSocketHandler{
		hub: h,
	}
}

func (h *WebSocketHandler) CreateRoom(c *gin.Context) {
	var request CreateRoomRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, exceptions.ErrorResponse(err))
		return
	}
	h.hub.Rooms[request.Id] = &Room{
		Id:      request.Id,
		Name:    request.Name,
		Clients: make(map[string]*Client),
	}
	c.JSON(http.StatusOK, request)
}

func (h *WebSocketHandler) JoinRoom(c *gin.Context) {
	connection, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, exceptions.ErrorResponse(err))
		return
	}
	var uriParam JoinRoomURIParam
	if err := c.ShouldBindUri(&uriParam); err != nil {
		c.JSON(http.StatusBadRequest, exceptions.ErrorResponse(err))
		return
	}
	var formParam JoinRoomFormParam
	if err := c.ShouldBind(&formParam); err != nil {
		c.JSON(http.StatusBadRequest, exceptions.ErrorResponse(err))
		return
	}
	client := &Client{
		Conn:     connection,
		Message:  make(chan *Message, 10),
		Id:       formParam.ClientId,
		RoomId:   uriParam.Id,
		Username: formParam.Username,
	}
	message := &Message{
		Content:  fmt.Sprintf("A new user named %s has joined the room", formParam.Username),
		RoomId:   uriParam.Id,
		Username: formParam.Username,
	}
	h.hub.Register <- client
	h.hub.Broadcast <- message
	go client.writeMessage()
	client.readMessage(&h.hub)
}

func (h *WebSocketHandler) GetRooms(c *gin.Context) {
	response := make([]GetRoomResponse, 0)
	for _, room := range h.hub.Rooms {
		roomResponse := &GetRoomResponse{
			Id:   room.Id,
			Name: room.Name,
		}
		response = append(response, *roomResponse)
	}
	c.JSON(http.StatusOK, response)
}

func (h *WebSocketHandler) GetClientsByRoomId(c *gin.Context) {
	var params GetClientByRoomIdParams
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, exceptions.ErrorResponse(err))
		return
	}
	response := make([]GetClientByRoomIdResponse, 0)
	if _, ok := h.hub.Rooms[params.Id]; !ok {
		c.JSON(http.StatusNotFound, nil)
		return
	}
	for _, client := range h.hub.Rooms[params.Id].Clients {
		clientResponse := &GetClientByRoomIdResponse{
			Id:       client.Id,
			Username: client.Username,
		}
		response = append(response, *clientResponse)
	}
	c.JSON(http.StatusOK, response)
}
