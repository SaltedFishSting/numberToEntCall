// num2ent_test
package main

import (
	"fmt"

	"testing"
	"time"

	"gopkg.in/mgo.v2"
)

func TestUser(t *testing.T) {
	session, err := mgo.Dial("103.25.23.104:20002")
	collection := session.DB("meetingDB").C("MeetingInfoSummary_FourTable")
	defer session.Close()

	if err != nil {
		fmt.Println(err)
	}
	tt := time.Now().Unix()
	fmt.Println(tt)

	getmeetqualifieds(collection, tt-300)

}
