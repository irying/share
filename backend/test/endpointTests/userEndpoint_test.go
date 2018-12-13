package endpointTests

import (
	"encoding/json"
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

// LoginData
type LoginData struct {
	Token string `json:"token"`
	Uid   int    `json:"uid"`
}

type LoginResponse struct {
	Code    string    `json:"code"`
	Message string    `json:"message"`
	Data    LoginData `json:"data"`
}

func TestOnLoginRequestForForm(t *testing.T) {
	// 初始化请求地址和请求参数
	uri := "/api/login"
	uid := 1

	param := make(map[string]string)
	param["username"] = "admin"
	param["password"] = "111111"

	// 发起post请求，以表单形式传递参数
	body := PostForm(uri, param, router)
	fmt.Printf("response:%v\n", string(body))

	// 解析响应，判断响应是否与预期一致
	response := &LoginResponse{}
	err := json.Unmarshal(body, response);
	Convey("Subject: Test Login Api", t, func() {
		Convey("Should Be Able To login", func() {
			So(err, ShouldEqual, nil)
		})
		Convey("Should Have Correct UID", func() {
			So(response.Data.Uid, ShouldEqual, uid)
		})
	})
}
