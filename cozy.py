import json
import auth
import markdown
import random
from fileinterface import myopen
from pathlib import Path
from bottle import Request, Response
import bleach
import webcolors
import colorsys

def main(username, request: Request, response: Response, *args):
    response = routes[args[0]](username, request, response, *args[1:])
    return response

def upload(username, request: Request, response: Response, *args):
    if request.method != 'PUT':
        response.status = 405
        return response
    
    return response


routes = {
    "post": post,
    "get": get,
    "listids": listids,
    
}