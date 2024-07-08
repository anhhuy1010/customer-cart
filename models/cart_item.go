package models

import (
	"context"
	"log"
	"time"

	// "log"

	"github.com/anhhuy1010/customer-cart/constant"
	"github.com/anhhuy1010/customer-cart/database"
	"github.com/anhhuy1010/customer-cart/helpers/util"

	// "github.com/anhhuy1010/customer-cart/helpers/util"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//"go.mongodb.org/mongo-driver/bson"

type CartItem struct {
	CartUuid     string    `json:"cart_uuid" bson:"cart_uuid"`
	ProductUuid  string    `json:"product_uuid" bson:"product_uuid"`
	CartItemUuid string    `json:"cart_item_uuid" bson:"cart_item_uuid"`
	ProductName  string    `json:"product_name" bson:"product_name"`
	ProductPrice float64   `json:"product_price" bson:"product_price"`
	Quantity     int       `json:"quantity" bson:"quantity"`
	ProductTotal float64   `json:"product_total" bson:"product_total"`
	IsDelete     int       `json:"is_delete" bson:"is_delete"`
	IsActive     int       `json:"is_active" bson:"is_active"`
	CreatedAt    time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" bson:"updated_at"`
}

func (u *CartItem) Model() *mongo.Collection {
	db := database.GetInstance()
	return db.Collection("cart_items")
}

func (u *CartItem) Find(conditions map[string]interface{}, opts ...*options.FindOptions) ([]*CartItem, error) {
	coll := u.Model()

	conditions["is_delete"] = constant.UNDELETE
	cursor, err := coll.Find(context.TODO(), conditions, opts...)
	if err != nil {
		return nil, err
	}

	var carts []*CartItem
	for cursor.Next(context.TODO()) {
		var elem CartItem
		err := cursor.Decode(&elem)
		if err != nil {
			return nil, err
		}

		carts = append(carts, &elem)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}
	_ = cursor.Close(context.TODO())

	return carts, nil
}

func (u *Carts) Pagination(ctx context.Context, conditions map[string]interface{}, modelOptions ...ModelOption) ([]*Carts, error) {
	coll := u.Model()

	conditions["is_delete"] = constant.UNDELETE

	modelOpt := ModelOption{}
	findOptions := modelOpt.GetOption(modelOptions)
	cursor, err := coll.Find(context.TODO(), conditions, findOptions)
	if err != nil {
		return nil, err
	}

	var items []*Carts
	for cursor.Next(context.TODO()) {
		var elem Carts
		err := cursor.Decode(&elem)
		if err != nil {
			log.Println("[Decode] PopularCuisine:", err)
			log.Println("-> #", elem.CartUuid)
			continue
		}

		items = append(items, &elem)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}
	_ = cursor.Close(context.TODO())

	return items, nil
}

func (u *CartItem) Distinct(conditions map[string]interface{}, fieldName string, opts ...*options.DistinctOptions) ([]interface{}, error) {
	coll := u.Model()

	conditions["is_delete"] = constant.UNDELETE

	values, err := coll.Distinct(context.TODO(), fieldName, conditions, opts...)
	if err != nil {
		return nil, err
	}

	return values, nil
}

func (u *CartItem) FindOne(conditions map[string]interface{}) (*CartItem, error) {
	coll := u.Model()

	conditions["is_delete"] = constant.UNDELETE
	err := coll.FindOne(context.TODO(), conditions).Decode(&u)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (u *CartItem) Insert() (interface{}, error) {
	coll := u.Model()

	resp, err := coll.InsertOne(context.TODO(), u)
	if err != nil {
		return 0, err
	}

	return resp, nil
}

func (u *CartItem) InsertMany(Items []interface{}) ([]interface{}, error) {
	coll := u.Model()

	resp, err := coll.InsertMany(context.TODO(), Items)
	if err != nil {
		return nil, err
	}

	return resp.InsertedIDs, nil
}

func (u *CartItem) Update() (int64, error) {
	coll := u.Model()

	condition := make(map[string]interface{})
	condition["cart_uuid"] = u.CartUuid

	u.UpdatedAt = util.GetNowUTC()
	updateStr := make(map[string]interface{})
	updateStr["$set"] = u

	resp, err := coll.UpdateOne(context.TODO(), condition, updateStr)
	if err != nil {
		return 0, err
	}

	return resp.ModifiedCount, nil
}

func (u *CartItem) UpdateByCondition(condition map[string]interface{}, data map[string]interface{}) (int64, error) {
	coll := u.Model()

	resp, err := coll.UpdateOne(context.TODO(), condition, data)
	if err != nil {
		return 0, err
	}

	return resp.ModifiedCount, nil
}

func (u *CartItem) UpdateMany(conditions map[string]interface{}, updateData map[string]interface{}) (int64, error) {
	coll := u.Model()
	resp, err := coll.UpdateMany(context.TODO(), conditions, updateData)
	if err != nil {
		return 0, err
	}

	return resp.ModifiedCount, nil
}

func (u *CartItem) Count(ctx context.Context, condition map[string]interface{}) (int64, error) {
	coll := u.Model()

	condition["is_delete"] = constant.UNDELETE

	total, err := coll.CountDocuments(ctx, condition)
	if err != nil {
		return 0, err
	}

	return total, nil
}
