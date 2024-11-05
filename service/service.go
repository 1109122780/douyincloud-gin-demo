/*
Copyright (year) Bytedance Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package service

import (
	"context"
	"douyincloud-gin-demo/component"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/volcengine/ve-tos-golang-sdk/v2/tos"
	"net/http"
	"strings"
)

func UploadHandler(c *gin.Context) {
	// 直接在代码中定义变量
	accessKey := "AKLTMmY0MzAxNGM3NDczNGIyY2I0NTBkMjYzMGU4NWY1ODg"
	secretKey := "T1dZek1UTmpNVGd4WldWaU5EUTVOV0U0TWpGa1ltSXdZV0k0TmpRM1ptWQ=="
	endpoint := "https://tos-cn-beijing.volces.com"
	region := "cn-beijing"
	bucketName := "tt42cc3471ddf817dd12-env-g3hekwkvus" // 填写实际的 bucket 名称
	objectKey := "test/example_object.txt"
	content := "1234567890abcdefghijklmnopqrstuvwxyz~!@#$%^&*()_+<>?,./   :'1234567890abcdefghijklmnopqrstuvwxyz~!@#$%^&*()_+<>?,./   :'"

	// 初始化 TOS 客户端
	client, err := tos.NewClientV2(endpoint, tos.WithRegion(region), tos.WithCredentials(tos.NewStaticCredentials(accessKey, secretKey)))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize TOS client", "details": err.Error()})
		return
	}

	ctx := context.Background()
	body := strings.NewReader(content)

	// 上传内容
	output, err := client.PutObjectV2(ctx, &tos.PutObjectV2Input{
		PutObjectBasicInput: tos.PutObjectBasicInput{
			Bucket: bucketName,
			Key:    objectKey,
		},
		Content: body,
	})

	// 处理错误
	if err != nil {
		if serverErr, ok := err.(*tos.TosServerError); ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":               "Server error",
				"request_id":          serverErr.RequestID,
				"status_code":         serverErr.StatusCode,
				"response_error_code": serverErr.Code,
				"response_error_msg":  serverErr.Message,
			})
		} else if clientErr, ok := err.(*tos.TosClientError); ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Client error",
				"details": clientErr.Cause.Error(),
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unknown error", "details": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"requestID": output.RequestID,
		"message":   "File uploaded successfully",
	})
}

func Hello(ctx *gin.Context) {
	target := ctx.Query("target")
	if target == "" {
		Failure(ctx, fmt.Errorf("param invalid"))
		return
	}
	fmt.Printf("target= %s\n", target)
	hello, err := component.GetComponent(target)
	if err != nil {
		Failure(ctx, fmt.Errorf("param invalid"))
		return
	}

	name, err := hello.GetName(ctx, "name")
	if err != nil {
		Failure(ctx, err)
		return
	}
	Success(ctx, name)
}

func SetName(ctx *gin.Context) {
	var req SetNameReq
	err := ctx.Bind(&req)
	if err != nil {
		Failure(ctx, err)
		return
	}
	hello, err := component.GetComponent(req.Target)
	if err != nil {
		Failure(ctx, fmt.Errorf("param invalid"))
		return
	}
	err = hello.SetName(ctx, "name", req.Name)
	if err != nil {
		Failure(ctx, err)
		return
	}
	Success(ctx, "")
}

func Failure(ctx *gin.Context, err error) {
	resp := &Resp{
		ErrNo:  -1,
		ErrMsg: err.Error(),
	}
	ctx.JSON(200, resp)
}

func Success(ctx *gin.Context, data string) {
	resp := &Resp{
		ErrNo:  0,
		ErrMsg: "success",
		Data:   data,
	}
	ctx.JSON(200, resp)
}

type HelloResp struct {
	ErrNo  int    `json:"err_no"`
	ErrMsg string `json:"err_msg"`
	Data   string `json:"data"`
}

type SetNameReq struct {
	Target string `json:"target"`
	Name   string `json:"name"`
}

type Resp struct {
	ErrNo  int         `json:"err_no"`
	ErrMsg string      `json:"err_msg"`
	Data   interface{} `json:"data"`
}
