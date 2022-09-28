package arcface

import (
	"log"
	"fmt"
	"bytes"
	"image"
	"io/ioutil"
	"strconv"

	"github.com/disintegration/imaging"
	"github.com/jack139/go-infer/helper"
	"github.com/jack139/arcface-go/arcface"
)

const vecLen = 512 // 特征向量长度

/* 训练好的模型权重 */
var (
	initOK = bool(false)
	threshHold float32
)


/* 初始化模型 */
func initModel() error {
	var err error

	if !initOK { // 模型只装入一次
		if err = arcface.LoadOnnxModel(helper.Settings.Customer["ArcfaceModelPath"]); err!=nil {
			return err
		}
		log.Println("Arcface onnx model loaded from: ", helper.Settings.Customer["ArcfaceModelPath"])


		// 初始化参数
		value, err := strconv.ParseFloat(helper.Settings.Customer["FACE_DistanceThreshold"], 32)
		if err != nil {
			return err
		}
		threshHold = float32(value)

		// 初始化标记
		initOK = true

		// 模型热身
		warmup(helper.Settings.Customer["FACE_WARM_UP_IMAGES"])
	}

	return nil
}


func locateInfer(imageByte []byte) ([][]int, int, error){

	// 转换为 image.Image
	reader := bytes.NewReader(imageByte)

	img, err := imaging.Decode(reader)
	if err!=nil {
		return nil, 9201,err
	}

	// 检测人脸
	dets, _, err := arcface.FaceDetect(img)
	if err != nil {
		return nil, 9202, err
	}

	// 返回整数结果
	r2 := make([][]int, len(dets))
	for i:=0;i<len(dets);i++ {
		r2[i] = make([]int, 4)
		for j:=0;j<4;j++ {
			r2[i][j] = int(dets[i][j])
		}
	}

	return r2, 0, nil
}


func featuresInfer(imageByte []byte) ([]float32, []int, image.Image, int, error){

	// 转换为 image.Image
	reader := bytes.NewReader(imageByte)

	img, err := imaging.Decode(reader)
	if err!=nil {
		return nil, nil, nil, 9201, err
	}

	// 检测人脸
	dets, kpss, err := arcface.FaceDetect(img)
	if err != nil {
		return nil, nil, nil, 9202, err
	}

	if len(dets)==0 {
		log.Println("No face detected.")
		return nil, nil, nil, 0, nil
	}

	// 只返回第一个人脸的特征
	features, normFace, err := arcface.FaceFeatures(img, kpss[0])
	if err != nil {
		return nil, nil, nil, 9203, err
	}

	// 返回整数结果
	r2 := make([]int, 4)
	for j:=0;j<4;j++ {
		r2[j] = int(dets[0][j])
	}

	return features, r2, normFace, 0, nil
}

// 模型热身
func warmup(path string){
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Printf("warmup fail: %s", err.Error())
		return
	}

	for _, file := range files {
		if file.IsDir() { continue }
	
		img, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", path, file.Name()))
		if err != nil { continue }

		r, r2, _, _, err := featuresInfer(img)
		if err==nil {
			log.Printf("warmup: %s %v %v", file.Name(), len(r), len(r2))
		} else {
			log.Printf("warmup fail: %s %s", file.Name(), err.Error())
		}
	}
}
