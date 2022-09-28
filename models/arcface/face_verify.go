package arcface

import (
	"fmt"
	"log"
	"strconv"
	"encoding/base64"

	"github.com/jack139/go-infer/helper"
)

/*  定义模型相关参数和方法  */
type FaceVerify struct{}

func (x *FaceVerify) Init() error {
	return initModel()
}

func (x *FaceVerify) ApiPath() string {
	return "/face2/verify"
}

func (x *FaceVerify) ApiEntry(reqData *map[string]interface{}) (*map[string]interface{}, error) {
	log.Println("Api_FaceVerify")

	// 检查参数
	imageBase64_1, ok := (*reqData)["image1"].(string)
	if !ok {
		return &map[string]interface{}{"code":9001}, fmt.Errorf("need image1")
	}

	imageBase64_2, ok := (*reqData)["image2"].(string)
	if !ok {
		return &map[string]interface{}{"code":9001}, fmt.Errorf("need image2")
	}

	// 构建请求参数
	reqDataMap := map[string]interface{}{
		"image1": imageBase64_1,
		"image2": imageBase64_2,
	}

	return &reqDataMap, nil
}


// 推理
func (x *FaceVerify) Infer(requestId string, reqData *map[string]interface{}) (*map[string]interface{}, error) {
	log.Println("Infer_FaceVerify")

	imageBase64_1 := (*reqData)["image1"].(string)
	imageBase64_2 := (*reqData)["image2"].(string)

	// 解码base64
	image1, err  := base64.StdEncoding.DecodeString(imageBase64_1)
	if err!=nil {
		return &map[string]interface{}{"code":9901}, err
	}

	image2, err  := base64.StdEncoding.DecodeString(imageBase64_2)
	if err!=nil {
		return &map[string]interface{}{"code":9901}, err
	}

	// 检查图片大小
	maxSize, _ := strconv.Atoi(helper.Settings.Customer["FACE_MAX_IMAGE_SIZE"])
	if len(image1) > maxSize {
		return &map[string]interface{}{"code":9002}, fmt.Errorf("图片数据太大")
	}

	if len(image2) > maxSize {
		return &map[string]interface{}{"code":9002}, fmt.Errorf("图片数据太大")
	}

	// 模型推理
	r1, _, _, code, err := featuresInfer(image1)
	if err != nil {
		return &map[string]interface{}{"code":code}, err
	}

	if r1==nil {  // 未检测到人脸
		return &map[string]interface{}{"is_match":false, "score":0.0}, nil
	}

	r2, _, _, code, err := featuresInfer(image2)
	if err != nil {
		return &map[string]interface{}{"code":code}, err
	}

	if r2==nil {  // 未检测到人脸
		return &map[string]interface{}{"is_match":false, "score":0.0}, nil
	}

	score, err := cosine(r1, r2)
	if err != nil {
		return &map[string]interface{}{"code":9003}, err
	}

	// 保存请求图片和结果
	//saveBackLog(requestId, image1, []byte(fmt.Sprintf("%v", score)))

	return &map[string]interface{}{
		"is_match" : float32(score)<threshHold, 
		"score"    : score / 2 + 0.5,
	}, nil
}
