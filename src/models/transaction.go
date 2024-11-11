package models

import (
	"github.com/UnknownRori/lagra_server/src"

	"github.com/google/uuid"
)

type Transaction struct {
	Uuid        string `param:"uuid" query:"uuid" form:"uuid" json:"uuid" xml:"uuid"`
	Pay         string `json:"pay" xml:"pay"`
	PayType     string `json:"payType" xml:"payType"`
	DisplayUser `json:"displayuser"`
}

type DetailTransaction struct {
	Transaction            
	TransactionDisplayItem []TransactionDisplayItem `json:"transactionItem"`
}

type FindTransactionByUuid struct {
	Uuid string `param:"uuid" query:"uuid" form:"uuid" json:"uuid" xml:"uuid" validate:"required"`
}

type NewTransaction struct {
	PayType string `json:"payType" xml:"payType" validate:"required,alphanum"`
}

func CreateTransaction(db *src.DB, item NewTransaction, userCarts []Cart, user User) (uuid.UUID, error) {
	uuid := uuid.New()
	stmt, err := db.Prepare("INSERT INTO transactions (uuid, pay, pay_type, consumer_id) VALUES (?, ?, ?, ?)")
	defer stmt.Close()

	if err != nil {
		return uuid, err
	}

	total := 0
	for _, cart := range userCarts {
		total += int(cart.Price)
	}

	_, err = stmt.Exec(uuid, total, item.PayType, user.Uuid)

	if err != nil {
		return uuid, err
	}

	return uuid, nil
}

func FetchTransactionByUuid(db *src.DB, uuid string, user User) (Transaction, error) {
	var item Transaction
	stmt, err := db.Prepare(`
		SELECT transactions.uuid, transactions.pay, transactions.pay_type
		FROM transactions 
		WHERE transactions.uuid = ? AND transactions.consumer_id = ? LIMIT 1`,
	)
	if err != nil {
		return item, err
	}

	query := stmt.QueryRow(uuid, user.Uuid)

	if err := query.Scan(&item.Uuid, &item.Pay, &item.PayType); err != nil {
		return item, err
	}

	item.DisplayUser.Uuid = user.Uuid
	item.DisplayUser.Username = user.Username

	return item, err
}

func FetchTransactions(db *src.DB, user User) ([]Transaction, error) {
	var items []Transaction
	stmt, err := db.Prepare(`
		SELECT transactions.uuid, transactions.pay, transactions.pay_type
		FROM transactions 
		WHERE transactions.consumer_id = ?`,
	)
	if err != nil {
		return items, err
	}

	query, err := stmt.Query(user.Uuid)
	if err != nil {
		return items, err
	}
	defer query.Close()

	for query.Next() {
		var item Transaction
		item.DisplayUser.Uuid = user.Uuid
		item.DisplayUser.Username = user.Username
		if err := query.Scan(&item.Uuid, &item.Pay, &item.PayType); err != nil {
			return items, err
		}

		items = append(items, item)

	}

	if err := query.Err(); err != nil {
		return items, err
	}

	return items, err
}
