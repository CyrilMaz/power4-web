package main

import (
	"log"
	"net/http"

	"github.com/CyrilMaz/power4-web/handlers"
	"github.com/CyrilMaz/power4-web/theme"
)

func main() {
	http.HandleFunc("/", handlers.Home)
	http.HandleFunc("/play", handlers.Play)
	http.HandleFunc("/power", handlers.UsePower)
	http.HandleFunc("/reset", handlers.Reset)
	http.HandleFunc("/set-theme", theme.ToggleHandler)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	log.Println("Serveur sur http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
