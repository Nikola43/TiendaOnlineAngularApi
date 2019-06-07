package main

import (
	"database/sql"
	"fmt"
)

type Article struct {
	ID           uint   `json:"id"`
	Name         string   `json:"name"`
	Category     string `json:"category"`
	UnitPrice    float64 `json:"unit_price"`
	UnitsInStock float64 `json:"units_in_stock"`
}

func (o *Article) insertArticle(db *sql.DB) error {
	statement := fmt.Sprintf("INSERT INTO articles (name, category, unit_price, units_in_stock) "+
		"VALUES('%s', '%s', '%f', '%f')",
		o.Name, o.Category, o.UnitPrice, o.UnitsInStock)

	_, err := db.Exec(statement)
	if err != nil {
		return err
	}
	return nil
}

func (o *Article) getArticle(db *sql.DB) error {
	statement := fmt.Sprintf("SELECT * FROM articles WHERE id=%d", o.ID)
	return db.QueryRow(statement).Scan(&o.ID, &o.Name, &o.Category, &o.UnitPrice, &o.UnitsInStock)
}

func (o *Article) updateArticle(db *sql.DB) error {
	statement := fmt.Sprintf(
		"UPDATE articles "+
			"SET name='%s', category='%s', unit_price=%f, units_in_stock=%f"+
			"WHERE id=%d",
		o.Name, o.Category, o.UnitPrice, o.UnitsInStock, o.ID)
	_, err := db.Exec(statement)
	return err
}

func (o *Article) deleteArticle(db *sql.DB) error {
	statement := fmt.Sprintf(
		"DELETE FROM articles WHERE id=%d", o.ID)
	_, err := db.Exec(statement)
	return err
}

func getArticles(db *sql.DB) ([]Article, error) {
	var list []Article

	rows, err := db.Query("SELECT * FROM articles")

	if err != nil {
		return nil, err
	}

	defer func() {
		_ = rows.Close()
	}()

	for rows.Next() {
		var o Article
		if err := rows.Scan(&o.ID, &o.Name, &o.Category, &o.UnitPrice, &o.UnitsInStock); err != nil {
			return nil, err
		}
		list = append(list, o)
	}
	return list, nil
}
