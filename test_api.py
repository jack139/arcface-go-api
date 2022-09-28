# coding:utf-8

import sys, urllib3, json, base64, time, hashlib
from datetime import datetime

urllib3.disable_warnings()

# 生成参数字符串
def gen_param_str(param1):
    param = param1.copy()
    name_list = sorted(param.keys())
    if 'data' in name_list: # data 按 key 排序, 中文不进行性转义，与go保持一致
        param['data'] = json.dumps(param['data'], sort_keys=True, ensure_ascii=False, separators=(',', ':'))
    return '&'.join(['%s=%s'%(str(i), str(param[i])) for i in name_list if str(param[i])!=''])


def request(hostname, body, url):
    appid = '3EA25569454745D01219080B779F021F'
    unixtime = int(time.time())
    body['timestamp'] = unixtime
    body['appId'] = appid

    param_str = gen_param_str(body)
    sign_str = '%s&key=%s' % (param_str, '41DF0E6AE27B5282C07EF5124642A352')

    #print(sign_str)

    signature_str =  base64.b64encode(hashlib.sha256(sign_str.encode('utf-8')).hexdigest().encode('utf-8')).decode('utf-8')

    body['signData'] = signature_str

    body = json.dumps(body)
    #print(body)

    pool = urllib3.PoolManager(num_pools=2, timeout=180, retries=False)

    start_time = datetime.now()
    r = pool.urlopen('POST', url, body=body)
    print('[Time taken: {!s}]'.format(datetime.now() - start_time))

    return r



if __name__ == '__main__':
    if len(sys.argv)<4:
        print("usage: python3 %s <host> <api> <image_path>" % sys.argv[0])
        sys.exit(2)

    hostname = sys.argv[1]
    cate     = sys.argv[2]
    filepath = sys.argv[3]

    host = 'http://%s:5000'%hostname

    body = {
        'version'  : '1',
        'signType' : 'SHA256', 
        'encType'  : 'plain',
        'data'     : {},
    }

    if cate=="face_locate":
        url = host+'/face2/locate'
        with open(filepath, 'rb') as f:
            img_data = f.read()
        body['data']['image'] = base64.b64encode(img_data).decode('utf-8')
        body['data']['max_face_num'] = 2
    elif cate=="face_verify":
        url = host+'/face2/verify'
        filepath = filepath.split('#')
        print(filepath)
        with open(filepath[0], 'rb') as f:
            img_data1 = f.read()
        with open(filepath[1], 'rb') as f:
            img_data2 = f.read()
        body['data']['image1'] = base64.b64encode(img_data1).decode('utf-8')
        body['data']['image2'] = base64.b64encode(img_data2).decode('utf-8')
    elif cate=="face_features":
        url = host+'/face2/features'
        with open(filepath, 'rb') as f:
            img_data = f.read()
        body['data']['image'] = base64.b64encode(img_data).decode('utf-8')
    else:
        print("unknown <api>")
        sys.exit(2)

    print("-->", url)

    r = request(hostname, body, url)

    print(r.status)
    if r.status==200:
        print(json.loads(r.data.decode('utf-8')))
    else:
        print(r.data)
