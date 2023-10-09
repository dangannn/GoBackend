package server

import (
	"GoBackend/controllers"
	"GoBackend/models"
	"GoBackend/repositories"
	"GoBackend/services"
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
	csrf "github.com/utrack/gin-csrf"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type HttpServer struct {
	config *viper.Viper
	router *gin.Engine
}

func InitHttpServer(config *viper.Viper, dbHandler *gorm.DB) HttpServer {

	//comments layer
	commentsRepository := repositories.NewCommentRepository(dbHandler)
	commentsService := services.NewCommentService(commentsRepository)
	commentsController := controllers.NewCommentController(commentsService)

	//post layer
	postsRepository := repositories.NewPostRepository(dbHandler)
	postsService := services.NewPostService(postsRepository)
	postsController := controllers.NewPostController(postsService)

	//user layer
	usersRepository := repositories.NewUserRepository(dbHandler)
	usersService := services.NewUserService(usersRepository)
	usersController := controllers.NewUserController(usersService)

	//emails
	emailsRepository := repositories.NewEmailRepository(dbHandler)
	emailsService := services.NewEmailService(emailsRepository)
	emailsController := controllers.NewEmailController(emailsService)

	router := gin.Default()
	emailsService.TaskScheduling()

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

	//Session & CSRF
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

	//Websocket

	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			origin := r.Header.Get("Origin")
			return origin == "http://127.0.0.1:5173" || origin == "http://127.0.0.1:8080"
		},
	}
	var connections = make(map[*websocket.Conn]bool)

	router.GET("/ws", func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			return
		}

		defer conn.Close()

		connections[conn] = true

		//Redis
		redisClient := initRedis()
		defer redisClient.Close()
		// Получение значения из Redis

		key := "posts"
		var firstResponse []byte

		cachedValue, err := redisClient.Get(c, key).Result()
		if err != nil {
			posts, _ := postsService.GetAll()
			jsonData, err := json.Marshal(posts)
			if err != nil {
				log.Println(err)
				return
			}
			value := jsonData
			// Сохранение значения в Redis с заданным временем жизни
			err = redisClient.Set(c, key, value, 10*time.Minute).Err()
			if err != nil {
				fmt.Println("Error caching data:", err)
				return
			}
			firstResponse = jsonData
		} else {
			log.Println("using cached data")
			firstResponse = []byte(cachedValue)
		}

		err = conn.WriteMessage(websocket.TextMessage, firstResponse)
		if err != nil {
			log.Println(err)
			return
		}

		for {
			messageType, p, err := conn.ReadMessage()
			if err != nil {
				log.Println(err)
				return
			}

			// Обработка и анализ входящих данных от клиента
			var receivedData *models.Post

			err = json.Unmarshal(p, &receivedData)
			if err != nil {
				log.Println(err)
				return
			}

			// Выполнение вставки данных в базу данных
			//log.Println(receivedData.Id)
			//log.Println(receivedData.Id == 0)
			if receivedData.Id == 0 {
				response, responseErr := postsService.Create(receivedData)
				if responseErr != nil {
					c.AbortWithStatusJSON(responseErr.Status, responseErr)
					return
				}
				jsonData, err := json.Marshal(response)
				if err != nil {
					log.Println(err)
					return
				}

				currentValue, err := redisClient.Get(c, key).Result()
				if err != nil && err != redis.Nil {
					log.Fatalf("Failed to get current value: %v", err)
					return
				}

				var objects []*models.Post

				if currentValue != "" {
					if err := json.Unmarshal([]byte(currentValue), &objects); err != nil {
						log.Fatalf("Failed to unmarshal current value: %v", err)
					}
				}

				objects = append(objects, receivedData)

				updatedValue, err := json.Marshal(objects)
				if err != nil {
					log.Fatalf("Failed to marshal updated value: %v", err)
				}

				err = redisClient.Set(c, key, updatedValue, 10*time.Minute).Err()
				if err != nil {
					log.Fatalf("Failed to set new value: %v", err)
				}

				// Отправка данных обратно клиенту через вебсокет
				for conn := range connections {
					err := conn.WriteMessage(websocket.TextMessage, []byte("New post created"))
					if err != nil {
						fmt.Println(err)
						conn.Close()
						delete(connections, conn)
					}
					err = conn.WriteMessage(messageType, jsonData)
					if err != nil {
						fmt.Println(err)
						conn.Close()
						delete(connections, conn)
					}
				}
			} else {
				responseErr := postsService.Delete(strconv.Itoa(int(receivedData.Id)))
				if responseErr != nil {
					c.AbortWithStatusJSON(responseErr.Status, responseErr)
					return
				}
				log.Println("удалено")
				//TODO clear redis & response to client

				currentValue, err := redisClient.Get(c, key).Result()
				if err != nil && err != redis.Nil {
					log.Fatalf("Failed to get current value: %v", err)
					return
				}

				var objects []*models.Post

				if currentValue != "" {
					if err := json.Unmarshal([]byte(currentValue), &objects); err != nil {
						log.Fatalf("Failed to unmarshal current value: %v", err)
					}
				}
				for i, v := range objects {
					if v.Id == receivedData.Id {
						objects = append(objects[:i], objects[i+1:]...)
						break
					}
				}

				updatedValue, err := json.Marshal(objects)
				if err != nil {
					log.Fatalf("Failed to marshal updated value: %v", err)
				}

				err = redisClient.Set(c, key, updatedValue, 10*time.Minute).Err()
				if err != nil {
					log.Fatalf("Failed to set new value: %v", err)
				}

				for conn := range connections {
					err := conn.WriteMessage(websocket.TextMessage, []byte("Post deleted"))
					if err != nil {
						fmt.Println(err)
						conn.Close()
						delete(connections, conn)
					}
					//err = conn.WriteMessage(messageType, jsonData)
					//if err != nil {
					//	fmt.Println(err)
					//	conn.Close()
					//	delete(connections, conn)
					//}
				}
			}
		}
	})

	router.GET("/ws/post/:id/comments", func(c *gin.Context) {
		id := c.Param("id")
		log.Println(id)
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			return
		}

		defer conn.Close()

		connections[conn] = true
		comments, _ := postsService.GetComments(id)
		jsonData, err := json.Marshal(comments)
		if err != nil {
			log.Println(err)
			return
		}
		err = conn.WriteMessage(websocket.TextMessage, jsonData)
		if err != nil {
			log.Println(err)
			return
		}
		for {
			messageType, p, err := conn.ReadMessage()
			if err != nil {
				log.Println(err)
				return
			}

			type TmpComment struct {
				Id        uint
				Text      string
				PostId    string
				Approved  bool
				AuthorId  uint
				CreatedAt time.Time
			}
			var tmp TmpComment

			err = json.Unmarshal(p, &tmp)
			if err != nil {
				log.Println(err)
				return
			}
			tmpPostId, err := strconv.Atoi(tmp.PostId)
			if err != nil {
				fmt.Println(err)
				c.Abort()
				return
			}
			receivedData := models.Comment{
				Id:        tmp.Id,
				Text:      tmp.Text,
				PostId:    uint(tmpPostId),
				Approved:  tmp.Approved,
				AuthorId:  tmp.AuthorId,
				CreatedAt: tmp.CreatedAt,
			}
			log.Println("измененные", receivedData)

			// Выполнение вставки данных в базу данных
			response, responseErr := commentsService.Create(&receivedData)
			if responseErr != nil {
				c.AbortWithStatusJSON(responseErr.Status, responseErr)
				return
			}
			jsonData, err := json.Marshal(response)
			if err != nil {
				log.Println(err)
				return
			}

			// Отправка данных обратно клиенту через вебсокет
			for conn := range connections {
				err := conn.WriteMessage(websocket.TextMessage, []byte("New comment created"))
				if err != nil {
					fmt.Println(err)
					conn.Close()
					delete(connections, conn)
				}
				err = conn.WriteMessage(messageType, jsonData)
				if err != nil {
					fmt.Println(err)
					conn.Close()
					delete(connections, conn)
				}
			}
		}
	})

	//API
	api := router.Group("/api")

	//email routes
	api.GET("/email", emailsController.SendEmail)
	api.GET("/add_view", emailsController.AddView)
	api.GET("/add_new_comment", emailsController.AddNewComment)

	//CRUD post & get comments & pagination
	api.POST("/post", postsController.Create)
	api.GET("/post/:id", postsController.GetById)
	api.GET("/post/:id/delete", postsController.Delete)
	api.GET("/post/:id/update", postsController.Update)
	api.GET("/post/all", postsController.GetAll)
	api.GET("/post/:id/comments", postsController.GetComments)
	api.GET("/posts/:page/page", postsController.GetPage)
	//api.POST("/post/:id/checked", postsController.AddCheck)

	//CRUD comments
	//api.GET("/comment", commentsController.CreateComment)
	//api.GET("/comment", commentsController.GetAllComment)
	//api.GET("/comment", commentsController.GetCommentById)
	//api.GET("/comment/:id", usersController.DeleteComment)
	//api.GET("/comment/:id", usersController.UpdateComment)
	api.GET("/comment/unapproved", commentsController.GetAllUnapproved)
	api.POST("/comment/:id/moderate", commentsController.Moderate)

	//CRUD users
	api.GET("/user/:id", usersController.GetUserById)
	api.GET("/user", usersController.GetAllUsers)
	api.POST("/register", usersController.CreateUser)
	api.POST("/login", usersController.LoginUser)
	api.GET("/user/:id/posts", usersController.GetUserPosts)
	//api.GET("/user/:id", usersController.DeletePost)
	//api.GET("/user/:id", usersController.UpdatePost)

	userRoutes := api.Group("/user").Use(AuthUser())
	{
		userRoutes.GET("/user/:id", usersController.GetUserById)
	}

	adminRoutes := api.Group("/admin").Use(AuthAdmin())
	{
		adminRoutes.GET("/user/:id", usersController.GetUserById)
		adminRoutes.GET("/post", postsController.GetAll)
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
