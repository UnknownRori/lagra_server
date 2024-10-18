package models

import (
	"github.com/UnknownRori/lagra_server/src"

	"github.com/google/uuid"
)

type Category struct {
	Uuid string `param:"uuid" query:"uuid" form:"uuid" json:"uuid" xml:"uuid"`
	Name string `json:"name" xml:"name"`
}

type FindCategoryByUuid struct {
	Uuid string `param:"uuid" query:"uuid" form:"uuid" json:"uuid" xml:"uuid" validate:"required"`
}

type NewCategory struct {
	Name string `json:"name" xml:"name" validate:"required,alphanum"`
}

func CreateCategory(db *src.DB, category NewCategory) (uuid.UUID, error) {
	uuid := uuid.New()
	stmt, err := db.Prepare("INSERT INTO categories (uuid, name) VALUES (?, ?)")
	defer stmt.Close()

	if err != nil {
		return uuid, err
	}

	_, err = stmt.Exec(uuid, category.Name)

	if err != nil {
		return uuid, err
	}

	return uuid, nil
}

func FetchCategories(db *src.DB) ([]Category, error) {
	var categories []Category
	stmt, err := db.Prepare("SELECT * FROM categories")
	if err != nil {
		return categories, err
	}

	query, err := stmt.Query()
	if err != nil {
		return categories, err
	}
	defer query.Close()

	for query.Next() {
		var category Category
		if err := query.Scan(&category.Uuid, &category.Name); err != nil {
			return categories, err
		}
		categories = append(categories, category)

	}

	if err := query.Err(); err != nil {
		return categories, err
	}

	return categories, err
}
