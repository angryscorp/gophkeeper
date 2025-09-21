package domain

import "github.com/google/uuid"

type BankCard struct {
	ID         uuid.UUID
	Owner      string
	Number     string
	CVV        string
	ExpireDate string
	Note       string
}
