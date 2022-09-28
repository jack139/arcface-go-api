# arcface-go api service



## 测试



### 编译

```
make
```



### 启动 dispatcher

```
LD_LIBRARY_PATH=/usr/local/lib64:/usr/local/lib build/arcface-api server 0
```



### 启动 http

```
LD_LIBRARY_PATH=/usr/local/lib64:/usr/local/lib build/arcface-api http
```



### 测试脚本

```
python3 test_api.py 127.0.0.1 face_features data/1.jpg
```
