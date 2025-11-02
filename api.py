from gevent import monkey
_ = monkey.patch_all()

import uuid
import os
from pathlib import Path

import bleach
from bottle import (
    Bottle,
    FileUpload,
    Response,
    Request,
    redirect,
    request,
    response,
    static_file,
    WSGIHeaderDict,
)
from bottle_cors_plugin import cors_plugin

import blog
import auth
import userroutes as ur  # these are all the routes needing authorization
from fileinterface import clean, getjson, html, myopen, setjson


app = Bottle()

BASE_DIR = Path(__file__).parent / "resources"
# Frontend is sibling to Python files inside the PEX
FRONTEND_DIR = BASE_DIR


def run():
    _ = auth.getsecret("cs")
    _ = auth.getsecret("ss")
    app.run(host="0.0.0.0", port=25558)


@app.get("/")
def index():
    redirect("/index.html")


@app.post("/auth_page")
def auth_page() -> str:
    
    cookie = auth.checkcookie(request)
    print("Cookie check:", cookie)
    if cookie is None:
        if auth.getsecret("ss") != "":
            with myopen(BASE_DIR / "auth.html", "r") as f:
                return f.read().format(
                    secret="""
    <label for="secret">Signup secret (mandatory to sign up):</label><br>
    <input type="password" id="secret" name="secret" required><br><br>"""
                )
        else:
            with myopen(BASE_DIR / "auth.html", "r") as f:
                return f.read().format(secret="")
    return "Signed in"


@app.post("/getin")
def getin():
    form_source = request.headers.get("Referer", "") # pyright: ignore[reportAttributeAccessIssue]
    username = request.forms["username"].split("@")[0]
    username = "".join(c for c in username if c.isalnum() or c in ("_", "-")).strip()
    password = request.forms["password"]
    if auth.checkpass(username, password, response):
        return html("Signed in successfully.", form_source, cozy=("cozy" in form_source))
    return html("Failed to sign in.", form_source, cozy=("cozy" in form_source))

@app.post("/getup")
def getup():
    form_source = request.headers.get("Referer", "") # pyright: ignore[reportAttributeAccessIssue]
    username = request.forms["username"].split("@")[0]
    username = "".join(c for c in username if c.isalnum() or c in ("_", "-")).strip()
    password = request.forms["password"]
    secret = request.forms["secret"] if auth.getsecret("ss") != "" else ""
    if auth.createuser(username, password, secret, response):
        return html("User signed up and in successfully.", form_source, cozy=("cozy" in form_source))
    return html("Invalid token (if applicable) or user exists, I won't tell which.", form_source, cozy=("cozy" in form_source))


@app.route("/user/<route:path>", method=["GET", "POST"]) # type: ignore
def user(route: str):
    username = auth.checkcookie(request)
    if username is None:
        return "Not logged in or session expired."
    return ur.routes[route.split("/")[0]](
        username, request, response, *route.split("/")[1:]
    )


@app.get("/posts")
def get_posts():
    postids = blog.listids(reverse=True)
    posts_html = ""
    for pid in postids:
        post_dict = blog.get(pid, format="html", truncate=True)
        content_html = post_dict["content"]
        title = post_dict["title"]
        author = post_dict["author"]
        color = post_dict["color"]
        textcolor = blog.get_text_color(color)

        with myopen(BASE_DIR / "postcard.html", "r") as f:
            postcard_template = f.read()
        social = auth.get_prefs(author.split("<")[0]).get("social", "#")
        if "://" not in social:
            social = "http://" + social

        postcard_html = postcard_template.format(
            bgcolor=color,
            color=textcolor,
            id=pid,
            title=title,
            author=author,
            content=content_html,
            social=social,
        )
        posts_html += postcard_html

    return posts_html if posts_html else "<p>No posts available.</p>"


