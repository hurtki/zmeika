package http_handlers

import (
	"encoding/json"
	"net/http"

	"github.com/hurtki/zmeika/internal/domain"
)

type JoinHandler struct {
	usecase *domain.GameUsecase
}

func NewJoinHandler(usecase *domain.GameUsecase) *JoinHandler {
	return &JoinHandler{
		usecase: usecase,
	}
}

func (h *JoinHandler) Join(rw http.ResponseWriter, req *http.Request) {
	out, err := h.usecase.JoinRoom(req.Context())
	if err != nil {
		responseErrorJson(rw, http.StatusServiceUnavailable, err.Error())
		return
	}
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(JoinResponse{
		Token:    out.Token,
		PlayerID: out.PlayerID,
		MapSize:  out.MapSize,
	})
}
