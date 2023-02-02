package service

import (
	"github.com/jinzhu/gorm"
	"simpledouyin/src/common"
	"simpledouyin/src/dao"
	"simpledouyin/src/entity"
)

// 新建user
func UserRegister(username string, password string) (uid uint, err error) {
	if IsUserNameUnique(username) {
		return 0, common.UserNameNotUnique
	}
	newUser := entity.User{
		Name:     username,
		Password: password,
	}
	dao.SqlSession.Model(&entity.User{}).Create(&newUser)
	return newUser.ID, nil

}
func IsUserNameUnique(username string) (ok bool) {
	var userExist = &entity.User{}
	if err := dao.SqlSession.Model(&entity.User{}).Where("name=?", username).First(&userExist).Error; gorm.IsRecordNotFoundError(err) {
		return false
	}
	return true
}

// 查看用户名是否合法,查看密码是否合法
func IsUserPasswordLegal(userName string, password string) (err error) {
	if userName == "" {
		return common.UserNameNull
	}
	if len(userName) > common.MaxUsernameLength {
		return common.UserNametoolong
	}
	if password == "" {
		return common.PasswordNull
	}
	if len(password) > common.MaxPasswordLength {
		return common.Passwordtoolong
	}
	return nil
}

func IsPasswordRight(username string, password string) (uid uint, err error) {
	var user = &entity.User{}
	dao.SqlSession.Model(&entity.User{}).Where("name=?", username).First(&user)
	if user.Password == password {
		return user.ID, nil
	} else {
		return 0, common.PasswordWrong
	}

}
