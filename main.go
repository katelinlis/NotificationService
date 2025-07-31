package main

import (
	"context"

	"github.com/avenir/notification-service/internal/app"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	app.StartApp(context.Background())
}
