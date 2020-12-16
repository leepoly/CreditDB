from django.db import models
import requests
import base64

# Create your models here.

gin_addr = "http://localhost:8085"

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
    name_bytes = name.encode("utf-8")
    correct_key = base64.b64encode(name_bytes).decode("utf-8")
    print("correctkey:{} key:{}".format(correct_key, key))
    return correct_key == key

def AcquireKey(name):
    name_bytes = name.encode("utf-8")
    correct_key = base64.b64encode(name_bytes).decode("utf-8")
    return correct_key
