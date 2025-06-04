package main

import "go.uber.org/zap"

type application struct {
	config config
	logger *zap.SugaredLogger
}

type config struct {
	port string
	env  string
}

func main() {

	cfg := &application{
		config: config{
			port: "8080",
			env:  "development",
		},
	}

	logger := zap.NewExample().Sugar()
	defer logger.Sync()

	app := &application{
		config: cfg.config,
		logger: logger,
	}

	mux := app.routes()
	logger.Fatal(app.server(mux))
}
