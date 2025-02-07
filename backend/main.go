package main

import (
	"github.com/jjPlusPlus/task-manager/backend/api"
	"github.com/jjPlusPlus/task-manager/backend/config"
	"github.com/jjPlusPlus/task-manager/backend/jobs"
	"github.com/jjPlusPlus/task-manager/backend/logging"
	"github.com/jjPlusPlus/task-manager/backend/migrations"
	"github.com/jjPlusPlus/task-manager/backend/utils"
	"github.com/rs/zerolog/log"
)

// @title           General Task API
// @version         0.1
// @description     Making knowledge workers more productive
// @termsOfService  https://resonant-kelpie-404a42.netlify.app/terms-of-service

// @contact.name   Support
// @contact.email  support@resonant-kelpie-404a42.netlify.app

// @host      localhost:8080
// @BasePath  /
func main() {
	env := config.GetEnvironment()
	utils.ConfigureLogger(env)
	log.Info().Msgf("Starting server in %s environment", env)
	// TODO: Validate .env/config at server startup

	err := migrations.RunMigrations("migrations")
	logger := logging.GetSentryLogger()
	if err != nil {
		logger.Error().Err(err).Msg("error running migrations")
	}
	apiStruct, dbCleanup := api.GetAPIWithDBCleanup()
	defer dbCleanup()
	scheduler, err := jobs.GetScheduler()
	if err != nil {
		logger.Error().Err(err).Msg("error getting job scheduler")
	} else {
		scheduler.StartAsync()
	}
	err = api.GetRouter(apiStruct).Run()
	if err != nil {
		logger.Error().Err(err).Msg("error running router")
	}
}
