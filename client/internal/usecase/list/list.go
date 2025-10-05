package list

import (
	"context"
	"encoding/json"
	"errors"
	"gophkeeper/client/internal/domain"
)

type List struct {
	repo Repository
}

func NewList(repo Repository) *List {
	return &List{
		repo: repo,
	}
}

func (l List) GetAllRecords() ([]domain.Record, error) {
	rows, err := l.repo.GetAll(context.Background())
	if err != nil {
		return nil, err
	}

	result := make([]domain.Record, len(rows))
	for i := range rows {
		res, err := l.mapRecord(rows[i])
		if err != nil {
			return nil, err
		}
		result[i] = *res
	}
	return result, nil
}

func (l List) mapRecord(record RawRecord) (*domain.Record, error) {
	switch record.Kind {
	case domain.UserDataKindBankCard:
		var bankCard domain.BankCard
		err := json.Unmarshal(record.Data, &bankCard)
		if err != nil {
			return nil, err
		}
		res := bankCard.ToRecord()
		return &res, nil

	case domain.UserDataKindCredentials:
		var credentials domain.Credentials
		err := json.Unmarshal(record.Data, &credentials)
		if err != nil {
			return nil, err
		}
		res := credentials.ToRecord()
		return &res, nil

	case domain.UserDataKindTextData:
		var textData domain.UserTextData
		err := json.Unmarshal(record.Data, &textData)
		if err != nil {
			return nil, err
		}
		res := textData.ToRecord()
		return &res, nil

	case domain.UserDataKindBinaryData:
		var binaryData domain.UserBinaryData
		err := json.Unmarshal(record.Data, &binaryData)
		if err != nil {
			return nil, err
		}
		res := binaryData.ToRecord()
		return &res, nil
	}

	return nil, errors.New("unknown data kind")
}
