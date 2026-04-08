package http_handlers

type JoinResponse struct {
	Token    string `json:"token"`
	PlayerID int    `json:"player_id"`
	MapSize  int    `json:"map_size"`
}
