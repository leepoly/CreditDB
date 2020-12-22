from django.urls import path

from .views import *

urlpatterns = [
    path('', HomePageView, name='home'),
    path('about/', AboutPageView, name='about'),
    path('txlist/', ListTransaction, name='txlist'),
    path('createtx/', CreateTxView, name='createtx'),
    path('processcreatetx/', CreateTx),
    path('checkmytx/', ListRelatedTx, name='checkmytx'),
    path('auth/', AuthUser),
    path('login/', LoginView, name='login'),
    path('logout/', Logout, name='logout'),
    path('signup/', SignupView, name='signup'),
    path('processsignup/', ProcessSignup),
    path('.well-known/pki-validation/9397C7D639E184C70EB8B7F709355471.txt', HttpsAuth), # for https cert auth
]
