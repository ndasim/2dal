package models

import (
	"errors"
	"time"

	"github.com/catinello/base62"
	"github.com/getsentry/sentry-go"
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

func (link *Link) Create(origin string, alias string, user *User) error {
	// Make sure to pass the model by reference (to update the model's "updated_at", "created_at" and "id" fields by mgm).
	collection := mgm.Coll(link)

	if len(alias) == 0 {
		count, err := collection.CountDocuments(mgm.Ctx(), bson.D{})

		if err != nil {
			return err
		}

		link.Alias = Rhash(int(count))
	} else {
		link.Alias = alias
	}

	duration := time.Hour * time.Duration(24)
	if user.Subscription == "BASIC" {
		duration = time.Hour * time.Duration(24)
	} else if user.Subscription == "PRO" || user.Subscription == "ULTRA" || user.Subscription == "MEGA" || user.Subscription == "Buy me coffee" {
		duration = time.Hour * time.Duration(24*365*10)
	}

	// Set timestamp of link
	link.From_ts = time.Now().Format(time.RFC3339)
	link.To_ts = time.Now().Add(duration).Format(time.RFC3339)

	// Set origin of link
	link.Origin_url = origin
	link.User = *user

	//print(time.Now().UnixMicro())
	//print("\n")
	//print(int(time.Now().UnixNano() - core.StartTime))
	//print("\n")
	//print(base62.Encode(int(time.Now().UnixNano() - core.StartTime)))

	if err := collection.Create(link); err != nil {
		return err
	}

	return nil
}

func (link *Link) Update(alias string, user *User) error {
	collection := mgm.Coll(link)

	duration := time.Hour * time.Duration(24)
	if user.Subscription == "BASIC" {
		duration = time.Hour * time.Duration(24)
	} else if user.Subscription == "PRO" || user.Subscription == "ULTRA" || user.Subscription == "MEGA" || user.Subscription == "Buy me coffee" {
		duration = time.Hour * time.Duration(24*365*10)
	}

	print(duration)

	// Set timestamp of link
	link.From_ts = time.Now().Format(time.RFC3339)
	link.To_ts = time.Now().Add(duration).Format(time.RFC3339)

	if alias != "" {
		link.Alias = alias
	}

	if err := collection.Update(link); err != nil {
		return err
	}

	return nil
}

func (link *Link) FindLink(alias string) error {
	model := mgm.Coll(link)

	result := []Link{}

	err := model.SimpleFind(&result, bson.D{{Key: "alias", Value: alias}})
	if err != nil {
		return err
	}

	if len(result) > 1 {
		sentry.CaptureMessage("alias: " + alias + " has returned more than one link!")
	} else if len(result) == 0 {
		return errors.New("not found")
	}

	link.Origin_url = result[0].Origin_url
	link.Alias = result[0].Alias
	link.CreatedAt = result[0].CreatedAt
	link.From_ts = result[0].From_ts
	link.To_ts = result[0].To_ts
	link.User = result[0].User

	return nil
}

func (link *Link) FindOrigin(origin string) error {
	model := mgm.Coll(link)

	result := []Link{}

	err := model.SimpleFind(&result, bson.D{{Key: "origin_url", Value: origin}})
	if err != nil {
		return err
	}

	if len(result) > 1 {
		sentry.CaptureMessage("origin_url: " + origin + " has returned more than one link!")
	} else if len(result) == 0 {
		return errors.New("no data")
	}

	link.ID = result[0].ID
	link.Origin_url = result[0].Origin_url
	link.Alias = result[0].Alias
	link.CreatedAt = result[0].CreatedAt
	link.UpdatedAt = result[0].UpdatedAt
	link.From_ts = result[0].From_ts
	link.To_ts = result[0].To_ts
	link.User = result[0].User

	return nil
}

func Rhash(n int) string {
	return base62.Encode(n * 1881 % 100000)
}

func UNRhash(n int) int {
	return n * 387420489 % 4000000000
}
