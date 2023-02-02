package common

import "errors"

var (
	UserNameNull      = errors.New("用户名为空")
	UserNametoolong   = errors.New("用户名过长")
	PasswordNull      = errors.New("密码为空")
	Passwordtoolong   = errors.New("密码过长")
	UserNameNotUnique = errors.New("用户名重复")
	PasswordWrong     = errors.New("密码错误")
	UserNoexist       = errors.New("用户不存在")
)
