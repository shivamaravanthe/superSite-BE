package server

import (
	"net/http"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/shivamaravanthe/superSite-BE/handlerFunctions"
	"github.com/shivamaravanthe/superSite-BE/jwt"
)

const (
	writeTimeout = time.Second * 15
	readTimeout  = time.Second * 15
	idleTimeout  = time.Second * 60
)

type httpRouter struct {
	router *mux.Router
}

func SetUpRoutes() httpRouter {
	router := mux.NewRouter()
	noAuthrouter := router.PathPrefix("/api/v1/").Subrouter()

	noAuthrouter.HandleFunc("/ping", handlerFunctions.Ping)
	noAuthrouter.HandleFunc("/login", handlerFunctions.Login).Methods("POST")
	noAuthrouter.HandleFunc("/signUp", handlerFunctions.SingUp).Methods("POST")

	authRouter := noAuthrouter.PathPrefix("").Subrouter()
	authRouter.Use(jwt.Authenticate)
	authRouter.HandleFunc("/password/create", handlerFunctions.StorePassword).Methods("POST")
	authRouter.HandleFunc("/password/update", handlerFunctions.ChangePassword).Methods("POST")
	authRouter.HandleFunc("/password/show", handlerFunctions.GetPassword).Methods("POST")
	authRouter.HandleFunc("/password/list", handlerFunctions.ListPassword).Methods("GET")
	// noAuthrouter.HandleFunc("/password/change", handlerFunctions.StorePassword).Methods("GET")

	// authRouter.Use()

	return httpRouter{router}
}

func (h httpRouter) Serve() error {
	_credentials := handlers.AllowCredentials()
	_methods := handlers.AllowedMethods([]string{"GET", "HEAD", "PUT", "PATCH", "POST", "DELETE", "OPTIONS"})
	_origins := handlers.AllowedOrigins([]string{"http://localhost:3000", "http://localhost:3000/", "http://localhost", "http://192.168.0.115", "http://192.168.0.115:3000"})
	_headers := handlers.AllowedHeaders([]string{"Access-Control-Allow-Headers", "Access-Control-Allow-Origin",
		"Access-Control-Allow-Headers", "Origin", "Accept", "X-Requested-With", "Content-Type", "Authorization",
		"Access-Control-Request-Method", "Access-Control-Request-Headers", "token", "Access-Control-Expose-Headers"})

	server := &http.Server{
		Addr:         ":8087",
		WriteTimeout: writeTimeout, // To avoid Solaris attacks.
		ReadTimeout:  readTimeout,
		IdleTimeout:  idleTimeout,
		Handler:      handlers.CORS(_credentials, _methods, _origins, _headers)(h.router),
	}

	err := server.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}
