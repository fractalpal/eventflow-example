package store

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/fractalpal/eventflow-example/log"
	"github.com/fractalpal/eventflow-example/payment-query/app"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoCollection struct {
	fields     logrus.Fields
	collection *mongo.Collection
}

func NewMongoCollection(ctx context.Context, sourceUrl, database, collection string, timeout time.Duration) (col *MongoCollection, err error) {
	client, err := connect(ctx, sourceUrl, timeout)
	if err != nil {
		return
	}
	fields := logrus.Fields{}
	fields["store"] = []string{"mongo"}
	fields["collection"] = []string{collection}
	return &MongoCollection{
		fields:     fields,
		collection: client.Database(database).Collection(collection),
	}, err
}

func (s *MongoCollection) Insert(ctx context.Context, payment app.Payment) (err error) {
	defer log.AddFieldsForErr(ctx, s.fields, err)
	// ensure it's not exists
	if _, err = s.findOne(ctx, payment.ID); err != nil {
		// if no results, continue
		if err != mongo.ErrNoDocuments {
			return
		}
	} else {
		return app.ErrExists
	}

	result, err := s.collection.InsertOne(ctx, &payment)
	if err != nil {
		return
	}

	l := log.FromContext(ctx)
	l = l.WithFields(s.fields).WithField("op", []interface{}{"created", result.InsertedID})
	l.Debug("created in store")
	return
}

func (s *MongoCollection) UpdateThirdParty(ctx context.Context, thirdParty app.ThirdParty, partyKey string) (err error) {
	defer log.AddFieldsForErr(ctx, s.fields, err)

	var result *mongo.UpdateResult
	filter := bson.D{{"id", thirdParty.PaymentID}, {"last_update", bson.M{"$lt": thirdParty.Time}}}
	update := bson.M{"$set": bson.M{"attributes." + partyKey: thirdParty.ThirdParty, "last_update": thirdParty.Time}}
	result, err = s.collection.UpdateOne(ctx, filter, update, &options.UpdateOptions{})
	if err != nil {
		return
	}
	l := log.FromContext(ctx)
	fields := logrus.Fields{}
	fields["op"] = "updated"
	fields["party"] = "party"
	fields["count"] = result.ModifiedCount
	fields["match_count"] = result.MatchedCount
	l = l.WithFields(s.fields).WithFields(fields)

	if result.ModifiedCount == 0 {
		l.Debug("not updated. Event from the past?")
		return
	}

	l.Debug("updated in store")
	return
}

func (s *MongoCollection) FindByID(ctx context.Context, id string) (p app.Payment, err error) {
	defer log.AddFieldsForErr(ctx, s.fields, err)

	if p, err = s.findOne(ctx, id); err != nil {
		if err == mongo.ErrNoDocuments {
			err = app.ErrNoResults
			return
		}
	}

	l := log.FromContext(ctx)
	fields := logrus.Fields{}
	fields["op"] = "find_by_id"
	fields["id"] = id
	l.Debug("found in store")
	return
}

func (s *MongoCollection) FindAll(ctx context.Context, page int64, limit int64) (p []app.Payment, err error) {
	defer log.AddFieldsForErr(ctx, s.fields, err)

	lim := limit
	skip := int64((page - 1) * limit)
	opts := &options.FindOptions{
		Limit: &lim,
		Skip:  &skip,
		//Sort: nil,
	}
	cursor, err := s.collection.Find(ctx, bson.D{}, opts)
	if err != nil {
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var payment app.Payment
		if err = cursor.Decode(&payment); err != nil {
			if err == mongo.ErrNoDocuments {
				err = app.ErrNoResults
				return
			}
		}
		p = append(p, payment)
	}

	l := log.FromContext(ctx)
	fields := logrus.Fields{}
	fields["op"] = "find_all"
	fields["count"] = len(p)
	l = l.WithFields(s.fields).WithFields(fields)

	l.Debug("found list in store")
	return
}

func (s *MongoCollection) Delete(ctx context.Context, id string) (err error) {
	defer log.AddFieldsForErr(ctx, s.fields, err)

	result, err := s.collection.DeleteOne(ctx, bson.M{"id": id})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			err = app.ErrNoResults
			return
		}
		return
	}

	l := log.FromContext(ctx)
	fields := logrus.Fields{}
	fields["op"] = "delete"
	fields["count"] = result.DeletedCount
	l = l.WithFields(s.fields).WithFields(fields)

	l.Debug("deleted from store")
	return
}

func connect(ctx context.Context, sourceUrl string, timeout time.Duration) (client *mongo.Client, err error) {
	client, err = mongo.NewClient(options.Client().ApplyURI(sourceUrl))
	if err != nil {
		return
	}
	conCtx, _ := context.WithTimeout(ctx, timeout)
	if err = client.Connect(conCtx); err != nil {
		return
	}
	if err = client.Ping(ctx, nil); err != nil {
		return
	}
	return
}

func (s *MongoCollection) findOne(ctx context.Context, id string) (p app.Payment, err error) {
	findRes := s.collection.FindOne(ctx, bson.M{"id": id})
	if findRes != nil {
		if err = findRes.Err(); err != nil {
			return
		}
		if err = findRes.Decode(&p); err != nil {
			return
		}
	}
	return
}
