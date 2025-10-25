import json
import auth
import markdown
import random
from fileinterface import myopen, getjson, setjson, html
from pathlib import Path
from bottle import Request, Response, FileUpload
import bleach
import webcolors
import colorsys
from dbm import sqlite3
import uuid
import os

BASE_DIR = Path(__file__).parent / 'resources'
# sibling to Python files inside the PEX
FRONTEND_DIR = BASE_DIR

def main(username, request: Request, response: Response, *args):
    response = routes[args[0]](username, request, response, *args[1:])
    return response

def postimg(username, request: Request, response: Response, *args):
    if request.method != 'POST':
        response.status = 405
        response.body = "Method Not Allowed"
        return response
    if request.content_type != 'multipart/form-data':
        response.status = 400
        response.body = "Bad Request: Expected multipart/form-data"
        return response
    upload: FileUpload = request.files.get("image")
    title = request.forms.get("title", "")
    if not upload:
        response.status = 400
        response.body = "Bad Request: Missing image file"
        return response
    
    postid = (uuid.uuid4().hex + "_" + upload.filename)
    save_path = Path("cozy") / username / "images" / postid
    cozy_data = getjson(Path("cozy") / "posts.json")
    if "posts" not in cozy_data:
        cozy_data["posts"] = []
    cozy_data["posts"].append({
        "id": postid,
        "title": title,
        "path": str(save_path)
    })
    setjson(Path("cozy") / "posts.json", cozy_data)
    with myopen(save_path, 'wb') as f:
        f.write(upload.file.read())
    response.status = 201
    response.body = "Image uploaded successfully."
    return response

def getpending(username, request: Request, response: Response, *args):
    cozy_data = getjson(Path("cozy") / username / "pending.json")
    pending = cozy_data.get("pending", [])
    
    for i in pending:
        with myopen(BASE_DIR / "confirm-card.html", 'r') as f:
            response.body += f.read().format(postid=i["id"])

routes = {
    "post": postimg,
    "get": getimg,
    "listids": listids,
    
}