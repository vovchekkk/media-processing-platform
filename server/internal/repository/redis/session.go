package redis

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	"media-processing-platform/server/internal/domain"
	"media-processing-platform/server/internal/repository"
)

type redisSession struct {
	db *redis.Client
}

func NewRedisSessionRepository(db *redis.Client) repository.Session {
	return &redisSession{db: db}
}

var _ repository.Session = (*redisSession)(nil)

func (r *redisSession) CreateSession(ctx context.Context, session *domain.Session, ttl time.Duration) error {
	key := "session:" + session.ID.String()

	return r.db.Set(
		ctx,
		key,
		session.UserID.String(),
		ttl,
	).Err()
}

func (r *redisSession) GetSessionByID(ctx context.Context, sessionID uuid.UUID) (*domain.Session, error) {
	key := "session:" + sessionID.String()

	userID, err := r.db.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, domain.ErrSessionNotFound
		}

		return nil, err
	}

	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	return &domain.Session{
		ID:     sessionID,
		UserID: parsedUserID,
	}, nil
}

func (r *redisSession) DeleteSession(ctx context.Context, sessionID uuid.UUID) error {
	key := "session:" + sessionID.String()

	return r.db.Del(ctx, key).Err()
}
