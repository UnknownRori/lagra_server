package models

import (
	"github.com/UnknownRori/lagra_server/src"

	"github.com/google/uuid"
)

type Item struct {
	Uuid     string `param:"uuid" query:"uuid" form:"uuid" json:"uuid" xml:"uuid"`
	Name     string `json:"name" xml:"name"`
	Price    int32  `json:"price" xml:"price"`
	Category `json:"category"`
}

type FindItemByUuid struct {
	Uuid string `param:"uuid" query:"uuid" form:"uuid" json:"uuid" xml:"uuid" validate:"required"`
}

type NewItem struct {
	Name       string `json:"name" xml:"name" validate:"required,alphanum"`
	Price      int32  `json:"price" xml:"price" validate:"required,numeric"`
	CategoryId string `json:"categoryId" xml:"categoryId" validate:"required"`
}

func CreateItem(db *src.DB, item NewItem) (uuid.UUID, error) {
	uuid := uuid.New()
	stmt, err := db.Prepare("INSERT INTO items (uuid, name, price, category_id) VALUES (?, ?, ?, ?)")
	defer stmt.Close()

	if err != nil {
		return uuid, err
	}

	_, err = stmt.Exec(uuid, item.Name, item.Price, item.CategoryId)

	if err != nil {
		return uuid, err
	}

	return uuid, nil
}

func FetchItemByUuid(db *src.DB, uuid string) ([]Item, error) {
	var items []Item
	stmt, err := db.Prepare(`
		SELECT items.uuid, items.name, items.price, categories.uuid as categories_uuid, categories.name 
		FROM items 
		INNER JOIN categories ON items.category_id = categories.uuid WHERE items.uuid = ? LIMIT 1`,
	)
	if err != nil {
		return items, err
	}

	query := stmt.QueryRow(uuid)

	var item Item
	var category Category
	if err := query.Scan(&item.Uuid, &item.Name, &item.Price, &category.Uuid, &category.Name); err != nil {
		return items, err
	}

	item.Category = category
	items = append(items, item)

	if err := query.Err(); err != nil {
		return items, err
	}

	return items, err
}

func FetchItems(db *src.DB) ([]Item, error) {
	var items []Item
	stmt, err := db.Prepare(`
		SELECT items.uuid, items.name, items.price, categories.uuid as categories_uuid, categories.name 
		FROM items 
		INNER JOIN categories ON items.category_id = categories.uuid`,
	)
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
		if err := query.Scan(&item.Uuid, &item.Name, &item.Price, &item.Category.Uuid, &item.Category.Name); err != nil {
			return items, err
		}

		items = append(items, item)

	}

	if err := query.Err(); err != nil {
		return items, err
	}

	return items, err
}
