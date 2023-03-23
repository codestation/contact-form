// Copyright 2022 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package api

import (
	"github.com/golang-jwt/jwt/v4"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

type JwtCustomClaims struct {
	jwt.RegisteredClaims
	UserID string `json:"user_id"`
}

func (api *API) UseJWT(g *echo.Group) {
	jwtConfig := echojwt.Config{
		SigningKey: api.app.Config().GeneralSettings.JwtSecret,
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return &JwtCustomClaims{}
		},
	}

	g.Use(echojwt.WithConfig(jwtConfig))
}

func (api *API) GetUserID(c echo.Context) string {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(*JwtCustomClaims)
	return claims.UserID
}
