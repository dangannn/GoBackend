package server

import (
	"GoBackend/controllers"
	"GoBackend/repositories"
	"GoBackend/services"
	"encoding/json"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
	csrf "github.com/utrack/gin-csrf"
	"gorm.io/gorm"
	"log"
	"net/http"
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
	usersController := controllers.NewUserController(usersService)

	emailsRepository := repositories.NewEmailRepository(dbHandler)
	emailsService := services.NewEmailService(emailsRepository)
	emailsController := controllers.NewEmailController(emailsService)

	router := gin.Default()
	// Middleware CORS
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "http://127.0.0.1:5173")
		c.Header("Access-Control-Allow-Methods", "*")
		c.Header("Access-Control-Allow-Headers", "Authorization, Origin, Content-Type, Cookie, X-CSRF-Token, Access-Control-Allow-Headers, Origin,Accept, X-Requested-With, Content-Type, Access-Control-Request-Method, Access-Control-Request-Headers")
		c.Header("Access-Control-Allow-Credentials", "true") // Разрешить отправку кук с запросом

		// Позволяем предварительные запросы (preflight) OPTIONS
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}

		// Продолжаем выполнение обработчика
		c.Next()
	})
	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("mysession", store))
	router.Use(csrf.Middleware(csrf.Options{
		Secret: "secrete-key",
		ErrorFunc: func(c *gin.Context) {
			c.String(400, "CSRF token mismatch")
			c.Abort()
		},
	}))

	router.GET("/protected", func(c *gin.Context) {
		c.String(200, csrf.GetToken(c))
	})

	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			origin := r.Header.Get("Origin")
			return origin == "http://127.0.0.1:5173" || origin == "http://127.0.0.1:8080"
		},
	}

	router.GET("/ws", func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			return
		}
		defer conn.Close()
		posts, _ := postsService.GetAllPosts()
		jsonData, err := json.Marshal(posts)
		if err != nil {
			log.Println(err)
			return
		}
		//err = conn.WriteMessage(websocket.TextMessage, jsonData)
		if err != nil {
			log.Println(err)
			return
		}
		i := 0
		for i < 2 {
			i++
			err = conn.WriteMessage(websocket.TextMessage, jsonData)
			time.Sleep(time.Second * 5)
		}
	})

	api := router.Group("/api")

	//CRUD post
	api.POST("/post", postsController.CreatePost)
	api.GET("/email", emailsController.SendEmail)
	//api.GET("/post", postsController.RetrieveAllPosts)
	//get all comments related to one post
	api.GET("/post/:id/comments", postsController.GetComments)
	//post's pagination
	api.GET("/posts/:page", postsController.GetPostPage)

	api.GET("/user", usersController.GetAllUsers)
	//router.GET("/user/:id", usersController.GetUserById)
	api.POST("/register", usersController.CreateUser)

	api.POST("/login", usersController.LoginUser)

	api.GET("/user/:id", usersController.GetUserById)
	api.GET("/user/:id/posts", usersController.GetUserPosts)
	userRoutes := api.Group("/user").Use(AuthUser())
	{
		userRoutes.GET("/user/:id", usersController.GetUserById)
	}

	adminRoutes := api.Group("/admin").Use(AuthAdmin())
	{
		adminRoutes.GET("/user/:id", usersController.GetUserById)
		adminRoutes.GET("/post", postsController.GetAllPosts)
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

func AuthUser() gin.HandlerFunc {
	return func(context *gin.Context) {
		type CustomClaims struct {
			Role string `json:"role"`
			Id   int    `json:"id"`
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

		if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid && claims.Role == "user" {
			log.Println("role" + claims.Role)
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

func AuthAdmin() gin.HandlerFunc {
	return func(context *gin.Context) {
		type CustomClaims struct {
			Role string `json:"role"`
			Id   int    `json:"id"`
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

		if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid && claims.Role == "admin" {
			log.Println("role" + claims.Role)
		} else {
			context.JSON(401, gin.H{"error": "wrong role"})
			context.Abort()
		}
		if err != nil {
			context.JSON(401, gin.H{"error": err.Error()})
			context.Abort()
			return
		}
		context.Next()
	}
}
