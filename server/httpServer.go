package server

import (
	"GoBackend/controllers"
	"GoBackend/repositories"
	"GoBackend/services"
	"gorm.io/gorm"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type HttpServer struct {
	config *viper.Viper
	router *gin.Engine
}

func InitHttpServer(config *viper.Viper, dbHandler *gorm.DB) HttpServer {
	//post layer
	postsRepository := repositories.NewPostRepository(dbHandler)
	postsService := services.NewPostService(postsRepository)
	postsController := controllers.NewPostController(postsService)

	//user layer
	usersRepository := repositories.NewUserRepository(dbHandler)
	usersService := services.NewUserService(usersRepository)
	usersController := controllers.NewUsersController(usersService)

	router := gin.Default()
	//CRUD post
	router.POST("/post", postsController.CreatePost)
	router.GET("/post", postsController.RetrieveAllPosts)
	//get all comments related to one post
	router.GET("/post/:id/comments", postsController.GetComments)
	//post's pagination
	router.GET("/posts/:page", postsController.GetPostPage)

	router.GET("/user", usersController.GetAllUsers)
	router.POST("/register", usersController.CreateUser)

	return HttpServer{
		config: config,
		router: router,
	}
}

func (hs HttpServer) Start() {
	err := hs.router.Run(hs.config.GetString("http.server_address"))
	if err != nil {
		log.Fatalf("Error while starting HTTP server: %v", err)
	}
}
