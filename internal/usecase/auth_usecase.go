package usecase

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/srgjo27/e-learning/internal/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	FindByID(ctx context.Context, id string) (*entity.User, error)
	UpdatePassword(ctx context.Context, userID string, newHashed string) error
	UpdateEmail(ctx context.Context, userID string, newEmail string) error
	Delete(ctx context.Context, userID string) error
	ListAll(ctx context.Context) ([]*entity.User, error)
	UpdateRole(ctx context.Context, userID string, role entity.Role) error
	FindUsersByIDs(ctx context.Context, studentIDs []primitive.ObjectID) ([]*entity.User, error)
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

func (a *AuthUseCase) Register(ctx context.Context, email, password string, role entity.Role) error {
	existingUser, err := a.userRepo.FindByEmail(ctx, email)
	if err == nil && existingUser != nil {
		return entity.ErrEmailExists
	}

	if role != entity.RoleAdmin && role != entity.RoleTeacher && role != entity.RoleStudent {
		role = entity.RoleStudent
	}

	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &entity.User{
		Email:		email,
		Password:	string(hashedPwd),
		Role: 		role,
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
		"role": string(user.Role),
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenStr, err := token.SignedString(a.jwtSecret)

	return tokenStr, err
}

func (a *AuthUseCase) ParseToken(tokenStr string) (string, string, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error){
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, entity.ErrInvalidToken
		}
		return a.jwtSecret, nil
	}) 

	if err != nil || !token.Valid {
		return "", "", entity.ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", entity.ErrInvalidToken
	}

	userID, ok := claims["userId"].(string)
	if !ok {
		return "", "", entity.ErrInvalidToken
	}

	role, ok := claims["role"].(string)
	if !ok {
		role = string(entity.RoleStudent)
	}

	return userID, role, nil
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

func (a *AuthUseCase) GetProfile(ctx context.Context, userID string) (*entity.User, error) {
	return a.userRepo.FindByID(ctx, userID)
}

func (a *AuthUseCase) UpdateProfile(ctx context.Context, userID, newEmail, newPassword string) error {
	user, err := a.userRepo.FindByID(ctx, userID)
	if err != nil {
		 return err
	}

	if newEmail != "" && newEmail != user.Email {
		err := a.userRepo.UpdateEmail(ctx, userID, newEmail)
		if err != nil {
			return err
		}
	}

	if newPassword != "" {
		hashedPwd, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		err = a.userRepo.UpdatePassword(ctx, userID, string(hashedPwd))
		if err != nil {
			return err
		}

	}

	return nil
}

func (a *AuthUseCase) DeleteUser(ctx context.Context, userID string) error {
	return a.userRepo.Delete(ctx, userID)
}

func (a *AuthUseCase) ListAllUsers(ctx context.Context) ([]*entity.User, error) {
	return a.userRepo.ListAll(ctx)
}

// func (a *AuthUseCase) UpdateUserRole(ctx context.Context, userID string, role entity.Role) error {
// 	user, err := a.userRepo.FindByID(ctx, userID)
// 	if err != nil {
// 		return err
// 	}

// 	if role != entity.RoleAdmin && role != entity.RoleTeacher && role != entity.RoleStudent {
// 		return errors.New("invalid role")
// 	}

// 	if user.Role == role {
// 		return nil
// 	}

// 	return a.userRepo.UpdateEmail(ctx, userID, user.Email)
// }

func (a *AuthUseCase) UpdateUserRole(ctx context.Context, userID string, role entity.Role) error {
	return a.userRepo.UpdateRole(ctx, userID, role)
}