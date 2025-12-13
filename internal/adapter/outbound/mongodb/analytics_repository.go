package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/okinn/service-presensi/internal/domain/entity"
	"github.com/okinn/service-presensi/internal/domain/repository"
)

type AnalyticsRepository struct {
	collection *mongo.Collection
}

func NewAnalyticsRepository(db *mongo.Database) repository.AnalyticsRepository {
	return &AnalyticsRepository{
		collection: db.Collection("presensi"),
	}
}

func (r *AnalyticsRepository) GetSummary(ctx context.Context, filter entity.AnalyticsFilter) (*entity.AttendanceSummary, error) {
	matchStage := bson.M{}

	if filter.UserID != "" {
		matchStage["user_id"] = filter.UserID
	}
	if !filter.StartDate.IsZero() || !filter.EndDate.IsZero() {
		dateFilter := bson.M{}
		if !filter.StartDate.IsZero() {
			dateFilter["$gte"] = filter.StartDate
		}
		if !filter.EndDate.IsZero() {
			dateFilter["$lte"] = filter.EndDate
		}
		matchStage["tanggal"] = dateFilter
	}

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: matchStage}},
		{{Key: "$group", Value: bson.M{
			"_id":            nil,
			"total_records":  bson.M{"$sum": 1},
			"total_hadir":    bson.M{"$sum": bson.M{"$cond": bson.A{bson.M{"$eq": bson.A{"$status", "hadir"}}, 1, 0}}},
			"total_terlambat": bson.M{"$sum": bson.M{"$cond": bson.A{bson.M{"$eq": bson.A{"$status", "terlambat"}}, 1, 0}}},
			"total_izin":     bson.M{"$sum": bson.M{"$cond": bson.A{bson.M{"$eq": bson.A{"$status", "izin"}}, 1, 0}}},
			"total_sakit":    bson.M{"$sum": bson.M{"$cond": bson.A{bson.M{"$eq": bson.A{"$status", "sakit"}}, 1, 0}}},
			"total_alpha":    bson.M{"$sum": bson.M{"$cond": bson.A{bson.M{"$eq": bson.A{"$status", "alpha"}}, 1, 0}}},
		}}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	summary := &entity.AttendanceSummary{}
	if len(results) > 0 {
		result := results[0]
		summary.TotalRecords = int(result["total_records"].(int32))
		summary.TotalHadir = int(result["total_hadir"].(int32))
		summary.TotalTerlambat = int(result["total_terlambat"].(int32))
		summary.TotalIzin = int(result["total_izin"].(int32))
		summary.TotalSakit = int(result["total_sakit"].(int32))
		summary.TotalAlpha = int(result["total_alpha"].(int32))
		summary.CalculatePercentage()
	}

	return summary, nil
}

func (r *AnalyticsRepository) GetDailySummary(ctx context.Context, date string) (*entity.DailySummary, error) {
	// Parse date string YYYY-MM-DD
	parsedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, err
	}

	// Set date range for the entire day
	startOfDay := time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(), 0, 0, 0, 0, parsedDate.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	filter := entity.AnalyticsFilter{
		StartDate: startOfDay,
		EndDate:   endOfDay,
	}

	summary, err := r.GetSummary(ctx, filter)
	if err != nil {
		return nil, err
	}

	breakdown, err := r.GetStatusBreakdown(ctx, filter)
	if err != nil {
		return nil, err
	}

	return &entity.DailySummary{
		Date:    parsedDate,
		Summary: *summary,
		Details: breakdown,
	}, nil
}

