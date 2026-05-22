package models

import "time"


type TransactionType string
// means: "create a new type called TransactionType which is based on the string type. 
// This allows us to define specific values for TransactionType, such as 'expense' and 'income', 
// while still treating it as a string in our code."

const (
	Expense TransactionType = "expense"
	Income  TransactionType = "income"
)

type Transaction struct {
	ID          string          `json:"id"`
	UserID      string          `json:"user_id"`
	Type        TransactionType `json:"type"`
	Amount      float64         `json:"amount"`
	Category    string          `json:"category"`
	Description string          `json:"description"`
	Date        time.Time       `json:"date"`
	CreatedAt   time.Time       `json:"created_at"`
}