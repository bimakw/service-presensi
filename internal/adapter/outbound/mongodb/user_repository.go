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

type userDocument struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Email     string             `bson:"email"`
	Password  string             `bson:"password"`
	Nama      string             `bson:"nama"`
	Role      string             `bson:"role"`
	IsActive  bool               `bson:"is_active"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}

type UserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) repository.UserRepository {
	return &UserRepository{
		collection: db.Collection("users"),
	}
}

func (r *UserRepository) Create(ctx context.Context, user *entity.User) error {
	doc := toUserDocument(user)
	result, err := r.collection.InsertOne(ctx, doc)
	if err != nil {
		return err
	}

	user.ID = result.InsertedID.(primitive.ObjectID).Hex()
	return nil
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*entity.User, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var doc userDocument
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&doc)
	if err != nil {
		return nil, err
	}

	return toUserEntity(&doc), nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	var doc userDocument
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&doc)
	if err != nil {
		return nil, err
	}

	return toUserEntity(&doc), nil
}

func (r *UserRepository) Update(ctx context.Context, user *entity.User) error {
	objectID, err := primitive.ObjectIDFromHex(user.ID)
	if err != nil {
		return err
	}

	doc := toUserDocument(user)
	doc.ID = objectID
	doc.UpdatedAt = time.Now()

	_, err = r.collection.ReplaceOne(ctx, bson.M{"_id": objectID}, doc)
	return err
}

func (r *UserRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}

func toUserDocument(u *entity.User) *userDocument {
	return &userDocument{
		Email:     u.Email,
		Password:  u.Password,
		Nama:      u.Nama,
		Role:      string(u.Role),
		IsActive:  u.IsActive,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func toUserEntity(doc *userDocument) *entity.User {
	return &entity.User{
		ID:        doc.ID.Hex(),
		Email:     doc.Email,
		Password:  doc.Password,
		Nama:      doc.Nama,
		Role:      entity.UserRole(doc.Role),
		IsActive:  doc.IsActive,
		CreatedAt: doc.CreatedAt,
		UpdatedAt: doc.UpdatedAt,
	}
}
