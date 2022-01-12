package main

import (
	app "TP_DB/internal/app"
	"github.com/valyala/fasthttp"
	"log"
	"os"
)

const logMessage = "main:"

func main() {
	application, err := app.NewApp()
	if err != nil {
		os.Exit(1)
	}

	router := application.CreateRouter()
	log.Println("server running at 5050")
	panic(fasthttp.ListenAndServe(":5050", router.HandleRequest))
}
