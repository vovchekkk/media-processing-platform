package service

import (
	"context"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"media-processing-platform/server/internal/domain"
	"media-processing-platform/server/internal/dto"
	"media-processing-platform/server/internal/repository"
)

type AuthService struct {
	userRepository    repository.User
	sessionRepository repository.Session
}

func NewAuthService(userRepo repository.User, sessionRepo repository.Session) *AuthService {
	return &AuthService{
		userRepository:    userRepo,
		sessionRepository: sessionRepo,
	}
}

func (authService *AuthService) Register(ctx context.Context, userDTO *dto.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userDTO.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	id, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	user := &domain.User{
		ID:       id,
		Username: userDTO.Username,
		Password: string(hashedPassword),
	}

	return authService.userRepository.CreateUser(ctx, user)
}

func (authService *AuthService) Login(ctx context.Context, userDTO *dto.User) (uuid.UUID, error) {
	user, err := authService.userRepository.GetByUsername(ctx, userDTO.Username)

	if err != nil {
		return uuid.Nil, err
	}

	if user == nil {
		return uuid.Nil, domain.ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userDTO.Password)); err != nil {
		return uuid.Nil, domain.ErrInvalidCredentials
	}

	sessionID, err := uuid.NewRandom()
	if err != nil {
		return uuid.Nil, err
	}

	session := &domain.Session{
		ID:     sessionID,
		UserID: user.ID,
	}

	if err := authService.sessionRepository.CreateSession(ctx, session); err != nil {
		return uuid.Nil, err
	}

	return sessionID, nil
}

func (authService *AuthService) Logout(ctx context.Context, sessionID uuid.UUID) error {
	return authService.sessionRepository.DeleteSession(ctx, sessionID)
}

func (authService *AuthService) ValidateToken(ctx context.Context, token uuid.UUID) (*domain.Session, error) {
	return authService.sessionRepository.GetSessionByID(ctx, token)
}
