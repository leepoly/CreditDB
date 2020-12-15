from django.urls import path

from .views import *

urlpatterns = [
    path('', HomePageView.as_view(), name='home'),
    path('about/', AboutPageView, name='about'),
    path('txlist/', ListTransaction, name='txlist'),
    path('createtx/', CreateTxView, name='createtx'),
    path('processcreatetx/', CreateTx),
    path('listrelatedtx/', ListRelatedTx),
    path('login/', ListRelatedTxView, name='login'),
]
