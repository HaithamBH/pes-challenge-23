package server

import (
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

func Authentificate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		if err := next(ctx); err != nil {
			ctx.Error(err)
		}

		token := ctx.Request().Header.Values("token")
		if token == nil {
			ctx.String(http.StatusUnauthorized, "unauthorized")
		}

		claims, err := PasreToken(strings.Join(token, ""))
		if err != nil {
			ctx.String(http.StatusUnauthorized, "unauthorized")
		}

		err = claims.StandardClaims.Valid()
		if err != nil {
			ctx.String(http.StatusUnauthorized, "unauthorized")
		}

		isValid := claims.StandardClaims.VerifyExpiresAt(time.Now().Unix(), true)
		if !isValid {
			ctx.String(http.StatusUnauthorized, "unauthorized")
		}

		ctx.Request().Header.Set("user", claims.UserID)
		return next(ctx)
	}
}
