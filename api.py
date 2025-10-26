import uuid
from gevent import monkey
monkey.patch_all()
from bottle import request, response, Bottle, static_file, redirect, template, FileUpload, HTTPResponse
from bottle_cors_plugin import cors_plugin

from fileinterface import getjson, myopen, html, setjson
import blog
import auth
import userroutes as ur # these are all the routes needing authorization
import auth
import os
import sys
from pathlib import Path

app = Bottle()

BASE_DIR = Path(__file__).parent / 'resources'
# Frontend is sibling to Python files inside the PEX
FRONTEND_DIR = BASE_DIR

def run():
    auth.getsecret("cs")
    auth.getsecret("ss")
    app.run(host='0.0.0.0', port=25557)



@app.get('/')
def index():
    redirect('/index.html')
    
@app.post('/auth_page')
def auth_page():
    cookie = auth.checkcookie(request)
    print("Cookie check:", cookie)
    if cookie is None:
        if auth.getsecret("ss") != "":
            with myopen(BASE_DIR / 'auth.html', 'r') as f:
                return f.read().format(secret="""
    <label for="secret">Signup secret (mandatory to sign up):</label><br>
    <input type="password" id="secret" name="secret" required><br><br>""")
        else:
            with myopen(BASE_DIR / 'auth.html', 'r') as f:
                return f.read().format(secret="")
    return "Signed in"



@app.post('/getin')
def getin():
    username = request.forms.username.split('@')[0]
    username = ''.join(c for c in username if c.isalnum() or c in ('_', '-')).strip()
    password = request.forms.password
    if auth.checkpass(username, password, response):
        return "Signed in successfully. Go back to <a href='/admin.html'>the dashboard</a>."
    return "Failed to sign in"

@app.post('/getup')
def getup():
    username = request.forms.username.split('@')[0]
    username = ''.join(c for c in username if c.isalnum() or c in ('_', '-')).strip()
    password = request.forms.password
    secret = request.forms.secret if auth.getsecret("ss") != "" else ""
    if auth.createuser(username, password, secret, response):
        return html("User signed up and in successfully. Go back to <a href='/admin.html'>the dashboard</a>.")
    return html("Invalid token (if applicable) or user exists, I won't tell which.")

@app.route('/user/<route:path>', method=['GET', 'POST'])
def user(route):
    username = auth.checkcookie(request)
    if username is None:
        return "Not logged in or session expired."
    return ur.routes[route.split('/')[0]](username, request, response, *route.split('/')[1:])


@app.get('/posts')
def get_posts():
    postids = blog.listids(reverse=True)
    posts_html = ""
    for pid in postids:
        post_dict = blog.get(pid, format='html', truncate=True)
        content_html = post_dict['content']
        title = post_dict['title']
        author = post_dict['author']
        color = post_dict['color']
        textcolor = blog.get_text_color(color)

        with myopen(BASE_DIR / 'postcard.html', 'r') as f:
            postcard_template = f.read()
        social = auth.get_prefs(author.split("<")[0]).get('social', '#')
        if "://" not in social:
            social = "http://"+social
        
        postcard_html = postcard_template.format(bgcolor=color, color=textcolor, id=pid, title=title, author=author, content=content_html, social=social)
        posts_html += postcard_html
        
    return posts_html if posts_html else "<p>No posts available.</p>"

@app.get('/post/<postid>')
def get_post(postid):
    post_dict = blog.get(postid, format='html')
    if post_dict:
        content_html = post_dict['content']
        title = post_dict['title']
        author = post_dict['author']
        color = post_dict['color']
        textcolor = blog.get_text_color(color)

        with myopen(BASE_DIR / 'post.html', 'r') as f:
            post_template = f.read()
        social = auth.get_prefs(author.split("<")[0]).get('social', '#')
        if "://" not in social:
            social = "http://"+social

        return post_template.format(bgcolor=color, color=textcolor, title=title, title1=title, author=author, content=content_html, social=social)
    return html("Post not found.")

@app.post('/cozy/postimg')
def cozy_postimg():
    if request.params.get("token") != auth.get_prefs(request.params.get("user")).get("cozy_token", ""):
        response.status = 401
        return response
    upload: FileUpload = request.files.get("image")
    if not upload:
        response.status = 400
        response.body = "Bad Request: Missing image file"
        return response
    postid = (uuid.uuid4().hex + "_" + upload.filename)
    save_path = Path("cozy") / request.params.get("user") / "pending" / postid
    ext = os.path.splitext(upload.filename)[1]
    if ext.lower() not in [".jpg", ".jpeg", ".png", ".gif", ".webp"]:
        response.status = 400
        response.body = "Unsupported file type"
        return response
    with myopen(save_path, 'wb') as f:
        f.write(upload.file.read())
    upload.file.close()
    pending = getjson(Path("cozy") / request.params.get("user") / "pending.json")
    if "pending" not in pending:
        pending["pending"] = []
    pending["pending"].append({
        "id": postid,
        "path": str(save_path)
    })
    setjson(Path("cozy") / request.params.get("user") / "pending.json", pending)
    response.status = 201
    response.body = "Image uploaded successfully."
    return response

@app.get('/cozy/posts')
def cozy_getposts():
    cozy_data = getjson(Path("cozy") / "posts.json")
    posts = cozy_data.get("posts", [])
    posts_html = ""
    for post in posts:
        with myopen(BASE_DIR / "imgcard.html", 'r') as f:
            card_template = f.read()
        filename = post["id"]
        author = post["path"].split('/')[1]
        title = f'<h2><a href="/cozy/img/{filename}" target="_blank">{post["title"]}</a></h2>' if post["title"] else ""
        card_html = card_template.format(filename=filename, title=title, author=author)
        posts_html += card_html
    return posts_html if posts_html else "<p>No images available.</p>"

@app.get('/cozy/img/<filename>')
def cozy_getimg(filename):
    posts = getjson(Path("cozy") / "posts.json").get("posts", [])
    for i in posts:
        if i["id"] == filename:
            filepath = i["path"]
            break
    else:
        filepath = "notfound"
    if Path(filepath).exists():
        return static_file(filepath, root=os.getcwd())
    response.status = 404
    return "Image not found."

@app.get('/<filepath>')
def server_static(filepath):
    return static_file(filepath, root=str(FRONTEND_DIR))
    

app.install(cors_plugin('*'))
if __name__ == "__main__":
    app.run(host='localhost', port=8080, debug=True)


