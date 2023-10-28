package api

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/Courtcircuits/HackTheCrous.api/storage"
	"github.com/Courtcircuits/HackTheCrous.api/types"
	"github.com/Courtcircuits/HackTheCrous.api/util"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

var server *Server

type Credentials struct {
	Mail     string `json:"mail"`
	Password string `json:"password"`
	Remember bool   `json:"remember"`
}
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
	config.AllowAllOrigins = true
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

	key := "Secret-session-key" // Replace with your SESSION_SECRET or similar
	maxAge := 86400 * 30        // 30 days
	isProd := false             // Set to true when serving over https

	storeSession := sessions.NewCookieStore([]byte(key))
	storeSession.MaxAge(maxAge)
	storeSession.Options.Path = "/"
	storeSession.Options.HttpOnly = true // HttpOnly should always be enabled
	storeSession.Options.Secure = isProd

	gothic.Store = storeSession

	goth.UseProviders(
		google.New(util.Get("GOOGLE_CLIENT_ID"), util.Get("GOOGLE_AUTH_SECRET"), util.Get("FULL_SERVER_URL")+"/auth/callback?provider=google", "email"),
	)

	return server
}

func (s *Server) Start() error {

	critical_route := s.router.Group("/")
	critical_route.Use(JWTAuth())
	critical_route.Use(GinContextMiddleware())

	s.router.POST("/login", s.Login)
	s.router.POST("/signup", s.Signup)
	s.router.POST("/logout", s.Logout)
	s.router.GET("/auth", s.GoogleAuth)
	s.router.GET("/auth/callback", s.GoogleAuthCallback)

	critical_route.POST("/graphql", s.GraphQLHandler)
	critical_route.POST("/mail/confirm", s.SendConfirmationMail)
	critical_route.POST("/mail/code", s.ConfirmMail)

	s.router.Run(s.listenAddr)
	return http.ListenAndServe(s.listenAddr, nil)
}

func (s *Server) GoogleAuth(c *gin.Context) {
	gothic.BeginAuthHandler(c.Writer, c.Request)
}

func (s *Server) GoogleAuthCallback(c *gin.Context) {
	google_user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		c.AbortWithStatusJSON(401, gin.H{
			"error": err.Error(),
		})
		return
	}
	log.Printf("Google user logged : %q\n", google_user.Email)

	credentials := Credentials{
		Mail:     google_user.Email,
		Password: "",
		Remember: true,
	}

	user, err := s.Store.GetUserByAuthCustomName(credentials.Mail)

	if err == storage.ErrUserNotFound {
		user, err = s.Store.CreateGoogleUser(credentials.Mail)
	}

	if err != nil {
		log.Println(err)
		c.AbortWithStatus(401)
		return
	}
	err = s.SendAuthToken(&user, credentials)
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(401)
		return
	}

	var mailVerified string
	if !user.Nonce.Valid {
		mailVerified = "true"
	} else {
		mailVerified = "false"
	}

	c.Redirect(301, util.Get("CLIENT_URL")+"?mailVerified="+mailVerified+"&token="+user.Auth_token.String+"&refreshToken="+user.Refresh_token.String+"&mail="+user.Email.String)
	return
}

func (s *Server) Logout(c *gin.Context) {
	c.String(200, "Logged out")
}

func (s *Server) Login(c *gin.Context) {

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

	err = s.SendAuthToken(&user, credentials)
	if err != nil {
		c.AbortWithStatus(401)
		return
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
	cred := Credentials{}

	if err := c.BindJSON(&cred); err != nil {
		fmt.Printf("error Signup : %q", err)
		c.AbortWithStatus(400)
		return
	}
	cred.Remember = true

	fmt.Printf("mail : %q, passsword :%q\n", cred.Mail, cred.Password)

	pg_storage := storage.NewPostgresDatabase()

	user, err := pg_storage.CreateLocalUser(cred.Mail, cred.Password)

	if err != nil {
		c.AbortWithStatusJSON(401, gin.H{
			"message": err.Error(),
		})
		return
	}

	s.SendAuthToken(&user, cred)

	c.JSON(200, gin.H{
		"type":         "success",
		"message":      "Signed up",
		"token":        user.Auth_token.String,
		"refreshToken": user.Refresh_token.String,
	})
}

func (s *Server) SendConfirmationMail(c *gin.Context) {
	type Payload struct {
		Student_mail string `json:"student_mail"`
	}

	var payload Payload
	if err := c.BindJSON(&payload); err != nil {
		fmt.Printf("wrong request")
		c.AbortWithStatus(400)
		return
	}

	id := c.GetInt("id")

	user, err := s.Store.GetUserByID(id)
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(401)
		return
	}
	if match, _ := regexp.MatchString("^.*@etu\\.umontpellier\\.fr$", payload.Student_mail); !match {
		c.AbortWithStatusJSON(400, gin.H{
			"message": "wrong format for student mail",
		})
		return
	}

	err = util.SendConfirmationMail(payload.Student_mail, user.Nonce.String)
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(401)
		return
	}
	err = s.Store.UpdateMail(int(user.ID.Int32), payload.Student_mail)
	c.JSON(200, gin.H{
		"type":    "Success",
		"message": "mail sent",
	})
	return
}

func (s *Server) ConfirmMail(c *gin.Context) {
	type Payload struct {
		Nonce string `json:"nonce"`
	}

	var payload Payload
	if err := c.BindJSON(&payload); err != nil {
		fmt.Printf("wrong request")
		c.AbortWithStatus(400)
	}

	id := c.GetInt("id")
	err := s.Store.ConfirmMail(id, payload.Nonce)
	if err != nil {
		c.JSON(400, gin.H{
			"type": "Error",
			"message": err,
		})
		return
	}
	c.JSON(200, gin.H{
		"type": "Success",
		"message": "Mail confirmed",
	})
}

// set auth token inside the user passed by reference
func (s *Server) SendAuthToken(user *types.User, credentials Credentials) error {
	fmt.Println(user)
	tokens, err := user.GetTokens(credentials.Remember)

	user.Auth_token = sql.NullString{String: tokens.Auth_token, Valid: true}
	if err == types.ErrRefreshTokenNeedUpdate && credentials.Remember {
		refresh_token := s.Store.UpdateRefreshToken(int(user.ID.Int32))
		user.Refresh_token = sql.NullString{String: refresh_token, Valid: true} //not gonna lie, this conversion seems wrong
	} else if err != nil {
		return errors.New("unauthorized")
	}

	if !credentials.Remember {
		user.Refresh_token = sql.NullString{
			String: "",
			Valid:  false,
		}
	}

	return nil
}
