
## API 文档



### 1. API清单

| url                    | 功能             |
| ---------------------- | ---------------- |
| /face2/locate           | 人脸定位         |
| /face2/verify           | 人脸对比         |
| /face2/features         | 人脸特征         |



### 2. 全局接口定义

输入参数

| 参数      | 类型   | 说明                          | 示例        |
| --------- | ------ | ----------------------------- | ----------- |
| appId     | string | 应用渠道编号                  |             |
| version   | string | 版本号                        |             |
| signType  | string | 签名算法，目前使用国密SM2算法 | SM2或SHA256 |
| signData  | string | 签名数据，具体算法见下文      |             |
| encType   | string | 接口数据加密算法，目前不加密  | plain       |
| timestamp | int    | unix时间戳（秒）              |             |
| data      | json   | 接口数据，详见各接口定义      |             |

> 签名/验签算法：
>
> 1. 筛选，获取参数键值对，剔除signData、encData、extra三个参数。data参数按key升序排列进行json序列化。
> 2. 排序，按key升序排序。
> 3. 拼接，按排序好的顺序拼接请求参数
>
> ```key1=value1&key2=value2&...&key=appSecret```，key=appSecret固定拼接在参数串末尾，appSecret需替换成应用渠道所分配的appSecret。
>
> 4. 签名，使用制定的算法进行加签获取二进制字节，使用 16进制进行编码Hex.encode得到签名串，然后base64编码。
> 5. 验签，对收到的参数按1-4步骤签名，比对得到的签名串与提交的签名串是否一致。

签名示例：

```json
请求参数：
{
    "appId":"19E179E5DC29C05E65B90CDE57A1C7E5",
    "version": "1",
    "signType": "SHA256",
    "signData": "...",
    "encType": "plain",
    "timestamp":1591943910,
    "data": {
    	"user_id":"gt",
    	"face_id":"5ed21b1c262daabe314048f5"
    }
}

密钥：
appSecret="D91CEB11EE62219CD91CEB11EE62219C"
SM2_privateKey="JShsBOJL0RgPAoPttEB1hgtPAvCikOl0V1oTOYL7k5U="

待加签串：
appId=19E179E5DC29C05E65B90CDE57A1C7E5&data={"face_id":"5ed21b1c262daabe314048f5","user_id":"gt"}&encType=plain&signType=SM2&timestamp=1591943910&version=1&key=D91CEB11EE62219CD91CEB11EE62219C

SHA256加签结果：
"2072bd8afb678c03ce9be14202e47b12031aa42a0a8c8593723d7027007ef804"

base64后结果：
"MjA3MmJkOGFmYjY3OGMwM2NlOWJlMTQyMDJlNDdiMTIwMzFhYTQyYTBhOGM4NTkzNzIzZDcwMjcwMDdlZjgwNA=="

SM2加签结果（每次不同）：
"LXgGBQNsXwofSXr+uXYiw0al7MFNNdUl0OyjpxHGKSPjJAr1N5oO6Tq3WL0C8UVX1pmDNH/GZK1Q0h+VvzKiEg=="

```



返回结果

| 参数      | 类型   | 说明                                                         |
| --------- | ------ | ------------------------------------------------------------ |
| code      | string | 接口返回状态代码                                             |
| timestamp | int    | unix时间戳                                                   |
| data      | json   | 成功时返回结果数据。出错时，data.msg返回错误说明。（人脸识别接口有此data.msg字段） |
| msg       | string | 出错时，msg返回错误说明。（特征库接口有此字段。）            |

> 成功时：code为0， msg为"success"，data内容见各接口定义；
>
> 出错时：code返回错误代码，具体定义建各接口说明

返回示例

```json
{
    "code": 0, 
    "timestamp": 1591943910,
    "data": {
       "msg": "success", 
       "...."
    }
}
```

全局出错代码

| 编码 | 说明                               |
| ---- | ---------------------------------- |
| 9800 | 无效签名                           |
| 9801 | 签名参数有错误                     |
| 9802 | 调用时间错误，unixtime超出接受范围 |



### 3. Arcface api

#### （1）人脸定位

> 检测图片中人脸并返回位置

请求URL

> http://127.0.0.1:5000/face2/locate

请求方式

> POST

