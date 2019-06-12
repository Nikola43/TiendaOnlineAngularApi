package main

import (
	"database/sql"
	"fmt"
	"time"
)

type Invoice struct {
	ID                    uint   `json:"id"`
	UserID                uint   `json:"user_id"`
	CreationDate          string `json:"creation_date"`
	ShippingMethod        string `json:"shipping_method"`
	PaymentMethod         string `json:"payment_method"`
	EstimatedDeliveryDate string `json:"estimated_delivery_date"`
	UserName              string `json:"user_name"`
}

func (o *Invoice) insertInvoice(db *sql.DB) error {
	date := fmt.Sprintf(time.Now().Format("2006-01-02 15:04:05"))
	statement := fmt.Sprintf("INSERT INTO invoices (user_id, creation_date, shipping_method, payment_method, estimated_delivery_date) "+
		"VALUES('%d', '%s', '%s', '%s', '%s')",
		o.UserID, date, o.ShippingMethod, o.PaymentMethod, o.EstimatedDeliveryDate)

	_, err := db.Exec(statement)
	if err != nil {
		return err
	}
	return nil
}

func (o *Invoice) getInvoice(db *sql.DB) error {
	statement := fmt.Sprintf("SELECT * FROM invoices WHERE id=%d", o.ID)
	return db.QueryRow(statement).Scan(&o)
}

func (o *Invoice) updateInvoice(db *sql.DB) error {
	statement := fmt.Sprintf(
		"UPDATE invoices "+
			"SET shipping_method='%s', estimated_delivery_date='%s' "+
			"WHERE id=%d",
		o.ShippingMethod, o.EstimatedDeliveryDate, o.ID)
	_, err := db.Exec(statement)
	return err
}

func (o *Invoice) deleteInvoice(db *sql.DB) error {
	statement := fmt.Sprintf(
		"DELETE FROM invoices WHERE id=%d", o.ID)
	_, err := db.Exec(statement)
	return err
}

func getInvoices(db *sql.DB) ([]Invoice, error) {
	var list []Invoice

	rows, err := db.Query("SELECT i.id, i.user_id, i.creation_date, i.shipping_method, i.payment_method, i.estimated_delivery_date, u.name FROM invoices as i INNER JOIN users as u ON i.user_id = u.id")

	if err != nil {
		return nil, err
	}

	defer func() {
		_ = rows.Close()
	}()

	for rows.Next() {
		var o Invoice
		if err := rows.Scan(&o.ID, &o.UserID, &o.CreationDate, &o.ShippingMethod, &o.PaymentMethod, &o.EstimatedDeliveryDate, &o.UserName); err != nil {
			return nil, err
		}
		list = append(list, o)
	}
	return list, nil
}
