package common

import "errors"

var (
	UserNameNull           = errors.New("用户名为空")
	UserNametoolong        = errors.New("用户名过长")
	PasswordNull           = errors.New("密码为空")
	Passwordtoolong        = errors.New("密码过长")
	Passwordtooshot        = errors.New("密码过短")
	UserNameNotUnique      = errors.New("用户名重复")
	PasswordWrong          = errors.New("密码错误")
	UserNoexist            = errors.New("该用户不存在")
	PasswordEncryptWrong   = errors.New("密码加密错误")
	UserLiveWrong          = errors.New("请求的用户信息不存在")
	VideoGetWrong          = errors.New("视频列表获取失败")
	ActionTypeWrong        = errors.New("用户动作违规")
	FollowActionDuplicate  = errors.New("用户操作重复")
	FollowActionWrong      = errors.New("不能自己关注自己")
	FavoriteExist          = errors.New("点赞操作重复")
	FavoriteNoExist        = errors.New("点赞关系不存在")
	SendMessageActionWrong = errors.New("自己给自己发消息")
)
