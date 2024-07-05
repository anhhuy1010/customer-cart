package models

import (
	"context"
	// "log"

	"github.com/anhhuy1010/customer-cart/constant"
	"github.com/anhhuy1010/customer-cart/database"

	// "github.com/anhhuy1010/customer-cart/helpers/util"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//"go.mongodb.org/mongo-driver/bson"

type Carts struct {
	Uuid  string  `json:"uuid" bson:"uuid"`
	Total float64 `json:"total" bson:"total"`
}

func (u *Carts) Model() *mongo.Collection {
	db := database.GetInstance()
	return db.Collection("carts")
}

func (u *Carts) Find(conditions map[string]interface{}, opts ...*options.FindOptions) ([]*Carts, error) {
	coll := u.Model()

	conditions["is_delete"] = constant.UNDELETE
	cursor, err := coll.Find(context.TODO(), conditions, opts...)
	if err != nil {
		return nil, err
	}

	var carts []*Carts
	for cursor.Next(context.TODO()) {
		var elem Carts
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

// func (u *Carts) Pagination(ctx context.Context, conditions map[string]interface{}, modelOptions ...ModelOption) ([]*Carts, error) {
// 	coll := u.Model()

// 	conditions["is_delete"] = constant.UNDELETE

// 	modelOpt := ModelOption{}
// 	findOptions := modelOpt.GetOption(modelOptions)
// 	cursor, err := coll.Find(context.TODO(), conditions, findOptions)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var items []*Carts
// 	for cursor.Next(context.TODO()) {
// 		var elem Carts
// 		err := cursor.Decode(&elem)
// 		if err != nil {
// 			log.Println("[Decode] PopularCuisine:", err)
// 			log.Println("-> #", elem.Uuid)
// 			continue
// 		}

// 		items = append(items, &elem)
// 	}

// 	if err := cursor.Err(); err != nil {
// 		return nil, err
// 	}
// 	_ = cursor.Close(context.TODO())

// 	return items, nil
// }

func (u *Carts) Distinct(conditions map[string]interface{}, fieldName string, opts ...*options.DistinctOptions) ([]interface{}, error) {
	coll := u.Model()

	conditions["is_delete"] = constant.UNDELETE

	values, err := coll.Distinct(context.TODO(), fieldName, conditions, opts...)
	if err != nil {
		return nil, err
	}

	return values, nil
}

func (u *Carts) FindOne(conditions map[string]interface{}) (*Carts, error) {
	coll := u.Model()

	conditions["is_delete"] = constant.UNDELETE
	err := coll.FindOne(context.TODO(), conditions).Decode(&u)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (u *Carts) Insert() (interface{}, error) {
	coll := u.Model()

	resp, err := coll.InsertOne(context.TODO(), u)
	if err != nil {
		return 0, err
	}

	return resp, nil
}

func (u *Carts) InsertMany(Items []interface{}) ([]interface{}, error) {
	coll := u.Model()

	resp, err := coll.InsertMany(context.TODO(), Items)
	if err != nil {
		return nil, err
	}

	return resp.InsertedIDs, nil
}

// func (u *Carts) Update() (int64, error) {
// 	coll := u.Model()

// 	condition := make(map[string]interface{})
// 	condition["uuid"] = u.Uuid

// 	u.UpdatedAt = util.GetNowUTC()
// 	updateStr := make(map[string]interface{})
// 	updateStr["$set"] = u

// 	resp, err := coll.UpdateOne(context.TODO(), condition, updateStr)
// 	if err != nil {
// 		return 0, err
// 	}

// 	return resp.ModifiedCount, nil
// }

func (u *Carts) UpdateByCondition(condition map[string]interface{}, data map[string]interface{}) (int64, error) {
	coll := u.Model()

	resp, err := coll.UpdateOne(context.TODO(), condition, data)
	if err != nil {
		return 0, err
	}

	return resp.ModifiedCount, nil
}

func (u *Carts) UpdateMany(conditions map[string]interface{}, updateData map[string]interface{}) (int64, error) {
	coll := u.Model()
	resp, err := coll.UpdateMany(context.TODO(), conditions, updateData)
	if err != nil {
		return 0, err
	}

	return resp.ModifiedCount, nil
}

func (u *Carts) Count(ctx context.Context, condition map[string]interface{}) (int64, error) {
	coll := u.Model()

	condition["is_delete"] = constant.UNDELETE

	total, err := coll.CountDocuments(ctx, condition)
	if err != nil {
		return 0, err
	}

	return total, nil
}
