package users

import (
	"OracleGo/internal/entities"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

var (
	accessSign  = []byte("my-secret-access-token")
	refreshSign = []byte("my-secret-refresh-token")
)

func (um *UsersManager) CreateAccessToken(user entities.User) (string, error) {
	unsignedToken := jwt.NewWithClaims(jwt.SigningMethodHS256, entities.Claims{
		UserId: user.Id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 5).Unix(),
		},
	})

	token, err := unsignedToken.SignedString(accessSign)
	if err != nil {
		return "", errors.Wrap(err, "signing token")
	}

	return token, nil
}

func (um *UsersManager) ValidateAccessToken(token string) (*entities.Claims, error) {
	claims := &entities.Claims{}
	t, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return accessSign, nil
	})

	if err != nil {
		return nil, errors.Wrap(err, "parsing token")
	}

	if !t.Valid {
		return claims, errors.New("invalid token")
	}

	return claims, nil
}

func (um *UsersManager) CreateRefreshToken(user entities.User) (string, error) {
	unsignedToken := jwt.NewWithClaims(jwt.SigningMethodHS256, entities.Claims{
		UserId:  user.Id,
		Version: "0",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24 * 7).Unix(),
		},
	})

	token, err := unsignedToken.SignedString(refreshSign)
	if err != nil {
		return "", errors.Wrap(err, "signing token")
	}

	err = um.repo.AddRows("tokens", map[string][]string{"user_id": {user.Id}, "version": {"0"}})
	if err != nil {
		return "", errors.Wrap(err, "adding token to database")
	}

	return token, nil
}

func (um *UsersManager) ValidateRefreshToken(token string) (*entities.Claims, error) {
	claims := &entities.Claims{}
	t, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return refreshSign, nil
	})

	if err != nil {
		return nil, errors.Wrap(err, "parsing token")
	}

	if !t.Valid {
		return nil, errors.New("invalid token")
	}

	rows, err := um.repo.GetRows("tokens", []string{"version"}, map[string][]string{"user_id": {claims.UserId}})
	if err != nil {
		return nil, errors.Wrap(err, "getting token version from database")
	}

	if claims.Version != rows[0][0] {
		return nil, errors.New("versions of tokens didn't match")
	}

	return claims, nil
}

func (um *UsersManager) RefreshAccessToken(refreshToken string) (string, *entities.Claims, error) {
	claims, err := um.ValidateRefreshToken(refreshToken)
	if err != nil {
		return "", nil, err
	}

	token, err := um.CreateAccessToken(entities.User{Id: claims.UserId})
	if err != nil {
		return "", nil, err
	}

	return token, claims, nil
}
