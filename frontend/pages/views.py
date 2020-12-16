from django.views.generic import TemplateView
from django.shortcuts import redirect, render
from django.contrib import messages
from django.http import HttpResponse, HttpResponseRedirect
from django.urls import reverse
from .models import *

def HomePageView(request):
    context = checkLogin(request)
    return render(request, "pages/home.html", context)

def AboutPageView(request):
    r = requests.get(gin_addr + "/hello")
    context = r.json()
    context.update(checkLogin(request))
    return render(request, "pages/about.html", context)

def ListTransaction(request):
    context = {}
    context.update(checkLogin(request))
    if "user_is_authenticated" in context:
        context.update(AcquireScore(request, context["user_name"]))
    context.update(AcquireAllTx(request))
    context["display_gist"] = True
    # print(context)
    return render(request, "pages/txlist.html", context)

def CreateTxView(request):
    context = checkLogin(request)
    return render(request, "pages/createtx.html", context)

def LoginView(request):
    context = checkLogin(request)
    return render(request, "pages/login.html", context)

def AuthUser(request):
    err_context = {}
    if request.POST:
        name = request.POST["name"]
        key = request.POST["key"]
        if ValidateKey(name, key):
            context = {}
            context["user_is_authenticated"] = True
            context["user_name"] = name
            context.update(AcquireRelatedTx(request, name))
            context.update(AcquireScore(request, name))
            rep = render(request, "pages/txlist.html", context)
            rep.set_cookie("login_username", name)
            return rep
        else:
            err_context["err"] = "Please check your name and key"
    else:
        err_context["err"] = "Incorrect parameters"
    return render(request, "pages/login.html", err_context)

def ListRelatedTx(request):
    context = {}
    context.update(checkLogin(request))
    if "user_is_authenticated" in context:
        context.update(AcquireRelatedTx(request, context["user_name"]))
        context.update(AcquireScore(request, context["user_name"]))
        rep = render(request, "pages/txlist.html", context)
        return rep
    return HttpResponseRedirect(reverse("home"))

def CreateTx(request):
    post_data = {}
    txid = MaxTxId()
    if request.POST:
        post_data["SenderName"] = request.COOKIES.get("login_username")
        post_data["RecverName"] = request.POST["recvername"]
        post_data["value"] = request.POST["value"]
        post_data["id"] = txid + 1
        response = requests.post(gin_addr + "/createTx", data=post_data)
        result = response.json()
        messages.error(request, result["Err"])

    return HttpResponseRedirect(reverse("checkmytx"))

def Logout(request):
    rep = redirect(reverse("home"))
    rep.delete_cookie("login_username")
    return rep