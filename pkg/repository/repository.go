package repository

import (
	"context"
	"lumel/pkg/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type SalesRepository interface {
	BulkUploadData(records []models.SaleRecord) error
	TotalRevenue(startDate, endDate time.Time) (float64, error)
	TotalRevenueByProduct(startDate, endDate time.Time) (map[string]float64, error)
	TotalRevenueByCategory(startDate, endDate time.Time) (map[string]float64, error)
}

type repo struct {
	collection *mongo.Collection
}

func (r repo) BulkUploadData(records []models.SaleRecord) error {
	var models []mongo.WriteModel
	for _, record := range records {
		models = append(models, mongo.NewInsertOneModel().SetDocument(record))
		if len(models) >= 1000 { // Batch size
			_, err := r.collection.BulkWrite(context.TODO(), models)
			if err != nil {
				return err
			}
			models = nil // Reset for the next batch
		}
	}
	// Insert any remaining records
	if len(models) > 0 {
		_, err := r.collection.BulkWrite(context.TODO(), models)
		return err
	}
	return nil
}

func (r repo) TotalRevenue(startDate, endDate time.Time) (float64, error) {
	filter := bson.M{
		"date_of_sale": bson.M{
			"$gte": startDate,
			"$lt":  endDate,
		},
	}

	pipeline := mongo.Pipeline{
		bson.D{{"$match", filter}},
		bson.D{{"$group", bson.M{
			"_id":   nil,
			"total": bson.M{"$sum": bson.M{"$multiply": []interface{}{"$quantity_sold", bson.M{"$subtract": []interface{}{"$unit_price", "$discount"}}}}}}},
		},
	}

	cursor, err := r.collection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return 0, err
	}
	defer cursor.Close(context.TODO())

	var result struct {
		Total float64 `bson:"total"`
	}
	if cursor.Next(context.TODO()) {
		if err := cursor.Decode(&result); err != nil {
			return 0, err
		}
	}

	return result.Total, nil
}

// TotalRevenueByProduct calculates total revenue by product for a given date range.
func (r repo) TotalRevenueByProduct(startDate, endDate time.Time) (map[string]float64, error) {
	filter := bson.M{
		"date_of_sale": bson.M{
			"$gte": startDate,
			"$lt":  endDate,
		},
	}

	pipeline := mongo.Pipeline{
		{{"$match", filter}},
		{{"$group", bson.M{
			"_id":   "$product_id",
			"total": bson.M{"$sum": bson.M{"$multiply": []interface{}{"$quantity_sold", bson.M{"$subtract": []interface{}{"$unit_price", "$discount"}}}}}},
		}},
	}

	cursor, err := r.collection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	revenueByProduct := make(map[string]float64)
	for cursor.Next(context.TODO()) {
		var result struct {
			ProductID string  `bson:"_id"`
			Total     float64 `bson:"total"`
		}
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}
		revenueByProduct[result.ProductID] = result.Total
	}

	return revenueByProduct, nil
}

// TotalRevenueByCategory calculates total revenue by category for a given date range.
func (r repo) TotalRevenueByCategory(startDate, endDate time.Time) (map[string]float64, error) {
	filter := bson.M{
		"date_of_sale": bson.M{
			"$gte": startDate,
			"$lt":  endDate,
		},
	}

	pipeline := mongo.Pipeline{
		bson.D{{"$match", filter}},
		bson.D{{"$group", bson.M{
			"_id":   "$category",
			"total": bson.M{"$sum": bson.M{"$multiply": []interface{}{"$quantity_sold", bson.M{"$subtract": []interface{}{"$unit_price", "$discount"}}}}}},
		}},
	}

	cursor, err := r.collection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	revenueByCategory := make(map[string]float64)
	for cursor.Next(context.TODO()) {
		var result struct {
			Category string  `bson:"_id"`
			Total    float64 `bson:"total"`
		}
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}
		revenueByCategory[result.Category] = result.Total
	}

	return revenueByCategory, nil
}

func NewRepo(db *mongo.Database) SalesRepository {
	return repo{
		collection: db.Collection("sales"),
	}
}
