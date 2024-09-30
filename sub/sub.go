package main

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func main() {
	ctx := context.Background()
	// Подключение к Redis и настройка Config
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	// Подписка на канал
	pubsub := client.Subscribe(ctx, "test")

	defer pubsub.Close()

	// получение сообщений
	// for {
	// 	msg, err := pubsub.ReceiveMessage(ctx)
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	fmt.Println(msg.Channel, msg.Payload)
	// }

	// Но самый простой метод подписки, использовать канал
	ch := pubsub.Channel()

	for msg := range ch {
		fmt.Println(msg.Channel, msg.Payload)
	}
}
