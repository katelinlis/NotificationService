package main

import (
	"context"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/katelinlis/NotificationService/internal/app"
)

func main() {
	app.StartApp(context.Background())
}
