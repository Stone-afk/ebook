//go:build k8s

// 使用 k8s 这个编译标签
package config

var Config = ebookConfig{
	DB: DBConfig{
		// 本地连接
		DSN: "root:root@tcp(localhost:13316)/ebook",
	},
	Redis: RedisConfig{
		Addr: "ebook-live-redis:11479",
	},
}
