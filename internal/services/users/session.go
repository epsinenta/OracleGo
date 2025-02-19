package users

import (
	"OracleGo/internal/entities"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (um *UsersManager) StartSession(user entities.User) (sessionId string, err error) {
	sessionId = uuid.NewString()

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	user.Password = string(hashedPass)

	err = um.cache.CacheData("session:"+sessionId, user, time.Minute*15)
	if err != nil {
		return "", err
	}

	err = um.cache.CacheData("session:"+user.Email, sessionId, time.Minute*15)
	if err != nil {
		um.cache.DeleteCachedData([]string{"session:" + sessionId})
		return "", err
	}

	return
}

func (um *UsersManager) GetSessionIdByEmail(email string) (sessionId string, err error) {
	err = um.cache.GetCachedData("session:"+email, &sessionId)
	return
}

func (um *UsersManager) GetSessionById(sessionId string) (user entities.User, err error) {
	err = um.cache.GetCachedData("session:"+sessionId, &user)
	return
}

func (um *UsersManager) EndSession(sessionId string, email string) error {
	err := um.cache.DeleteCachedData([]string{"session:" + sessionId, "session:" + email})
	return err
}
