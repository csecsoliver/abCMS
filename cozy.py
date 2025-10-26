import json
import auth
import markdown
import random
from fileinterface import myopen, getjson, setjson, html
from pathlib import Path
from bottle import Request, Response, FileUpload, static_file
import bleach
import webcolors
import colorsys
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
    # if request.content_type != 'multipart/form-data':
    #     response.status = 400
    #     response.body = "Bad Request: Expected multipart/form-data"
    #     print("print:", response.body)
    #     return response
    upload: FileUpload = request.files.get("image")
    title = request.forms.get("title", "")
    if not upload:
        response.status = 400
        response.body = "Bad Request: Missing image file"
        print("print:", response.body)
        return response
    
    postid = (uuid.uuid4().hex + "_" + upload.filename)
    save_path = Path("cozy") / username / "images" / postid

    ext = os.path.splitext(upload.filename)[1]
    if ext.lower() not in [".jpg", ".jpeg", ".png", ".gif", ".webp"]:
        response.status = 400
        response.body = "Unsupported file type"
        return response
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
    upload.file.close()
    response.status = 201
    response.body = "Image uploaded successfully."
    return response

def getpending(username, request: Request, response: Response, *args):
    cozy_data = getjson(Path("cozy") / username / "pending.json")
    pending = cozy_data.get("pending", [])
    
    for i in pending:
        with myopen(BASE_DIR / "confirm-card.html", 'r') as f:
            response.body += f.read().format(postid=i["id"])
    
    return response

def getimg(username, request: Request, response: Response, *args):
    save_path = Path("cozy") / username / "images" / args[0]
    if not save_path.exists():
        response.status = 404
        response.body = "Image not found."
        return response
    return static_file(save_path, root=os.getcwd())
def dashboard(username, request: Request, response: Response, *args):
    with myopen(BASE_DIR / 'imgdashboard-sect.html', 'r') as f:
        dashboard_template = f.read()
    response.body = dashboard_template.format(username=username)
    return response
routes = {
    "postimg": postimg, # manual form submission endpoint
    "getimg": getimg, # needs another route section to point to an image
    "getpending": getpending, # html for pending images
    "dashboard": dashboard
}