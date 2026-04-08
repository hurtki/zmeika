package ws

import (
	"sync"

	"github.com/hurtki/zmeika/internal/app"
	"github.com/hurtki/zmeika/internal/domain"
	"golang.org/x/net/websocket"
)

type sessionEntry struct {
	playerID int
	conn     *websocket.Conn

	closeOnce *sync.Once
}

type Server struct {
	sessions map[string]sessionEntry

	usecase *domain.GameUsecase

	mu sync.Mutex
}

func NewServer(usecase *domain.GameUsecase) *Server {
	return &Server{
		sessions: make(map[string]sessionEntry),
		usecase:  usecase,
	}
}

func (s *Server) HandleWS(conn *websocket.Conn, token string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.sessions[token]
	if !ok {
		conn.Close()
		return
	}
	if entry.conn != nil {
		conn.Close()
		return
	}
	entry.conn = conn
	s.sessions[token] = entry
	go s.readLoop(conn, token)
	go s.writeLoop(conn, token)
}

func (s *Server) readLoop(conn *websocket.Conn, token string) {
	buf := make([]byte, 1)
	n, err := conn.Read(buf)
	if n != 1 || err != nil {
		s.closeSession(token)
		return
	}

	motion, err := app.NewDirection(uint8(buf[0]))
	if err != nil {
		s.closeSession(token)
		return
	}
	s.mu.Lock()
	entry, ok := s.sessions[token]
	s.mu.Unlock()
	if !ok {
		s.closeSession(token)
		return
	}

	playerID := entry.playerID

	_ = s.usecase.Move(conn.Request().Context(), motion, playerID)
}

func (s *Server) writeLoop(conn *websocket.Conn, _ string) {
	for {
		plot, _ := s.usecase.GetMap(conn.Request().Context())
		conn.Write(SerializePlot(plot))
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
