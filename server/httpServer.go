package server

import (
	"GoBackend/controllers"
	"GoBackend/repositories"
	"GoBackend/services"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
	"log"
	"time"

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
	//router.GET("/user/:id", usersController.GetUserById)
	router.POST("/register", usersController.CreateUser)
	router.POST("/login", usersController.LoginUser)

	secured := router.Group("/secured").Use(Auth())
	{
		secured.GET("/user/:id", usersController.GetUserById)
	}

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

func Auth() gin.HandlerFunc {
	return func(context *gin.Context) {
		type CustomClaims struct {
			Role string `json:"claims"`
			jwt.RegisteredClaims
		}
		tokenString := context.GetHeader("Authorization")
		if tokenString == "" {
			context.JSON(401, gin.H{"error": "request does not contain an access token"})
			context.Abort()
			return
		}
		token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte("secrete-key"), nil
		}, jwt.WithLeeway(2*time.Second))

		if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
			log.Printf(claims.Role)
		} else {
			log.Println(err)
		}
		if err != nil {
			context.JSON(401, gin.H{"error": err.Error()})
			context.Abort()
			return
		}
		context.Next()
	}
}