func (r *AnalyticsRepository) GetMonthlySummary(ctx context.Context, month string) (*entity.MonthlySummary, error) {
	// Parse month string YYYY-MM
	parsedMonth, err := time.Parse("2006-01", month)
	if err != nil {
		return nil, err
	}

	// Set date range for the entire month
	startOfMonth := time.Date(parsedMonth.Year(), parsedMonth.Month(), 1, 0, 0, 0, 0, parsedMonth.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, 0)

	filter := entity.AnalyticsFilter{
		StartDate: startOfMonth,
		EndDate:   endOfMonth,
	}

	summary, err := r.GetSummary(ctx, filter)
	if err != nil {
		return nil, err
	}

	// Get daily stats within the month
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{
			"tanggal": bson.M{
				"$gte": startOfMonth,
				"$lt":  endOfMonth,
			},
		}}},
		{{Key: "$group", Value: bson.M{
			"_id": bson.M{
				"$dateToString": bson.M{"format": "%Y-%m-%d", "date": "$tanggal"},
			},
			"count": bson.M{"$sum": 1},
		}}},
		{{Key: "$sort", Value: bson.M{"_id": 1}}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var dailyResults []bson.M
	if err := cursor.All(ctx, &dailyResults); err != nil {
		return nil, err
	}

	dailyStats := make([]entity.DailyStats, 0, len(dailyResults))
	for _, result := range dailyResults {
		dailyStats = append(dailyStats, entity.DailyStats{
			Date:  result["_id"].(string),
			Count: int(result["count"].(int32)),
		})
	}

	return &entity.MonthlySummary{
		Month:      month,
		Summary:    *summary,
		DailyStats: dailyStats,
	}, nil
}

func (r *AnalyticsRepository) GetUserSummary(ctx context.Context, userID string, filter entity.AnalyticsFilter) (*entity.UserSummary, error) {
	// Set user filter
	filter.UserID = userID

	summary, err := r.GetSummary(ctx, filter)
	if err != nil {
		return nil, err
	}

	statusBreakdown, err := r.GetStatusBreakdown(ctx, filter)
	if err != nil {
		return nil, err
	}

	// Get user name from the first matching record
	matchFilter := bson.M{"user_id": userID}
	var doc struct {
		Nama string `bson:"nama"`
	}
	err = r.collection.FindOne(ctx, matchFilter).Decode(&doc)
	userName := ""
	if err == nil {
		userName = doc.Nama
	}

	// Format period
	period := "all-time"
	if !filter.StartDate.IsZero() && !filter.EndDate.IsZero() {
		period = filter.StartDate.Format("2006-01-02") + " to " + filter.EndDate.Format("2006-01-02")
	} else if !filter.StartDate.IsZero() {
		period = "from " + filter.StartDate.Format("2006-01-02")
	} else if !filter.EndDate.IsZero() {
		period = "until " + filter.EndDate.Format("2006-01-02")
	}

	return &entity.UserSummary{
		UserID:       userID,
		UserName:     userName,
		Period:       period,
		Summary:      *summary,
		StatusDetail: statusBreakdown,
	}, nil
}

func (r *AnalyticsRepository) GetStatusBreakdown(ctx context.Context, filter entity.AnalyticsFilter) ([]entity.StatusBreakdown, error) {
	matchStage := bson.M{}

	if filter.UserID != "" {
		matchStage["user_id"] = filter.UserID
	}
	if !filter.StartDate.IsZero() || !filter.EndDate.IsZero() {
		dateFilter := bson.M{}
		if !filter.StartDate.IsZero() {
			dateFilter["$gte"] = filter.StartDate
		}
		if !filter.EndDate.IsZero() {
			dateFilter["$lte"] = filter.EndDate
		}
		matchStage["tanggal"] = dateFilter
	}

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: matchStage}},
		{{Key: "$group", Value: bson.M{
			"_id":   "$status",
			"count": bson.M{"$sum": 1},
		}}},
		{{Key: "$sort", Value: bson.M{"count": -1}}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	breakdown := make([]entity.StatusBreakdown, 0, len(results))
	for _, result := range results {
		breakdown = append(breakdown, entity.StatusBreakdown{
			Status: result["_id"].(string),
			Count:  int(result["count"].(int32)),
		})
	}

	return breakdown, nil
}
