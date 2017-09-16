package core

import (
	"time"

	mgo "gopkg.in/mgo.v2"
)

// Rate rate info
type Rate struct {
	Neo *Neo
	Btc *Btc
}

// Neo Neo info
type Neo struct {
	Btc     float32   `bson:"btc"`
	BtcRate float32   `bson:"btcrate"`
	Usd     float32   `bson:"usd"`
	UsdRate float32   `bson:"usdrate"`
	Date    time.Time `bson:"date"`
}

// Btc Btc info
type Btc struct {
	Buy  float32   `json:"buy"`
	Sell float32   `json:"sell"`
	Date time.Time `bson:"date"`
}

// RateAccessor accessor
type RateAccessor struct {
}

// Description description of accessor
func (accessor *RateAccessor) Description() string {
	return "Neo and Btc info accessor"
}

// GetLatestRateInfo list latest rate info
func (accessor *RateAccessor) GetLatestRateInfo(db interface{}) (*Rate, error) {
	rate := &Rate{}
	var neo *Neo
	var btc *Btc
	collection := db.(*mgo.Database).C("rates")
	err := collection.Find(nil).Sort("-date").One(&neo)
	if err != nil {
		return nil, err
	}

	rate.Neo = neo

	collection = db.(*mgo.Database).C("btcrates")
	err = collection.Find(nil).Sort("-date").One(&btc)
	if err != nil {
		return nil, err
	}

	rate.Btc = btc
	return rate, nil
}
