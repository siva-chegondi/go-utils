package firebase

import (
	"context"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	"github.com/siva-chegondi/go-utils/logger"
	"google.golang.org/api/option"
	"net/http"
	"strings"
)

var log = logger.DefaultLogger
var (
	app        *firebase.App
	AuthClient *auth.Client
)

// InitFirebase initializes Firebase Admin SDK
func InitFirebase() error {
	var err error

	// Deployment should mount account info
	// as secret having the following information.
	opt := option.WithCredentialsFile("/firebase/service-account.json")
	app, err = firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return err
	}

	AuthClient, err = app.Auth(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func VerifyToken(ctx *gin.Context) {
	bearerToken := ctx.Request.Header.Get("Authorization")
	bearerToken = strings.TrimPrefix(bearerToken, "Bearer ")

	token, err := AuthClient.VerifyIDToken(ctx, bearerToken)
	if err != nil {
		log.Error().Err(err).Msg("Failed to verify token")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Debug().Str("uid", token.UID).Msg("User verified")
	ctx.Set("uid", token.UID)
	ctx.Next()
}
