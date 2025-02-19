package users

import (
	"OracleGo/internal/entities"
	"OracleGo/internal/log"
	"database/sql"
	"fmt"
	"net/smtp"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

var (
	smtpHost    = "smtp.gmail.com"
	smtpPort    = "587"
	senderEmail = "example@example.com" // почта
	appPassword = "exam plea pppa sswd" // пароль к приложению почты
)

func (um *UsersManager) Login(user entities.User) (token string, err error) {
	userFromDb, err := um.getUser(user.Email)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(userFromDb.Password), []byte(user.Password))
	if err != nil {
		return "", err
	}

	return um.CreateRefreshToken(userFromDb)
}

// требует заполнения полей email password
func (um *UsersManager) StartRegistration(user entities.User) (sessionId string, err error) {
	_, err = um.GetSessionIdByEmail(user.Email)
	if err != nil {
		if errors.Cause(err) != redis.Nil {
			return "", err
		}
	} else {
		return "", entities.SessionAlreadyStartedError
	}

	_, err = um.getUser(user.Email)
	if err != nil {
		if errors.Cause(err) != sql.ErrNoRows {
			return "", err
		}
	} else {
		return "", entities.AlreadyRegisteredError
	}

	sessionId, err = um.StartSession(user)
	if err != nil {
		return "", err
	}

	um.sendConfirmationEmail(sessionId, user.Email)

	return
}

func (um *UsersManager) CompleteRegistration(sessionId string) error {
	user, err := um.GetSessionById(sessionId)
	if err != nil {
		return err
	}

	err = um.addUsers([]entities.User{user})
	if err != nil {
		return err
	}

	err = um.EndSession(sessionId, user.Email)
	if err != nil {
		return err
	}

	return nil
}

func (um *UsersManager) ResendEmail(sessionId string) error {
	type Result struct {
		Email string
	}
	res := Result{}
	err := um.cache.GetCachedData("confirmation:"+sessionId, &res)
	if err != nil {
		return err
	}

	um.sendConfirmationEmail(sessionId, res.Email)

	return nil
}

func (um *UsersManager) sendConfirmationEmail(sessionId string, email string) {
	go func() {
		subject := "Subject: Confirmation message\n"
		body := fmt.Sprint("your confirmation link is ", "http://localhost:8080/confirm?sessionId=", sessionId)

		msg := []byte(subject + "\n" + body)

		auth := smtp.PlainAuth("", senderEmail, appPassword, smtpHost)

		err := smtp.SendMail(smtpHost+":"+smtpPort, auth, senderEmail, []string{email}, msg)
		if err != nil {
			log.New().Log(errors.Wrap(err, "sending confirmation email"))
		}
	}()
}

func (um *UsersManager) getUser(email string) (entities.User, error) {
	userRows, err := um.repo.GetRows("users", []string{"email", "password"}, map[string][]string{"email": {email}})
	if err != nil {
		return entities.User{}, errors.Wrap(err, "getting rows from db")
	}
	if len(userRows) == 0 {
		return entities.User{}, errors.Wrap(sql.ErrNoRows, "getting rows from db")
	}

	return entities.User{
		Email:    userRows[0][0],
		Password: userRows[0][1],
	}, nil
}

func (um *UsersManager) addUsers(users []entities.User) error {
	ids := make([]string, len(users))
	emails := make([]string, len(users))
	passwords := make([]string, len(users))
	for i := 0; i < len(users); i++ {
		ids[i] = uuid.NewString()
		emails[i] = users[i].Email
		passwords[i] = users[i].Password
	}

	err := um.repo.AddRows("users", map[string][]string{"id": ids, "email": emails, "password": passwords})
	if err != nil {
		return err
	}

	return nil
}
