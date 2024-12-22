package main

import (
	"clickerProj/main/handler"
	"clickerProj/main/pg"
	"clickerProj/main/service"
	"context"
	"fmt"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func main() {
	// Инициализация БД
	pgDB, err := pg.InitDB() // Для PostgreSQL
	if err != nil {
		fmt.Errorf("pgdb.Dial failed: %w", err)
	}

	// Миграции
	if pgDB != nil {
		log.Println("Launched PostgreSQL migrations")
		/*if err := runPgMigrations(); err != nil {
			log.Printf("runPgMigrations failed: %w", err)
		}*/
	}

	go service.StartStatsSync(pgDB)

	// Создаем канал для обработки сигналов завершения
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	go func() {

		router := fasthttprouter.New()
		router.GET("/counter/:bannerID", handler.IncrementClickHandler) // Читаем параметр :bannerID из урла
		router.POST("/stats/:bannerID", func(ctx *fasthttp.RequestCtx) {
			bannerIDStr := ctx.UserValue("bannerID").(string)
			bannerID, err := strconv.Atoi(bannerIDStr)

			if err != nil {
				ctx.Error("Invalid banner ID", fasthttp.StatusBadRequest)
				return
			}

			handler.GetStatsHandler(ctx, pgDB, bannerID)
		})

		// Запуск сервера
		log.Println("Server is running on port 8080")
		if err := fasthttp.ListenAndServe(":8080", router.Handler); err != nil {
			log.Fatalf("Error starting server: %v", err)
		}

	}()

	// Ожидаем сигнала завершения
	<-done

	log.Println("Received termination signal, shutting down server...")

	// Завершаем работу сервера с тайм-аутом
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Здесь можно выполнить дополнительные задачи перед завершением (например, закрыть соединения с базой)
	if err := shutdownServer(shutdownCtx); err != nil {
		log.Printf("Error shutting down server: %v", err)
	} else {
		log.Println("Server gracefully stopped")
	}
}

// Завершение работы сервера с тайм-аутом
func shutdownServer(ctx context.Context) error {
	// Здесь можно добавить дополнительные действия перед остановкой, например, ожидание завершения всех горутин

	// Закрытие соединений с БД, если необходимо
	// db.Close() // Пример для PostgreSQL

	// Ожидание завершения всех активных соединений
	// Например, можно использовать pool.Close() для подключения к базе данных PostgreSQL
	return nil
}
