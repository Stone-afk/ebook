package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

//func main() {
//	initViperWatch()
//	// 这里暂时随便搞一下
//	// 搞成依赖注入
//	app := Init()
//	for _, c := range app.Consumers {
//		err := c.Start()
//		if err != nil {
//			panic(err)
//		}
//	}
//	//go func() {
//	//	err := app.webAdmin.Start()
//	//	if err != nil {
//	//		log.Println(err)
//	//		panic(err)
//	//	}
//	//}()
//	err := app.server.Serve()
//	if err != nil {
//		log.Println(err)
//		panic(err)
//	}
//}

func initViperWatch() {
	cfile := pflag.String("config",
		"config/config.yaml", "配置文件路径")
	pflag.Parse()
	// 直接指定文件路径
	viper.SetConfigFile(*cfile)
	// 实时监听配置变更
	viper.WatchConfig()
	// 只能告诉你文件变了，不能告诉你，文件的哪些内容变了
	viper.OnConfigChange(func(in fsnotify.Event) {
		// 比较好的设计，它会在 in 里面告诉你变更前的数据，和变更后的数据
		// 更好的设计是，它会直接告诉你差异。
		fmt.Println(in.Name, in.Op)
		fmt.Println(viper.GetString("db.dsn"))
	})
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}
