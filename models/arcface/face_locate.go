package arcface

import (
	"fmt"
	"log"
	"strconv"
	"encoding/base64"

	"github.com/jack139/go-infer/helper"
)

/*  定义模型相关参数和方法  */
type FaceLocate struct{}

func (x *FaceLocate) Init() error {
	return initModel()
}

func (x *FaceLocate) ApiPath() string {
	return "/face2/locate"
}

func (x *FaceLocate) ApiEntry(reqData *map[string]interface{}) (*map[string]interface{}, error) {
	log.Println("Api_FaceLocate")

	// 检查参数
	imageBase64, ok := (*reqData)["image"].(string)
	if !ok {
		return &map[string]interface{}{"code":9001}, fmt.Errorf("need image")
	}

	var maxFace float64
	maxFace, ok = (*reqData)["max_face_num"].(float64)
	if !ok {
		maxFace = 1
	}

	// 构建请求参数
	reqDataMap := map[string]interface{}{
		"image": imageBase64,
		"max_face_num" : maxFace,
	}

	return &reqDataMap, nil
}


// 推理
func (x *FaceLocate) Infer(requestId string, reqData *map[string]interface{}) (*map[string]interface{}, error) {
	log.Println("Infer_FaceLocate")

	imageBase64 := (*reqData)["image"].(string)
	maxFace := (*reqData)["max_face_num"].(float64)

	log.Println("maxFace: ", maxFace)

	// 解码base64
	image, err  := base64.StdEncoding.DecodeString(imageBase64)
	if err!=nil {
		return &map[string]interface{}{"code":9901}, err
	}

	// 检查图片大小
	maxSize, _ := strconv.Atoi(helper.Settings.Customer["FACE_MAX_IMAGE_SIZE"])
	if len(image) > maxSize {
		return &map[string]interface{}{"code":9002}, fmt.Errorf("图片数据太大")
	}

	// 模型推理
	r, code, err := locateInfer(image)
	if err != nil {
		return &map[string]interface{}{"code":code}, err
	}

	log.Println("face num--> ", len(r), int(maxFace))

	// 最多返回 maxFace 个数据
	if len(r) > int(maxFace) {
		r = r[:int(maxFace)]
	}

	return &map[string]interface{}{"locations":r, "face_num":len(r)}, nil
}
