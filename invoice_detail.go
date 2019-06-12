package main

import (
	"database/sql"
	"fmt"
)

type InvoiceDetail struct {
	InvoiceID   uint   `json:"invoice_id"`
	ArticleID   uint   `json:"article_id"`
	ArticleName string `json:"article_name"`
	Quantity    uint   `json:"quantity"`
}

func (o *InvoiceDetail) insertInvoiceDetail(db *sql.DB) error {
	statement := fmt.Sprintf("INSERT INTO invoice_details (invoice_id, article_id, quantity) "+
		"VALUES('%d', '%d', '%d')",
		o.InvoiceID, o.ArticleID, o.Quantity)
	_, err := db.Exec(statement)
	fmt.Println(statement)
	if err != nil {
		return err
	}
	return nil
}

func (o *InvoiceDetail) getInvoiceDetail(db *sql.DB) error {
	statement := fmt.Sprintf("SELECT * FROM invoice_details WHERE invoice_id=%d AND article_id=%d", o.InvoiceID, o.ArticleID)
	return db.QueryRow(statement).Scan(&o.InvoiceID, &o.ArticleID, &o.Quantity)
}

func (o *InvoiceDetail) getInvoiceDetailsByInvoiceID(db *sql.DB) ([]InvoiceDetail, error) {
	statement := fmt.Sprintf("SELECT id.invoice_id, id.article_id, a.name, id.quantity "+
		"FROM invoice_details as id "+
		"INNER JOIN articles as a "+
		"ON id.article_id = a.id "+
		"WHERE invoice_id=%d", o.InvoiceID)
	rows, err := db.Query(statement)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = rows.Close()
	}()

	var list []InvoiceDetail

	for rows.Next() {
		var o InvoiceDetail
		if err := rows.Scan(&o.InvoiceID, &o.ArticleID, &o.ArticleName, &o.Quantity); err != nil {
			return nil, err
		}
		list = append(list, o)
	}
	return list, nil
}

func (o *InvoiceDetail) getInvoiceDetailByArticleID(db *sql.DB) error {
	statement := fmt.Sprintf("SELECT * FROM invoice_details WHERE article_id=%d", o.ArticleID)
	return db.QueryRow(statement).Scan(&o.InvoiceID, &o.ArticleID, &o.Quantity)
}

func (o *InvoiceDetail) updateInvoiceDetail(db *sql.DB) error {
	statement := fmt.Sprintf("UPDATE invoice_details SET quantity=%d WHERE invoice_id=%d AND article_id=%d",
		o.Quantity, o.InvoiceID, o.ArticleID)
	_, err := db.Exec(statement)
	return err
}

func (o *InvoiceDetail) deleteInvoiceDetail(db *sql.DB) error {
	statement := fmt.Sprintf("DELETE FROM invoice_details WHERE invoice_id=%d AND article_id=%d", o.InvoiceID, o.ArticleID)
	fmt.Println(statement)
	_, err := db.Exec(statement)
	return err
}

func getInvoiceDetails(db *sql.DB) ([]InvoiceDetail, error) {
	statement := fmt.Sprintf("SELECT * FROM invoice_details")
	rows, err := db.Query(statement)

	if err != nil {
		return nil, err
	}

	defer func() {
		_ = rows.Close()
	}()

	var list []InvoiceDetail

	for rows.Next() {
		var o InvoiceDetail
		if err := rows.Scan(&o.InvoiceID, &o.ArticleID, &o.Quantity); err != nil {
			return nil, err
		}
		list = append(list, o)
	}
	return list, nil
}
