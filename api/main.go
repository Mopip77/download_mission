package api

import "onedrive/serializer"

func ErrorResponse(err error) serializer.Response {
	return serializer.Response{
		Status: 400,
		Data:   nil,
		Msg:    "错误",
		Error:  err.Error(),
	}
}
