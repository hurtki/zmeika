package ws

import (
	"context"
	"log/slog"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/hurtki/zmeika/internal/app"
	"github.com/hurtki/zmeika/internal/domain"
)

type sessionEntry struct {
	playerID int
	conn     *websocket.Conn

	closeOnce *sync.Once
}

type Server struct {
	sessions map[string]sessionEntry

	usecase *domain.GameUsecase

	mu     sync.Mutex
	logger *slog.Logger
}

func NewServer(usecase *domain.GameUsecase, logger *slog.Logger) *Server {
	return &Server{
		sessions: make(map[string]sessionEntry),
		usecase:  usecase,
		logger:   logger,
	}
}

func (s *Server) HandleWS(conn *websocket.Conn, token string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.sessions[token]
	if !ok {
		s.logger.Debug("no session found, closing conn", "tok", token, "addr", conn.RemoteAddr())
		conn.Close()
		return
	}
	if entry.conn != nil {
		s.logger.Debug("conn already exists, closing new one", "tok", token, "addr", conn.RemoteAddr())
		conn.Close()
		return
	}
	entry.conn = conn
	s.sessions[token] = entry
	go s.readLoop(conn, token)
	go s.writeLoop(conn, token)
}

func (s *Server) readLoop(conn *websocket.Conn, token string) {
	for {
		_, buf, err := conn.ReadMessage()
		if err != nil {
			s.logger.Debug("can't read message, closing session", "err", err, "tok", token)
			s.closeSession(token)
			return
		}

		motion, err := app.NewDirection(uint8(buf[0]))
		if err != nil {
			s.closeSession(token)
			return
		}
		s.logger.Debug("got new move", "direction", motion)
		s.mu.Lock()
		entry, ok := s.sessions[token]
		s.mu.Unlock()
		if !ok {
			s.closeSession(token)
			return
		}

		playerID := entry.playerID

		_ = s.usecase.Move(context.Background(), motion, playerID)

	}
}

func (s *Server) writeLoop(conn *websocket.Conn, _ string) {
	for {
		plot, _ := s.usecase.GetMap(context.Background())
		s.logger.Debug("writing to conn new plot")
		conn.WriteMessage(websocket.BinaryMessage, SerializePlot(plot))
	}
}

func (s *Server) closeSession(token string) {
	s.mu.Lock()
	entry, ok := s.sessions[token]
	if ok {
		delete(s.sessions, token)
	}
	s.mu.Unlock()

	if !ok {
		return
	}

	entry.closeOnce.Do(func() {
		_ = entry.conn.Close()
	})
}
