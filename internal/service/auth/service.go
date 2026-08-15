package auth

import (
	"context"
	"time"

	"github.com/Grisha1Kadetov/TeamTaskTrackerService/internal/model/user"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	secretKey []byte
	userRepo  userRepo
}

type claims struct {
	UserID uuid.UUID `json:"user_id"`
	jwt.RegisteredClaims
}

func New(secretKey []byte, userRepo userRepo) *AuthService {
	return &AuthService{
		secretKey: secretKey,
		userRepo:  userRepo,
	}
}

func (s *AuthService) Register(ctx context.Context, email, password, name string) (user.User, string, error) {
	ok, err := s.userRepo.EmailExists(ctx, email)
	if err != nil {
		return user.User{}, "", err
	}
	if ok {
		return user.User{}, "", user.ErrEmailAlreadyExists
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return user.User{}, "", err
	}
	u := user.User{
		ID:           uuid.New(),
		Name:         name,
		Email:        email,
		PasswordHash: string(hash),
		CreatedAt:    time.Now().UTC(),
	}

	if err := s.userRepo.Create(ctx, u); err != nil {
		return user.User{}, "", err
	}

	token, err := generateToken(u.ID, s.secretKey)
	if err != nil {
		return user.User{}, "", err
	}

	return u, token, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	u, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", err
	}

	if bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)) != nil {
		return "", user.ErrIncorrectPassword
	}

	tokenString, err := generateToken(u.ID, s.secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *AuthService) ParseToken(ctx context.Context, tokenString string) (uuid.UUID, error) {
	claims := &claims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		return s.secretKey, nil
	})

	if err != nil {
		return uuid.Nil, err
	}

	return claims.UserID, nil
}

func generateToken(uuid uuid.UUID, secretKey []byte) (string, error) {
	now := time.Now().UTC()

	claims := claims{
		UserID: uuid,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   uuid.String(),
			ExpiresAt: jwt.NewNumericDate(now.Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(secretKey)
}
