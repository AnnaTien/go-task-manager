package cli

import (
	"fmt"
	"net/http"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"go-task-manager/api"
	"go-task-manager/database"
	"go-task-manager/internal/task"

	"go-task-manager/internal/middleware"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	// Configure zerolog for human-readable output in development
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

// config holds the application's configuration settings.
var config struct {
	Port     string
	Database struct {
		DSN string
	}
}

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "go-task-manager",
	Short: "A simple task manager application",
}

// apiCmd is the subcommand to run the REST API server.
var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "Runs the REST API server",
	Run: func(cmd *cobra.Command, args []string) {
		// Build the DSN from environment variables.
		dbHost := os.Getenv("DB_HOST")
		dbUser := os.Getenv("DB_USER")
		dbPass := os.Getenv("DB_PASSWORD")
		dbName := os.Getenv("DB_NAME")
		dbPort := os.Getenv("DB_PORT")

		// Check if environment variables exist.
		if dbHost != "" && dbUser != "" {
			config.Database.DSN = fmt.Sprintf(
				"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
				dbHost, dbUser, dbPass, dbName, dbPort,
			)
		} else {
			// If not, read from the config file.
			viper.SetConfigName("config")
			viper.SetConfigType("yaml")
			viper.AddConfigPath("./configs")
			if err := viper.ReadInConfig(); err != nil {
				log.Fatal().Err(err).Msg("Error reading config file")
			}
			if err := viper.Unmarshal(&config); err != nil {
				log.Fatal().Err(err).Msg("Error unmarshalling config")
			}
		}

		db := database.InitDB(config.Database.DSN)
		err := db.AutoMigrate(&task.Task{})
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to auto migrate database")
		}

		storage := task.NewGormStorage(db)
		apiHandler := api.APIHandler{Storage: storage}

		r := mux.NewRouter()
		authRouter := r.PathPrefix("/").Subrouter()
		authRouter.Use(middleware.AuthMiddleware)

		// API routes
		r.HandleFunc("/tasks", apiHandler.GetTasksHandler).Methods("GET")
		r.HandleFunc("/tasks/search", apiHandler.SearchTasksHandler).Methods("GET")
		r.HandleFunc("/tasks/{id}", apiHandler.GetTaskByIDHandler).Methods("GET")
		authRouter.HandleFunc("/tasks", apiHandler.AddTaskHandler).Methods("POST")
		authRouter.HandleFunc("/tasks/{id}", apiHandler.UpdateTaskHandler).Methods("PUT")
		authRouter.HandleFunc("/tasks/{id}", apiHandler.DeleteTaskHandler).Methods("DELETE")

		// Apply the middleware to the router
		loggedRouter := middleware.LoggingMiddleware(r)

		log.Info().Msgf("API server started on port %s", config.Port)
		// Pass the wrapped router to the server
		log.Fatal().Err(http.ListenAndServe(":"+config.Port, loggedRouter)).Msg("Server failed to start")
	},
}

// Execute adds all child commands to the root command and sets flags.
func Execute() {
	rootCmd.AddCommand(apiCmd)
	apiCmd.Flags().StringVarP(&config.Port, "port", "p", "8080", "Port for the API server")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
