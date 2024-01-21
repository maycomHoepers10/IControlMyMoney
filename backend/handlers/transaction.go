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

type TransactionHandler struct {
	db       *gorm.DB
	validate *validator.Validate
}

func NewTransactionHandler(db *gorm.DB) *TransactionHandler {
	return &TransactionHandler{
		db:       db,
		validate: validator.New(),
	}
}

func (h *TransactionHandler) CreateTransaction(c echo.Context) error {

	transaction := new(models.Transaction)
	if err := c.Bind(transaction); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	// Validação do nome da transação usando o validator
	if err := h.validate.Struct(transaction); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	// Define o campo user_id no objeto Transaction
	transaction.UserID = c.Get("userID").(int)

	if err := h.db.Create(transaction).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, transaction)
}

func (h *TransactionHandler) UpdateTransaction(c echo.Context) error {
	userID := c.Get("userID").(int)

	// Extrair o ID da categoria da URL
	transactionIDStr := c.Param("id")
	transactionID, err := strconv.Atoi(transactionIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "ID de categoria inválido")
	}

	// Verificar se a transação pertence ao usuário atual
	var existingTransaction models.Transaction
	if err := h.db.Where("transaction_id = ? AND user_id = ?", transactionID, userID).First(&existingTransaction).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, "Transação não encontrada!")
		}
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	// Bind dos dados da transação da solicitação
	newTransaction := new(models.Transaction)
	if err := c.Bind(newTransaction); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	// Atualize os campos da transação existente com os dados da nova transação
	existingTransaction.AccountID = newTransaction.AccountID
	existingTransaction.Date = newTransaction.Date
	existingTransaction.TransactionType = newTransaction.TransactionType
	existingTransaction.CategoryID = newTransaction.CategoryID
	existingTransaction.Amount = newTransaction.Amount
	existingTransaction.Description = newTransaction.Description
	existingTransaction.Status = newTransaction.Status

	// Salvar a transação atualizada no banco de dados
	if err := h.db.Save(&existingTransaction).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, existingTransaction)
}

func (h *TransactionHandler) DeleteTransaction(c echo.Context) error {
	transactionID := c.Param("id")
	var transaction models.Transaction

	userID := c.Get("userID").(int)
	// Define o campo user_id no objeto Transação
	transaction.UserID = userID

	// Adicione uma cláusula WHERE para garantir que a transação seja excluída apenas se corresponder ao user_id
	if err := h.db.Where("transaction_id = ? AND user_id = ?", transactionID, userID).Delete(&transaction).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, "Transação excluída com sucesso!")
}

func (h *TransactionHandler) ListTransactions(c echo.Context) error {
	userID := c.Get("userID").(int)

	// Consulte o banco de dados para obter todas as transações do user_id especificado
	var transactions []models.Transaction
	if err := h.db.Where("user_id = ?", userID).Find(&transactions).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, transactions)
}

func (h *TransactionHandler) GetTransaction(c echo.Context) error {
	// Obtenha o ID da transação da URL
	transactionID := c.Param("id")

	// Consulte o banco de dados para obter a transação pelo ID
	var transaction models.Transaction

	userID := c.Get("userID").(int)

	// Define o campo user_id no objeto Category
	transaction.UserID = userID

	if err := h.db.Where("transaction_id = ? AND user_id = ?", transactionID, userID).First(&transaction).Error; err != nil {
		return c.JSON(http.StatusNotFound, "Transação não encontrada")
	}

	return c.JSON(http.StatusOK, transaction)
}

func (h *TransactionHandler) ApproveTransaction(c echo.Context) error {
	// Criar uma estrutura para representar os dados da solicitação
	var requestData struct {
		TransactionIds []int `json:"transactionIds"`
	}

	// Fazer o binding do corpo da solicitação para a estrutura
	if err := c.Bind(&requestData); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Erro ao ler o corpo da solicitação"})
	}

	// Obter o user_id do contexto
	userID := c.Get("userID").(int)

	// Atualizar o status das transações usando uma única consulta SQL
	if err := h.db.Model(&models.Transaction{}).
		Where("transaction_id IN (?) AND user_id = ?", requestData.TransactionIds, userID).
		Update("status", 2).
		Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Erro ao atualizar o status das transações"})
	}

	// Retorne uma resposta de sucesso
	return c.JSON(http.StatusOK, map[string]string{"message": "Status das transações atualizado com sucesso"})
}
