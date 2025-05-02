package usecase

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/srgjo27/e-learning/internal/entity"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	UpdatePassword(ctx context.Context, userID string, newHashed string) error
	UpdateEmail(ctx context.Context, userID string, newEmail string) error
}

type AuthUseCase struct {
	userRepo UserRepository
	jwtSecret []byte
}

func NewAuthUseCase(repo UserRepository, jwtSecret []byte) *AuthUseCase {
	return &AuthUseCase{
		userRepo: repo,
		jwtSecret: jwtSecret,
	}
}

var resetTokens = make(map[string]string)

func (a *AuthUseCase) Register(ctx context.Context, email, password string) error {
	existingUser, err := a.userRepo.FindByEmail(ctx, email)
	if err == nil && existingUser != nil {
		return entity.ErrEmailExists
	}

	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &entity.User{
		Email:		email,
		Password:	string(hashedPwd),
		CreatedAt: 	time.Now(),
	}

	return a.userRepo.Create(ctx, user)
}

func (a *AuthUseCase) Login(ctx context.Context, email, password string) (string, error) {
	user, err := a.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return "", entity.ErrUserNotFound
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", entity.ErrInvalidPassword
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": user.ID.Hex(),
		"email": user.Email,
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenStr, err := token.SignedString(a.jwtSecret)

	return tokenStr, err
}

func (a *AuthUseCase) ParseToken(tokenStr string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error){
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, entity.ErrInvalidToken
		}
		return a.jwtSecret, nil
	}) 

	if err != nil || !token.Valid {
		return "", entity.ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", entity.ErrInvalidToken
	}

	userID, ok := claims["userId"].(string)
	if !ok {
		return "", entity.ErrInvalidToken
	}

	return userID, nil
}

func (a *AuthUseCase) RequestPasswordReset(ctx context.Context, email string) (string, error) {
	user, err := a.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return "", entity.ErrUserNotFound
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": user.Email,
		"exp": time.Now().Add(time.Hour * 1).Unix(),
	})

	tokenStr, err := token.SignedString(a.jwtSecret)
	if err != nil {
		return "", err
	}

	resetTokens[tokenStr] = user.Email

	return tokenStr, nil
}

func (a *AuthUseCase) ResetPassword(ctx context.Context, tokenStr, newPassword string) error {
	email, ok := resetTokens[tokenStr]
	if !ok {
		return entity.ErrInvalidToken
	}

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error){
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, entity.ErrInvalidToken
		}
		return a.jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return entity.ErrInvalidToken
	}

	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user, err := a.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return entity.ErrUserNotFound
	}

	err = a.userRepo.UpdatePassword(ctx, user.ID.Hex(), string(hashedPwd))
	if err == nil {
		delete(resetTokens, tokenStr)
	}

	return err
}

