package handler

import (
	"fmt"
	"log"
	"lumel/pkg/repository"
	"lumel/pkg/services"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type handler struct {
	service services.Service
	repo    repository.SalesRepository
}

func (h handler) RefreshData(c *gin.Context) {
	go func() {
		log.Println("Starting data refresh in the background...")
		err := h.service.UploadData()
		if err != nil {
			log.Printf("Failed to refresh data: %v", err)
			return
		}
		log.Println("Data successfully refreshed!")
	}()

	c.JSON(http.StatusAccepted, gin.H{"status": "Syncing in background"})
}

func (h handler) GetTotalRevenue(c *gin.Context) {
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")

	startDateTime, endDateTime, err := parseDateTime(startDate, endDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid date formate",
		})
		return
	}

	revenue, err := h.repo.TotalRevenue(startDateTime, endDateTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "error fetching revenue details",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    revenue,
		"message": "success",
	})

}

func (h handler) GetTotalRevenueByProduct(c *gin.Context) {
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")

	startDateTime, endDateTime, err := parseDateTime(startDate, endDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid date formate",
		})
		return
	}

	revenue, err := h.repo.TotalRevenueByProduct(startDateTime, endDateTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "error fetching revenue details",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    revenue,
		"message": "success",
	})

}

func (h handler) GetTotalRevenueByCategory(c *gin.Context) {
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")

	startDateTime, endDateTime, err := parseDateTime(startDate, endDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid date formate",
		})
		return
	}

	revenue, err := h.repo.TotalRevenueByCategory(startDateTime, endDateTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "error fetching revenue details",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    revenue,
		"message": "success",
	})

}

func parseDateTime(startDate, endDate string) (startTime, endTime time.Time, err error) {
	layout := "2006-01-02"

	startTime, err = time.Parse(layout, startDate)
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return startTime, endTime, err
	}

	endTime, err = time.Parse(layout, endDate)
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return startTime, endTime, err
	}

	return startTime, endTime, nil
}

func NewHandler(service services.Service, repo repository.SalesRepository) handler {
	return handler{
		service: service,
		repo:    repo,
	}
}
