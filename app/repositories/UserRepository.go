package repositories

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/InfinityIkh/app/validation"
)

type User struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Date     string `json:"date_registration"`
}

type Body struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (ur *UserRepository) CreateUser(req *http.Request) (User, error) {
	var body Body
	err := json.NewDecoder(req.Body).Decode(&body)
	if err != nil {
		return User{}, err
	}
	err = validation.Validate(body.Username, body.Email)
	if err != nil {
		return User{}, err
	}

	res, err := ur.db.Exec("INSERT INTO users(username ,email) VALUES (? ,?)", body.Username, body.Email)
	if err != nil {
		return User{}, err
	}

	userId, _ := res.LastInsertId()
	createdUser := User{
		Id:       int(userId),
		Username: body.Username,
		Email:    body.Email,
	}

	return createdUser, nil
}

func (ur *UserRepository) GetUsers() ([]User, error) {
	rows, err := ur.db.Query("SELECT * FROM users")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var users []User = []User{}
	for rows.Next() {
		var user User
		err := rows.Scan(
			&user.Id,
			&user.Username,
			&user.Email,
			&user.Date,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

func (ur *UserRepository) GetUser(req *http.Request) (User ,error) {
	var user User
	var id string
	id = req.PathValue("id")
	stmt ,err := ur.db.Prepare("SELECT * FROM users WHERE id = ?")
	if err != nil {
		return User{} ,err
	}
	defer stmt.Close()
	err = stmt.QueryRow(id).Scan(&user.Id ,&user.Username ,&user.Email ,&user.Date)
	if err != nil {
		return User{} ,err
	}
	return user ,nil
}