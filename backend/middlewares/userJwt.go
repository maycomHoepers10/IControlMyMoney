package middlewares

import (
	"fmt"

	"home_money/models"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func getUserIDFromJWT(c echo.Context, db *gorm.DB) (int, error) {

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	// Extrair o email do JWT
	email, ok := claims["email"].(string)
	if !ok {
		return 0, fmt.Errorf("email não encontrado no JWT")
	}

	// Consultar o banco de dados para obter o user_id com base no email
	var userDb models.User
	if err := db.Table("users").Where("email = ?", email).First(&userDb).Error; err != nil {
		return 0, err
	}

	return userDb.UserID, nil
}

func ExtractUserIDMiddleware(db *gorm.DB) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userID, err := getUserIDFromJWT(c, db)
			if err != nil {
				// Se ocorrer um erro ao obter o userID, defina como 0 ou outro valor padrão.
				userID = 0
			}

			// Adicionar userID ao contexto para que esteja disponível nos handlers subsequentes
			c.Set("userID", userID)

			return next(c)
		}
	}
}
