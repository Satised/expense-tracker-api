package main

import (
	"log"
	"net/http"

	"expense-tracker-api/internal/handler"
	"expense-tracker-api/internal/repository"
	"expense-tracker-api/internal/service"
)

func main() {
	store := repository.NewExpenseStore()
	svc := service.NewExpenseService(store)
	h := handler.New(svc)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", h.Healthz)
	mux.HandleFunc("POST /expenses", h.CreateExpense)
	mux.HandleFunc("GET /expenses", h.ListExpenses)

	log.Println("сервер слушает на http://localhost:8080")

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}
