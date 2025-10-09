package domain

// UserDataKind defines the type of user data stored in the vault.
// It is used as an enum to categorize different record types.
type UserDataKind int64

const (
	// UserDataKindBankCard represents a bank card entry.
	UserDataKindBankCard UserDataKind = iota + 1
	// UserDataKindCredentials represents login/password credentials.
	UserDataKindCredentials
	// UserDataKindTextData represents plain text or notes.
	UserDataKindTextData
	// UserDataKindBinaryData represents arbitrary binary data or files.
	UserDataKindBinaryData
)

// Title returns a human-friendly title for the UserDataKind.
// It may include emoji to visually distinguish data types.
func (k UserDataKind) Title() string {
	switch k {
	case UserDataKindBankCard:
		return "💳 Bank Card"
	case UserDataKindCredentials:
		return "🔑 Credentials"
	case UserDataKindTextData:
		return "📝 Text Data"
	case UserDataKindBinaryData:
		return "📁 Binary Data"
	default:
		panic("unknown user data kind")
	}
}

// Description returns a short textual description for the UserDataKind.
func (k UserDataKind) Description() string {
	switch k {
	case UserDataKindBankCard:
		return "Credit/Debit card information"
	case UserDataKindCredentials:
		return "Username and password"
	case UserDataKindTextData:
		return "Notes and text information"
	case UserDataKindBinaryData:
		return "Files and binary information"
	default:
		panic("unknown user data kind")
	}
}
