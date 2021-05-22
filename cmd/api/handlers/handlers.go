package handlers

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/shaneu/indahaus/graph"
	"github.com/shaneu/indahaus/graph/generated"
	"github.com/shaneu/indahaus/internal/data/ipresult"
	"github.com/shaneu/indahaus/internal/iplookup"
	"github.com/shaneu/indahaus/internal/mid"
	"github.com/shaneu/indahaus/pkg/auth"
)

// API binds our HTTP routes and applies our middleware and returns our http.Handler interface
func API(build string, shutdown chan os.Signal, a auth.Auth, db *sqlx.DB, log *log.Logger) http.Handler {
	e := echo.New()

	// global middlewares to be applied to each request
	e.Use(
		echo.WrapMiddleware(mid.InsertValues()),
		echo.WrapMiddleware(mid.Logger(log)),
		echo.WrapMiddleware(mid.Metrics()),
		middleware.Recover(),
	)

	basicAuth := middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		return a.Authenticate(username, password), nil
	})

	gqlResolver := graph.Resolver{
		IPResult: ipresult.New(log, db),
		IPLookup: iplookup.New(log),
	}
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &gqlResolver}))

	srv.SetRecoverFunc(func(ctx context.Context, err interface{}) error {
		return errors.New("internal server error")
	})

	gqlGrp := graphqlGroup{
		srv: srv,
	}
	e.GET("/", gqlGrp.playground, basicAuth)
	e.POST("/graphql", gqlGrp.graphql, basicAuth)

	checkGroup := checkGroup{
		build: build,
		db:    db,
	}
	e.GET("/readiness", checkGroup.readiness)
	e.GET("/liveness", checkGroup.liveness)

	return e
}
