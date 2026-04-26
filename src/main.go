package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	_ "time/tzdata"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lmittmann/tint"
	"golang.org/x/crypto/bcrypt"
	"wynglet.chimbori.dev/conf"
	"wynglet.chimbori.dev/core"
	"wynglet.chimbori.dev/dashboard"
	"wynglet.chimbori.dev/db"
	"wynglet.chimbori.dev/embedfs"
	"wynglet.chimbori.dev/github"
	"wynglet.chimbori.dev/linkpreviews"
	"wynglet.chimbori.dev/qrcode"
	"wynglet.chimbori.dev/rating"
	"wynglet.chimbori.dev/slogdb"
)

func main() {
	tintHandler := tint.NewHandler(os.Stderr, &tint.Options{TimeFormat: "2006-01-02 15:04:05.000"})
	slog.SetDefault(slog.New(tintHandler))
	slog.Info(conf.AppName, "build-timestamp", conf.BuildTimestamp)

	bcryptFlag := flag.Bool("bcrypt", false, "print bcrypt hash for given password & exit")
	healthCheckFlag := flag.Bool("healthcheck", false, "verify health of running service & exit")
	configYmlFlag := flag.String("config", "wynglet.yml", "path to wynglet.yml")
	flag.Parse()

	// If run with “--bcrypt”, read a password via the terminal, output a bcrypt hash, and exit.
	if *bcryptFlag {
		password := core.ReadPassword()
		bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			slog.Error("Failed to generate password hash", tint.Err(err))
			os.Exit(1)
		}
		fmt.Println(password)
		fmt.Println(string(bytes))
		os.Exit(0)
	}

	// Read config before any routine maintenance is performed.
	var err error
	if conf.Config, err = conf.ReadConfig(*configYmlFlag); err != nil {
		slog.Error("Failed to parse config", tint.Err(err))
		os.Exit(1)
	}

	if *healthCheckFlag {
		os.Exit(core.VerifyHealthCheck(conf.Config.Web.Port))
	}

	// If debug mode was turned on in the config file, print logs at DEBUG or above.
	if conf.Config.Debug {
		tintHandler = tint.NewHandler(os.Stderr, &tint.Options{
			Level:      slog.LevelDebug,
			TimeFormat: "2006-01-02 15:04:05.000",
		})
		slog.SetDefault(slog.New(tintHandler))
	}

	if conf.Config.Database.Url == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	// Run migrations using [database/sql] before connecting to the DB using [pgxpool.Pool].
	if err := core.RunMigrations(conf.Config.Database.Url, db.EmbedMigrations); err != nil {
		slog.Error("Error running critical migrations", tint.Err(err))
		os.Exit(1)
	}
	db.Pool, err = pgxpool.New(context.Background(), conf.Config.Database.Url)
	if err != nil {
		slog.Error("Unable to connect to database", tint.Err(err))
		os.Exit(1)
	}
	slog.Info("Connected to database successfully")

	// Now that the database is connected, wrap the console handler with the DB handler
	// so that all error-level logs are also written to the database.
	slog.SetDefault(slog.New(slogdb.NewDBHandler(tintHandler, db.Pool)))
	slog.Info("Database error logging enabled")

	// Set up a graceful cleanup for when the process is terminated.
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	// Set up the Web server and start serving.
	mux := http.NewServeMux()
	core.SetupHealthCheck(mux)
	core.ServeWebManifest(mux, conf.AppName, "/dashboard", "#2575fc")
	embedfs.ServeStaticFS(mux)
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, req *http.Request) {
		IndexTempl().Render(req.Context(), w)
	})
	linkpreviews.Init(mux)
	qrcode.Init(mux)
	github.Init(mux)
	rating.Init(mux)
	dashboard.Init(mux)

	// Set up cron task for routine maintenance.
	go func() {
		// Do a one-off cleanup before scheduling a recurring task.
		performMaintenance()
		ticker := time.Tick(2 * time.Hour)
		for {
			<-ticker
			performMaintenance()
		}
	}()

	addr := net.JoinHostPort("", strconv.Itoa(conf.Config.Web.Port))
	server := &http.Server{
		Addr:    addr,
		Handler: core.SecurityHeaders(mux),
	}

	go func() {
		<-signalCh
		fmt.Println()
		slog.Info("Shutting down…")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			slog.Error("HTTP server shutdown error", tint.Err(err))
		}
		db.Pool.Close()
		slog.Info("Shutdown successfully!")
		os.Exit(0)
	}()

	slog.Info("Listening", "url", "http://localhost"+addr) // Not "https://", since this app does not terminate SSL.
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
