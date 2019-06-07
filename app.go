package main

import (
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/context"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strings"
)


type App struct {
	Router *mux.Router
	DB     *sql.DB
	ArticleController *ArticleController
	UserController *UserController
	InvoiceController *InvoiceController
	InvoiceDetailController *InvoiceDetailController
}

func (a *App) Initialize(user, password, dbname string) {
	connectionString := fmt.Sprintf("%s:%s@/%s", user, password, dbname)

	var err error
	a.DB, err = sql.Open("mysql", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	a.Router = mux.NewRouter()
	a.ArticleController = NewArticleController(a.DB)
	a.UserController = NewUserController(a.DB)
	a.InvoiceController = NewInvoiceController(a.DB)
	a.InvoiceDetailController = NewInvoiceDetailController(a.DB)
	a.initializeRoutes()
}

func (a *App) Run(addr string) {
	allowedHeaders := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	allowedOrigins := handlers.AllowedOrigins([]string{"*"})
	allowedMethods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})
	log.Fatal(http.ListenAndServe(addr, handlers.CORS(allowedHeaders, allowedOrigins, allowedMethods)(a.Router)))
}

func (a *App) initializeRoutes() {
	// AUTH JWT
	a.Router.HandleFunc("/login", a.getToken).Methods("POST")

	// ARTICLES
	a.Router.HandleFunc("/articles", AuthenticationMiddleware(a.ArticleController.getArticles)).Methods("GET")
	a.Router.HandleFunc("/article", AuthenticationMiddleware(a.ArticleController.createArticle)).Methods("POST")
	a.Router.HandleFunc("/article/{id:[0-9]+}", AuthenticationMiddleware(a.ArticleController.getArticle)).Methods("GET")
	a.Router.HandleFunc("/article/{id:[0-9]+}", AuthenticationMiddleware(a.ArticleController.updateArticle)).Methods("PUT")
	a.Router.HandleFunc("/article/{id:[0-9]+}", AuthenticationMiddleware(a.ArticleController.deleteArticle)).Methods("DELETE")

	// INVOICES
	a.Router.HandleFunc("/invoice", AuthenticationMiddleware(a.InvoiceController.getInvoices)).Methods("GET")
	a.Router.HandleFunc("/invoice", AuthenticationMiddleware(a.InvoiceController.createInvoice)).Methods("POST")
	a.Router.HandleFunc("/invoice/{id:[0-9]+}", AuthenticationMiddleware(a.InvoiceController.getInvoice)).Methods("GET")
	a.Router.HandleFunc("/invoice/{id:[0-9]+}", AuthenticationMiddleware(a.InvoiceController.updateInvoice)).Methods("PUT")
	a.Router.HandleFunc("/invoice/{id:[0-9]+}", AuthenticationMiddleware(a.InvoiceController.deleteInvoice)).Methods("DELETE")


	// INVOICE DETAILS
	a.Router.HandleFunc("/invoice_details", AuthenticationMiddleware(a.InvoiceDetailController.getInvoiceDetails)).Methods("GET")
	a.Router.HandleFunc("/invoice_details", AuthenticationMiddleware(a.InvoiceDetailController.createInvoiceDetail)).Methods("POST")
	a.Router.HandleFunc("/invoice_details/{invoice_id:[0-9]+}/{article_id:[0-9]+}", AuthenticationMiddleware(a.InvoiceDetailController.getInvoiceDetail)).Methods("GET")
	a.Router.HandleFunc("/invoice_details", AuthenticationMiddleware(a.InvoiceDetailController.updateInvoiceDetail)).Methods("PUT")
	a.Router.HandleFunc("/invoice_details/{invoice_id:[0-9]+}/{article_id:[0-9]+}", AuthenticationMiddleware(a.InvoiceDetailController.deleteInvoiceDetail)).Methods("DELETE")



	// USER
	a.Router.HandleFunc("/user", AuthenticationMiddleware(a.UserController.getUsers)).Methods("GET")
	a.Router.HandleFunc("/user", AuthenticationMiddleware(a.UserController.createUser)).Methods("POST")
	a.Router.HandleFunc("/user/{id:[0-9]+}", AuthenticationMiddleware(a.UserController.getUser)).Methods("GET")
	a.Router.HandleFunc("/user/{id:[0-9]+}", AuthenticationMiddleware(a.UserController.updateUser)).Methods("PUT")
	a.Router.HandleFunc("/user/{id:[0-9]+}", AuthenticationMiddleware(a.UserController.deleteUser)).Methods("DELETE")
}


func (a *App) getRefreshToken(w http.ResponseWriter, r *http.Request) {

}

// AUTH JWT
func (a *App) getToken(w http.ResponseWriter, r *http.Request) {
	var user User
	_ = json.NewDecoder(r.Body).Decode(&user)
	passwordHash := sha256.Sum256([]byte(user.Password))


	if err := user.login(a.DB, user.Username, string(passwordHash[:])); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusUnauthorized, "username or password don't match")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"password": passwordHash,
	})

	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		fmt.Println(err)
	}

	u := User{Username: user.Username, Token:tokenString}
	if err := u.updateUserTokenByUsername(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := u.getUserByUsername(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, u)
}

func AuthenticationMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		authorizationHeader := req.Header.Get("Authorization")
		if authorizationHeader != "" {
			bearerToken := strings.Split(authorizationHeader, " ")
			if len(bearerToken) == 2 {
				token, err := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, fmt.Errorf("there was an error")
					}
					return []byte("secret"), nil
				})
				if err != nil {
					respondWithError(w, http.StatusBadRequest, err.Error())
					return
				}
				if token.Valid {
					log.Println("TOKEN WAS VALID")
					context.Set(req, "decoded", token.Claims)
					next(w, req)
				} else {
					respondWithError(w, http.StatusForbidden, "Invalid authorization token")
				}
			}
		} else {
			respondWithError(w, http.StatusForbidden, "An authorization header is required")
		}
	})
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, _ = w.Write(response)
}
