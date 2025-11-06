package routes

import (
	"net/http"

	"github.com/jfernsio/slotswapper/internals/handlers"
	"github.com/jfernsio/slotswapper/internals/middleware"
)

func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("SlotSwapper backend (refactored)"))
	})
	mux.HandleFunc("/health", handlers.HealthHandler)

	// Auth routes
	mux.HandleFunc("/api/signup", handlers.SignupHandler)
	// Note: login route will be added in the next step
	mux.HandleFunc("/api/login",handlers.Login)
	//Protected routes
	mux.Handle("/profile",middleware.AuthMiddleware(http.HandlerFunc(handlers.Dashboard)))
	mux.Handle("/api/create/event",middleware.AuthMiddleware(http.HandlerFunc(handlers.CreateEvent)))
	mux.Handle("/api/events",middleware.AuthMiddleware(http.HandlerFunc(handlers.ListEvents)))
	mux.Handle("/update-events/{id}",middleware.AuthMiddleware(http.HandlerFunc(handlers.UpdateEvent)))
	mux.Handle("/delet-events/{id}",middleware.AuthMiddleware(http.HandlerFunc(handlers.DeletEvent)))
	mux.Handle("/api/swappable-slots",middleware.AuthMiddleware(http.HandlerFunc(handlers.GetSwappableSlots)))
	mux.Handle("/api/swap-req",middleware.AuthMiddleware(http.HandlerFunc(handlers.CreateSwapRequest)))
	mux.Handle("/api/swap-res",middleware.AuthMiddleware(http.HandlerFunc(handlers.RespondToSwap)))






}
