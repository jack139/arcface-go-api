# Arcface-go API Service

使用Arcface的预训练模型实现的API服务示例：

| API             | 功能     |
| --------------- | -------- |
| /face2/locate   | 人脸定位 |
| /face2/verify   | 人脸对比 |
| /face2/features | 人脸特征 |



代码中使用了[arcface-go](https://github.com/jack139/arcface-go)和[go-infer](https://github.com/jack139/go-infer)，也可作为使用这两个项目的示例。




- 编译

```
make
```



- 启动推理服务

```
build/arcface-api server 0
```



- 启动 http

```
build/arcface-api http
```



- 测试脚本

```
python3 test_api.py 127.0.0.1 face_features data/1.jpg
```



- 说明
	- 运行时需要一个redis，配置在 ```config/settings.yaml```里
	- arcface的预训练模型使用的是 ["**buffalo_l**"](https://insightface.cn-sh2.ufileos.com/models/buffalo_l.zip)， 在[这里](https://github.com/deepinsight/insightface/tree/master/model_zoo)
	- 如果找不到opencv的运行库，可能需要设置```LD_LIBRARY_PATH```指向opencv的安装路径
	- API文档在[这里](doc/API.md)

