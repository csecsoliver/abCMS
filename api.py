from gevent import monkey
monkey.patch_all()
from bottle import request, response, Bottle, static_file, redirect, template
from bottle_cors_plugin import cors_plugin

from fileinterface import myopen
import blog
import auth
import userroutes as ur # these are all the routes needing authorization
import auth
import os
import sys
from pathlib import Path

app = Bottle()

# Determine base directory for resources  
# For PEX: files are extracted to a temp directory, use __file__ parent
BASE_DIR = Path(__file__).parent / 'resources'
# Frontend is sibling to Python files in PEX
FRONTEND_DIR = BASE_DIR

# Helper function to find HTML files in backend directory
def get_html_path(filename):
    """Find HTML file in BASE_DIR, works for both dev and PEX environments"""
    return BASE_DIR / filename

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
        return static_file('auth.html', root=str(BASE_DIR))
    return "Signed in"



@app.post('/getin')
def getin():
    username = request.forms.username.split('@')[0]
    password = request.forms.password
    if auth.checkpass(username, password, response):
        return "Signed in successfully. Go back to <a href='/admin.html'>the dashboard</a>."
    return "Failed to sign in"

@app.post('/getup')
def getup():
    username = request.forms.username.split('@')[0]
    password = request.forms.password
    secret = request.forms.secret
    if auth.createuser(username, password, secret, response):
        return "User signed up and in successfully. Go back to <a href='/admin.html'>the dashboard</a>."
    return "Invalid token or user exists, I won't tell which."

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
        post_dict = blog.get(pid, format='html')
        content_html = post_dict['content']
        title = post_dict['title']
        author = post_dict['author']
        color = post_dict['color']
        
        with myopen(get_html_path('postcard.html'), 'r') as f:
            postcard_template = f.read()

        postcard_html = postcard_template.format(color=color, id=pid, title=title, author=author, content=content_html)
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

        with myopen(get_html_path('post.html'), 'r') as f:
            post_template = f.read()

        return post_template.format(color=color, title=title, title1=title, author=author, content=content_html)
    return "<p>Post not found.</p>"

@app.get('/<filepath>')
def server_static(filepath):
    return static_file(filepath, root=str(FRONTEND_DIR))

app.install(cors_plugin('*'))
if __name__ == "__main__":
    app.run(host='localhost', port=8080, debug=True)
