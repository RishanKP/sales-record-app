package services

import (
	"encoding/csv"
	"log"
	"lumel/library/config"
	"lumel/pkg/models"
	"lumel/pkg/repository"
	"os"
	"strconv"
	"sync"
	"time"
)

type Service interface {
	UploadData() error
}

type service struct {
	repo repository.SalesRepository
}

func (s service) UploadData() error {
	file, err := os.Open(config.PATH_TO_CSV)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	batchSize := 1000
	var currentBatch []models.SaleRecord

	for _, record := range records[1:] {
		quantitySold, _ := strconv.Atoi(record[5])
		unitPrice, _ := strconv.ParseFloat(record[6], 64)
		discount, _ := strconv.ParseFloat(record[7], 64)
		shippingCost, _ := strconv.ParseFloat(record[8], 64)
		dateOfSale, _ := time.Parse("2006-01-02", record[9])

		saleRecord := models.SaleRecord{
			OrderID:       record[0],
			ProductID:     record[1],
			CustomerID:    record[2],
			QuantitySold:  quantitySold,
			UnitPrice:     unitPrice,
			Discount:      discount,
			ShippingCost:  shippingCost,
			PaymentMethod: record[10],
			DateOfSale:    dateOfSale,
		}

		currentBatch = append(currentBatch, saleRecord)
		if len(currentBatch) >= batchSize {
			wg.Add(1)
			go func(batch []models.SaleRecord) {
				defer wg.Done()
				if err := s.repo.BulkUploadData(batch); err != nil {
					log.Printf("Error inserting batch: %v", err)
				}
			}(currentBatch)
			currentBatch = nil // Reset current batch
		}
	}

	return nil
}

func NewService(repo repository.SalesRepository) Service {
	return service{
		repo: repo,
	}
}
