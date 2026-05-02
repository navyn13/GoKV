package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	gokv "github.com/navyn13/GoKV/cmd"
)

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
func main() {
	loadEnv()
	server := gokv.NewServer(gokv.Config{
		ListenAddr: ":8080",
		Username:   os.Getenv("USERNAME"),
		Password:   os.Getenv("PASSWORD"),
	})
	go func() {
		log.Fatal(server.Start())
	}()
	select {}

	// quit := make(chan os.Signal, 1)
	// signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	// <-quit

	// slog.Info("Shutting down server gracefully...")
	// server.Shutdown()
	// slog.Info("Server stopped.")

}
