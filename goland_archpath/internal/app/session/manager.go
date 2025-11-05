package session

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type SessionData struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
}

type Manager struct {
	client *redis.Client
	ttl    time.Duration
}

func NewManager(redisAddr string, ttl time.Duration) (*Manager, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "",
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &Manager{
		client: client,
		ttl:    ttl,
	}, nil
}

func (m *Manager) CreateSession(ctx context.Context, userID uint, role string) (string, error) {
	sessionID, err := generateSessionID()
	if err != nil {
		return "", err
	}

	data := SessionData{
		UserID: userID,
		Role:   role,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	key := fmt.Sprintf("session:%s", sessionID)
	if err := m.client.Set(ctx, key, jsonData, m.ttl).Err(); err != nil {
		return "", err
	}

	return sessionID, nil
}

func (m *Manager) GetSession(ctx context.Context, sessionID string) (*SessionData, error) {
	key := fmt.Sprintf("session:%s", sessionID)
	val, err := m.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, fmt.Errorf("session not found")
	}
	if err != nil {
		return nil, err
	}

	var data SessionData
	if err := json.Unmarshal([]byte(val), &data); err != nil {
		return nil, err
	}

	return &data, nil
}

func (m *Manager) DeleteSession(ctx context.Context, sessionID string) error {
	key := fmt.Sprintf("session:%s", sessionID)
	return m.client.Del(ctx, key).Err()
}

func generateSessionID() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}