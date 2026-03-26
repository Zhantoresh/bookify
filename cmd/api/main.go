package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	
	"bookify/internal/database"
	"bookify/internal/handlers"
	"bookify/internal/usecase"

	
	_ "github.com/lib/pq" 
)

func main() {
	
	connStr := "user=postgres password=yourpass dbname=bookify sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	
	userRepo := database.NewUserRepository(db)
	
	
	userUsecase := usecase.NewUserUsecase(userRepo)
	
	
	authHandler := handlers.NewAuthHandler(userUsecase)

	
	http.HandleFunc("/register", authHandler.Register)
	http.HandleFunc("/login", authHandler.Login)

	
	fmt.Println("Server starting on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}