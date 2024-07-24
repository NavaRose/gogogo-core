package core

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"strings"
)

func InitEngine(RouteCreator func() func(engine *gin.Engine)) *gin.Engine {
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

	// Set trusted proxies
	trustedProxies := []string{"127.0.0.1"}
	if os.Getenv("TRUSTED_PROXIES") != "" {
		trustedProxies = strings.Split(os.Getenv("TRUSTED_PROXIES"), ",")
	}
	err = engine.SetTrustedProxies(trustedProxies)
	if err != nil {
		log.Fatal("Error setting trusted proxies")
	}

	if os.Getenv("IS_USE_LOGGER") == "true" {
		engine.Use(gin.Logger())
		engine.Use(gin.Recovery())
	}

	InitRoute(engine, RouteCreator())
	return engine
}

func InitRoute(engine *gin.Engine, routeCreator func(*gin.Engine)) {
	routeCreator(engine)
}

func AutoMigrate(schemas []interface{}) {
	db := InitDatabaseWithoutEngine()
	if os.Getenv("ALLOW_AUTO_MIGRATE") == "true" {
		//model.Migrate(db)
		for _, schema := range schemas {
			err := db.AutoMigrate(&schema)
			if err != nil {
				ErrorHandle(err)
			}
		}
	}
	CloseDatabase(db)
}

func InitDatabaseWithoutEngine() *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s user=%s dbname=%s sslmode=disable password=%s TimeZone=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("TIME_ZONE"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Error opening database connection")
	}

	return db
}

func CloseDatabase(db *gorm.DB) {
	defer func() {
		sqlDB, _ := db.DB()
		if err := sqlDB.Close(); err != nil {
			log.Fatal("Failed to close database connection: ", err)
		}
	}()
}
