package main

import (
	"log/slog"
	"net"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()
	// r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	if address := os.Getenv("LISTEN_SOCKET_PATH"); len(address) > 0 {
		// 监听 Unix socket
		listener, err := net.Listen("unix", address)
		if err != nil {
			slog.Error("Error listening on socket: " + err.Error())
			return
		}
		defer listener.Close()

		// 确保程序退出时删除套接字文件
		defer os.Remove(address)

		slog.Info("Listening on socket: " + listener.Addr().String())

		// 使用 Gin 路由器处理来自 socket 的请求
		if err := http.Serve(listener, r); err != nil {
			slog.Error("Error serving on socket: " + err.Error())
		}
	} else {
		var port = os.Getenv("PORT")
		if port == "" {
			port = "3000"
		}
		if err := http.ListenAndServe(":"+port, r); err != nil {
			slog.Error("failed to start HTTP server: " + err.Error())
		}
	}
}
