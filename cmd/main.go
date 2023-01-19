package main

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"mentat-backend/cmd/setup"
	"os"
	"os/signal"
	"sync"
	"time"
)

func TrapSignal(e *echo.Echo) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	sig := <-c
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	e.Shutdown(ctx)
	logrus.Println("Gracefully shutting down the server: ", sig)
}

func main() {
	var wg sync.WaitGroup
	wg.Add(1)

	setup.ConfigureLogrus()

	go func() {
		port := 8080
		err, e := setup.InitializeEcho()
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"component": "main",
			}).WithError(err).Error("API server running failed")
			return
		}

		err = e.Start(":" + fmt.Sprint(port))
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"component": "main",
			}).WithError(err).Error("API server running failed")
		}

		fmt.Println("Server is now running at the port :" + fmt.Sprint(port))
		logrus.WithFields(logrus.Fields{
			"component": "main",
			"port":      port,
		}).Info("API server running")

		TrapSignal(e)
	}()

	wg.Wait()
}
