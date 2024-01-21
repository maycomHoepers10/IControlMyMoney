package models

import "time"

type User struct {
	UserID     int    `gorm:"column:user_id;primaryKey" json:"id"`
	Name       string `gorm:"column:name" json:"name" validate:"required"`
	Email      string `gorm:"column:email;unique" json:"email" validate:"required,email"`
	Password   string `gorm:"column:password" json:"password" validate:"required"`
	Permission int    `gorm:"column:permission" json:"permission" validate:"required"`
}

type Credentials struct {
	Email    string `json:"email" form:"email" validate:"required"`
	Password string `json:"password" form:"password" validate:"required"`
}

type Category struct {
	CategoryID   int    `gorm:"column:category_id;primaryKey" json:"id"`
	CategoryName string `gorm:"column:category_name;not null" json:"categoryName" validate:"required"`
	UserID       int    `gorm:"column:user_id" json:"-"`
}

type Transaction struct {
	TransactionID   int       `gorm:"primaryKey" json:"transaction_id"`
	UserID          int       `gorm:"column:user_id" json:"user_id"`
	AccountID       int       `gorm:"column:account_id" json:"account_id"`
	Date            time.Time `gorm:"column:date;type:date" json:"date"`
	TransactionType string    `gorm:"column:transaction_type;type:varchar(10)" json:"transaction_type"`
	CategoryID      int       `gorm:"column:category_id" json:"category_id"`
	Amount          float64   `gorm:"column:amount;type:numeric(10,2)" json:"amount"`
	Description     string    `gorm:"column:description" json:"description"`
	Status          int       `gorm:"column:status" json:"status"` /* 1 - Análise, 2 - Aprovado, 3 - Rejeitado, 4 - Insersão manual */
}

type Budget struct {
	BudgetID     int     `gorm:"primaryKey" json:"budget_id"`
	UserID       int     `gorm:"column:user_id" json:"user_id"`
	CategoryID   int     `gorm:"column:category_id" json:"category_id"`
	BudgetAmount float64 `gorm:"column:budget_amount;type:numeric(10,2)" json:"budget_amount"`
	Year         int     `gorm:"column:year" json:"year"`
	Month        int     `gorm:"column:month" json:"month"`
}

type Financialaccount struct {
	AccountID   int     `gorm:"primaryKey" json:"account_id"`
	UserID      int     `gorm:"column:user_id" json:"user_id"`
	AccountName string  `gorm:"column:account_name" json:"account_name"`
	Balance     float64 `gorm:"column:balance;type:numeric(10,2)" json:"balance"`
}

type DescriptionsByCategories struct {
	CategoryID  int    `gorm:"column:category_id" json:"category_id"`
	Description string `gorm:"column:description" json:"description"`
}
