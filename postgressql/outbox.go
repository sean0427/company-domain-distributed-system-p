// TODO might be move to util repository
package postgressql

import (
	"context"
	"encoding/json"

	"github.com/sean0427/company-domain-distributed-system-p/model"
	"gorm.io/gorm"
)

type outbox struct {
	Query string `json:"query"`
	//TODO
}

func getOutboxQuery(data interface{}) (string, error) {
	b, err := json.Marshal(data)
	return string(b), err
}

func TransactionWithOutboxMsg(ctx context.Context, db *gorm.DB, data *model.Company, queryFunc func(tx *gorm.DB) error) error {
	msg, err := getOutboxQuery(data)
	if err != nil {
		return err
	}

	outbox := outbox{Query: msg}

	err = db.WithContext(ctx).Transaction(func(_tx *gorm.DB) error {
		err := queryFunc(_tx)
		if err != nil {
			return err
		}

		ret := _tx.Model(&outbox).Create(&outbox)
		if ret.Error != nil {
			return ret.Error
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil

}
