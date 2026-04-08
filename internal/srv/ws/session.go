package ws

import (
	"context"
	"sync"
)

func (s *Server) CreateSession(ctx context.Context, token string, playerID int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.sessions[token] = sessionEntry{
		playerID:  playerID,
		conn:      nil,
		closeOnce: &sync.Once{},
	}

	return nil
}
