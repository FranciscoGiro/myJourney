package main

import(
	"os"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/FranciscoGiro/myJourney/backend/src/routes"
	"github.com/FranciscoGiro/myJourney/backend/src/database"
)


func main() {

	fmt.Println("Vou começar aqui")

	mydir, err := os.Getwd()
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println("MY DIRECTORY:",mydir)

	err = godotenv.Load(".env")
	if err != nil {
		fmt.Printf("Error loading environment variables")
		os.Exit(1)
	}

	port := os.Getenv("PORT")

	defer database.Disconnect()

	

	app := gin.New()
	router := app.Group("/api/")
	routes.AllRoutes(router) // pass client if needed


    app.Run(":"+port)

	fmt.Println("Server started on port ", os.Getenv("SERVER_PORT"))
}