package users

// import (
// 	"OracleGo/internal/entities"
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// )

// func TestCreateAndValidateAccessToken(t *testing.T) {
// 	expected := entities.User{
// 		Email:    entities.Email{Value: "test@test.com"},
// 		Password: entities.Password{Value: "test123"},
// 	}
// 	token, err := CreateAccessToken(expected)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	claims, err := ValidateAccessToken(token)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	assert.Equal(t, expected.Email.GetValue(), claims.Email)
// }

// func TestCreateAndValidateRefreshToken(t *testing.T) {
// 	token, err := CreateRefreshToken()
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	_, err = ValidateRefreshToken(token)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// }
