package route

import (
	"github.com/ayanami77/pecopeco-service/internal/db"
	"github.com/ayanami77/pecopeco-service/internal/domain/user"
	"github.com/ayanami77/pecopeco-service/internal/infrastructure/repository"
	hh "github.com/ayanami77/pecopeco-service/internal/presentation/health_handler"
	rh "github.com/ayanami77/pecopeco-service/internal/presentation/restaurant"
	uh "github.com/ayanami77/pecopeco-service/internal/presentation/user"
	mymiddleware "github.com/ayanami77/pecopeco-service/internal/server/middleware"
	ru "github.com/ayanami77/pecopeco-service/internal/usecase/restaurant"
	uu "github.com/ayanami77/pecopeco-service/internal/usecase/user"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func InitRoute(r *chi.Mux) {
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/v1", func(r chi.Router) {
		r.Get("/health-check", hh.HealthCheck)
		userRoute(r)
		restaurantRoute(r)
	})
}

func userRoute(r chi.Router) chi.Router {
	db := db.GetDB()
	userRepository := repository.NewUserRepository(db)
	h := uh.NewHandler(
		uu.NewFindUserUseCase(userRepository),
		uu.NewLoginUsecase(
			user.NewUserDomainService(userRepository),
			userRepository,
		),
	)
	return r.Route("/users", func(r chi.Router) {
		r.Post("/login", h.Login)
		r.With(mymiddleware.Auth).Get("/me", h.FindUser)
	})
}

func restaurantRoute(r chi.Router) chi.Router {
	db := db.GetDB()
	restaurantRepository := repository.NewRestaurantRepository(db)
	h := rh.NewHandler(
		ru.NewListRestaurantsUseCase(restaurantRepository),
		ru.NewSaveRestaurantUseCase(restaurantRepository),
		ru.NewDeleteRestaurantsUseCase(restaurantRepository),
	)
	return r.Route("/restaurants", func(r chi.Router) {
		r.Use(mymiddleware.Auth)
		r.Get("/", h.ListRestaurants)
		r.Post("/", h.SaveRestaurant)
		r.Delete("/{id}", h.DeleteRestaurant)
	})
}
