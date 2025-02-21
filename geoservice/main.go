package main

import (
	"fmt"
	_ "geoservice/docs"
	"geoservice/handlers"
	"geoservice/middleware"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Geo Service API
// @version 1.0
// @description Сервис для поиска и геокодирования адресов
// @host localhost:8080
// @BasePath /api
func main() {
	r := chi.NewRouter()

	// Создаем JWT Auth
	tokenAuth := jwtauth.New("HS256", []byte("secret"), nil)

	// Добавляем эндпоинт для метрик Prometheus
	r.Handle("/metrics", promhttp.Handler())

	// Публичные маршруты
	r.Group(func(r chi.Router) {
		r.Post("/api/login", handlers.LoginHandler(tokenAuth))
		r.Post("/api/register", handlers.RegisterHandler)
	})

	// Защищенные маршруты
	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(tokenAuth))
		r.Post("/api/address/search", handlers.SearchHandler)
		r.Post("/api/address/geocode", handlers.GeoHandler)
	})

	// Swagger endpoints
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	time.Sleep(3 * time.Second)
	fmt.Println(`
==================
  HELLO ANDREY!
==================
Starting application...`)

	time.Sleep(2 * time.Second)

	log.Printf("Server starting on :8080")
	log.Printf("Swagger UI available at http://localhost:8080/swagger/index.html")
	log.Printf("Prometheus metrics available at http://localhost:9090")
	log.Printf("Metrics available at http://localhost:8080/metrics")
	log.Printf("Grafana available at http://localhost:3000")

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)

	}

}