@app.get("/post/<postid>")
def get_post(postid: str):
    
    post_dict = blog.get(postid, format="html")
    if post_dict:
        content_html = post_dict["content"]
        title = post_dict["title"]
        author = post_dict["author"]
        color = post_dict["color"]
        textcolor = blog.get_text_color(color)

        with myopen(BASE_DIR / "post.html", "r") as f:
            post_template = f.read()
        social = auth.get_prefs(author.split("<")[0]).get("social", "#")
        if "://" not in social:
            social = "http://" + social

        return post_template.format(
            bgcolor=color,
            color=textcolor,
            title=title,
            title1=title,
            author=author,
            content=content_html,
            social=social,
        )
    return html("Post not found.")


@app.post("/cozy/postimg")
def cozy_postimg(): # this is the endpoint for uploading images to cozy with an app. No need to return anything other than text.
    print("Cozy postimg request from", request.params["user"])
    if request.params["token"] != auth.get_prefs(request.params["user"])["cozy_token"]:
        response.status = 401
        return response
    upload: FileUpload = request.files["image"]
    if not upload:
        response.status = 400
        response.body = "Bad Request: Missing image file"
        return response
    filename = clean(str(upload.filename))
    postid = uuid.uuid4().hex + "_" + str(filename)
    save_path = Path("cozy") / request.params["user"] / "pending" / postid
    ext = os.path.splitext(str(filename))[1]
    if ext.lower() not in [".jpg", ".jpeg", ".png", ".gif", ".webp"]:
        response.status = 400
        response.body = "Unsupported file type"
        return response
    myopen(save_path, "wb").close()
    _ = upload.save(str(save_path), overwrite=True)
    pending = getjson(Path("cozy") / request.params["user"] / "pending.json")
    if "pending" not in pending:
        pending["pending"] = []
    _ = pending["pending"].append({"id": postid, "path": str(save_path)})
    setjson(Path("cozy") / request.params["user"] / "pending.json", pending)
    response.status = 201
    response.body = "Image uploaded successfully."
    return response


@app.get("/cozy/posts")
def cozy_getposts():
    username = request.query.get("username", "") # pyright: ignore[reportAttributeAccessIssue]
    print("Cozy getposts request, filter by user:", username)
    cozy_data = getjson(Path("cozy") / "posts.json")
    posts = cozy_data.get("posts", [])
    _ = posts.reverse()
    posts_html = ""
    for post in posts:
        with myopen(BASE_DIR / "imgcard.html", "r") as f:
            card_template = f.read()
        filename = post["id"]
        author = post["path"].split("/")[1]
        title = post["title"].replace("[spookify]", "")
        title = bleach.clean(title)
        title = (
            f'<h2><a>{bleach.clean(title)}</a></h2>'
            if title
            else ""
        )
        card_html = card_template.format(filename=filename, title=title, author=author)
        if author == username or username == "":
            posts_html += card_html
    return posts_html if posts_html else "<p>No images available.</p>"





@app.get("/cozy/img/<filename>")
def cozy_getimg(filename: str) -> Response:
    form_source = request.headers.get("Referer", "") # pyright: ignore[reportAttributeAccessIssue]
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
    response.body = html("Image not found.", form_source, cozy=True)
    return response

@app.get("/cozy/post/<filename>")
def cozy_getpost(filename: str) -> Response:
    posts: list = getjson(Path("cozy") / "posts.json").get("posts", [])
    postdata = None
    for i in posts:
        if i["id"] == filename:
            postdata = i
            break
    with myopen(BASE_DIR / "imgpost.html", "r") as f:
        post_template = f.read()
    if postdata:
        spookify = "false"
        if "[spookify]" in postdata["title"]:
            postdata["title"] = postdata["title"].replace("[spookify]", "")
            spookify = "true"
        response.body = post_template.format(title = bleach.clean(postdata["title"]), postid=postdata["id"], author=postdata["path"].split("/")[1], spookify=spookify)
    else:
        response.status = 404
        response.body = html("Post not found.", "/cozy.html", cozy=True)
    return response

@app.get("/<filepath>")
def server_static(filepath: str) -> Response:
    return static_file(filepath, root=str(FRONTEND_DIR))


_ = app.install(cors_plugin("*"))
if __name__ == "__main__":
    app.run(host="localhost", port=8080, debug=True)
