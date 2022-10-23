package models

import "github.com/kamva/mgm/v3"

type Activity struct {
	mgm.DefaultModel `bson:",inline"`
	activity string         // activiy name
	ts       int            // activity timestamp
	extra    map[string]int // extra data of activity
}
