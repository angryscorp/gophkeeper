package domain

// UserInfoSaver defines the interface for persisting
// different types of user-provided data in the vault.
// It abstracts saving logic for multiple domain entities.
type UserInfoSaver interface {
	// SaveBankCard stores a bank card record.
	SaveBankCard(bankCard BankCard) error
	// SaveCredentials stores a credentials record.
	SaveCredentials(credentials Credentials) error
	// SaveUserBinaryData stores a binary data record.
	SaveUserBinaryData(userBinaryData UserBinaryData) error
	// SaveUserTextData stores a text data record.
	SaveUserTextData(userTextData UserTextData) error
}
