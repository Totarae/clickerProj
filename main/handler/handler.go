package handler

import (
	"clickerProj/main/model"
	"clickerProj/main/pg"
	"clickerProj/main/service"
	"encoding/json"
	"github.com/valyala/fasthttp"
	"log"
	"strconv"
)

func IncrementClickHandler(ctx *fasthttp.RequestCtx) {
	bannerIDStr := ctx.UserValue("bannerID").(string)
	bannerID, err := strconv.Atoi(bannerIDStr)
	if err != nil {
		ctx.Error("Invalid banner ID", fasthttp.StatusBadRequest)
		return
	}

	service.IncrementClick(bannerID) // Увеличиваем счетчик кликов
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody([]byte("Click counted"))
}

func GetStatsHandler(ctx *fasthttp.RequestCtx, db *pg.DB, bannerID int) {

	var statsRequest model.StatsRequest
	if err := json.Unmarshal(ctx.PostBody(), &statsRequest); err != nil {
		ctx.Error("Invalid request body", fasthttp.StatusBadRequest)
		return
	}

	stats, err := pg.GetStats(db.DB, bannerID, statsRequest.TsFrom, statsRequest.TsTo)
	if err != nil {
		log.Printf("Error retrieving statistics: %v", err)
		ctx.Error("Error retrieving statistics", fasthttp.StatusInternalServerError)
		return
	}

	// Ответ с поминутной статистикой
	response, _ := json.Marshal(stats)
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody(response)
}
