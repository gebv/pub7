package tarantool

import (
	"log"
	"os"
	"time"

	tarantool "github.com/tarantool/go-tarantool"
)

var (
	DefaultTimeout   time.Duration = 500 * time.Millisecond
	DefaultReconnect time.Duration = 5 * time.Second
)

// SetupFromENV инициализация клиента
// TARANTOOL_SERVER
// TARANTOOL_USER_NAME
// TARANTOOL_USER_PASSWORD
func SetupFromENV() (client *tarantool.Connection, err error) {
	// 127.0.0.1:3013
	server := os.Getenv("TARANTOOL_SERVER")
	user := os.Getenv("TARANTOOL_USER_NAME")
	pwd := os.Getenv("TARANTOOL_USER_PASSWORD")

	log.Println("tarantool opts", server, user, pwd)

	opts := tarantool.Opts{
		Timeout:       DefaultTimeout,
		Reconnect:     DefaultReconnect,
		MaxReconnects: 100,
		User:          user,
		Pass:          pwd,
	}

	client, err = tarantool.Connect(server, opts)
	if err != nil {
		return client, err
	}

	_, err = client.Ping()

	return client, err
}
