package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"golang.org/x/sync/errgroup"
	"net/http"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
	"time"

	"github.com/labstack/echo-contrib/echoprometheus"
	sharedconfig "github.com/otterize/network-mapper/src/shared/config"
	"github.com/otterize/network-mapper/src/sniffer/pkg/mapperclient"
	"github.com/otterize/network-mapper/src/sniffer/pkg/sniffer"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	if viper.GetBool(sharedconfig.DebugKey) {
		logrus.SetLevel(logrus.DebugLevel)
	}
	mapperClient := mapperclient.NewMapperClient(viper.GetString(sharedconfig.MapperApiUrlKey))

	healthServer := echo.New()
	healthServer.GET("/healthz", func(c echo.Context) error {
		err := mapperClient.Health(c.Request().Context())
		if err != nil {
			return err
		}
		return c.NoContent(http.StatusOK)
	})

	metricsServer := echo.New()

	metricsServer.GET("/metrics", echoprometheus.NewHandler())
	errgrp, errGroupCtx := errgroup.WithContext(signals.SetupSignalHandler())
	errgrp.Go(func() error {
		return metricsServer.Start(fmt.Sprintf(":%d", viper.GetInt(sharedconfig.PrometheusMetricsPortKey)))
	})
	errgrp.Go(func() error {
		return healthServer.Start(":9090")
	})

	errgrp.Go(func() error {
		snifferInstance := sniffer.NewSniffer(mapperClient)
		return snifferInstance.RunForever(errGroupCtx)
	})

	err := errgrp.Wait()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		logrus.WithError(err).Fatal("Error when running server or HTTP server")
	}

	timeoutCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = healthServer.Shutdown(timeoutCtx)
	if err != nil {
		logrus.WithError(err).Fatal("Error when shutting down")
	}

	err = metricsServer.Shutdown(timeoutCtx)
	if err != nil {
		logrus.WithError(err).Fatal("Error when shutting down")
	}
}
