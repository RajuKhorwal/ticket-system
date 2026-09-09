package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

// createTicketRequest is the expected JSON body for POST /tickets.
type createTicketRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

// HandleCreateTicket creates a new ticket owned by the logged-in user.
// Every new ticket starts with status "open".
func (app *App) HandleCreateTicket(w http.ResponseWriter, r *http.Request) {
	userID := userIDFromContext(r)

	var req createTicketRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	req.Title = strings.TrimSpace(req.Title)
	req.Description = strings.TrimSpace(req.Description)
	if req.Title == "" {
		writeError(w, http.StatusBadRequest, "title is required")
		return
	}

	ticket := app.Store.CreateTicket(userID, req.Title, req.Description)
	writeJSON(w, http.StatusCreated, ticket)
}

// HandleListTickets returns every ticket that belongs to the logged-in user.
// Tickets created by other users are never included.
func (app *App) HandleListTickets(w http.ResponseWriter, r *http.Request) {
	userID := userIDFromContext(r)
	tickets := app.Store.ListTicketsByUser(userID)
	writeJSON(w, http.StatusOK, tickets)
}

// HandleGetTicket returns a single ticket by ID, but only if it belongs
// to the logged-in user. Otherwise, it responds as if the ticket does not
// exist at all (404), so no information leaks about other users' tickets.
func (app *App) HandleGetTicket(w http.ResponseWriter, r *http.Request) {
	userID := userIDFromContext(r)
	ticketID := r.PathValue("id")

	ticket, ok := app.Store.GetTicketByID(ticketID)
	if !ok || ticket.UserID != userID {
		writeError(w, http.StatusNotFound, ErrTicketNotFound.Error())
		return
	}

	writeJSON(w, http.StatusOK, ticket)
}

// updateStatusRequest is the expected JSON body for PATCH /tickets/{id}/status.
type updateStatusRequest struct {
	Status string `json:"status"`
}

// HandleUpdateTicketStatus moves a ticket to a new status, following the
// required flow: open -> in_progress -> closed. A closed ticket can never
// move to any other status again.
func (app *App) HandleUpdateTicketStatus(w http.ResponseWriter, r *http.Request) {
	userID := userIDFromContext(r)
	ticketID := r.PathValue("id")

	ticket, ok := app.Store.GetTicketByID(ticketID)
	if !ok || ticket.UserID != userID {
		writeError(w, http.StatusNotFound, ErrTicketNotFound.Error())
		return
	}

	var req updateStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	newStatus := strings.TrimSpace(strings.ToLower(req.Status))
	if newStatus != StatusOpen && newStatus != StatusInProgress && newStatus != StatusClosed {
		writeError(w, http.StatusBadRequest, "status must be one of: open, in_progress, closed")
		return
	}

	// Check that this exact transition is allowed from the ticket's current status.
	if allowedTransitions[ticket.Status] != newStatus {
		writeError(w, http.StatusBadRequest, ErrInvalidTransition.Error()+
			": cannot move from '"+ticket.Status+"' to '"+newStatus+"'")
		return
	}

	app.Store.UpdateTicketStatus(ticket, newStatus)
	writeJSON(w, http.StatusOK, ticket)
}
