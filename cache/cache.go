package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

// UserService сервис для работы с пользователями
type UserService struct {
	db  *sql.DB
	rdb *redis.Client
}

// User структура для демонстрации
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// GetUser получает информацию о пользователе по его ID
func (s *UserService) GetUser(ctx context.Context, id int) (*User, error) {
	// Проверяем, есть ли пользователь в кеше
	cachedUser, err := s.rdb.Get(ctx, fmt.Sprintf("user:%d", id)).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}

	if cachedUser != "" {
		var user User
		err = json.Unmarshal([]byte(cachedUser), &user)
		if err != nil {
			return nil, err
		}
		return &user, nil
	}

	// Если пользователя нет в кеше, получаем его из базы данных
	user, err := s.getUserFromDB(ctx, id)
	if err != nil {
		return nil, err
	}

	// Кешируем пользователя
	userJSON, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}
	err = s.rdb.Set(ctx, fmt.Sprintf("user:%d", id), userJSON, 0).Err()
	if err != nil {
		return nil, err
	}

	return user, nil
}

// getUserFromDB получает информацию о пользователе из базы данных
func (s *UserService) getUserFromDB(ctx context.Context, id int) (*User, error) {
	row := s.db.QueryRowContext(ctx, "SELECT * FROM users WHERE id = $1", id)
	var user User
	err := row.Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func main() {
	// Создаем клиент Redis
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // адрес Redis
		Password: "",               // пароль Redis
		DB:       0,                // номер базы данных
	})

	// Создаем подключение к базе данных
	db, err := sql.Open("postgres", "user=myuser dbname=mydb sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	// Создаем сервис для работы с пользователями
	userService := &UserService{db: db, rdb: rdb}

	// Получаем информацию о пользователе
	user, err := userService.GetUser(ctx, 1)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(user) // выводит &{1 John john@example.com}
}

// Мы используем Redis в качестве кеша, чтобы уменьшить нагрузку на базу данных и увеличить скорость ответа.
// Если пользователь уже находится в кеше, мы возвращаем его из кеша,
// в противном случае мы получаем его из базы данных и кешируем его.
