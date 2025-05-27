package usecase

import (
	"encoding/json"
	"fmt"
	"net/smtp"
	"os"
	"time"
	"user-service/domain"
	local_redis "user-service/internal/redis"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

type UserUsecase struct {
	UserRepo domain.UserRepository
	redis    *redis.Client
}

func NewUserUsecase(userRepo domain.UserRepository, redisClient *redis.Client) *UserUsecase {
	return &UserUsecase{
		UserRepo: userRepo,
		redis:    redisClient,
	}
}

func (uc *UserUsecase) RegisterUser(username, password string) (string, error) {
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return "", err
	}
	user := &domain.User{
		UserID:   GenerateUserID(),
		Username: username,
		Password: hashedPassword,
	}
	SendEmail("Hello New User", "Welcme for my project", username)
	return uc.UserRepo.Create(user)
}

func (uc *UserUsecase) AuthenticateUser(username, password string) (string, error) {
	return uc.UserRepo.Authenticate(username, password)
}

func (uc *UserUsecase) GetUserProfile(userID string) (*domain.User, error) {
	key := fmt.Sprintf("user:%d", userID)

	// Check cache
	val, err := uc.redis.Get(local_redis.Ctx, key).Result()
	if err == nil {
		var cachedUser domain.User
		json.Unmarshal([]byte(val), &cachedUser)
		return &cachedUser, nil
	}
	user, err := uc.UserRepo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}
	jsonData, _ := json.Marshal(user)
	uc.redis.Set(local_redis.Ctx, key, jsonData, time.Minute*5)
	return user, nil
}

func GenerateUserID() string {
	return fmt.Sprintf("user-%d", time.Now().UnixNano())
}

func SendEmail(subject, message string, recipient string) error {
	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("Error loading .env file")
	}

	from := os.Getenv("EMAIL")
	password := os.Getenv("EMAIL_PASSWORD")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	msg := []byte(fmt.Sprintf("Subject: %s\n\n%s", subject, message))

	auth := smtp.PlainAuth("", from, password, smtpHost)

	return smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{recipient}, msg)
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
