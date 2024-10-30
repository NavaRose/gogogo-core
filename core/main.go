package core

import (
	"github.com/NavaRose/gogogo-core/constants"
	"github.com/NavaRose/gogogo-core/database"
	"github.com/NavaRose/gogogo-core/exception"
	"github.com/NavaRose/gogogo-core/middleware"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type Interface interface{}

func InitEngine(RouteCreator func(engine *gin.Engine)) *gin.Engine {
	// Load environment variables
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Set app mode
	appMode := gin.DebugMode
	if os.Getenv("GIN_MODE") != "" {
		appMode = os.Getenv("GIN_MODE")
	}
	gin.SetMode(appMode)

	// Init engine
	engine := gin.New()

	// Setup middleware
	engine.Use(middleware.HandleCORS())
	engine.Use(middleware.ErrorHandler())

	// Set trusted proxies
	trustedProxies := []string{"127.0.0.1"}
	if os.Getenv("TRUSTED_PROXIES") != "" {
		trustedProxies = strings.Split(os.Getenv("TRUSTED_PROXIES"), ",")
	}

	//
	err = engine.SetTrustedProxies(trustedProxies)
	if err != nil {
		log.Fatal("Error setting trusted proxies")
	}

	if os.Getenv("IS_USE_LOGGER") == "true" {
		engine.Use(gin.Logger())
		engine.Use(gin.Recovery())
	}

	RouteCreator(engine)
	// Set middlewares
	return engine
}

func AutoMigrate(schemas []Interface) {
	db := database.InitDatabaseWithoutEngine()
	if os.Getenv("ALLOW_AUTO_MIGRATE") == "true" {
		//model.Migrate(db)
		for _, schema := range schemas {
			err := db.AutoMigrate(&schema)
			if err != nil {
				ErrorHandle(err)
			}
		}
	}
	database.CloseDatabase(db)
}

func ErrorChecking(err error) bool {
	if err != nil {
		ErrorHandle(err)
		return true
	}
	return false
}

func CreateError(statusCode int, message string, metadata string) exception.Http {
	return exception.NewHttpError(message, metadata, statusCode)
}

func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, exist := ctx.Get(constants.TokenKey)
		if exist == false || token == nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Invalid token",
				"data":    nil,
			})
			ctx.Abort()
			return
		}
	}
}

func SetTokenCookie(accessToken string, ctx *gin.Context) {
	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.SetCookie(constants.TokenKey, accessToken, int(time.Now().Add(time.Hour*24).Unix()), "", "", false, false)
}

func DestroyTokenCookie(ctx *gin.Context) {
	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.SetCookie(constants.TokenKey, "", -1, "", "", false, false)
}
