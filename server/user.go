package server

import (
	"my_mange_system/common"
	"my_mange_system/model"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type UserList struct {
	Username string `json:"username"`
	Roleid   int    `json:"roleid"`
	Userid   uint   `json:"userid"`
}

func CheckOutUser(ctx *gin.Context, username string, password string) (bool, model.User) {
	var user model.User
	DB := model.DB.Model(&model.User{})
	DB.Where("username = ?", username).First(&user)
	if user.Password == password {
		common.SetSession(ctx, "user", user)
		return true, user
	}
	return false, model.User{}
}

func UpdateLoginInfo(city string, username string) {
	DB := model.DB.Model(&model.User{})
	timestamp := time.Now().Unix()
	DB.Where("username = ?", username).Updates(model.User{City: city, LastLoginTime: timestamp})
}

func GetUsetList(username string, roleid int, offset int, limit int) ([]UserList, int64) {
	var users []model.User
	var new_users []UserList
	var total int64
	DB := model.DB.Model(&model.User{})
	if username != "" {
		DB = DB.Where("username LIKE ?", "%"+username+"%")
	}
	if roleid > 0 {
		DB = DB.Where("roleid = ?", roleid)
	}
	DB.Count(&total)
	DB.Limit(limit).Offset(offset).Find(&users)
	for _, user := range users {
		row := UserList{
			Username: user.Username,
			Roleid:   user.RoleId,
			Userid:   user.ID,
		}

		new_users = append(new_users, row)
	}
	return new_users, total
}

func DeleteUserList(userids []string, roleids []string, userinfo model.User) (bool, string) {
	var ids []uint
	for index, id := range userids {
		if roleids[index] == "1" {
			return false, "无法删除管理员用户"
		}
		_id, err := strconv.Atoi(id)
		if err != nil {
			return false, "用户数据解析失败"
		}
		if uint(_id) == userinfo.ID {
			return false, "删除用户种包含自己"
		}
		ids = append(ids, uint(_id))
	}
	DB := model.DB.Model(&model.User{})
	DB.Delete(&model.User{}, ids)
	return true, "删除成功"
}

func UpdateUserList(userid uint, username string, password string) (bool, string) {
	var user model.User
	DB := model.DB.Model(&model.User{})
	if username == "" {
		username = user.Username
	} else {
		DB.Where("username = ?", username).First(&user)
		if user.ID != userid {
			return false, "用户名已存在"
		}
	}
	if password == "" {
		password = user.Password
	}
	DB.Where("id = ?", userid).Updates(model.User{Username: username, Password: password})
	return true, "更新成功"
}

func ChangeUserPassword(userid uint, oldpassword string, newpassword string) (bool, string) {
	var user model.User
	DB := model.DB.Model(&model.User{})
	DB.Where("id = ?", userid).First(&user)
	if user.Password != oldpassword {
		return false, "密码不正确"
	}
	DB.Where("id = ?", userid).Updates(model.User{Password: newpassword})
	return true, "修改成功"
}
