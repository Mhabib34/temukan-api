package helper

import (
	"os"

	"github.com/gin-gonic/gin"
)

func isProduction() bool {
	return os.Getenv("APP_ENV") == "production"
}

func SetAuthCookies(ctx *gin.Context, accessToken, refreshToken string) {
	secure := isProduction()
	ctx.SetCookie("access_token", accessToken, 24*60*60, "/", "", secure, true)
	ctx.SetCookie("refresh_token", refreshToken, 7*24*60*60, "/", "", secure, true)
}

func SetAccessTokenCookie(ctx *gin.Context, accessToken string) {
	secure := isProduction()
	ctx.SetCookie("access_token", accessToken, 24*60*60, "/", "", secure, true)
}

func ClearAuthCookies(ctx *gin.Context) {
	secure := isProduction()
	ctx.SetCookie("access_token", "", -1, "/", "", secure, true)
	ctx.SetCookie("refresh_token", "", -1, "/", "", secure, true)
}
