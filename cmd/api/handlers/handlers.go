package handlers

import (
	"context"
	"errors"
	"log"
	"net/http"

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
	"github.com/vektah/gqlparser/v2/gqlerror"
)

// API binds our HTTP routes and applies our middleware and returns our http.Handler interface
func API(build string, a auth.Auth, db *sqlx.DB, log *log.Logger) http.Handler {
	e := echo.New()

	// global middlewares to be applied to each request
	e.Use(
		echo.WrapMiddleware(mid.InsertValues()),
		echo.WrapMiddleware(mid.Logger(log)),
		echo.WrapMiddleware(mid.Metrics()),
		middleware.Recover(),
	)

	customHTTPErrorHandler := func(err error, c echo.Context) {
		v := c.Request().Context().Value(mid.RequestValueKey).(*mid.RequestValues)
		log.Printf("%s : ERROR    : %v", v.TraceID, err)

		msg := map[string]string{
			"message": http.StatusText(http.StatusInternalServerError),
		}

		if !c.Response().Committed {
			err = c.JSON(http.StatusInternalServerError, msg)
			if err != nil {
				log.Printf("%s : ERROR    : %v", v.TraceID, err)
			}
		}
	}
	// global http error handling
	e.HTTPErrorHandler = customHTTPErrorHandler

	basicAuth := middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		return a.Authenticate(username, password), nil
	})

	gqlResolver := graph.Resolver{
		IPResult: ipresult.New(log, db),
		IPLookup: iplookup.New(log),
	}
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &gqlResolver}))

	// global graphql panic handling
	srv.SetRecoverFunc(func(ctx context.Context, err interface{}) error {
		return errors.New("internal server error")
	})

	// global graphql error handling
	srv.SetErrorPresenter(func(ctx context.Context, err error) *gqlerror.Error {
		v := ctx.Value(mid.RequestValueKey).(*mid.RequestValues)

		log.Printf("%s : ERROR    : %v", v.TraceID, err)

		return gqlerror.Errorf("graphql error")
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
