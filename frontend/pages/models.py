from django.db import models
import requests
import base64

# Create your models here.

gin_addr = "http://localhost:8085"
facesdk_access_token="***REMOVED***"
facesdk_search_url = "https://aip.baidubce.com/rest/2.0/face/v3/search"
facesdk_add_url = "https://aip.baidubce.com/rest/2.0/face/v3/faceset/user/add"

def checkLogin(request):
    context = {}
    if request.COOKIES.get("login_username"):
        context["user_is_authenticated"] = True
        context["user_name"] = request.COOKIES.get("login_username")
    return context

def MaxTxId():
    r = requests.get(gin_addr + "/listTx")
    tx_list = r.json()
    max_id = 0
    for tx_item in tx_list:
        id_i = int(tx_item["Key"].replace("TX", ""))
        if id_i > max_id:
            max_id = id_i
    return max_id

def AcquireAllTx(request):
    r = requests.get(gin_addr + "/listTx")
    context = {}
    context["TxList"] = r.json()
    context["txFilter"] = "everyone"
    return context

def AcquireScore(request, name):
    post_data = {}
    post_data["name"] = name
    r = requests.post(gin_addr + "/scoreUser", data=post_data)
    score = r.json()["Score"]
    context = {}
    context["creditScore"] = score
    context["normalizedCreditScore"] = int(score) / 10
    return context

def AcquireRelatedTx(request, name):
    post_data = {}
    post_data["name"] = name
    r = requests.post(gin_addr + "/queryUser", data=post_data)
    context = {}
    context["TxList"] = r.json()
    context["txFilter"] = post_data["name"]
    return context

def ValidateKey(name, key) -> bool:
    if (len(name) == 0):
        return False
    name_bytes = name.encode("utf-8")
    correct_key = base64.b64encode(name_bytes).decode("utf-8")
    print("correctkey:{} key:{}".format(correct_key, key))
    return correct_key == key

def AcquireKey(name):
    name_bytes = name.encode("utf-8")
    correct_key = base64.b64encode(name_bytes).decode("utf-8")
    return correct_key

def AppendFace(name, imgEnc):
    imgEnc = imgEnc.replace("data:image/png;base64,", "")
    name_key = AcquireKey(name)
    add_params = {}
    add_params["image"] = imgEnc
    add_params["image_type"] = "BASE64"
    add_params["group_id"] = "common"
    add_params["user_id"] = name
    add_params["quality_control"] = "NONE"
    add_params["liveness_control"] = "NONE"
    add_req = facesdk_add_url + "?access_token=" + facesdk_access_token
    headers = {'content-type': 'application/json'}
    response = requests.post(add_req, data=add_params, headers=headers)
    if response:
        result = response.json()
        print(result)
        if result["error_msg"] == "SUCCESS":
            return True
    return False

def SearchFace(imgEnc):
    imgEnc = imgEnc.replace("data:image/png;base64,", "")
    # name_key = AcquireKey(name)
    search_params = {}
    search_params["image"] = imgEnc
    search_params["image_type"] = "BASE64"
    search_params["group_id_list"] = "common"
    search_params["quality_control"] = "NONE"
    search_params["liveness_control"] = "NONE"
    search_req = facesdk_search_url + "?access_token=" + facesdk_access_token
    headers = {'content-type': 'application/json'}
    response = requests.post(search_req, data=search_params, headers=headers)
    max_score = 0
    max_user = ""
    if response:
        result = response.json()
        print(result)
        if not result["error_msg"] == "SUCCESS":
            return max_user
        result = result["result"]
        if "face_list" in result:
            for face_item in result["face_list"]:
                for user_item in face_item["user_list"]:
                    if int(user_item["score"]) > max_score:
                        max_score = int(user_item["score"])
                        max_user = user_item["user_id"]
        else:
            for user_item in result["user_list"]:
                if int(user_item["score"]) > max_score:
                    max_score = int(user_item["score"])
                    max_user = user_item["user_id"]
    print("Search {} score: {}".format(max_user, max_score))
    return max_user