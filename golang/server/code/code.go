package code

import (
	"errors"
	"fmt"
	"reflect"
)

var (
	ResDupErr      = errors.New(Messages[CodeResourceDuplicated])
	ResNotFoundErr = errors.New(Messages[CodeResourceNotFound])
)

type CodeCount int

const (
	CodeOk CodeCount = 0 + iota
	CodeSystemError
	CodeParamWrong
	CodeResourceDuplicated
	CodeResourceNotFound
	CodeResourceOccupied
	CodeIllegeField

	CodeWrongPassword = 1000 + iota
	CodeWrongTicket
	CodeNotLogin
	CodeWrongExpired
	CodeIllegeOP
	CodeParamInvalid
	CodeTokenInvalid
)

var Messages = map[CodeCount]string{
	CodeOk:                 "成功",
	CodeSystemError:        "系统错误",
	CodeParamWrong:         "参数错误",
	CodeResourceDuplicated: "资源重复",
	CodeResourceNotFound:   "资源未找到",
	CodeResourceOccupied:   "资源被占用",
	CodeIllegeField:        "格式不对",

	CodeWrongPassword: "密码错误",
	CodeNotLogin:      "未登陆",
	CodeWrongTicket:   "账号在其他设备登录",
	CodeWrongExpired:  "更新会话失败",
	CodeIllegeOP:      "非法操作",
	CodeParamInvalid:  "参数非法",
	CodeTokenInvalid:  "token非法",
}

type CodeInfo struct {
	Code    CodeCount `json:"code"`
	Message string    `json:"message"`
}

func GetCodeInfo(ci *CodeInfo) string {
	return fmt.Sprintf(`{"code":%v, "message": "%v"}`, ci.Code, ci.Message)
}

// NewCode construct CodeCount and string
// if msg is nil, return a default message
// return a type of CodeInfo
func NewCode(code CodeCount, msg string) *CodeInfo {
	if msg == "" {
		msg = Messages[code]
	}
	return &CodeInfo{
		Code:    code,
		Message: msg,
	}
}

// Response used for controllers` APIs response
func Response(v interface{}, ci *CodeInfo) interface{} {
	s := reflect.ValueOf(v).Elem()
	t := s.FieldByName("CodeInfo")
	t.FieldByName("Code").SetInt(int64(ci.Code))
	t.FieldByName("Message").SetString(Messages[ci.Code])
	if ci.Message != "" {
		t.FieldByName("Message").SetString(ci.Message)
	}
	return v
}

func GetCodeMessage(code CodeCount) string {
	return Messages[code]
}

func FillCodeInfo(v interface{}, ci *CodeInfo) interface{} {
	ele := reflect.ValueOf(v).Elem()
	field := ele.FieldByName("CodeInfo")

	if ci.Message == "" {
		ci.Message = GetCodeMessage(ci.Code)
	}

	// set field
	field.FieldByName("Code").SetInt(int64(ci.Code))
	field.FieldByName("Message").SetString(ci.Message)

	return v
}

func GetTheCodeInfo(code CodeCount) *CodeInfo {
	return &CodeInfo{
		Code:    code,
		Message: GetCodeMessage(code),
	}
}
