package domain

type UserDataKind int64

const (
	UserDataKindBankCard UserDataKind = iota + 1
	UserDataKindCredentials
	UserDataKindTextData
	UserDataKindBinaryData
)

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
