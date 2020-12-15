from django.views.generic import TemplateView
from django.shortcuts import render
from django.contrib import messages
from django.http import HttpResponse, HttpResponseRedirect
import requests
from django.urls import reverse


class HomePageView(TemplateView):
    template_name = "pages/home.html"

gin_addr = "http://localhost:8085"

def AboutPageView(request):
    r = requests.get(gin_addr + "/hello")
    context = r.json()
    return render(request, "pages/about.html", context)

def ListTransaction(request):
    r = requests.get(gin_addr + "/listTx")
    context = {}
    context["TxList"] = r.json()
    # print(context)
    return render(request, "pages/txlist.html", context)

def MaxTxId():
    r = requests.get(gin_addr + "/listTx")
    tx_list = r.json()
    max_id = 0
    for tx_item in tx_list:
        id_i = int(tx_item["Key"].replace("TX", ""))
        if id_i > max_id:
            max_id = id_i
    return max_id

def CreateTxView(request):
    return render(request, "pages/createtx.html", {})

def ListRelatedTxView(request):
    return render(request, "pages/login.html", {})

def ListRelatedTx(request):
    post_data = {}
    if request.POST:
        post_data["name"] = request.POST["name"]
        context = {}
        r = requests.post(gin_addr + "/queryUser", data=post_data)
        context["TxList"] = r.json()
        return render(request, "pages/txlist.html", context)
    return HttpResponseRedirect(reverse("home"))

def CreateTx(request):
    post_data = {}
    txid = MaxTxId()
    if request.POST:
        post_data["SenderName"] = request.POST["identity"]
        post_data["RecverName"] = request.POST["recvername"]
        post_data["value"] = request.POST["value"]
        post_data["id"] = txid + 1
        response = requests.post(gin_addr + "/createTx", data=post_data)
        result = response.json()
        messages.error(request, result["Err"])

    return HttpResponseRedirect(reverse("txlist"))