package main

import (
	"fmt"
	"log"

	"github.com/alfiannnas/go-gin-starter/internals/config"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	server()
}

func server() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Println("Error loading .env file:", err)
	}

	supabaseConfig, err := config.InitSupabaseClient()
	if err != nil {
		log.Fatalf("Failed to initialize Supabase: %v", err)
	}

	supabaseClient := supabaseConfig.GetClient()
	fmt.Println("Supabase client successfully initialized:", supabaseClient != nil)

	g := gin.Default()
	g.Use(gin.Recovery())

	// usersGroup := g.Group("/users")
	// gorm := config.NewGormPostgres()
	// userRepo := repository.NewUserQuery(gorm)
	// userSvc := service.NewUserService(userRepo)
	// userHdl := handlers.NewUserHandler(userSvc)
	// userRouter := routes.NewUserRouter(usersGroup, userHdl, gorm)
	// userRouter.Mount()

	g.Run(":8080")
}
