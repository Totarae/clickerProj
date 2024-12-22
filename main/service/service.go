package service

import (
	"clickerProj/main/pg"
	"log"
	"sync"
	"time"
)

var clickStats = struct {
	mu    sync.Mutex
	stats map[int]int // Ключ: bannerID, значение: количество кликов
}{
	stats: make(map[int]int),
}

// Инкремент кликов
func IncrementClick(bannerID int) {
	clickStats.mu.Lock()
	defer clickStats.mu.Unlock()

	if clickCount, ok := clickStats.stats[bannerID]; ok {
		clickStats.stats[bannerID] = clickCount + 1
	} else {
		clickStats.stats[bannerID] = 1
	}
}

// Отправка статистики в базу данных
func sendStatsToDB(db *pg.DB) {
	clickStats.mu.Lock()
	defer clickStats.mu.Unlock()

	now := time.Now().Truncate(time.Minute)
	for bannerID, clickCount := range clickStats.stats {
		// Отправляем данные в базу данных
		err := db.SaveClickData(bannerID, now, clickCount)
		if err != nil {
			log.Printf("Error saving click data for banner %d: %v", bannerID, err)
		} else {
			log.Printf("Saved click data for banner %d: %d clicks", bannerID, clickCount) // Log successful save
		}
	}

	// Очищаем после синхронизации
	clickStats.stats = make(map[int]int)
}

// Фоновая синхронизация статистики
func StartStatsSync(db *pg.DB) {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			log.Println("Starting stats synchronization...") // Log before sending stats
			sendStatsToDB(db)
			log.Println("Stats synchronization completed.") // Log after stats sent
		}
	}
}
