package server

import (
	"github.com/SemmiDev/go-product/internal/app/handler"
	"github.com/SemmiDev/go-product/internal/app/repository"
	"github.com/SemmiDev/go-product/internal/app/service"
	"github.com/SemmiDev/go-product/internal/config"
	"github.com/SemmiDev/go-product/internal/db/mysql"
	"github.com/SemmiDev/go-product/internal/db/redis"
	"github.com/SemmiDev/go-product/internal/security/middleware"
	"net/http"

	"github.com/go-chi/chi"
	chimiddleware "github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
)

func NewRouter(mysqlClient mysql.Client, redisClient redis.Client) *chi.Mux {
	router := chi.NewRouter()

	router.Use(httprate.LimitByIP(
		config.Cfg().HttpRateLimitRequest,
		config.Cfg().HttpRateLimitTime,
	))

	router.Use(cors.AllowAll().Handler)
	router.Use(chimiddleware.Logger)
	router.Use(chimiddleware.Recoverer)

	merchantRepository := repository.NewMerchantRepository(mysqlClient, redisClient)
	productRepository := repository.NewProductRepository(mysqlClient, redisClient)

	authService := service.NewAuthService(merchantRepository)
	merchantService := service.NewMerchantService(merchantRepository)
	productService := service.NewProductService(productRepository)

	authHandler := handler.NewAuthHandler(authService)
	merchantHandler := handler.NewMerchantHandler(merchantService)
	productHandler := handler.NewProductHandler(productService)

	router.Options("/*", func(w http.ResponseWriter, r *http.Request) {})
	api := router.Route("/v1", func(router chi.Router) {})

	api.Route("/merchants", func(r chi.Router) {
		r.Post("/auth", authHandler.Login())

		r.Post("/", merchantHandler.Create())
		r.Get("/", merchantHandler.List())
		r.Get("/{merchant_id}", merchantHandler.Get())
		r.With(middleware.JWTVerifier).Put("/{merchant_id}", merchantHandler.Update())
		r.With(middleware.JWTVerifier).Put("/{merchant_id}/password", merchantHandler.UpdatePassword())
		r.With(middleware.JWTVerifier).Delete("/{merchant_id}", merchantHandler.Delete())
	})

	api.Route("/products", func(r chi.Router) {
		r.With(middleware.JWTVerifier).Post("/", productHandler.Create())
		r.Get("/", productHandler.List())
		r.Get("/{post_id}", productHandler.Get())
		r.With(middleware.JWTVerifier).Put("/{post_id}", productHandler.Update())
		r.With(middleware.JWTVerifier).Delete("/{post_id}", productHandler.Delete())
	})

	return router
}
