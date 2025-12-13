/*
 * Copyright (c) 2024 Bima Kharisma Wicaksana
 * GitHub: https://github.com/bimakw
 *
 * Licensed under MIT License with Attribution Requirement.
 * See LICENSE file for details.
 */

package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/okinn/service-presensi/internal/domain/entity"
	"github.com/okinn/service-presensi/internal/domain/repository"
)

// auditLogDocument adalah representasi MongoDB document untuk audit log
type auditLogDocument struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	EntityType string             `bson:"entity_type"`
	EntityID   string             `bson:"entity_id"`
	Action     string             `bson:"action"`
	OldValue   string             `bson:"old_value,omitempty"`
	NewValue   string             `bson:"new_value,omitempty"`
	Changes    string             `bson:"changes,omitempty"`
	UserID     string             `bson:"user_id"`
	UserName   string             `bson:"user_name"`
	UserRole   string             `bson:"user_role"`
	IPAddress  string             `bson:"ip_address"`
	UserAgent  string             `bson:"user_agent,omitempty"`
	RequestID  string             `bson:"request_id,omitempty"`
	StatusCode int                `bson:"status_code,omitempty"`
	ErrorMsg   string             `bson:"error_msg,omitempty"`
	Duration   int64              `bson:"duration_ms,omitempty"`
	CreatedAt  time.Time          `bson:"created_at"`
}

// AuditLogRepository implements repository.AuditLogRepository
type AuditLogRepository struct {
	collection *mongo.Collection
}

// NewAuditLogRepository creates a new audit log repository
func NewAuditLogRepository(db *mongo.Database) repository.AuditLogRepository {
	collection := db.Collection("audit_logs")

	// Create indexes for better query performance
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "entity_type", Value: 1}, {Key: "entity_id", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "user_id", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "created_at", Value: -1}},
		},
		{
			Keys: bson.D{{Key: "action", Value: 1}},
		},
		{
			Keys:    bson.D{{Key: "created_at", Value: 1}},
			Options: options.Index().SetExpireAfterSeconds(90 * 24 * 60 * 60), // 90 days TTL
		},
	}

	collection.Indexes().CreateMany(ctx, indexes)

	return &AuditLogRepository{
		collection: collection,
	}
}

func (r *AuditLogRepository) Create(ctx context.Context, auditLog *entity.AuditLog) error {
	doc := toAuditLogDocument(auditLog)
	result, err := r.collection.InsertOne(ctx, doc)
	if err != nil {
		return err
	}

	auditLog.ID = result.InsertedID.(primitive.ObjectID).Hex()
	return nil
}

func (r *AuditLogRepository) GetByID(ctx context.Context, id string) (*entity.AuditLog, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var doc auditLogDocument
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&doc)
	if err != nil {
		return nil, err
	}

	return toAuditLogEntity(&doc), nil
}

func (r *AuditLogRepository) GetAll(ctx context.Context, filter repository.AuditLogFilter, page, limit int) ([]entity.AuditLog, int64, error) {
	bsonFilter := buildAuditLogFilter(filter)

	total, err := r.collection.CountDocuments(ctx, bsonFilter)
	if err != nil {
		return nil, 0, err
	}

	skip := int64((page - 1) * limit)
	opts := options.Find().
		SetSkip(skip).
		SetLimit(int64(limit)).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.collection.Find(ctx, bsonFilter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var docs []auditLogDocument
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, 0, err
	}

	auditLogs := make([]entity.AuditLog, len(docs))
	for i, doc := range docs {
		auditLogs[i] = *toAuditLogEntity(&doc)
	}

	return auditLogs, total, nil
}

func (r *AuditLogRepository) GetByEntityID(ctx context.Context, entityType, entityID string) ([]entity.AuditLog, error) {
	filter := bson.M{
		"entity_type": entityType,
		"entity_id":   entityID,
	}

	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var docs []auditLogDocument
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, err
	}

	auditLogs := make([]entity.AuditLog, len(docs))
	for i, doc := range docs {
		auditLogs[i] = *toAuditLogEntity(&doc)
	}

	return auditLogs, nil
}

func (r *AuditLogRepository) GetByUserID(ctx context.Context, userID string, page, limit int) ([]entity.AuditLog, int64, error) {
	filter := bson.M{"user_id": userID}

	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	skip := int64((page - 1) * limit)
	opts := options.Find().
		SetSkip(skip).
		SetLimit(int64(limit)).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var docs []auditLogDocument
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, 0, err
	}

	auditLogs := make([]entity.AuditLog, len(docs))
	for i, doc := range docs {
		auditLogs[i] = *toAuditLogEntity(&doc)
	}

	return auditLogs, total, nil
}

func (r *AuditLogRepository) DeleteOldLogs(ctx context.Context, olderThan time.Time) (int64, error) {
	result, err := r.collection.DeleteMany(ctx, bson.M{
		"created_at": bson.M{"$lt": olderThan},
	})
	if err != nil {
		return 0, err
	}
	return result.DeletedCount, nil
}

// Helper functions

func buildAuditLogFilter(filter repository.AuditLogFilter) bson.M {
	bsonFilter := bson.M{}

	if filter.EntityType != "" {
		bsonFilter["entity_type"] = filter.EntityType
	}
	if filter.EntityID != "" {
		bsonFilter["entity_id"] = filter.EntityID
	}
	if filter.Action != "" {
		bsonFilter["action"] = string(filter.Action)
	}
	if filter.UserID != "" {
		bsonFilter["user_id"] = filter.UserID
	}
	if filter.IPAddress != "" {
		bsonFilter["ip_address"] = filter.IPAddress
	}
	if !filter.StartDate.IsZero() || !filter.EndDate.IsZero() {
		dateFilter := bson.M{}
		if !filter.StartDate.IsZero() {
			dateFilter["$gte"] = filter.StartDate
		}
		if !filter.EndDate.IsZero() {
			dateFilter["$lte"] = filter.EndDate
		}
		bsonFilter["created_at"] = dateFilter
	}

	return bsonFilter
}

func toAuditLogDocument(a *entity.AuditLog) *auditLogDocument {
	return &auditLogDocument{
		EntityType: a.EntityType,
		EntityID:   a.EntityID,
		Action:     string(a.Action),
		OldValue:   a.OldValue,
		NewValue:   a.NewValue,
		Changes:    a.Changes,
		UserID:     a.UserID,
		UserName:   a.UserName,
		UserRole:   a.UserRole,
		IPAddress:  a.IPAddress,
		UserAgent:  a.UserAgent,
		RequestID:  a.RequestID,
		StatusCode: a.StatusCode,
		ErrorMsg:   a.ErrorMsg,
		Duration:   a.Duration,
		CreatedAt:  a.CreatedAt,
	}
}

func toAuditLogEntity(doc *auditLogDocument) *entity.AuditLog {
	return &entity.AuditLog{
		ID:         doc.ID.Hex(),
		EntityType: doc.EntityType,
		EntityID:   doc.EntityID,
		Action:     entity.AuditAction(doc.Action),
		OldValue:   doc.OldValue,
		NewValue:   doc.NewValue,
		Changes:    doc.Changes,
		UserID:     doc.UserID,
		UserName:   doc.UserName,
		UserRole:   doc.UserRole,
		IPAddress:  doc.IPAddress,
		UserAgent:  doc.UserAgent,
		RequestID:  doc.RequestID,
		StatusCode: doc.StatusCode,
		ErrorMsg:   doc.ErrorMsg,
		Duration:   doc.Duration,
		CreatedAt:  doc.CreatedAt,
	}
}
