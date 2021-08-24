import os
import wmi
import sys
from urllib import request, parse
import json
from http import cookiejar
import socket
import re
import psutil

is_active = None
base_url = "http://10.1.32.188"
user_add_api = "{base}/api/user/add".format(base=base_url)
user_json_url = '{base}/api/user/json'.format(base=base_url)
user_login_api = "{base}/login".format(base=base_url)

c = wmi.WMI()    
my_system = c.Win32_ComputerSystem()[0]

ip_address = socket.gethostbyname(socket.gethostname())

print("host ip: " + ip_address)

used_user = os.getenv('username')

user_data = {
    "username": os.getenv('username'),
    "email": os.getenv("username") + '@' + "dominos.com.cn",
    "is_superuser": 0,
    "is_staff": 1,
    "is_active": 1,
    "is_delete": 0,
}

cookie = cookiejar.CookieJar()
cookie_handler = request.HTTPCookieProcessor(cookie)

form_data = {
    "username": "admin",
    "password": "admin@123"
}

data = parse.urlencode(form_data).encode("utf-8")
opener = request.build_opener(cookie_handler)
req = request.Request(url=user_login_api, data=data)
resp = opener.open(req)


def do_post_api(url, req_dict):
    return opener.open(request.Request(url=url, data=parse.urlencode(req_dict).encode("utf-8")))


def do_get_json(url):
    return json.loads(opener.open(url).read().decode('utf-8'))


user_response_list = do_get_json(url=user_json_url)


def detect_user_is_exist(query_list, username):
    ret_msg = dict()
    for element in query_list:
        if re.search(username, element['UserName'], re.IGNORECASE):
            ret_msg["status"] = True
            ret_msg["id"] = element["Id"]
        else:
            continue

    return ret_msg


user_ret_dict = detect_user_is_exist(query_list=user_response_list, username=used_user)
if user_ret_dict:
    print("%s is exist, id: %d" % (used_user, user_ret_dict["id"]))
else:
    do_post_api(url=user_add_api, req_dict=user_data)
    user_ret_dict = detect_user_is_exist(query_list=do_get_json(url=user_json_url), username=used_user)


# 设备类型
deviceType_api_url = "{base}/api/assets/device_type/add".format(base=base_url)
deviceType_json_url = '{base}/api/assets/device_type/json'.format(base=base_url)

deviceType_data = {
    "deviceType": "Server",
    "desc": "设备类型",
    "is_active": 1,
    "is_delete": 0
}


device_type_response = do_get_json(url=deviceType_json_url)


def detect_device_type(res_list):
    ret_msg = dict()
    for device in res_list:
        if device["DeviceType"] == deviceType_data["deviceType"]:
            ret_msg["status"] = True
            ret_msg["id"] = device["Id"]
        else:
            continue
    return ret_msg


device_type_dict = detect_device_type(res_list=device_type_response)

if device_type_dict:
    print("device type already exist, id: %d" % device_type_dict["id"])
else:
    resp = do_post_api(url=deviceType_api_url, req_dict=deviceType_data)
    print(resp.read().decode('utf-8'))
    # device_type_dict = detect_device_type(res_list=do_get_json(url=deviceType_json_url))


# 品牌
brand_api_url = "{base}/api/assets/brand/add".format(base=base_url)
brand_json_url = '{base}/api/assets/brand/json'.format(base=base_url)

brand_data = {
    "brand_name": my_system.Manufacturer.split(" ")[0],
    "is_delete": 0,
    "is_active": 1,
    "desc": "品牌名称",
}

brand_resp = do_get_json(url=brand_json_url)


def detect_brands_exists(res_list):
    ret_msg = dict()
    for b in res_list:
        if b["Name"] == brand_data["brand_name"]:
            ret_msg["status"] = True
            ret_msg["id"] = b["Id"]
        else:
            continue
    return ret_msg


# 检测品牌是否存在
brands_dict = detect_brands_exists(res_list=brand_resp)

if brands_dict:
    print("brand id: %d" % brands_dict["id"])
else:
    resp = do_post_api(url=brand_api_url, req_dict=brand_data)
    print(resp.read().decode('utf-8'))
    # brands_dict = detect_brands_exists(res_list=do_get_json(url=brand_json_url))

# 主机

pc_api_url = "{base}/api/assets/pc/add".format(base=base_url)
pc_json_url = '{base}/api/assets/pc/json'.format(base=base_url)

if my_system.Status == "OK":
    is_active = 1


cpu = psutil.cpu_count()
memory = psutil.virtual_memory().total
partitions = psutil.disk_partitions()

p_list = list()
for p in partitions:
    if p.fstype == '':
        continue
    else:
        p_tag = str(p.device).strip(':\\')
        size = psutil.disk_usage('%s' % p.device).total
        p_list.append(p_tag + "盘: " + str(size))

# print(p_list)
disk = ','.join(p_list)
# print(disk)

pc_data = {
    "location": "shanghai",
    "asset_no": my_system.Name,
    "model": my_system.Model,
    "name": my_system.Name,
    "is_active": is_active,
    "is_delete": 0,
    "brand": brands_dict["id"],
    "deviceType": device_type_dict["id"],
    "users": user_ret_dict["id"],
    "remark": os.getenv('username') + "的办公电脑",
    "ip_address": ip_address,
    "cpu": cpu,
    "memory": memory,
    "disk": disk
}

print(pc_data)

pc_resp = do_get_json(url=pc_json_url)

# print(pc_resp)


def detect_pc_exists(res_list):
    ret_msg = dict()
    for pc in res_list:
        if pc["Name"] == pc_data["name"]:
            ret_msg["status"] = True
            ret_msg["id"] = pc["Id"]
        else:
            continue
    return ret_msg


pc_dict = detect_pc_exists(res_list=pc_resp)
# 如果主机名存在，输出已存在，否则添加
if pc_dict:
    print("pc already exist pc_id: %d" % pc_dict["id"])
else:
    res = do_post_api(url=pc_api_url, req_dict=pc_data)
    print(res.read().decode('utf-8'))
    # pc_dict = detect_pc_exists(res_list=do_get_json(url=pc_json_url))

sys.exit()
