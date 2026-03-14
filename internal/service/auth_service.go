package service

import (
	"database/sql"
	"errors"
	"hospital/internal/models"
	"hospital/internal/repository"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte("supersecret")

func Register(email, password string) (*models.User, error) {
	existingUser, err := repository.GetUserByEmail(email)
	if err == nil && existingUser != nil {
		return nil, errors.New("email already exists")
	}
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		ID:           uuid.New().String(),
		Email:        email,
		PasswordHash: string(hash),
		VerificationToken: sql.NullString{
			String: uuid.New().String(),
			Valid:  true,
		},
		IsVerified: false,
		Status:     "inactive",
	}

	err = repository.CreateUser(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func Verify(token string) error {
	return repository.VerifyUser(token)
}

func Login(email, password string) (string, string, error) {
	user, err := repository.GetUserByEmail(email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", "", errors.New("invalid credentials")
		}
		return "", "", err
	}

	if !user.IsVerified {
		return "", "", errors.New("account not verified")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", "", errors.New("invalid credentials")
	}

	accessToken := generateJWT(user.ID)
	refreshToken := uuid.New().String()

	err = repository.SaveRefreshToken(user.ID, refreshToken)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func generateJWT(userID string) string {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, _ := token.SignedString(jwtSecret)
	return signed
}

func RefreshToken(oldRefreshToken string) (string, string, error) {
	// 1. Lấy thông tin token từ DB
	tokenRecord, err := repository.GetRefreshToken(oldRefreshToken)
	if err != nil {
		return "", "", errors.New("invalid refresh token")
	}

	// 2. Kiểm tra xem token đã bị revoke chưa
	if tokenRecord.RevokedAt.Valid {
		// CẢNH BÁO BẢO MẬT: Token đã bị hủy mà vẫn mang ra dùng -> Có thể user bị hack
		// Ở đây ta chỉ trả về lỗi, nhưng trong thực tế nên log lại IP này.
		return "", "", errors.New("token has been revoked")
	}

	// 3. Kiểm tra hết hạn
	if time.Now().After(tokenRecord.ExpiresAt) {
		return "", "", errors.New("refresh token expired")
	}

	// 4. Token Rotation: Hủy token cũ để tránh dùng lại
	err = repository.RevokeRefreshToken(oldRefreshToken)
	if err != nil {
		return "", "", err
	}

	// 5. Tạo cặp token mới
	newAccessToken := generateJWT(tokenRecord.UserID)
	newRefreshToken := uuid.New().String()

	// 6. Lưu token mới vào DB
	err = repository.SaveRefreshToken(tokenRecord.UserID, newRefreshToken)
	if err != nil {
		return "", "", err
	}

	return newAccessToken, newRefreshToken, nil
}
func Logout(refreshToken string) error {
	return repository.RevokeRefreshToken(refreshToken)
}
