package models

import (
	"github.com/UnknownRori/lagra_server/src"

	"github.com/google/uuid"
)

type Cart struct {
	Uuid  string `param:"uuid" query:"uuid" form:"uuid" json:"uuid" xml:"uuid"`
	Total int32  `json:"name" xml:"name"`
	Item
}

type FindCartByUuid struct {
	Uuid string `param:"uuid" query:"uuid" form:"uuid" json:"uuid" xml:"uuid" validate:"required"`
}

type NewCart struct {
	Name   string `json:"name" xml:"name" validate:"required,alphanum"`
	Total  int32  `json:"total" xml:"total" validate:"required,numeric"`
	ItemId string `json:"itemId" xml:"itemId" validate:"required"`
}

func CreateCart(db *src.DB, cart NewCart) (uuid.UUID, error) {
	uuid := uuid.New()
	stmt, err := db.Prepare("INSERT INTO carts (uuid, total, item_id) VALUES (?, ?, ?)")
	defer stmt.Close()

	if err != nil {
		return uuid, err
	}

	_, err = stmt.Exec(uuid, cart.Total, cart.ItemId)

	if err != nil {
		return uuid, err
	}

	return uuid, nil
}

func FetchCarts(db *src.DB) ([]Item, error) {
	var items []Item
	stmt, err := db.Prepare("SELECT * FROM items")
	if err != nil {
		return items, err
	}

	query, err := stmt.Query()
	if err != nil {
		return items, err
	}
	defer query.Close()

	for query.Next() {
		var item Item
		if err := query.Scan(&item.Uuid, &item.Name); err != nil {
			return items, err
		}
		items = append(items, item)

	}

	if err := query.Err(); err != nil {
		return items, err
	}

	return items, err
}
