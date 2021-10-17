package repository

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"sber/types"
)

const DbEmployees = "employees"
const DbTableUsers = "users"

type UserMongo struct {
	db *mongo.Database
}

func NewUserMongo(client *mongo.Client) *UserMongo {
	return &UserMongo{
		db: client.Database(DbEmployees),
	}
}

func (u *UserMongo) GetById(ctx context.Context, userId string) (*types.User, error) {

	var users []types.User

	pipeline := mongo.Pipeline{
		u.getMatch(userId),
		u.getGraphLookup(),
		u.getProject(),
		u.getUnwind(),
		u.getSort(),
		u.getGroup(),
		u.getAddFields(bson.D{
			{"result", u.getReduce()},
		}),
		u.getAddFields(bson.D{
			{"result", "$result.currentLevelEmployees"},
		}),
	}

	cursor, err := u.db.Collection(DbTableUsers).Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("aggregate error: %v", err)
	}
	defer cursor.Close(ctx)

	err = cursor.All(ctx, &users)
	if err != nil {
		return nil, err
	}

	if len(users) != 1 {
		if len(users) == 0 {
			return nil, fmt.Errorf("Employee not found. ")
		}
		return nil, fmt.Errorf("Result array`s length is not 1. ")
	}

	return &users[0], err
}

func (u *UserMongo) getMatch(userId string) bson.D {
	return bson.D{{
		"$match", bson.D{
			{"_id", userId},
		},
	}}
}

func (u *UserMongo) getGraphLookup() bson.D {
	return bson.D{{
		"$graphLookup", bson.D{
			{"from", DbTableUsers},
			{"startWith", "$_id"},
			{"connectFromField", "_id"},
			{"connectToField", "parent"},
			{"as", "result"},
			{"depthField", "level"},
		},
	}}
}
func (u *UserMongo) getProject() bson.D {
	return bson.D{{
		"$project", bson.D{
			{"name", 1},
			{"employedat", 1},
			{"job", 1},
			{"parent", 1},
			{"result.name", 1},
			{"result.employedat", 1},
			{"result.job", 1},
			{"result.parent", 1},
			{"result.level", 1},
			{"result._id", 1},
		},
	}}
}
func (u *UserMongo) getUnwind() bson.D {
	return bson.D{{
		"$unwind", bson.D{
			{"path", "$result"},
			{"preserveNullAndEmptyArrays", true},
		},
	}}
}
func (u *UserMongo) getSort() bson.D {
	return bson.D{{
		"$sort", bson.D{
			{"result.level", -1},
		},
	}}
}
func (u *UserMongo) getGroup() bson.D {
	return bson.D{{
		"$group", bson.D{
			{"_id", "$_id"},
			{"parent", bson.D{
				{"$first", "$parent"},
			}},
			{"name", bson.D{
				{"$first", "$name"},
			}},
			{"job", bson.D{
				{"$first", "$job"},
			}},
			{"employedat", bson.D{
				{"$first", "$employedat"},
			}},
			{"result", bson.D{
				{"$push", "$result"},
			}},
		},
	}}
}
func (u *UserMongo) getAddFields(bs bson.D) bson.D {
	return bson.D{{
		"$addFields", bs,
	}}
}
func (u *UserMongo) getVars() bson.E {
	return bson.E{Key: "vars", Value: bson.D{
		{"prev", bson.D{
			{"$cond", []interface{}{
				bson.D{
					{"$eq", []interface{}{
						"$$value.currentLevel", "$$this.level",
					}},
				},
				"$$value.previousLevelEmployees",
				"$$value.currentLevelEmployees",
			}},
		}},
		{"current", bson.D{
			{"$cond", []interface{}{
				bson.D{
					{"$eq", []interface{}{
						"$$value.currentLevel", "$$this.level",
					}},
				},
				"$$value.currentLevelEmployees",
				[]interface{}{},
			}},
		}},
	}}
}
func (u *UserMongo) getReduce() bson.D {
	return bson.D{
		{"$reduce", bson.D{
			{"input", "$result"},
			{"initialValue", bson.D{
				{"currentLevel", -1},
				{"currentLevelEmployees", []interface{}{}},
				{"previousLevelEmployees", []interface{}{}},
			}},
			{"in", bson.D{
				{"$let", bson.D{
					u.getVars(),
					{"in", bson.D{
						{"currentLevel", "$$this.level"},
						{"previousLevelEmployees", "$$prev"},
						{"currentLevelEmployees", u.getCurrentLevel()},
					}},
				}},
			}},
		}},
	}
}

func (u *UserMongo) getCurrentLevel() bson.D {
	return bson.D{
		{"$concatArrays", []interface{}{
			"$$current",
			[]interface{}{
				bson.D{
					{"$mergeObjects", []interface{}{
						"$$this",
						bson.D{
							{"result", bson.D{
								{"$filter", bson.D{
									{"input", "$$prev"},
									{"as", "e"},
									{"cond", bson.D{
										{"$eq", []interface{}{
											"$$e.parent",
											"$$this._id",
										}},
									}},
								}},
							}},
						},
					}},
				},
			},
		}},
	}
}
