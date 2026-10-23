package session

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"time"

	"github.com/redis/go-redis/v9"
)

// Session structure
type Session struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
	Role  string `json:"role"`
	Photo string `json:"photo"`
}

type SessionManager struct {
	redisClient *redis.Client
}

func NewSessionManager(redisCLient *redis.Client) *SessionManager {
	return &SessionManager{redisClient: redisCLient}
}

// CreateSession stores session in Redis
func (sm *SessionManager) CreateSession(ctx context.Context, sessionID string, sessionData Session, ttl time.Duration) error {
	sessionJSON, err := json.Marshal(sessionData)
	if err != nil {
		return err
	}

	return sm.redisClient.Set(ctx, sessionID, sessionJSON, ttl).Err()
}

// GetSession retrieves session from Redis
func (sm *SessionManager) GetSession(ctx context.Context, sessionID string) (*Session, error) {
	sessionJSON, err := sm.redisClient.Get(ctx, sessionID).Result()
	if err == redis.Nil {
		return nil, fmt.Errorf("session not found")
	} else if err != nil {
		return nil, err
	}

	var session Session
	if err := json.Unmarshal([]byte(sessionJSON), &session); err != nil {
		return nil, err
	}

	return &session, nil
}

// DeleteSession removes session from Redis
func (sm *SessionManager) DeleteSession(ctx context.Context, sessionID string) error {
	return sm.redisClient.Del(ctx, sessionID).Err()
}

// InvalidateAllSessions clears all sessions (Admin action)
func (sm *SessionManager) InvalidateAllSessions(ctx context.Context) error {
	return sm.redisClient.FlushDB(ctx).Err() // Clears all sessions
}

func GetUserSession(c *fiber.Ctx) (*Session, error) {
	userSession := c.Locals("user").(*Session)
	if userSession == nil {
		return nil, c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "User not found"})
	}
	user := Session{
		ID:    userSession.ID,
		Email: userSession.Email,
		Name:  userSession.Name,
		Role:  userSession.Role,
		Photo: userSession.Photo,
	}

	return &user, nil
}
