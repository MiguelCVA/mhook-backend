package main

import (
	"github.com/MiguelCVA/mhook-backend/internal/config"
	"github.com/MiguelCVA/mhook-backend/internal/database"
	"github.com/MiguelCVA/mhook-backend/internal/server"
)

func main() {
	/**########## DOTENV ##########**/
	config.ConfigDotenv()

	/**########## DATABASE ##########**/
	database.StartDB()

	/**########## SERVER ##########**/
	server := server.NewServer()
	server.Run()
}
