package dal

import (
	"log"
	"neo/config"
	"time"

	mgo "gopkg.in/mgo.v2"
)

var (
	mongoDBDialInfo *mgo.DialInfo
)

func getDBName(config *config.Config) string {
	return config.MongoDB.Database
}

func getDBHost(config *config.Config) string {
	return config.MongoDB.Host
}

// Open open database connection
func Open(config *config.Config) (interface{}, error) {
	// no auth needed
	mongoDBDialInfo = &mgo.DialInfo{
		Addrs:    []string{getDBHost(config)},
		Timeout:  15 * time.Second,
		Database: getDBName(config),
	}

	mongoSession, err := mgo.DialWithInfo(mongoDBDialInfo)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	// Optional. Switch the session to a monotonic behavior.
	mongoSession.SetMode(mgo.Monotonic, true)
	return mongoSession.DB(getDBName(config)), nil
}
