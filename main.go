package main

import (
	"context"

	"github.com/avenir/notification-service/internal/app"
)

func main() {
	app.StartApp(context.Background())
}
