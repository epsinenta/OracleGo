package middlewares

import (
	"OracleGo/internal/entities"
	"OracleGo/internal/log"
	"OracleGo/internal/services"
	"context"
	"fmt"
	"net/http"
	"time"
)

type MiddlewareManager struct {
	services *services.ServicesManager
	logger   *log.Logger
}

func NewMiddlewareManager(servicesManager *services.ServicesManager, logger *log.Logger) *MiddlewareManager {
	return &MiddlewareManager{services: servicesManager, logger: logger}
}

type Middleware func(handler http.Handler) http.Handler

// пока не рабочая штука
func (mm *MiddlewareManager) Authenticator(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(entities.AccessTokenCookieName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			mm.logger.Log(err)
			return
		}

		token := cookie.Value
		claims, err := mm.services.Users.ValidateAccessToken(token)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			mm.logger.Log(err)
			return
		}

		ctx := context.WithValue(r.Context(), entities.AuthClaimsContextKey{}, claims)
		r = r.WithContext(ctx)

		handler.ServeHTTP(w, r)
	})
}

func (mm *MiddlewareManager) AuthenticationChecker(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		refresh := false
		isLoggedIn := true
		accessTokenCookie, err := r.Cookie(entities.AccessTokenCookieName)
		if err != nil {
			refresh = true
		}

		var claims *entities.Claims
		if !refresh {
			token := accessTokenCookie.Value
			claims, err = mm.services.Users.ValidateAccessToken(token)
			if err != nil {
				refresh = true
			}
		}

		if refresh {
			refreshTokenCookie, err := r.Cookie(entities.RefreshTokenCookieName)
			if err != nil {
				isLoggedIn = false
			} else {
				token, newClaims, err := mm.services.Users.RefreshAccessToken(refreshTokenCookie.Value)
				if err != nil {
					isLoggedIn = false
				} else {
					claims = newClaims

					http.SetCookie(w, &http.Cookie{
						Name:     entities.AccessTokenCookieName,
						Value:    token,
						HttpOnly: true,
						Secure:   false,
						SameSite: http.SameSiteStrictMode,
						Path:     "/",
					})
				}
			}
		}

		ctx := context.WithValue(context.WithValue(r.Context(), entities.AuthStatusContextKey{}, isLoggedIn), entities.AuthClaimsContextKey{}, claims)
		r = r.WithContext(ctx)

		handler.ServeHTTP(w, r)
	})
}

func (mm *MiddlewareManager) RegistrationChecker(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isInRegistration := true
		var claims *entities.Claims

		cookie, err := r.Cookie(entities.RegistrationCookieName)
		if err != nil {
			fmt.Println(err)
			isInRegistration = false
		}

		if isInRegistration {
			sessionId := cookie.Value
			_, err := mm.services.Users.GetSessionById(sessionId)
			if err != nil {
				fmt.Println(err)
				isInRegistration = false
			} else {
				claims = &entities.Claims{SessionId: sessionId}
			}
		}

		ctx := context.WithValue(context.WithValue(r.Context(), entities.RegistrationStatusContextKey{}, isInRegistration), entities.AuthClaimsContextKey{}, claims)
		r = r.WithContext(ctx)

		handler.ServeHTTP(w, r)
	})
}

func (mm *MiddlewareManager) Recoverer(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				mm.logger.Log(err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()

		handler.ServeHTTP(w, r)
	})
}

type responseWriterMock struct {
	http.ResponseWriter
	code int
}

func (r *responseWriterMock) WriteHeader(statusCode int) {
	r.code = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

func (mm *MiddlewareManager) Logger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wMock := &responseWriterMock{
			ResponseWriter: w,
			code:           http.StatusOK,
		}

		handler.ServeHTTP(wMock, r)

		mm.logger.Info().Fields(map[string]interface{}{
			"Host":     r.Host,
			"Method":   r.Method,
			"URL":      r.URL.Path,
			"Duration": time.Since(start).String(),
			"Status":   wMock.code,
		}).Msg("")
	})
}

// Иерархия идет слева направо, т.е. миддлвар слева оборачивает тот что справа
func (mm *MiddlewareManager) BuildChain(middlewares ...Middleware) Middleware {
	return func(handler http.Handler) http.Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			handler = middlewares[i](handler)
		}
		return handler
	}
}
