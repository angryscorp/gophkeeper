package list

import (
	"context"
	"encoding/json"
	"errors"
	"gophkeeper/client/internal/domain"
)

type List struct {
	repo      Repository
	decryptor func([]byte) ([]byte, error)
}

func NewList(
	repo Repository,
	decryptor func([]byte) ([]byte, error),
) *List {
	return &List{
		repo:      repo,
		decryptor: decryptor,
	}
}

func (l List) GetAllRecords() ([]domain.Record, error) {
	rows, err := l.repo.GetAll(context.Background())
	if err != nil {
		return nil, err
	}

	result := make([]domain.Record, len(rows))
	for i := range rows {
		res, err := l.decryptRecord(rows[i])
		if err != nil {
			return nil, err
		}
		result[i] = *res
	}
	return result, nil
}

func (l List) decryptRecord(record RawRecord) (*domain.Record, error) {
	data, err := l.decryptor(record.Payload)
	if err != nil {
		return nil, err
	}

	switch record.Kind {
	case domain.UserDataKindBankCard:
		var bankCard domain.BankCard
		err := json.Unmarshal(data, &bankCard)
		if err != nil {
			return nil, err
		}
		res := bankCard.ToRecord()
		return &res, nil

	case domain.UserDataKindCredentials:
		var credentials domain.Credentials
		err := json.Unmarshal(data, &credentials)
		if err != nil {
			return nil, err
		}
		res := credentials.ToRecord()
		return &res, nil

	case domain.UserDataKindTextData:
		var textData domain.UserTextData
		err := json.Unmarshal(data, &textData)
		if err != nil {
			return nil, err
		}
		res := textData.ToRecord()
		return &res, nil

	case domain.UserDataKindBinaryData:
		var binaryData domain.UserBinaryData
		err := json.Unmarshal(data, &binaryData)
		if err != nil {
			return nil, err
		}
		res := binaryData.ToRecord()
		return &res, nil
	}

	return nil, errors.New("unknown data kind")
}
