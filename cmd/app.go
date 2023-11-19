package main

import (
	"ebook/cmd/internal/events"
	"github.com/gin-gonic/gin"
)

type App struct {
	server    *gin.Engine
	consumers []events.Consumer
}
