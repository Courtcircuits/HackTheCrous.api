package api

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/Courtcircuits/HackTheCrous.api/storage"
	"github.com/Courtcircuits/HackTheCrous.api/types"
	"github.com/Courtcircuits/HackTheCrous.api/util"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var server *Server

type Server struct {
	listenAddr     string
	router         *gin.Engine
	graphqlHandler *handler.Server
	Store          storage.PostgresDatabase
	Cache          *storage.RedisCache
}

func GetServer() *Server {
	return server
}

func NewServer(listenAddr string, store storage.PostgresDatabase, cache *storage.RedisCache, h *handler.Server) *Server {
	r := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{util.Get("CLIENT_URL")}
	config.AllowMethods = []string{"GET", "POST"}
	config.AllowHeaders = []string{"Authorization", "Content-Type", "Origin"}

	r.Use(cors.New(config))

	server = &Server{
		listenAddr:     listenAddr,
		Store:          store,
		Cache:          cache,
		router:         r,
		graphqlHandler: h,
	}

	return server
}

func (s *Server) Start() error {
	critical_route := s.router.Group("/")
	critical_route.Use(JWTAuth())
	critical_route.Use(GinContextMiddleware())

	s.router.POST("/login", s.Login)
	s.router.POST("/signup", s.Signup)
	critical_route.POST("/graphql", s.GraphQLHandler)

	s.router.Run(s.listenAddr)
	return http.ListenAndServe(s.listenAddr, nil)
}

func (s *Server) Login(c *gin.Context) {
	type Credentials struct {
		Mail     string `json:"mail"`
		Password string `json:"password"`
		Remember bool   `json:"remember"`
	}

	var credentials Credentials

	if err := c.BindJSON(&credentials); err != nil {
		c.AbortWithStatus(400)
		return
	}

	user, err := s.Store.GetUserByEmail(credentials.Mail)

	log.Printf("user auth_token : %q\n", user.Auth_token.String)

	if err != nil {
		c.AbortWithStatus(401)
		return
	}

	if !user.CheckPassword(credentials.Password) {
		c.AbortWithStatus(401)
		return
	}

	tokens, err := user.GetTokens(credentials.Remember)

	user.Auth_token = sql.NullString{String: tokens.Auth_token, Valid: true}

	if err == types.ErrRefreshTokenNeedUpdate && credentials.Remember {
		refresh_token := s.Store.UpdateRefreshToken(int(user.ID.Int32))
		user.Refresh_token = sql.NullString{String: refresh_token, Valid: true} //not gonna lie, this conversion seems wrong
	} else if err != nil {
		c.AbortWithStatus(401)
		return
	}

	if !credentials.Remember {
		user.Refresh_token = sql.NullString{
			String: "",
			Valid:  false,
		}
	}

	c.JSON(200, gin.H{
		"type":         "success",
		"message":      "Logged in",
		"token":        user.Auth_token.String,
		"refreshToken": user.Refresh_token.String,
		"mail":         user.Email.String,
	})
}

func (s *Server) GraphQLHandler(c *gin.Context) {

	s.graphqlHandler.ServeHTTP(c.Writer, c.Request)
}

func (s *Server) Signup(c *gin.Context) {
	type Credentials struct {
		Mail     string `json:"mail"`
		Password string `json:"password"`
	}

	cred := Credentials{}

	if err := c.BindJSON(&cred); err != nil {
		fmt.Printf("error Signup : %q", err)
		c.AbortWithStatus(400)
		return
	}

	fmt.Printf("mail : %q, passsword :%q\n", cred.Mail, cred.Password)

	pg_storage := storage.NewPostgresDatabase()

	user, err := pg_storage.CreateUser(cred.Mail, cred.Password)

	if err != nil {
		c.AbortWithStatusJSON(401, gin.H{
			"message": err.Error(),
		})
		return
	}

	tokens, err := user.GetTokens(true)

	log.Printf("generated auth token : %q \n", tokens.Auth_token)

	if err != nil {
		log.Printf("error caught : %q", err)
		if err == types.ErrRefreshTokenNeedUpdate {
			refresh_token := s.Store.UpdateRefreshToken(int(user.ID.Int32))
			user.Refresh_token = sql.NullString{String: refresh_token, Valid: true}
		} else {
			c.AbortWithStatus(500)
		}
	}

	c.JSON(200, gin.H{
		"type":         "success",
		"message":      "Signed up",
		"token":        tokens.Auth_token,
		"refreshToken": user.Refresh_token.String,
	})
}
