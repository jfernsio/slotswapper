package main

import (
	"log"
	"net/http"
	// "time"

	// "github.com/jfernsio/slotswapper/internals/config"
	"github.com/jfernsio/slotswapper/internals/database"
	"github.com/jfernsio/slotswapper/internals/middleware"
	"github.com/jfernsio/slotswapper/internals/routes"
)

func main() {
	database.Init()

	mux := http.NewServeMux()
	routes.RegisterRoutes(mux)
	

	log.Printf("🚀 Server running on http://localhost:", 8080)
	http.ListenAndServe(":8080", middleware.EnableCORS(mux))

}
