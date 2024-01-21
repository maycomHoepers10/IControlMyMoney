package handlers

import (
	"net/http"
	"time"

	"github.com/go-playground/locales/pt_BR"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	models "home_money/models"
)

var (
	customValidator *validator.Validate
	translator      ut.Translator
)

type UserHandler struct {
	db        *gorm.DB
	jwtSecret []byte
}

func NewUserHandler(db *gorm.DB, jwtSecret []byte) *UserHandler {
	customValidator = NewValidator()

	pt := pt_BR.New()
	uni := ut.New(pt, pt)

	translator, _ = uni.GetTranslator("pt_BR")

	validate := validator.New()
	validate.RegisterTranslation("required", translator, func(ut ut.Translator) error {
		return ut.Add("required", "{0} deve ser preenchido", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required", fe.Field())
		return t
	})

	customValidator = validate

	return &UserHandler{
		db:        db,
		jwtSecret: jwtSecret,
	}
}

func NewValidator() *validator.Validate {
	validate := validator.New()
	return validate
}

func (h *UserHandler) RegisterUser(c echo.Context) error {
	u := new(models.User)
	if err := c.Bind(u); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	// Valide a estrutura de dados usando o validador personalizado
	if err := customValidator.Struct(u); err != nil {
		// Erros de validação ocorreram
		var validationErrors []string
		for _, e := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, e.Translate(translator))
		}
		return c.JSON(http.StatusBadRequest, validationErrors)
	}

	// Verificar se o usuário já existe
	exists, err := h.userExists(u.Name)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	if exists {
		return c.JSON(http.StatusConflict, "Usuário já cadastrado")
	}

	// Hash da senha
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Erro ao gerar o hash da senha")
	}

	// Salvar o usuário no banco de dados usando GORM
	u.Password = string(hashedPassword) // Defina a senha hash no objeto User
	if err := h.db.Create(u).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, u)
}

func (h *UserHandler) userExists(name string) (bool, error) {
	var count int64
	result := h.db.Model(&models.User{}).Where("name = ?", name).Count(&count)
	if result.Error != nil {
		return false, result.Error
	}

	return count > 0, nil
}

func (h *UserHandler) Login(c echo.Context) error {
	credentials := new(models.Credentials)
	if err := c.Bind(credentials); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	// Verificar se o nome de usuário existe na tabela de usuários
	user, err := h.getUserByEmail(credentials.Email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	if user == nil {
		return c.JSON(http.StatusUnauthorized, "Credenciais inválidas")
	}

	// Comparar a senha fornecida com a senha armazenada
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password))
	if err != nil {
		return c.JSON(http.StatusUnauthorized, "Credenciais inválidas")
	}

	// Gerar o token JWT
	token, err := h.generateJWT(user.Name, user.Email, user.Permission)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	// Retornar o token JWT como resposta
	return c.JSON(http.StatusOK, echo.Map{
		"token": token,
	})
}

func (h *UserHandler) getUserByEmail(email string) (*models.User, error) {
	var user models.User
	result := h.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}

	return &user, nil
}

func (h *UserHandler) ValidateToken(c echo.Context) error {
	// Estrutura para receber o JSON que contém o token
	type TokenRequest struct {
		Token string `json:"token"`
	}

	// Bind do corpo da requisição para a estrutura TokenRequest
	var tokenRequest TokenRequest
	if err := c.Bind(&tokenRequest); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"isValid": false, "message": "Erro ao processar a requisição"})
	}

	// Obter o token da estrutura deserializada
	tokenString := tokenRequest.Token

	if tokenString == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"isValid": false, "message": "Token não fornecido"})
	}

	// Obter a chave secreta do arquivo de configuração usando viper
	key := []byte(h.jwtSecret)

	// Função de validação personalizada
	validationFunc := func(token *jwt.Token) (interface{}, error) {
		return key, nil
	}

	// Parse do token
	token, err := jwt.Parse(tokenString, validationFunc)

	if err != nil || !token.Valid {
		return c.JSON(http.StatusOK, map[string]interface{}{"isValid": false, "message": "Token inválido"})
	}

	// Token é válido
	return c.JSON(http.StatusOK, map[string]interface{}{"isValid": true, "message": "Token válido"})
}

func (h *UserHandler) generateJWT(name string, email string, permission int) (string, error) {
	// Definir as claims (dados) do token
	claims := jwt.MapClaims{
		"name":       name,
		"email":      email,                                 //fazer uma criptografia aqui
		"permission": permission,                            // Use 'permission' em vez de 'role'
		"exp":        time.Now().Add(time.Hour * 24).Unix(), // Expiração do token após 24 horas
	}

	// Criar o token JWT com as claims e a assinatura usando a chave secreta
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(h.jwtSecret)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
