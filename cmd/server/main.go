package main

import (
	"context"
	"fmt"
	"kunkun-go/internal/repository"
	"kunkun-go/internal/router"
	"kunkun-go/pkg/config"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/viper"
)

func main() {
	// 初始化配置
	config.InitConfig()

	// 初始化数据库
	repository.InitDB()

	// 初始化 Redis
	repository.InitRedis()

	// 初始化路由
	route := router.InitRouter()

	// 设置端口号
	port := viper.GetString("server.port")

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: route,
	}

	// 启动服务器 goroutine 中运行
	go func() {
		log.Println("Server is running on port", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	//优雅重启 graceful shutdown  让主 goroutine 阻塞等待「退出信号」，收到后再走后面的优雅关闭流程  因为刚刚把 http 服务放在另一个 goroutine 服务里面跑所以不用阻塞主 goroutine
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	//设定 5 秒的超时时间
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v\n", err)
	}

	if err := repository.CloseRedis(); err != nil {
		log.Printf("redis close: %v\n", err)
	}

	log.Println("Server exiting")
}
