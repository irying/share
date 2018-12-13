package user

import (
	"backend/api/middleware"
	"backend/api/proxy"
	"backend/pb/user"
	"backend/system/exception"
	"backend/system/upload"
	"backend/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
	"strconv"
)
type loginRequest struct {
	Username string `form:"username"  binding:"required"`
	Password string	`form:"password"  binding:"required"`
}

func Login(c *gin.Context) {
	data := map[string]interface{}{}
	req := &loginRequest{}
	if err := c.ShouldBindWith(req, binding.Form); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    exception.INVALID_PARAMS,
			"message": exception.GetMsg(exception.INVALID_PARAMS),
			"data":    data,
		})
		return
	}
	username := c.PostForm("username")
	password := c.PostForm("password")
	res, _ := pb.NewUserServiceClient(proxy.NewRPCConn()).FindUserByUsernameAndPassword(c, &pb.ByUserInfoRequest{Username: username, Password: password})

	
	if 0 == res.Uid {
		c.JSON(http.StatusOK, gin.H{
			"code":    exception.ERROR_USER_LOGIN_FAIL,
			"message": exception.GetMsg(exception.ERROR_USER_LOGIN_FAIL),
			"data":    data,
		})
	} else {
		token, err := middleware.GenerateToken(username, password)
		if err != nil {
			utils.LogPrintError(err)
		}
		saveToken(c, token, username)

		data["token"] = token
		data["uid"] = res.Uid
		c.JSON(http.StatusOK, gin.H{
			"code":    exception.SUCCESS,
			"message": "success",
			"data":    data,
		})
	}
}

func saveToken(c *gin.Context, token string, username string) {
	pb.NewUserServiceClient(proxy.NewRPCConn()).SaveToken(c, &pb.TokenRequest{Username: username, Token:token})
}

func Index(c *gin.Context) {
	queryUid := c.Query("uid")
	token := c.Query("token")
	uid, err := strconv.Atoi(queryUid)
	utils.LogPrintError(err)
	res, _ := pb.NewUserServiceClient(proxy.NewRPCConn()).FindUserByUID(c, &pb.ByUidRequest{Uid: int64(uid)})

	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"uid":					queryUid,
		"username":            res.Username,
		"nickname":            res.Nickname,
		"profile_picture_url": res.ProfilePictureUrl,
		"token":               token,
	})
}

func Profile(c *gin.Context) {
	queryUid := c.Query("uid")
	//token := c.Query("token")
	uid, err := strconv.Atoi(queryUid)
	utils.LogPrintError(err)
	res, _ := pb.NewUserServiceClient(proxy.NewRPCConn()).FindUserByUID(c, &pb.ByUidRequest{Uid: int64(uid)})

	c.JSON(http.StatusOK, gin.H{
		"code":    exception.SUCCESS,
		"message": "success",
		"data":    res,
	})
}


func UpdateProfile(c *gin.Context) {
	queryUid := c.Query("uid")
	token := c.Query("token")
	uid, err := strconv.Atoi(queryUid)
	nickname := c.PostForm("nickname")

	imageUrl, _, err := getUrlAfterImageUploaded(c)
	backUrl := fmt.Sprintf("/backend/index?uid=%d&token=%s", uid, token)
	if err != nil {
		utils.LogPrintError(err)

		c.Redirect(http.StatusFound, backUrl)
	} else {
		pb.NewUserServiceClient(proxy.NewRPCConn()).SaveProfile(c, &pb.ProfileUpdateRequest{Uid: int64(uid), Nickname:nickname, PictureUrl:imageUrl})
		c.Redirect(http.StatusFound, backUrl)
	}

}
func getUrlAfterImageUploaded(c *gin.Context) (string, int, error) {
	url := ""

	file, image, err := c.Request.FormFile("profile_picture")

	if err != nil {
		return url, exception.ERROR_UPLOAD_IMAGE_FAIL, err
	}
	if image == nil {
		return url, exception.SUCCESS, nil
	}
	imageName := upload.GetImageName(image.Filename)
	src := upload.GetImagePath() + imageName

	if ! upload.CheckImageSize(file) {
		return url, exception.ERROR_UPLOAD_CHECK_IMAGE_FORMAT, nil
	}
	if err := c.SaveUploadedFile(image, src); err != nil {
		fmt.Println(err)
		return url, exception.ERROR_UPLOAD_SAVE_IMAGE_FAIL, nil
	}

	return imageName, exception.SUCCESS, nil
}
