package model

import "time"

// Структура запроса на статистику
type StatsRequest struct {
	TsFrom time.Time `json:"tsFrom"`
	TsTo   time.Time `json:"tsTo"`
}

// Структура ответа на статистику
type StatsResponse struct {
	Minute     time.Time `json:"minute"`
	ClickCount int       `json:"clickCount"`
}

type Stats struct {
	Minute time.Time `json:"minute"` // Время округленное до минуты
	Count  int       `json:"count"`  // Количество кликов
}
