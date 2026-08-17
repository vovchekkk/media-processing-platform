package service

import (
	"context"

	"golang.org/x/crypto/bcrypt"
	"github.com/google/uuid"

	"media-processing-platform/internal/repository"
	"media-processing-platform/internal/domain"
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

func (authService *AuthService) Register(ctx context.Context, userDTO *domain.UserDTO) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userDTO.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &domain.User{
		Username: userDTO.Username,
		Password: string(hashedPassword),
	}

	return authService.userRepository.CreateUser(ctx, user)
}

func (authService *AuthService) Login(ctx context.Context, userDTO *domain.UserDTO) (uuid.UUID, error) {
	user, err := authService.userRepository.GetByUsername(ctx, userDTO.Username)
	if err != nil || user == nil {
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

func (authService *AuthService) ValidateToken(ctx context.Context, token uuid.UUID) (*domain.Session, error) {
	return authService.sessionRepository.GetSessionByID(ctx, token)
}