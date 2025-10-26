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
    # Save the uploaded file
    upload.save(str(save_path), overwrite=True)
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

def confirm(username, request: Request, response: Response, *args):
    postid = args[0]
    pending_path = Path("cozy") / username / "pending.json"
    pending = getjson(pending_path)
    pending_list = pending.get("pending", [])
    for i in pending_list:
        if i["id"] == postid:
            image_path = Path(i["path"])
            final_path = Path("cozy") / username / "images" / postid
            if request.forms.get("approve") == "reject":
                os.remove(image_path)
                pending_list.remove(i)
                response.body = "Image rejected and removed."
                response.status = 200
                return response
            elif request.forms.get("approve") == "approve":
                os.rename(image_path, final_path)
            else:
                response.status = 400
                response.body = "Bad Request: Missing or invalid approval action"
                return response
            pending_list.remove(i)
            break
    pending["pending"] = pending_list
    setjson(pending_path, pending)
    
    posts = getjson(Path("cozy") / "posts.json")
    if "posts" not in posts:
        posts["posts"] = []
    title = request.forms.get("title", "")
    posts["posts"].append({
        "id": postid,
        "title": title,
        "path": str(final_path)
    })
    setjson(Path("cozy") / "posts.json", posts)
    
    response.body = "Image confirmed and moved to posts."
    return response

def getsettings(username, request: Request, response: Response, *args):
    settings = auth.get_prefs(username)
    with myopen(FRONTEND_DIR / 'settingscozy.html', 'r') as f:
        html = f.read().format(
            cozy_token=settings.get('cozy_token', '')
        )
    response.content_type = 'text/html'
    return html
def savesettings(username, request: Request, response: Response, *args):
    settings = auth.get_prefs(username)
    settings["cozy_token"] = request.forms.cozy_token
    auth.save_prefs(username, settings)
    response.add_header('HX-Trigger', "postlist")
    return 'Settings saved!'

routes = {
    "postimg": postimg, # manual form submission endpoint
    "getimg": getimg, # needs another route section to point to an image
    "getpending": getpending, # html for pending images
    "confirm": confirm,
    "dashboard": dashboard,
    "getsettings": getsettings, # html for settings
    "savesettings": savesettings
}