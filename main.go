package main

import (
	"fmt"

	"github.com/shivamaravanthe/superSite-BE/server"
	"github.com/shivamaravanthe/superSite-BE/store"
)

func main() {
	store.DBConnect()
	if store.DB == nil {
		fmt.Println("Database not connected")
		return
	}
	routes := server.SetUpRoutes()
	fmt.Println("Server Running")
	err := routes.Serve()
	if err != nil {
		fmt.Printf("Errored Starting Server %v\n", err.Error())
		return
	}

}
