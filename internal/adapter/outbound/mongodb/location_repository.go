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

	"github.com/okinn/service-presensi/internal/domain/entity"
	"github.com/okinn/service-presensi/internal/domain/repository"
)

// allowedLocationDocument adalah representasi MongoDB document
type allowedLocationDocument struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	Name         string             `bson:"name"`
	Latitude     float64            `bson:"latitude"`
	Longitude    float64            `bson:"longitude"`
	RadiusMeters float64            `bson:"radius_meters"`
	Address      string             `bson:"address,omitempty"`
	IsActive     bool               `bson:"is_active"`
	CreatedAt    time.Time          `bson:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at"`
}

type AllowedLocationRepository struct {
	collection *mongo.Collection
}

func NewAllowedLocationRepository(db *mongo.Database) repository.AllowedLocationRepository {
	return &AllowedLocationRepository{
		collection: db.Collection("allowed_locations"),
	}
}

func (r *AllowedLocationRepository) Create(ctx context.Context, location *entity.AllowedLocation) error {
	doc := toLocationDocument(location)
	result, err := r.collection.InsertOne(ctx, doc)
	if err != nil {
		return err
	}

	location.ID = result.InsertedID.(primitive.ObjectID).Hex()
	return nil
}

func (r *AllowedLocationRepository) GetByID(ctx context.Context, id string) (*entity.AllowedLocation, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var doc allowedLocationDocument
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&doc)
	if err != nil {
		return nil, err
	}

	return toLocationEntity(&doc), nil
}

func (r *AllowedLocationRepository) GetAll(ctx context.Context) ([]entity.AllowedLocation, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var docs []allowedLocationDocument
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, err
	}

	locations := make([]entity.AllowedLocation, len(docs))
	for i, doc := range docs {
		locations[i] = *toLocationEntity(&doc)
	}

	return locations, nil
}

func (r *AllowedLocationRepository) GetAllActive(ctx context.Context) ([]entity.AllowedLocation, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"is_active": true})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var docs []allowedLocationDocument
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, err
	}

	locations := make([]entity.AllowedLocation, len(docs))
	for i, doc := range docs {
		locations[i] = *toLocationEntity(&doc)
	}

	return locations, nil
}

func (r *AllowedLocationRepository) Update(ctx context.Context, location *entity.AllowedLocation) error {
	objectID, err := primitive.ObjectIDFromHex(location.ID)
	if err != nil {
		return err
	}

	doc := toLocationDocument(location)
	doc.ID = objectID
	doc.UpdatedAt = time.Now()

	_, err = r.collection.ReplaceOne(ctx, bson.M{"_id": objectID}, doc)
	return err
}

func (r *AllowedLocationRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}

// Helper functions untuk konversi antara entity dan document

func toLocationDocument(l *entity.AllowedLocation) *allowedLocationDocument {
	return &allowedLocationDocument{
		Name:         l.Name,
		Latitude:     l.Latitude,
		Longitude:    l.Longitude,
		RadiusMeters: l.RadiusMeters,
		Address:      l.Address,
		IsActive:     l.IsActive,
		CreatedAt:    l.CreatedAt,
		UpdatedAt:    l.UpdatedAt,
	}
}

func toLocationEntity(doc *allowedLocationDocument) *entity.AllowedLocation {
	return &entity.AllowedLocation{
		ID:           doc.ID.Hex(),
		Name:         doc.Name,
		Latitude:     doc.Latitude,
		Longitude:    doc.Longitude,
		RadiusMeters: doc.RadiusMeters,
		Address:      doc.Address,
		IsActive:     doc.IsActive,
		CreatedAt:    doc.CreatedAt,
		UpdatedAt:    doc.UpdatedAt,
	}
}
