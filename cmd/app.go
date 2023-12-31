package main

import (
	"ebook/cmd/pkg/saramax"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
)

type App struct {
	server    *gin.Engine
	consumers []saramax.Consumer
	cron      *cron.Cron
}
