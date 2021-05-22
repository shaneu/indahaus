package handlers

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/labstack/echo/v4"
)

type graphqlGroup struct {
	srv *handler.Server
}

func (*graphqlGroup) playground(c echo.Context) error {
	h := playground.Handler("GraphQL playground", "/graphql")
	h.ServeHTTP(c.Response(), c.Request())
	return nil
}

func (g *graphqlGroup) graphql(c echo.Context) error {
	g.srv.ServeHTTP(c.Response(), c.Request())
	return nil
}
