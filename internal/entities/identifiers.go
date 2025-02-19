package entities

var (
	AccessTokenCookieName  = "access_token"
	RefreshTokenCookieName = "refresh_token"
	RegistrationCookieName = "registration_token"

	IsLoggedInKey   = "IsLoggedIn"
	ErrorMessageKey = "ErrorMessage"
)

type (
	AuthStatusContextKey         struct{}
	RegistrationStatusContextKey struct{}
	AuthClaimsContextKey         struct{}
)
