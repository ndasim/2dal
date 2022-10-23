package models

import (
	"2dal/core"
	"2dal/shortener/db"
	"time"

	"github.com/catinello/base62"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
)

type Link struct {
	mgm.DefaultModel `bson:",inline"`
	Origin_url       string
	Alias            string
	From_ts          string
	To_ts            string
	User             User
}

func (link *Link) Create(origin string, alias string, user *User) {
	db.GetConnection()

	// Make sure to pass the model by reference (to update the model's "updated_at", "created_at" and "id" fields by mgm).
	collection := mgm.Coll(link)

	if len(alias) == 0 {
		count, err := collection.CountDocuments(mgm.Ctx(), bson.D{})

		if err != nil {
			panic(err)
		}

		link.Alias = Rhash(int(count))
	} else {
		link.Alias = alias
	}
	
	// Set timestamp of link
	link.From_ts = time.Now().Format(time.RFC3339)
	link.To_ts = time.Now().Add(time.Duration(time.Duration.Hours(1))).Format(time.RFC3339)

	// Set origin of link
	link.Origin_url = origin
	link.User = *user

	print(time.Now().UnixMicro())
	print("\n")
	print(int(time.Now().UnixNano() - core.StartTime))
	print("\n")
	print(base62.Encode(int(time.Now().UnixNano() - core.StartTime)))

	err := collection.Create(link)

	if err != nil {
		panic(err)
	}
}

func Rhash(n int) string {
	return base62.Encode(n * 1881 % 100000)
}

func UNRhash(n int) int {
	return n * 387420489 % 4000000000
}
