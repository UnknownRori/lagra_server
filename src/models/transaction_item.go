package models

import (
	"fmt"
	"strings"

	"github.com/UnknownRori/lagra_server/src"

	"github.com/google/uuid"
)

type TransactionItem struct {
	Uuid        string `param:"uuid" query:"uuid" form:"uuid" json:"uuid" xml:"uuid"`
	Total       string `json:"total" xml:"total"`
	Transaction `json:"transaction"`
	Item        `json:"item"`
}

type TransactionDisplayItem struct {
	Uuid  string `param:"uuid" query:"uuid" form:"uuid" json:"uuid" xml:"uuid"`
	Total string `json:"total" xml:"total"`
	Item  `json:"item"`
}

func CreateTransactionItemsFromCarts(db *src.DB, transaction Transaction, userCarts []Cart) error {
	valString := []string{}
	valArgs := []interface{}{}
	for _, cart := range userCarts {
		uuid := uuid.New()
		valString = append(valString, "(?, ?, ?, ?)")

		valArgs = append(valArgs, uuid)
		valArgs = append(valArgs, cart.Total)
		valArgs = append(valArgs, transaction.Uuid)
		valArgs = append(valArgs, cart.Item.Uuid)
	}

	query := fmt.Sprintf("INSERT INTO transaction_item (uuid, total, transaction_id, item_id) VALUES %s", strings.Join(valString, ","))
	stmt, err := db.Prepare(query)
	defer stmt.Close()

	if err != nil {
		return err
	}

	_, err = stmt.Exec(valArgs...)

	if err != nil {
		return err
	}

	return err
}

func FetchDetailTransactions(db *src.DB, uuid string, user User) (DetailTransaction, error) {
	var items []TransactionDisplayItem
	var detail DetailTransaction
	transaction, err := FetchTransactionByUuid(db, uuid, user)
	detail.Transaction = transaction

	if err != nil {
		return detail, err
	}

	stmt, err := db.Prepare(`
		SELECT 
			transaction_item.uuid, transaction_item.total,
			items.uuid, items.name, items.price, categories.uuid, categories.name
		FROM transaction_item
		INNER JOIN
			items ON items.uuid = transaction_item.item_id
		INNER JOIN
			categories ON items.category_id = categories.uuid
		WHERE transaction_item.transaction_id = ?`,
	)
	if err != nil {
		return detail, err
	}

	query, err := stmt.Query(uuid)
	if err != nil {
		return detail, err
	}
	defer query.Close()

	for query.Next() {
		var item TransactionDisplayItem
		if err := query.Scan(&item.Uuid, &item.Total, &item.Item.Uuid, &item.Item.Name,
			&item.Item.Price, &item.Category.Uuid, &item.Category.Name); err != nil {
			return detail, err
		}

		items = append(items, item)

	}

	if err := query.Err(); err != nil {
		return detail, err
	}

	detail.TransactionDisplayItem = items
	return detail, err
}
