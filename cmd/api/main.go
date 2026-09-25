package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"

	"expense-tracker-api/internal/handler"
	"expense-tracker-api/internal/repository"
	"expense-tracker-api/internal/service"
)

func main() {
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("переменная DB_URL не задана")
	}

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatal("не удалось создать пул соединений: ", err)
	}
	defer pool.Close()

	err = pool.Ping(ctx)
	if err != nil {
		log.Fatal("база не отвечает: ", err)
	}
	log.Println("подключились к базе")

	store := repository.NewPostgresStore(pool)
	svc := service.NewExpenseService(store)
	h := handler.New(svc)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", h.Healthz)
	mux.HandleFunc("POST /expenses", h.CreateExpense)
	mux.HandleFunc("GET /expenses", h.ListExpenses)

	log.Println("сервер слушает на http://localhost:8080")

	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}
