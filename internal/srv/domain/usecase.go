package domain

import (
	"context"
	"math/rand"

	"github.com/hurtki/ascii-snake/internal/srv/app"

	"golang.org/x/sync/singleflight"
)

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const baseTokenSymbolsCount = 50

func genRandomString(length int) string {
	str := make([]rune, length)
	for i := range length {
		str[i] = rune(alphabet[int(rand.Int31())%len(alphabet)])
	}
	return string(str)
}

type GameUsecase struct {
	game *app.Game
	sm   SessionManager

	sg singleflight.Group
}

type SessionManager interface {
	CreateSession(ctx context.Context, token string, playerID int) error
}

func NewGameUsecase(game *app.Game, sm SessionManager) *GameUsecase {
	return &GameUsecase{
		game: game,
		sm:   sm,
	}
}

func (u *GameUsecase) JoinRoom(ctx context.Context) (JoinOut, error) {
	playerID, err := u.game.AddPlayer()
	if err != nil {
		return JoinOut{}, err
	}
	token := genRandomString(baseTokenSymbolsCount)
	size := u.game.GetMapSize()

	u.sm.CreateSession(ctx, token, playerID)

	return JoinOut{
		Token:    token,
		MapSize:  size,
		PlayerID: playerID,
	}, nil
}

func (u *GameUsecase) Move(ctx context.Context, motion app.Direction, playerID int) error {
	u.game.AddMove(app.Move{PlayerID: playerID, Direction: motion})
	return nil
}

func (u *GameUsecase) GetMap(ctx context.Context) ([][]app.Cell, error) {
	plot, _, _ := u.sg.Do("", func() (any, error) {
		return u.game.GetTickMap(), nil
	})

	field, _ := plot.([][]app.Cell)
	return field, nil
}
