package models

import (
	"github.com/UnknownRori/lagra_server/src"

	"github.com/google/uuid"
)

type Cart struct {
	Uuid        string `param:"uuid" query:"uuid" form:"uuid" json:"uuid" xml:"uuid"`
	Total       int32  `json:"total" xml:"total"`
	Item        `json:"item" xml:"item"`
	DisplayUser `json:"user" xml:"user"`
}

type FindCartByUuid struct {
	Uuid string `param:"uuid" query:"uuid" form:"uuid" json:"uuid" xml:"uuid" validate:"required"`
}

type NewCart struct {
	Total  int32  `json:"total" xml:"total" validate:"required,numeric"`
	ItemId string `json:"itemId" xml:"itemId" validate:"required"`
}

func CreateCart(db *src.DB, cart NewCart, user User) (uuid.UUID, error) {
	uuid := uuid.New()
	stmt, err := db.Prepare("insert into carts (uuid, total, item_id, user_id) values (?, ?, ?, ?)")
	defer stmt.Close()

	if err != nil {
		return uuid, err
	}

	_, err = stmt.Exec(uuid, cart.Total, cart.ItemId, user.Uuid)

	if err != nil {
		return uuid, err
	}

	return uuid, nil
}

func CleanCarts(db *src.DB, user User) error {
	stmt, err := db.Prepare("DELETE FROM carts WHERE carts.user_id = ?")
	defer stmt.Close()

	if err != nil {
		return err
	}

	_, err = stmt.Exec(user.Uuid)

	return err
}

func FetchCartsByUuid(db *src.DB, uuid string, user User) (Cart, error) {
	var carts Cart
	stmt, err := db.Prepare(`
		SELECT 
			carts.uuid, carts.total, 
			items.uuid as items_uuid, items.name, items.price, 
			categories.uuid, categories.name
		FROM carts
		INNER JOIN items ON carts.item_id = items.uuid
		INNER JOIN categories ON items.category_id = categories.uuid
		WHERE carts.uuid = ? AND carts.user_id = ?
		`)
	if err != nil {
		return carts, err
	}

	query := stmt.QueryRow(uuid, user.Uuid)

	if err := query.Scan(&carts.Uuid, &carts.Total, &carts.Item.Uuid, &carts.Item.Name, &carts.Item.Price, &carts.Item.Category.Uuid, &carts.Item.Category.Name); err != nil {
		return carts, err
	}
	carts.DisplayUser.Uuid = user.Uuid
	carts.DisplayUser.Username = user.Username
	if err := query.Err(); err != nil {
		return carts, err
	}

	return carts, err
}

func FetchCarts(db *src.DB, user User) ([]Cart, error) {
	var carts []Cart
	stmt, err := db.Prepare(`
		SELECT 
			carts.uuid, carts.total, 
			items.uuid as items_uuid, items.name, items.price, 
			categories.uuid, categories.name
		FROM carts
		INNER JOIN items ON carts.item_id = items.uuid
		INNER JOIN categories ON items.category_id = categories.uuid
		WHERE carts.user_id = ?
		`)
	if err != nil {
		return carts, err
	}

	query, err := stmt.Query(user.Uuid)
	if err != nil {
		return carts, err
	}
	defer query.Close()

	for query.Next() {
		var cart Cart

		cart.DisplayUser.Uuid = user.Uuid
		cart.DisplayUser.Username = user.Username
		if err := query.Scan(&cart.Uuid, &cart.Total, &cart.Item.Uuid, &cart.Item.Name, &cart.Item.Price, &cart.Category.Uuid, &cart.Category.Name); err != nil {
			return carts, err
		}
		carts = append(carts, cart)

	}

	if err := query.Err(); err != nil {
		return carts, err
	}

	return carts, err
}
