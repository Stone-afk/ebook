package main

import (
	"ebook/cmd/internal/events"
	"github.com/gin-gonic/gin"
)

type App struct {
	web       *gin.Engine
	consumers []events.Consumer
}
