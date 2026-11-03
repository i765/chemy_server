package main

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
)

//go:embed frontend/dist*
var assets embed.FS

func main() {

	const (
		address = "127.0.0.1:9765"
	)

	frontendFS, err := fs.Sub(assets, "frontend/dist")

	if err != nil {
		panic(err)
	}

	// Создаем обработчик для статических файлов
	fileServer := http.FileServer(http.FS(frontendFS))

	// Назначаем обработчик на главный роут
	http.Handle("/", fileServer)

	fmt.Printf("Сервер запущен на http://%s", address)
	http.ListenAndServe(address, nil)

}
