package pg

import (
	"clickerProj/main/config"
	"clickerProj/main/model"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"time"
)

// Таймаут для БД
const Timeout = 5

// структура-обертка длля DB
type DB struct {
	*gorm.DB
}

type ClickData struct {
	BannerID   int       `gorm:"primaryKey;column:banner_id" json:"banner_id"` // Banner ID
	Timestamp  time.Time `gorm:"primaryKey;column:timestamp" json:"timestamp"` // Timestamp
	ClickCount int       `gorm:"column:click_count" json:"click_count"`        // Click count for that minute
}

// Подключаемся, проверяем DB
func InitDB() (*DB, error) {
	cfg := config.Get()

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Etc/GMT-3",
		cfg.Database.Host, cfg.Database.User, cfg.Database.Password, cfg.Database.DbName, cfg.Database.Port,
	)

	var db *gorm.DB
	var err error

	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	//db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	// Тестовый запрос за несколько попыток
	var attempt uint

	const maxAttempts = 10

	for {
		attempt++

		log.Printf("[PostgreSQL.Dial] (Ping attempt %d) SELECT 1\n", attempt)

		//_, err := pgDB.Exec("SELECT 1")
		db, err = gorm.Open(postgres.Open(dsn), gormConfig)
		if err != nil {
			log.Printf("[PostgreSQL.Dial] (Ping attempt %d) error: %s\n", attempt, err)

			if attempt < maxAttempts {
				time.Sleep(1 * time.Second)

				continue
			}

			return nil, fmt.Errorf("pgDB.Exec failed: %w", err)
		}

		log.Printf("[PostgreSQL.Dial] (Ping attempt %d) OK\n", attempt)

		break
	}

	// Вторая автомиграция, если кто-то заглянет в проект))
	err = db.AutoMigrate(&ClickData{})
	if err != nil {
		return nil, fmt.Errorf("failed to auto migrate: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get raw SQL DB object: %w", err)
	}

	sqlDB.SetConnMaxLifetime(time.Second * time.Duration(Timeout))
	sqlDB.SetConnMaxLifetime(time.Minute * 10)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	return &DB{db}, nil
}

// Сохранение кликов в базу данных
func (db *DB) SaveClickData(bannerID int, timestamp time.Time, clickCount int) error {
	// Проверим баннер Id и время
	var existingClick ClickData
	err := db.DB.Where("banner_id = ? AND timestamp = ?", bannerID, timestamp).First(&existingClick).Error
	if err == nil {
		// Обновим карту
		existingClick.ClickCount += clickCount
		err = db.DB.Save(&existingClick).Error
	} else if err == gorm.ErrRecordNotFound {
		// Записи нет, создадим новую
		clickData := ClickData{
			BannerID:   bannerID,
			Timestamp:  timestamp,
			ClickCount: clickCount,
		}
		err = db.DB.Create(&clickData).Error
	} else {
		// Все остальные ошибки
		return fmt.Errorf("error checking or saving click data: %w", err)
	}

	if err != nil {
		return fmt.Errorf("error saving click data: %w", err)
	}

	return nil
}

// Функция для получения статистики по баннеру
func GetStats(db *gorm.DB, bannerID int, tsFrom, tsTo time.Time) ([]model.Stats, error) {
	var stats []model.Stats

	utcPlus3 := time.FixedZone("UTC+3", 3*60*60) // 3 hours offset from UTC

	// Convert current time to UTC+3
	tsFrom3 := time.Date(tsFrom.Year(), tsFrom.Month(), tsFrom.Day(), tsFrom.Hour(), tsFrom.Minute(),
		tsFrom.Second(), tsFrom.Nanosecond(), utcPlus3)
	tsTo3 := time.Date(tsTo.Year(), tsTo.Month(), tsTo.Day(), tsTo.Hour(), tsTo.Minute(),
		tsTo.Second(), tsTo.Nanosecond(), utcPlus3)

	// Выполним запрос для получения статистики по кликам по баннеру за указанный промежуток времени
	err := db.Table("click_data").
		Where("banner_id = ? AND timestamp BETWEEN ? AND ?", bannerID, tsFrom3, tsTo3).
		Group("timestamp").
		Select("timestamp AS minute, COUNT(*) AS count").
		Scan(&stats).Error

	if err != nil {
		return nil, fmt.Errorf("Error fetching stats: %v", err)
	}

	return stats, nil
}
