package api

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/Courtcircuits/HackTheCrous.api/storage"
	"github.com/Courtcircuits/HackTheCrous.api/types"
	"github.com/gin-gonic/gin"
)

type Server struct {
	listenAddr string
	router     *gin.Engine
	store      storage.PostgresDatabase
}

func NewServer(listenAddr string, store storage.PostgresDatabase) *Server {
	return &Server{
		listenAddr: listenAddr,
		store:      store,
		router:     gin.Default(),
	}
}

func (s *Server) Start() error {
	s.router.POST("/login", s.Login)
	s.router.POST("/signup", s.Signup)

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

	user, err := s.store.GetUserByEmail(credentials.Mail)

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

	if err == types.ErrRefreshTokenNeedUpdate {
		refresh_token := s.store.UpdateRefreshToken(int(user.ID.Int32))
		user.Refresh_token = sql.NullString{String: refresh_token, Valid: true} //not gonna lie, this conversion seems wrong
	} else if err != nil {
		c.AbortWithStatus(401)
		return
	}

	fmt.Println(user.Auth_token)
	fmt.Println(user.Refresh_token)

	c.JSON(200, gin.H{
		"type":         "success",
		"message":      "Logged in",
		"token":        user.Auth_token.String,
		"refreshToken": user.Refresh_token.String,
		"mail":         user.Email.String,
	})
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
			refresh_token := s.store.UpdateRefreshToken(int(user.ID.Int32))
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
