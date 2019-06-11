package main

import (
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type InvoiceDetailController struct {
	DB *sql.DB
}

func NewInvoiceDetailController(db *sql.DB) *InvoiceDetailController {
	o := &InvoiceDetailController{db}
	return o
}

func (a *InvoiceDetailController) getInvoiceDetails(w http.ResponseWriter, r *http.Request) {
	result, err := getInvoiceDetails(a.DB)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, result)
}

func (a *InvoiceDetailController) createInvoiceDetail(w http.ResponseWriter, r *http.Request) {
	var c InvoiceDetail
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&c); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}
	defer func() {
		_ = r.Body.Close()
	}()

	if err := c.insertInvoiceDetail(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *InvoiceDetailController) getInvoiceDetailsFromInvoiceID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	invoiceId, err := strconv.Atoi(vars["invoice_id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Invoince ID")
		return
	}

	c := InvoiceDetail{InvoiceID: uint(invoiceId)}
	result, err := c.getInvoiceDetailsByInvoiceID(a.DB)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, result)
}

func (a *InvoiceDetailController) getInvoiceDetail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	invoiceId, err := strconv.Atoi(vars["invoice_id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Invoince ID")
		return
	}

	articleId, err := strconv.Atoi(vars["article_id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Invoince ID")
		return
	}

	c := InvoiceDetail{InvoiceID: uint(invoiceId), ArticleID: uint(articleId)}
	if err := c.getInvoiceDetail(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Invoice Detail not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, c)
}

func (a *InvoiceDetailController) updateInvoiceDetail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	invoiceId, err := strconv.Atoi(vars["invoice_id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Invoince ID")
		return
	}

	articleId, err := strconv.Atoi(vars["article_id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Invoince ID")
		return
	}

	c := InvoiceDetail{InvoiceID: uint(invoiceId), ArticleID: uint(articleId)}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&c); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}
	defer func() {
		_ = r.Body.Close()
	}()

	if err := c.updateInvoiceDetail(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *InvoiceDetailController) deleteInvoiceDetail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	invoiceId, err := strconv.Atoi(vars["invoice_id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Invoince ID")
		return
	}

	articleId, err := strconv.Atoi(vars["article_id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Invoince ID")
		return
	}

	c := InvoiceDetail{InvoiceID: uint(invoiceId), ArticleID: uint(articleId)}
	if err := c.deleteInvoiceDetail(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}
