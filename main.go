package main

import (
	"log"
	"net/http"
	"os"
)

// App holds shared dependencies for the handlers.
type App struct {
	Store *Store
}

func main() {
	app := &App{
		Store: NewStore(),
	}

	mux := http.NewServeMux()

	// Public routes - no login required.
	mux.HandleFunc("GET /health", HandleHealth)
	mux.HandleFunc("POST /auth/register", app.HandleRegister)
	mux.HandleFunc("POST /auth/login", app.HandleLogin)

	// Protected routes - require a valid "Authorization: Bearer <token>" header.
	mux.HandleFunc("POST /tickets", RequireAuth(app.HandleCreateTicket))
	mux.HandleFunc("GET /tickets", RequireAuth(app.HandleListTickets))
	mux.HandleFunc("GET /tickets/{id}", RequireAuth(app.HandleGetTicket))
	mux.HandleFunc("PATCH /tickets/{id}/status", RequireAuth(app.HandleUpdateTicketStatus))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("ticket-system listening on port %s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}
}

// HandleHealth is a simple public endpoint used to verify the service is running.
func HandleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
