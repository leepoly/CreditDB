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
        imgEnc = request.POST["imgEnc"]
        face_rec_name = SearchFace(imgEnc)
        validated = False
        if len(face_rec_name) > 0:
            name = face_rec_name
            validated = True
        else:
            validated = ValidateKey(name, key)
        if validated:
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
    txid = MaxTxId()
    if request.POST:
        post_data = {}
        post_data["SenderName"] = request.POST["getloanfrom"]
        post_data["RecverName"] = request.COOKIES.get("login_username")
        post_data["value"] = request.POST["value"]
        post_data["id"] = txid + 1
        response = requests.post(gin_addr + "/createTx", data=post_data)
        result = response.json()
        err = result["Err"]
        if len(err)>0:
            err_context = {}
            err_context["err"] = err
            rep = render(request, "pages/txlist.html", err_context)
            return rep

    return HttpResponseRedirect(reverse("checkmytx"))

def Logout(request):
    rep = redirect(reverse("home"))
    rep.delete_cookie("login_username")
    return rep

def SignupView(request):
    context = {}
    context["signup"] = True
    context["err"] = "To sign up, you should make a transaction to our assigned accounts, then you will receive your personal key"
    return render(request, "pages/signup.html", context)

def HttpsAuth(request):
    return render(request, ".well-known/pki-validation/9397C7D639E184C70EB8B7F709355471.txt")

def ProcessSignup(request):
    txid = MaxTxId()
    err_context = {}
    if request.POST:
        post_data = {}
        name = request.POST["identity"]
        imgEnc = request.POST["imgEnc"]
        err = ""
        if imgEnc:
            if not AppendFace(name, imgEnc):
                err = "Append face fails. Please retake."
        if len(err) == 0:
            post_data["SenderName"] = name
            post_data["RecverName"] = "ProfXuTHU"
            post_data["value"] = "0.01"
            post_data["id"] = txid + 1
            response = requests.post(gin_addr + "/createTx", data=post_data)
            result = response.json()
            err = result["Err"]
        if len(err) > 0:
            err_context["err"] = err
        else:
            err_context["err"] = "Sign up successfully. Your key: \'" + AcquireKey(name) + "\'. Now please login with your identity."
    else:
        err_context["err"] = "Incorrect parameters"
    rep = render(request, "pages/signup.html", err_context)
    return rep
