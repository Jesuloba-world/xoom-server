package command

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humabunrouter"
	"github.com/redis/go-redis/v9"
	"github.com/rs/cors"
	"github.com/uptrace/bun"
	"github.com/uptrace/bunrouter"
	"github.com/uptrace/bunrouter/extra/reqlog"
	"github.com/urfave/cli/v2"

	activemeetings "github.com/Jesuloba-world/xoom-server/apps/activeMeetings"
	"github.com/Jesuloba-world/xoom-server/apps/cloudinary"
	logto "github.com/Jesuloba-world/xoom-server/apps/logtoApp"
	meetingservice "github.com/Jesuloba-world/xoom-server/services/meetingService"
	signallingserver "github.com/Jesuloba-world/xoom-server/services/signallingServer"
	userservice "github.com/Jesuloba-world/xoom-server/services/userService"
	"github.com/Jesuloba-world/xoom-server/util"

)

func HttpCommand(db *bun.DB, rdb *redis.Client) *cli.Command {
	return &cli.Command{
		Name:  "serve",
		Usage: "Start the HTTP server",
		Action: func(c *cli.Context) error {
			return startHTTPServer(db, rdb)
		},
	}
}

func startHTTPServer(db *bun.DB, rdb *redis.Client) error {
	port := "10001"
	config, err := util.GetConfig()
	if err != nil {
		slog.Error("Error reading config", "error", err)
		os.Exit(1)
	}
	humaConfig := huma.DefaultConfig("Xoom API", "1.0.0")
	router := bunrouter.New(
		bunrouter.Use(reqlog.NewMiddleware()),
	)

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000", "http://192.168.202.180:3000", "https://xoom-ui-development.up.railway.app"},
		AllowedMethods: []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders: []string{"Link"},
	})

	handler := corsMiddleware.Handler(router)

	// router.Use(middleware.ErrorLogger)
	api := humabunrouter.New(router, humaConfig)

	cloudinary, err := cloudinary.NewCloudinaryApp(config.CloudinaryApiKey, config.CloudinaryApiSecret, config.CloudinaryCloudName)
	if err != nil {
		slog.Error("Error initializing cloudinary", "error", err)
		os.Exit(1)
	}

	logto, err := logto.NewLogtoApp(config.LogtoEndpoint, config.LogtoApplicationId, config.LogtoApplicationSecret, cloudinary, api)
	if err != nil {
		slog.Error("Error initializing logto", "error", err)
		os.Exit(1)
	}

	activeMeetings := activemeetings.NewActiveMeetingService(rdb)

	user := userservice.NewUserService(api, logto)
	user.RegisterRoutes()

	meeting := meetingservice.NewMeetingService(api, logto, activeMeetings)
	meeting.RegisterRoutes()

	signalServer := signallingserver.NewSignallingServer(rdb, activeMeetings, logto)
	signalServer.RegisterRoute(router)

	slog.Info("Server starting", "port", port)
	err = http.ListenAndServe(":"+port, handler)
	if err != nil {
		slog.Error("Server failed to start", "error", err)
		os.Exit(1)
	}
	return nil
}
