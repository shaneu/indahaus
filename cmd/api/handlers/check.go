package handlers

import (
	"net/http"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/shaneu/indahaus/internal/mid"
)

type checkGroup struct {
	build string
	db    *sqlx.DB
}

func (cg checkGroup) readiness(c echo.Context) error {
	v := c.Request().Context().Value(mid.RequestValueKey).(*mid.RequestValues)

	status := "ok"
	statusCode := http.StatusOK

	if err := cg.db.Ping(); err != nil {
		status = "db not ready"
		statusCode = http.StatusInternalServerError
	}

	v.StatusCode = statusCode

	health := struct {
		Status string `json:"status"`
	}{
		Status: status,
	}

	return c.JSON(statusCode, health)
}

func (cg checkGroup) liveness(c echo.Context) error {
	v := c.Request().Context().Value(mid.RequestValueKey).(*mid.RequestValues)

	statusCode := http.StatusOK
	v.StatusCode = statusCode

	host, err := os.Hostname()
	if err != nil {
		host = "unavailable"
	}

	info := struct {
		Status    string `json:"status,omitempty"`
		Build     string `json:"build,omitempty"`
		Host      string `json:"host,omitempty"`
		Pod       string `json:"pod,omitempty"`
		PodIP     string `json:"podIP,omitempty"`
		Node      string `json:"node,omitempty"`
		Namespace string `json:"namespace,omitempty"`
	}{
		Status:    "up",
		Build:     cg.build,
		Host:      host,
		Pod:       os.Getenv("KUBERNETES_PODNAME"),
		PodIP:     os.Getenv("KUBERNETES_NAMESPACE_POD_IP"),
		Node:      os.Getenv("KUBERNETES_NODENAME"),
		Namespace: os.Getenv("KUBERNETES_NAMESPACE"),
	}

	return c.JSON(statusCode, info)
}
