package ws

type CreateRoomRequest struct {
	Id   string `json:"id" binding:"required"`
	Name string `json:"name" binding:"required"`
}

type GetRoomResponse struct {
	Id   string `json:"id" binding:"required"`
	Name string `json:"name"`
}

type GetClientByRoomIdParams struct {
	Id string `uri:"id" binding:"required"`
}

type GetClientByRoomIdResponse struct {
	Id       string `json:"id" binding:"required"`
	Username string `json:"username" binding:"required"`
}

type JoinRoomFormParam struct {
	ClientId string `form:"clientId" binding:"required"`
	Username string `form:"username" binding:"required"`
}

type JoinRoomURIParam struct {
	Id string `uri:"id" binding:"required"`
}
