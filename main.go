package main

import (
	"context"
	"embed"
	"os"
	"os/signal"
	"syscall"
)

//go:embed frontend/dist
var assets embed.FS

// ANSI Codes
const (
	CheckMark = "\033[32m✔\033[0m"
	ErrorMark = "\033[31m✖\033[0m"
)

func main() {
	// Создаем контекст, который отменится при получении SIGINT (Ctrl+C) или SIGTERM
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	app := NewApp()
	app.run(ctx)
}
