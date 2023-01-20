package main

import (
	"chatgpt-service/cmd/setup"
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
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

	//setup.ConfigureLogrus()

	go func() {
		port := 8080
		err, e := setup.InitializeEcho()
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"service":   "chatgpt-service",
				"component": "main",
			}).WithError(err).Error("ChatGPT Service API server running failed")
			fmt.Println("ChatGPT API server running failed")
			os.Exit(1)
		}

		err = e.Start(":" + fmt.Sprint(port))
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"service":   "chatgpt-service",
				"component": "main",
			}).WithError(err).Error("ChatGPT Service API server running failed")
			fmt.Println("ChatGPT API server running failed")
			os.Exit(1)
		}

		fmt.Println("ChatGPT Service Server is now running at the port :" + fmt.Sprint(port))
		logrus.WithFields(logrus.Fields{
			"service":   "chatgpt-service",
			"component": "main",
			"port":      port,
		}).Info("ChatGPT Service API server running")

		TrapSignal(e)
	}()

	wg.Wait()
}
