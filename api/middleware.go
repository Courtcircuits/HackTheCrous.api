package api

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/Courtcircuits/HackTheCrous.api/util"
	"github.com/gin-gonic/gin"
)

func AuthRewriteMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		cookies := ctx.Request.Cookies()
		var token string
		for _, cookie := range cookies {
			if strings.HasPrefix(cookie.Value, "token=") {
				token = strings.TrimPrefix(cookie.Value, "token=")
			}
		}
		ctx.Request.Header.Del("Authorization")
		ctx.Request.Header.Del("Cookie")
		ctx.Request.Header.Set("Authorization", "Bearer "+token)
		ctx.Next()
	}
}

func JWTAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		auth_header := ctx.GetHeader("Authorization")

		if auth_header == "" {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "Auth header is empty",
			})
			return
		}

		parts := strings.Split(auth_header, " ")

		if len(parts) != 2 || parts[0] != "Bearer" {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "Auth header must provide a bearer token",
			})
			return
		}

		parsed_token, err := util.VerifyJWT(parts[1])

		if err != nil {
			log.Print(err)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": err.Error(),
			})
			return
		}

		id := int(parsed_token["iduser"].(float64))

		log.Printf("Authenticated %d, %q\n", id, parsed_token["mail"])

		ctx.Set("id", id)
		ctx.Set("mail", parsed_token["mail"])

		ctx.Next()
	}
}

// enable to retrieve gin context in graphql resolvers
func GinContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), "GinContext", c)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

// WARNING : not a middleware
func GetGinContext(ctx context.Context) (*gin.Context, error) {
	gCtx := ctx.Value("GinContext")
	if gCtx == nil {
		return nil, errors.New("could not retive gin.Context")
	}

	gc, ok := gCtx.(*gin.Context)
	if !ok {
		return nil, errors.New("gin.Context has wrong type")
	}
	return gc, nil
}
