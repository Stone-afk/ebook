package ioc

import (
	"ebook/cmd/interactive/repository/dao"
	"ebook/cmd/pkg/ginx"
	"ebook/cmd/pkg/gormx/connpool"
	"ebook/cmd/pkg/logger"
	"ebook/cmd/pkg/migrator/events"
	"ebook/cmd/pkg/migrator/events/canal"
	"ebook/cmd/pkg/migrator/events/fixer"
	"ebook/cmd/pkg/migrator/scheduler"
	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/viper"
)

const topic = "migrator_interactive"

type InteractiveMySQLBinlogConsumer struct {
	*canal.MySQLBinlogConsumer[dao.Interactive]
}

func InitInteractiveMySQLBinlogConsumer(
	client sarama.Client,
	l logger.Logger,
	src SrcDB,
	dst DstDB,
	p events.Producer) *InteractiveMySQLBinlogConsumer {
	return &InteractiveMySQLBinlogConsumer{
		MySQLBinlogConsumer: canal.NewMySQLBinlogConsumer[dao.Interactive](
			client, l, src, dst, p,
			"interactive_binlog",
			"migrator_interactive", "ebook_interactive"),
	}
}

func InitFixDataConsumer(l logger.Logger,
	src SrcDB,
	dst DstDB,
	client sarama.Client) *fixer.Consumer[dao.Interactive] {
	res, err := fixer.NewConsumer[dao.Interactive](client, l,
		src, dst, topic)
	if err != nil {
		panic(err)
	}
	return res
}

func InitMigratorProducer(p sarama.SyncProducer) events.Producer {
	return events.NewSaramaProducer(p, topic)
}

func InitMigratorWeb(
	l logger.Logger,
	src SrcDB,
	dst DstDB,
	pool *connpool.DoubleWritePool,
	producer events.Producer,
) *ginx.Server {
	// 在这里，有多少张表，就初始化多少个 scheduler
	intrSch := scheduler.NewScheduler[dao.Interactive](l, src, dst, pool, producer)
	engine := gin.Default()
	ginx.InitCounter(prometheus.CounterOpts{
		Namespace: "ebook_interactive",
		Subsystem: "ebook_interactive_admin",
		Name:      "http_biz_code",
		Help:      "HTTP 的业务错误码",
	})
	intrSch.RegisterRoutes(engine.Group("/migrator"))
	//intrSch.RegisterRoutes(engine.Group("/migrator/interactive"))
	addr := viper.GetString("migrator.web.addr")
	return &ginx.Server{
		Addr:   addr,
		Engine: engine,
	}
}
