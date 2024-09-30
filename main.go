package main

import (
	"context"
	"encoding/json"
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

	// Работа со строками
	err := client.Set(ctx, "key", "value", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := client.Get(ctx, "key").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("key:", val)

	// Работа с map
	user := map[string]interface{}{
		"name":    "Joe Peach",
		"age":     38,
		"country": "USA, Maslochusetts",
	}

	err = client.HSet(ctx, "user:2", user).Err()
	if err != nil {
		panic(err)
	}

	user2, err := client.HGetAll(ctx, "user:2").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println(user2)

	// Работа со структурами

	type Info struct {
		Name    string `json:"name"`
		Age     int    `json:"age"`
		Country string `json:"country"`
	}

	info1 := Info{Name: "Joe Peach", Age: 38, Country: "USA, Maslochusetts"}
	jsoninfo, _ := json.Marshal(info1)

	err = client.Set(ctx, "userInfo", jsoninfo, 0).Err()
	if err != nil {
		panic(err)
	}

	userinfo, err := client.Get(ctx, "userInfo").Result()
	if err != nil {
		panic(err)
	}
	getinfo := new(Info)
	json.Unmarshal([]byte(userinfo), getinfo)
	fmt.Println(getinfo)

	// Redis как брокер сообщений PUB/SUB
	// Отправка сообщения
	err = client.Publish(ctx, "test", "message").Err()
	if err != nil {
		panic(err)
	}

	// Прием сообщения
	pubsub := client.Subscribe(ctx, "test")

	defer pubsub.Close()

	// получение сообщений
	for {
		msg, err := pubsub.ReceiveMessage(ctx)
		if err != nil {
			panic(err)
		}

		fmt.Println(msg.Channel, msg.Payload)
	}

	// Но самый простой метод подписки, использовать канал
	// ch := pubsub.Channel()

	// for msg := range ch {
	// 	fmt.Println(msg.Channel, msg.Payload)
	// }

}
