package main

import (
	"fmt"
	"lumel/library/db"
	handler "lumel/pkg/handler/sales"
	"lumel/pkg/repository"
	"lumel/pkg/services"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	db.Connect()
	defer db.Disconnect()

	repo := repository.NewRepo(db.Client.Database("dbname"))
	service := services.NewService(repo)
	salesHandler := handler.NewHandler(service, repo)

	go func() {
		err := service.UploadData()
		if err != nil {
			fmt.Println("failed to sync data")
		}

		//refresh data every 12 hours
		time.Sleep(12 * time.Hour)
	}()

	r := gin.Default()

	r.POST("/v1/refresh", salesHandler.RefreshData)
	r.GET("/v1/revenue", salesHandler.GetTotalRevenue)
	r.POST("/v1/revenue/product", salesHandler.GetTotalRevenueByProduct)
	r.POST("/v1/revenue/category", salesHandler.GetTotalRevenueByCategory)

	r.Run(":8080")
}
