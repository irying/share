package user

import (
	"backend/pb/user"
	"backend/storage"
	"backend/utils"
	"context"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

// TOKEN_PREFIX Redis key prefix for token
const TOKEN_PREFIX = "token_"

// User user row
type User struct {
	UID               int64  `db:"uid"`
	Username          string `db:"username"`
	Nickname          string `db:"nickname"`
	Password          string `db:"password"`
	CreateTime        int    `db:"create_time"`
	UpdateTime        int    `db:"update_time"`
	ProfilePictureURL string `db:"profile_picture_url"`
}

// FindUserByUID find by uid
func (u *User) FindUserByUID(ctx context.Context, in *pb.ByUidRequest) (*pb.UserResponse, error) {
	conn := storage.NewDBConn("")
	defer conn.Close()

	user := User{}
	//query := fmt.Sprintf("select * from user where uid = %s", in.UID)
	//err := conn.QueryRowx("select * from user where uid = ?", in.UID).StructScan(user)
	err := conn.Get(&user, "select * from user where uid = ?", in.Uid)
	utils.LogPrintError(err)
	return &pb.UserResponse{
		Uid:               user.UID,
		Username:          user.Username,
		Nickname:          user.Nickname,
		ProfilePictureUrl: user.ProfilePictureURL,
	}, nil
}

// FindUserByUsernameAndPassword find by username and password
func (u *User) FindUserByUsernameAndPassword(ctx context.Context, in *pb.ByUserInfoRequest) (*pb.UserResponse, error) {
	conn := storage.NewDBConn("")
	defer conn.Close()

	user := &User{}
	err := conn.QueryRowx("select * from user where username = ? and password = ?", in.Username, in.Password).StructScan(user)
	utils.LogPrintError(err)
	return &pb.UserResponse{
		Uid:               user.UID,
		Username:          user.Username,
		Nickname:          user.Nickname,
		ProfilePictureUrl: user.ProfilePictureURL,
	}, nil
}

// SaveProfile save profile
func (u *User) SaveProfile(ctx context.Context, in *pb.ProfileUpdateRequest) (*pb.SaveResponse, error) {
	conn := storage.NewDBConn("")
	defer conn.Close()

	query := getUpdateQueryByRequest(in)
	_, err := conn.Exec(query)
	utils.LogPrintError(err)
	if err != nil {
		return &pb.SaveResponse{
			Success: false,
		}, nil
	}

	return &pb.SaveResponse{
		Success: true,
	}, nil
}

// return update query
func getUpdateQueryByRequest(in *pb.ProfileUpdateRequest) string {
	query := ""
	if in.Nickname != "" && in.PictureUrl != "" {
		query = fmt.Sprintf("update user set `nickname` = '%s' ,`profile_picture_url` = '%s' where uid= %d", in.Nickname,
			in.PictureUrl, in.Uid)
	} else if in.Nickname != "" {
		query = fmt.Sprintf("update user set nickname = %s where uid= %s", in.Nickname, in.Uid)
	} else if in.PictureUrl != "" {
		query = fmt.Sprintf("update user set profile_picture_url = %s where uid= %s", in.PictureUrl, in.Uid)
	}

	return query
}

// IsExpiredToken judge the token ttl
func (u *User) IsExpiredToken(ctx context.Context, in *pb.TokenRequest) (*pb.TokenResponse, error) {
	client := storage.NewRedis()
	key := TOKEN_PREFIX + in.Username
	value , isExist:= client.Get(key)
	if isExist == false || value.(string) != in.Token {
		return &pb.TokenResponse{
			IsExpired:true,
		}, nil
	}

	client.Put(key, in.Token, time.Hour)
	return &pb.TokenResponse{
		IsExpired:false,
	}, nil
}

// SaveToken save once token
func (u *User) SaveToken(ctx context.Context, in *pb.TokenRequest) (*pb.SaveResponse, error) {
	client := storage.NewRedis()
	key := TOKEN_PREFIX + in.Username
	_ , isExist:= client.Get(key)
	if isExist == false {
		return &pb.SaveResponse{
			Success:true,
		}, nil
	}
	client.Put(key, in.Token, time.Hour)
	return &pb.SaveResponse{
		Success:true,
	}, nil
}
