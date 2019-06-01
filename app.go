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
	"strconv"
	"strings"
)


type App struct {
	Router *mux.Router
	DB     *sql.DB
	userController *UserController
}

func (a *App) Initialize(user, password, dbname string) {
	connectionString := fmt.Sprintf("%s:%s@/%s", user, password, dbname)

	var err error
	a.DB, err = sql.Open("mysql", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	a.Router = mux.NewRouter()
	a.userController = NewUserController(a.DB)
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

	/*

	// ARTICLES
	a.Router.HandleFunc("/articles", a.getCalculators).Methods("GET")
	a.Router.HandleFunc("/articles", a.createCalculator).Methods("POST")
	a.Router.HandleFunc("/articles/{id:[0-9]+}", a.getCalculator).Methods("GET")
	a.Router.HandleFunc("/articles/{id:[0-9]+}", a.updateCalculator).Methods("PUT")
	a.Router.HandleFunc("/articles/{id:[0-9]+}", a.deleteCalculator).Methods("DELETE")

	// INVOICES
	a.Router.HandleFunc("/invoices", AuthenticationMiddleware(a.getInvoices)).Methods("GET")
	a.Router.HandleFunc("/invoices", AuthenticationMiddleware(a.createInvoice)).Methods("POST")
	a.Router.HandleFunc("/invoices/{id:[0-9]+}", AuthenticationMiddleware(a.getInvoice)).Methods("GET")
	a.Router.HandleFunc("/invoices/{id:[0-9]+}", AuthenticationMiddleware(a.updateInvoice)).Methods("PUT")
	a.Router.HandleFunc("/invoices/{id:[0-9]+}", AuthenticationMiddleware(a.deleteInvoice)).Methods("DELETE")

	// USER PROMOS
	a.Router.HandleFunc("/user_promo/{id:[0-9]+}", AuthenticationMiddleware(a.getUserPromos)).Methods("GET")
	a.Router.HandleFunc("/user_promo", AuthenticationMiddleware(a.createUserPromo)).Methods("POST")
	a.Router.HandleFunc("/user_promo/emitter/{id:[0-9]+}", AuthenticationMiddleware(a.getUserPromoByReceiverID)).Methods("GET")
	a.Router.HandleFunc("/user_promo/receiver/{id:[0-9]+}", AuthenticationMiddleware(a.getUserPromoByReceiverID)).Methods("GET")
	a.Router.HandleFunc("/user_promo", AuthenticationMiddleware(a.updateUserPromo)).Methods("PUT")
	a.Router.HandleFunc("/user_promo", AuthenticationMiddleware(a.deleteUserPromo)).Methods("DELETE")

	*/

	// USER
	a.Router.HandleFunc("/user", AuthenticationMiddleware(a.userController.getUsers)).Methods("GET")
	a.Router.HandleFunc("/user", AuthenticationMiddleware(a.userController.createUser)).Methods("POST")
	a.Router.HandleFunc("/user/{id:[0-9]+}", AuthenticationMiddleware(a.userController.getUser)).Methods("GET")
	a.Router.HandleFunc("/user/{id:[0-9]+}", AuthenticationMiddleware(a.userController.updateUser)).Methods("PUT")
	a.Router.HandleFunc("/user/{id:[0-9]+}", AuthenticationMiddleware(a.userController.deleteUser)).Methods("DELETE")
}


// INVOICES ------------------------------------------------------------------------------------------------------------
func (a *App) getInvoices(w http.ResponseWriter, r *http.Request) {
	result, err := getInvoices(a.DB)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, result)
}

func (a *App) createInvoice(w http.ResponseWriter, r *http.Request) {
	var c Invoice
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&c); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}
	defer func() {
		_ = r.Body.Close()
	}()

	if err := c.insertInvoice(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *App) getInvoice(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Invoice ID")
		return
	}

	c := Invoice{ID: uint(id)}
	if err := c.getInvoice(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Invoice not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, c)
}

func (a *App) updateInvoice(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Invoice ID")
		return
	}

	var c Invoice
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&c); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}

	defer func() {
		_ = r.Body.Close()
	}()

	c.ID = uint(id)
	if err := c.updateInvoice(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *App) deleteInvoice(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Invoice ID")
		return
	}

	c := Invoice{ID: uint(id)}
	if err := c.deleteInvoice(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}
// INVOICES ------------------------------------------------------------------------------------------------------------


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
