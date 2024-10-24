package models

import (
	"github.com/UnknownRori/lagra_server/src"

	"github.com/google/uuid"
)

type Item struct {
	Uuid     string `param:"uuid" query:"uuid" form:"uuid" json:"uuid" xml:"uuid"`
	Name     string `json:"name" xml:"name"`
	Price    int32  `json:"price" xml:"price"`
	ImgUrl   string `json:"img" xml:"img"`
	Category `json:"category"`
}

type FindItemByUuid struct {
	Uuid string `param:"uuid" query:"uuid" form:"uuid" json:"uuid" xml:"uuid" validate:"required"`
}

type NewItem struct {
	Name       string `form:"name" validate:"required,alphanum"`
	Price      int32  `form:"price" validate:"required,numeric"`
	CategoryId string `form:"categoryId" xml:"categoryId" validate:"required"`
	ImgUrl     string
}

func CreateItem(db *src.DB, item NewItem) (uuid.UUID, error) {
	uuid := uuid.New()
	stmt, err := db.Prepare("INSERT INTO items (uuid, name, price, img_url, category_id) VALUES (?, ?, ?, ?, ?)")
	defer stmt.Close()

	if err != nil {
		return uuid, err
	}

	_, err = stmt.Exec(uuid, item.Name, item.Price, item.ImgUrl, item.CategoryId)

	if err != nil {
		return uuid, err
	}

	return uuid, nil
}

func FetchItemByUuid(db *src.DB, uuid string) (Item, error) {
	var item Item
	stmt, err := db.Prepare(`
		SELECT items.uuid, items.name, items.price, items.img_url, categories.uuid as categories_uuid, categories.name 
		FROM items 
		INNER JOIN categories ON items.category_id = categories.uuid WHERE items.uuid = ? LIMIT 1`,
	)
	if err != nil {
		return item, err
	}

	query := stmt.QueryRow(uuid)

	if err := query.Scan(&item.Uuid, &item.Name, &item.Price, &item.ImgUrl, &item.Category.Uuid, &item.Category.Name); err != nil {
		return item, err
	}

	return item, err
}

func FetchItems(db *src.DB) ([]Item, error) {
	var items []Item
	stmt, err := db.Prepare(`
		SELECT items.uuid, items.name, items.price, items.img_url, categories.uuid as categories_uuid, categories.name 
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
		if err := query.Scan(&item.Uuid, &item.Name, &item.Price, &item.ImgUrl, &item.Category.Uuid, &item.Category.Name); err != nil {
			return items, err
		}

		items = append(items, item)

	}

	if err := query.Err(); err != nil {
		return items, err
	}

	return items, err
}
