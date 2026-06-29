package repositories

import (
	"context"
	"errors"
	"time"

	"task-manager/internals/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TaskRepository struct {
	collection *mongo.Collection
}

func NewTaskRepository(client *mongo.Client, dbName, collectionName string) *TaskRepository {
	return &TaskRepository{
		collection: client.Database(dbName).Collection(collectionName),
	}
}

func NewMongoClient(ctx context.Context, uri string) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	return client, nil
}

func (r *TaskRepository) CreateTask(ctx context.Context, task *models.Task) error {
	if task.ID == "" {
		task.ID = primitive.NewObjectID().Hex()
	}
	if task.CreatedAt.IsZero() {
		task.CreatedAt = time.Now().UTC()
	}
	task.UpdatedAt = task.CreatedAt

	_, err := r.collection.InsertOne(ctx, task)
	return err
}

func (r *TaskRepository) ListTasks(ctx context.Context) ([]*models.Task, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tasks []*models.Task
	for cursor.Next(ctx) {
		var task models.Task
		if err := cursor.Decode(&task); err != nil {
			return nil, err
		}
		tasks = append(tasks, &task)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (r *TaskRepository) GetTask(ctx context.Context, id string) (*models.Task, error) {
	var task models.Task
	filter := bson.M{"id": id}
	if err := r.collection.FindOne(ctx, filter).Decode(&task); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &task, nil
}

func (r *TaskRepository) UpdateTask(ctx context.Context, id string, updates *models.UpdateTaskRequest) (*models.Task, error) {
	update := bson.M{"$set": bson.M{"updated_at": time.Now().UTC()}}
	if updates.Title != nil {
		update["$set"].(bson.M)["title"] = *updates.Title
	}
	if updates.Description != nil {
		update["$set"].(bson.M)["description"] = *updates.Description
	}
	if updates.Status != nil {
		update["$set"].(bson.M)["status"] = *updates.Status
	}

	result := r.collection.FindOneAndUpdate(ctx, bson.M{"id": id}, update, options.FindOneAndUpdate().SetReturnDocument(options.After))
	var task models.Task
	if err := result.Decode(&task); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &task, nil
}

func (r *TaskRepository) DeleteTask(ctx context.Context, id string) error {
	result, err := r.collection.DeleteOne(ctx, bson.M{"id": id})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}
