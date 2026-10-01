package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	database "github.com/InfinityIkh/database"
	repositories "github.com/InfinityIkh/app/Repositories"
	"github.com/joho/godotenv"
)

func main(){
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	db ,err := database.Connect()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	router := http.NewServeMux()
	router.HandleFunc("GET /users" ,func(w http.ResponseWriter, r *http.Request) {
		users ,err := repositories.NewUserRepository(db).GetUsers()
		if err != nil {
			panic(err)
		}
		w.Header().Set("Content-Type","application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"users" : users,
		})
	})

	router.HandleFunc("POST /users" ,func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		user ,err := repositories.NewUserRepository(db).CreateUser(r)
		if err != nil{
			w.Header().Set("Content-Type","application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]any{
				"message" : "Could not create user",
				"details": err.Error(),
			})
			return
		}
		w.Header().Set("Content-Type","application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]any{
			"user" : user,
			"message" : "User successfully created in database",
		})

	})
	
	fmt.Println("Listenning on 8080")
	http.ListenAndServe(":8080" ,router)
}