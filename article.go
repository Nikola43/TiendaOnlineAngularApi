package main

import (
	"database/sql"
	"fmt"
	"time"
)

type Article struct {
	ID                    uint   `json:"id"`
	UserID                uint   `json:"user_id"`
	CreationDate          string `json:"creation_date"`
	ShippingMethod        string `json:"shipping_method"`
	PaymentMethod         string `json:"payment_method"`
	EstimatedDeliveryDate string `json:"estimated_delivery_date"`
}

func (o *Article) insertArticle(db *sql.DB) error {
	date := fmt.Sprintf(time.Now().Format("2006-01-02 15:04:05"))
	statement := fmt.Sprintf("INSERT INTO articles (user_id, creation_date, shipping_method, payment_method, estimated_delivery_date) "+
		"VALUES('%d', '%s', '%s', '%s', '%s')",
		o.UserID, date, o.ShippingMethod, o.PaymentMethod, o.EstimatedDeliveryDate)

	_, err := db.Exec(statement)
	if err != nil {
		return err
	}
	return nil
}

func (o *Article) getArticle(db *sql.DB) error {
	statement := fmt.Sprintf("SELECT * FROM articles WHERE id=%d", o.ID)
	return db.QueryRow(statement).Scan(&o)
}

func (o *Article) updateArticle(db *sql.DB) error {
	statement := fmt.Sprintf(
		"UPDATE articles "+
			"SET shipping_method='%s', estimated_delivery_date='%s' " +
			"WHERE id=%d",
		o.ShippingMethod, o.EstimatedDeliveryDate, o.ID)
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
		if err := rows.Scan(&o); err != nil {
			return nil, err
		}
		list = append(list, o)
	}
	return list, nil
}
