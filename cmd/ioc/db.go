package ioc

import (
	"ebook/cmd/internal/repository/dao"
	prometheus2 "ebook/cmd/pkg/gormx/callbacks/prometheus"
	"ebook/cmd/pkg/logger"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
	"gorm.io/plugin/opentelemetry/tracing"
	"gorm.io/plugin/prometheus"
	"time"
)

//func InitMongoDB() *mongo.Database {
//	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
//	defer cancel()
//	monitor := &event.CommandMonitor{
//		Started: func(ctx context.Context,
//			startedEvent *event.CommandStartedEvent) {
//			fmt.Println(startedEvent.Command)
//		},
//	}
//	opts := options.Client().
//		ApplyURI("mongodb://root:example@localhost:27017/").
//		SetMonitor(monitor)
//	client, err := mongo.Connect(ctx, opts)
//	if err != nil {
//		panic(err)
//	}
//	mongoDB := client.Database("ebook")
//	return mongoDB
//}

func InitDB() *gorm.DB {
	type Config struct {
		DSN string `yaml:"dsn"`

		// 有些人的做法
		// localhost:13316
		//Addr string
		//// localhost
		//Domain string
		//// 13316
		//Port string
		//Protocol string
		//// root
		//Username string
		//// root
		//Password string
		//// webook
		//DBName string
	}
	var cfg = Config{
		DSN: "root:root@tcp(localhost:13316)/ebook_default",
	}
	// 看起来，remote 不支持 key 的切割
	err := viper.UnmarshalKey("db", &cfg)
	//dsn := viper.GetString("db.mysql")
	//println(dsn)
	//if err != nil {
	//	panic(err)
	//}
	db, err := gorm.Open(mysql.Open(cfg.DSN), &gorm.Config{
		// 缺了一个 writer
		Logger: glogger.New(nil, glogger.Config{
			// 慢查询阈值，只有执行时间超过这个阈值，才会使用
			// 50ms， 100ms
			// SQL 查询必然要求命中索引，最好就是走一次磁盘 IO
			// 一次磁盘 IO 是不到 10ms
			SlowThreshold: time.Millisecond * 10,
			// 忽略未找到记录的错误
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries:      true,
			LogLevel:                  glogger.Info,
		}),
	})
	if err != nil {
		// 我只会在初始化过程中 panic
		// panic 相当于整个 goroutine 结束
		// 一旦初始化过程出错，应用就不要启动了
		panic(err)
	}

	//dao.NewUserDAOV1(func() *gorm.DB {
	//viper.OnConfigChange(func(in fsnotify.Event) {
	//oldDB := db
	//db, err = gorm.Open(mysql.Open())
	//pt := unsafe.Pointer(&db)
	//atomic.StorePointer(&pt, unsafe.Pointer(&db))
	//oldDB.Close()
	//})
	// 要用原子操作
	//return db
	//})

	// 接入 prometheus
	err = db.Use(prometheus.New(prometheus.Config{
		DBName: "ebook",
		// 每 15 秒采集一些数据
		RefreshInterval: 15,
		MetricsCollector: []prometheus.MetricsCollector{
			&prometheus.MySQL{
				VariableNames: []string{"Threads_running"},
			},
		}, // user defined metrics
	}))
	if err != nil {
		panic(err)
	}
	err = db.Use(tracing.NewPlugin(tracing.WithoutMetrics()))
	if err != nil {
		panic(err)
	}
	prom := prometheus2.Callbacks{
		Namespace:  "server database",
		Subsystem:  "ebook",
		Name:       "gorm",
		InstanceID: "my-instance-1",
		Help:       "gorm DB 查询",
	}
	err = prom.Register(db)
	if err != nil {
		panic(err)
	}

	err = dao.InitTables(db)
	if err != nil {
		panic(err)
	}
	return db
}

type gormLoggerFunc func(msg string, fields ...logger.Field)

func (g gormLoggerFunc) Printf(msg string, args ...interface{}) {
	g(msg, logger.Field{Key: "args", Value: args})
}

type DoSomething interface {
	DoABC() string
}

type DoSomethingFunc func() string

func (d DoSomethingFunc) DoABC() string {
	return d()
}
