package main

import (
	"net/http"
	"safe_talk/config"
	"safe_talk/internal/handler"
	"safe_talk/internal/repository"
	"safe_talk/internal/usecase"
	"safe_talk/pkg/db/dbpostgres"
	"safe_talk/pkg/logger"
	"time"
)

func main() {

	l := logger.New()
	configs, err := config.New()
	if err != nil {
		l.Fatal(err)
	}
	l.Info("Успешно прочитаны конфиги")
	db, err := dbpostgres.New(configs.Postgres, l)
	if err != nil {
		l.Printf("Ошибка при прдключении к БД. Ошибка: %v", err)
		return
	}
	l.Info("Успешное подключение к БД")

	repos := repository.NewRepos(db, l)
	useCase := usecase.New(l, repos)
	h := handler.New(useCase, l)
	routes := handler.InitRoutes(h)
	server := http.Server{
		Addr:         configs.Srv.Host + configs.Srv.Port,
		ReadTimeout:  time.Duration(configs.Srv.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(configs.Srv.WriteTimeout) * time.Second,
		Handler:      routes,
	}
	l.Info("Сервер доступен по порту " + configs.Srv.Port)
	if err := server.ListenAndServe(); err != nil {
		l.Fatal(err)
	}

}
