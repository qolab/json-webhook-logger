package main

import (
	"context"
	"fmt"
	"github.com/Jeffail/gabs"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}

func Hook(c *gin.Context) {
	jsonData, err := c.GetRawData()
	if err != nil {
		c.IndentedJSON(http.StatusOK, gin.H{"message": "invalid json"})
	}
	jsonParsedObj, _ := gabs.ParseJSON(jsonData)

	// Output to console
	defer fmt.Println(jsonParsedObj.String())

	c.IndentedJSON(http.StatusOK, gin.H{"message": "processed"})
}

func NotFound(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{"code": "PAGE_NOT_FOUND", "message": "404 page not found"})
}

func main() {
	port := getEnv("PORT", "8080")
	address := fmt.Sprintf("0.0.0.0:%s", port)

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.POST("/hook", Hook)
	r.NoRoute(NotFound)
	r.Run(address)

	srv := &http.Server{
		Addr:    address,
		Handler: r,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error().Err(err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info().Msg("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Error().Err(err)
	}
	select {
	case <-ctx.Done():
		log.Info().Msg("timeout of 2 seconds.")
	}
	log.Info().Msg("Server exiting")
}
