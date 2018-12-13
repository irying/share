package serviceTests

import (
	"backend/pb/user"
	"backend/service/user"
	"context"
	"testing"
	. "github.com/smartystreets/goconvey/convey"
)

var service = &user.User{}
var ctx = context.Background()

func Test_FindUserByUid(t *testing.T) {
	request := &pb.ByUidRequest{Uid:1}

	username := "admin"
	record, err := service.FindUserByUID(ctx, request)

	Convey("Subject: Test User Service", t, func() {
		Convey("Should Be Able To FindByUid", func() {
			So(err, ShouldEqual, nil)
		})
		Convey("Should Have Correct Username", func() {
			So(record.Username, ShouldEqual, username)
		})
	})
}


func Test_FindUserByUsernameAndPassword(t *testing.T) {

	request := &pb.ByUserInfoRequest{Username:"admin", Password:"111111"}

	uid := 1
	record, err := service.FindUserByUsernameAndPassword(ctx, request)

	Convey("Subject: Test User Service", t, func() {
		Convey("Should Be Able To FindByUsernameAndPassword", func() {
			So(err, ShouldEqual, nil)
		})
		Convey("Should Have Correct UID", func() {
			So(record.Uid, ShouldEqual, uid)
		})
	})
}

func Test_SaveProfile(t *testing.T) {
	nickname := "我是test"
	pitcureUrl:= "test.png"

	// first update
	request := &pb.ProfileUpdateRequest{Uid:2, Nickname:nickname, PictureUrl:pitcureUrl}
	record, err := service.SaveProfile(ctx, request)

	// senond find
	uidRequest := &pb.ByUidRequest{Uid:2}
	newRecord, _ := service.FindUserByUID(ctx, uidRequest)

	Convey("Subject: Test User Service", t, func() {
		Convey("Should Be Able To SaveProfile", func() {
			So(err, ShouldEqual, nil)
			So(record.Success, ShouldEqual, true)
		})

		Convey("Should Have New Nickname", func() {
			So(newRecord.Nickname, ShouldEqual, nickname)
			So(newRecord.ProfilePictureUrl, ShouldEqual, pitcureUrl)
		})
	})
}

//func Test_UserWorkFlow(t *testing.T)  {
//	t.Run("Update", Test_SaveProfile)
//	t.Run("Find", Test_FindUserByUid)
//}