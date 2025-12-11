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
	"github.com/okinn/service-presensi/internal/domain/valueobject"
)

// presensiDocument adalah representasi MongoDB document
type presensiDocument struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	UserID     string             `bson:"user_id"`
	Nama       string             `bson:"nama"`
	Tanggal    time.Time          `bson:"tanggal"`
	JamMasuk   *time.Time         `bson:"jam_masuk,omitempty"`
	JamKeluar  *time.Time         `bson:"jam_keluar,omitempty"`
	Status     string             `bson:"status"`
	Keterangan string             `bson:"keterangan,omitempty"`
	Lokasi     *lokasiDocument    `bson:"lokasi,omitempty"`
	CreatedAt  time.Time          `bson:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at"`
}

type lokasiDocument struct {
	Latitude  float64 `bson:"latitude"`
	Longitude float64 `bson:"longitude"`
	Alamat    string  `bson:"alamat,omitempty"`
}

type PresensiRepository struct {
	collection *mongo.Collection
}

func NewPresensiRepository(db *mongo.Database) repository.PresensiRepository {
	return &PresensiRepository{
		collection: db.Collection("presensi"),
	}
}

func (r *PresensiRepository) Create(ctx context.Context, presensi *entity.Presensi) error {
	doc := toDocument(presensi)
	result, err := r.collection.InsertOne(ctx, doc)
	if err != nil {
		return err
	}

	presensi.ID = result.InsertedID.(primitive.ObjectID).Hex()
	return nil
}

func (r *PresensiRepository) GetByID(ctx context.Context, id string) (*entity.Presensi, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var doc presensiDocument
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&doc)
	if err != nil {
		return nil, err
	}

	return toEntity(&doc), nil
}

func (r *PresensiRepository) GetAll(ctx context.Context, filter repository.PresensiFilter, page, limit int) ([]entity.Presensi, int64, error) {
	bsonFilter := bson.M{}

	if filter.UserID != "" {
		bsonFilter["user_id"] = filter.UserID
	}
	if filter.Status != "" {
		bsonFilter["status"] = string(filter.Status)
	}
	if !filter.StartDate.IsZero() || !filter.EndDate.IsZero() {
		dateFilter := bson.M{}
		if !filter.StartDate.IsZero() {
			dateFilter["$gte"] = filter.StartDate
		}
		if !filter.EndDate.IsZero() {
			dateFilter["$lte"] = filter.EndDate
		}
		bsonFilter["tanggal"] = dateFilter
	}

	total, err := r.collection.CountDocuments(ctx, bsonFilter)
	if err != nil {
		return nil, 0, err
	}

	skip := int64((page - 1) * limit)
	opts := options.Find().
		SetSkip(skip).
		SetLimit(int64(limit)).
		SetSort(bson.D{{Key: "tanggal", Value: -1}})

	cursor, err := r.collection.Find(ctx, bsonFilter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var docs []presensiDocument
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, 0, err
	}

	presensiList := make([]entity.Presensi, len(docs))
	for i, doc := range docs {
		presensiList[i] = *toEntity(&doc)
	}

	return presensiList, total, nil
}

func (r *PresensiRepository) Update(ctx context.Context, presensi *entity.Presensi) error {
	objectID, err := primitive.ObjectIDFromHex(presensi.ID)
	if err != nil {
		return err
	}

	doc := toDocument(presensi)
	doc.ID = objectID
	doc.UpdatedAt = time.Now()

	_, err = r.collection.ReplaceOne(ctx, bson.M{"_id": objectID}, doc)
	return err
}

func (r *PresensiRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}

// Helper functions untuk konversi antara entity dan document

func toDocument(p *entity.Presensi) *presensiDocument {
	doc := &presensiDocument{
		UserID:     p.UserID,
		Nama:       p.Nama,
		Tanggal:    p.Tanggal,
		JamMasuk:   p.JamMasuk,
		JamKeluar:  p.JamKeluar,
		Status:     string(p.Status),
		Keterangan: p.Keterangan,
		CreatedAt:  p.CreatedAt,
		UpdatedAt:  p.UpdatedAt,
	}

	if p.Lokasi != nil {
		doc.Lokasi = &lokasiDocument{
			Latitude:  p.Lokasi.Latitude,
			Longitude: p.Lokasi.Longitude,
			Alamat:    p.Lokasi.Alamat,
		}
	}

	return doc
}

func toEntity(doc *presensiDocument) *entity.Presensi {
	p := &entity.Presensi{
		ID:         doc.ID.Hex(),
		UserID:     doc.UserID,
		Nama:       doc.Nama,
		Tanggal:    doc.Tanggal,
		JamMasuk:   doc.JamMasuk,
		JamKeluar:  doc.JamKeluar,
		Status:     valueobject.StatusPresensi(doc.Status),
		Keterangan: doc.Keterangan,
		CreatedAt:  doc.CreatedAt,
		UpdatedAt:  doc.UpdatedAt,
	}

	if doc.Lokasi != nil {
		p.Lokasi = &valueobject.Lokasi{
			Latitude:  doc.Lokasi.Latitude,
			Longitude: doc.Lokasi.Longitude,
			Alamat:    doc.Lokasi.Alamat,
		}
	}

	return p
}
