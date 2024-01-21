package middlewares

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
)

// CustomClaims define a estrutura dos claims do token JWT
type CustomClaims struct {
	UserID int    `json:"user_id"`
	Role   string `json:"role"`
	jwt.StandardClaims
}

// LoadJWTConfig carrega a configuração do JWT a partir do arquivo de configuração
func LoadJWTConfig() (*middleware.JWTConfig, error) {
	jwtSecret := viper.GetString("jwt.secret")
	return &middleware.JWTConfig{
		SigningKey: []byte(jwtSecret),
	}, nil
}

// JWTMiddleware retorna o middleware JWT para autenticação
func JWTMiddleware(jwtConfig *middleware.JWTConfig) echo.MiddlewareFunc {
	return middleware.JWTWithConfig(*jwtConfig)
}
