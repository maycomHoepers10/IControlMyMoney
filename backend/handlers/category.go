package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"home_money/models"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type CategoryHandler struct {
	db       *gorm.DB
	validate *validator.Validate
}

func NewCategoryHandler(db *gorm.DB) *CategoryHandler {
	return &CategoryHandler{
		db:       db,
		validate: validator.New(),
	}
}

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

func (h *CategoryHandler) CreateCategory(c echo.Context) error {
	category := new(models.Category)
	if err := c.Bind(category); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	// Validação do nome da categoria usando o validator
	if err := h.validate.Struct(category); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	// Extrair o user_id do banco de dados usando o email do token JWT
	userID, err := getUserIDFromJWT(c, h.db)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	// Define o campo user_id no objeto Category
	category.UserID = userID

	// Verificar se a categoria com o mesmo nome já existe
	if h.categoryExists(category.CategoryName) {
		return c.JSON(http.StatusConflict, "Categoria com o mesmo nome já existe")
	}

	if err := h.db.Create(category).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, category)
}

func (h *CategoryHandler) UpdateCategory(c echo.Context) error {
	// Extrair o user_id do banco de dados usando o email do token JWT
	userID, err := getUserIDFromJWT(c, h.db)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	// Extrair o ID da categoria da URL
	categoryIDStr := c.Param("id")
	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "ID de categoria inválido")
	}

	// Verificar se a categoria pertence ao usuário atual
	var existingCategory models.Category
	if err := h.db.Where("category_id = ? AND user_id = ?", categoryID, userID).First(&existingCategory).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, "Categoria não encontrada ou não pertence ao usuário")
		}
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	// Bind dos dados da categoria da solicitação
	newCategory := new(models.Category)
	if err := c.Bind(newCategory); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	// Atualize os campos da categoria existente com os dados da nova categoria
	existingCategory.CategoryName = newCategory.CategoryName

	// Validação do nome da categoria usando o validator
	if err := h.validate.Struct(existingCategory); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	// Verificar se a categoria com o mesmo nome já existe, excluindo a categoria atual
	if h.categoryExistsWithDifferentID(existingCategory.CategoryName, existingCategory.CategoryID) {
		return c.JSON(http.StatusConflict, "Categoria com o mesmo nome já existe")
	}

	// Salvar a categoria atualizada no banco de dados
	if err := h.db.Save(&existingCategory).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, existingCategory)
}

func (h *CategoryHandler) DeleteCategory(c echo.Context) error {
	categoryID := c.Param("id")
	var category models.Category

	// Extrair o user_id do banco de dados usando o email do token JWT
	userID, err := getUserIDFromJWT(c, h.db)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	// Define o campo user_id no objeto Category
	category.UserID = userID

	// Adicione uma cláusula WHERE para garantir que a categoria seja excluída apenas se corresponder ao user_id
	if err := h.db.Where("category_id = ? AND user_id = ?", categoryID, userID).Delete(&category).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, "Categoria excluída com sucesso")
}

func (h *CategoryHandler) ListCategories(c echo.Context) error {
	// Extrair o user_id do banco de dados usando o email do token JWT
	userID, err := getUserIDFromJWT(c, h.db)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	// Consulte o banco de dados para obter todas as categorias do user_id especificado
	var categories []models.Category
	if err := h.db.Where("user_id = ?", userID).Find(&categories).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, categories)
}

func (h *CategoryHandler) GetCategory(c echo.Context) error {
	// Obtenha o ID da categoria da URL
	categoryID := c.Param("id")

	// Consulte o banco de dados para obter a categoria pelo ID
	var category models.Category

	// Extrair o user_id do banco de dados usando o email do token JWT
	userID, err := getUserIDFromJWT(c, h.db)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	// Define o campo user_id no objeto Category
	category.UserID = userID

	if err := h.db.Where("category_id = ? AND user_id = ?", categoryID, userID).First(&category).Error; err != nil {
		return c.JSON(http.StatusNotFound, "Categoria não encontrada")
	}

	return c.JSON(http.StatusOK, category)
}

func (h *CategoryHandler) categoryExists(categoryName string) bool {
	var count int64
	h.db.Model(&models.Category{}).Where("category_name = ?", categoryName).Count(&count)
	return count > 0
}

func (h *CategoryHandler) categoryExistsWithDifferentID(categoryName string, categoryID int) bool {
	var count int64
	h.db.Model(&models.Category{}).Where("category_name = ? AND category_id != ?", categoryName, categoryID).Count(&count)
	return count > 0
}

func (h *CategoryHandler) GetCategoryIDByDescription(description string, userID int) (int, error) {
	var categoryID int

	if err := h.db.
		Table("categories").
		Joins("INNER JOIN descriptions_by_categories ON categories.category_id = descriptions_by_categories.category_id").
		Where("descriptions_by_categories.description = ? AND categories.user_id = ?", description, userID).
		Pluck("categories.category_id", &categoryID).
		Error; err != nil {
		return 0, err
	}

	return categoryID, nil
}

func (h *CategoryHandler) GetCategoryByName(categoryName string, userID int) int {

	// Consulte o banco de dados para obter a categoria pelo ID
	var category models.Category

	// Define o campo user_id no objeto Category
	category.UserID = userID

	if err := h.db.Where("category_name= ? AND user_id = ?", categoryName, userID).First(&category).Error; err != nil {
		fmt.Println("Categoria não encontrada")
	}

	return category.CategoryID
}
