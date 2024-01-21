package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"home_money/models"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type AccountHandler struct {
	db       *gorm.DB
	validate *validator.Validate
}

func NewAccountHandler(db *gorm.DB) *AccountHandler {
	return &AccountHandler{
		db:       db,
		validate: validator.New(),
	}
}

func (h *AccountHandler) CreateAccount(c echo.Context) error {

	account := new(models.Financialaccount)
	if err := c.Bind(account); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	// Validação do nome da conta usando o validator
	if err := h.validate.Struct(account); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	// Define o campo user_id no objeto Account
	account.UserID = c.Get("userID").(int)

	if err := h.db.Create(account).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, account)
}

func (h *AccountHandler) UpdateAccount(c echo.Context) error {
	userID := c.Get("userID").(int)

	// Extrair o ID da conta da URL
	accountIDStr := c.Param("id")
	accountID, err := strconv.Atoi(accountIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "ID de conta inválido")
	}

	// Verificar se a conta pertence ao usuário atual
	var existingAccount models.Financialaccount
	if err := h.db.Where("account_id = ? AND user_id = ?", accountID, userID).First(&existingAccount).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, "Conta não encontrada!")
		}
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	// Bind dos dados da conta da solicitação
	if err := c.Bind(&existingAccount); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	// Salvar a conta atualizada no banco de dados
	if err := h.db.Save(&existingAccount).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, existingAccount)
}

func (h *AccountHandler) DeleteAccount(c echo.Context) error {
	accountID := c.Param("id")
	var account models.Financialaccount

	userID := c.Get("userID").(int)
	// Define o campo user_id no objeto Transação
	account.UserID = userID

	// Adicione uma cláusula WHERE para garantir que a conta seja excluída apenas se corresponder ao user_id
	if err := h.db.Where("account_id = ? AND user_id = ?", accountID, userID).Delete(&account).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, "Transação excluída com sucesso!")
}

func (h *AccountHandler) ListAccounts(c echo.Context) error {
	userID := c.Get("userID").(int)

	// Consulte o banco de dados para obter todas as transações do user_id especificado
	var accounts []models.Financialaccount
	if err := h.db.Where("user_id = ?", userID).Find(&accounts).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, accounts)
}

func (h *AccountHandler) GetAccount(c echo.Context) error {
	// Obtenha o ID da conta da URL
	accountID := c.Param("id")

	// Consulte o banco de dados para obter a conta pelo ID
	var account models.Financialaccount

	userID := c.Get("userID").(int)

	// Define o campo user_id no objeto conta
	account.UserID = userID

	if err := h.db.Where("account_id = ? AND user_id = ?", accountID, userID).First(&account).Error; err != nil {
		return c.JSON(http.StatusNotFound, "Transação não encontrada")
	}

	return c.JSON(http.StatusOK, account)
}