输入参数

| 参数         | 必选 | 类型   | 说明                                              |
| ------------ | ---- | ------ | ------------------------------------------------- |
| image        | 是   | string | base64编码图片数据                                |
| max_face_num | 否   | int    | 最多定位的人脸数量，默认为1，仅检测面积最大的一个 |

请求示例

```json
{
    "image" : "....", 
    "max_face_num" : 5
}
```

返回结果

| 参数                           | 必选 | 类型  | 说明                 |
| ------------------------------ | ---- | ----- | -------------------- |
| face_num                       | 是   | int   | 检测到的图片人脸数量 |
| locations                      | 是   | array | 人脸位置坐标列表     |
| + [ top, right, bottom, left ] | 是   | array | 人脸位置             |

返回示例

```json
{
    "appId": "3EA25569454745D01219080B779F021F", 
    "code": 0, 
    "data": {
        "face_num": 2, 
        "locations": [
            [1145, 364, 1335, 607], 
            [764, 391, 947, 641]
        ], 
        "msg": "success", 
        "requestId": "202209281709289cfaac99b62c26f65346fe767c976c1a"
    }, 
    "encType": "plain", 
    "signType": "plain", 
    "success": true, 
    "timestamp": 1662013631
}
```

出错代码

| 编码 | 说明                              |
| ---- | --------------------------------- |
| 9001 | 缺少参数                          |
| 9002 | 图片数据太大，base64数据不大于2MB |
| 9901 | base64编码异常                    |



#### （2）人脸对比

> 比对两张照片中人脸的相似度（1:1），返回相似度分值

请求URL

> http://127.0.0.1:5000/face2/verify

请求方式

> POST

输入参数

| 参数   | 必选 | 类型   | 说明               |
| ------ | ---- | ------ | ------------------ |
| image1 | 是   | string | base64编码图片数据 |
| image2 | 是   | string | base64编码图片数据 |

请求示例

```json
{
    "image1": "....", 
    "image2": "....", 
}
```

返回结果

| 参数     | 必选 | 类型    | 说明                           |
| -------- | ---- | ------- | ------------------------------ |
| is_match | 是   | boolean | 是否同一人，TRUE 或 FALSE      |
| score    | 是   | float   | 相似度得分（值越小相似度越高） |

返回示例

```json
{
    "appId": "3EA25569454745D01219080B779F021F", 
    "code": 0, 
    "data": {
        "is_match": true, 
        "msg": "success", 
        "requestId": "202209281709289cfaac99b62c26f65346fe767c976c1a", 
        "score": 0.06643416019927911
    }, 
    "encType": "plain", 
    "signType": "plain", 
    "success": true, 
    "timestamp": 1662014458
}
```

出错代码

| 编码 | 说明                              |
| ---- | --------------------------------- |
| 9001 | 缺少参数                          |
| 9002 | 图片数据太大，base64数据不大于2MB |
| 9901 | base64编码异常                    |



#### （3）人脸特征

> 返回arcface模型计算的人脸特征值，长度512

请求URL

> http://127.0.0.1:5000/face2/features

请求方式

> POST

输入参数

| 参数         | 必选 | 类型   | 说明                                              |
| ------------ | ---- | ------ | ------------------------------------------------- |
| image        | 是   | string | base64编码图片数据                                |

请求示例

```json
{
    "image" : "....", 
}
```

返回结果

| 参数                           | 必选 | 类型  | 说明                 |
| ------------------------------ | ---- | ----- | -------------------- |
| features                      | 是   | array | 人脸特征值     |


返回示例

```json
{
    "appId": "3EA25569454745D01219080B779F021F", 
    "code": 0, 
    "data": {
        "features": [ 0.053627726, -0.040755205, 0.029853795, 0.038380787, -0.03890704, "...."], 
        "msg": "success", 
        "requestId": "202209281709289cfaac99b62c26f65346fe767c976c1a"
    }, 
    "encType": "plain", 
    "signType": "plain", 
    "success": true, 
    "timestamp": 1664356169
}
```

出错代码

| 编码 | 说明                              |
| ---- | --------------------------------- |
| 9001 | 缺少参数                          |
| 9002 | 图片数据太大，base64数据不大于2MB |
| 9901 | base64编码异常                    |
