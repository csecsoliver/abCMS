from gevent import monkey
monkey.patch_all()
from bottle import app, route, post, request, response, Bottle, static_file, redirect
from bottle_cors_plugin import cors_plugin

from fileinterface import myopen
import blog
import auth
import userroutes as ur # these are all the routes needing authorization
import auth
import os
from pathlib import Path

app = Bottle()

# Determine base directory for resources
# When running from PEX, __file__ will be in the extracted temp directory
BASE_DIR = Path(__file__).parent
FRONTEND_DIR = BASE_DIR / 'frontend' if (BASE_DIR / 'frontend').exists() else BASE_DIR.parent / 'frontend'

# Helper function to find HTML files in backend directory
def get_html_path(filename):
    """Find HTML file in BASE_DIR, works for both dev and PEX environments"""
    return BASE_DIR / filename

def run():
    auth.getsecret("cs")
    auth.getsecret("ss")
    app.run(host='0.0.0.0', port=25556)


@app.get('/<filepath>')
def server_static(filepath):
    return static_file(filepath, root=str(FRONTEND_DIR))
@app.get('/')
def index():
    redirect('/index.html')
    
@app.get('/auth')
def auth_page():
    if auth.checkcookie(request) is not None:
        return "Signed in"
    return static_file('auth.html', root=str(BASE_DIR))




# https://bottlepy.org/docs/dev/api.html#bottle.BaseResponse.set_cookie
@app.post('/getin')
def getin():
    username = request.forms.username
    password = request.forms.password
    if auth.checkpass(username, password, response):
        return "Signed in successfully"
    return "Failed to sign in"

@app.post('/getup')
def getup():
    username = request.forms.username
    password = request.forms.password
    secret = request.forms.secret
    if auth.createuser(username, password, secret, response):
        return "User signed up and in successfully"
    return "Invalid token"

@app.route('/user/<route:path>')
def user(route):
    username = auth.checkcookie(request)
    if username is None:
        response.status = 401
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
        
        with myopen(get_html_path('postcard.html'), 'r') as f:
            postcard_template = f.read()
            
        postcard_html = postcard_template.format(title=title, author=author, content=content_html)
        posts_html += postcard_html
        
    return posts_html if posts_html else "<p>No posts available.</p>"

app.install(cors_plugin('*'))
if __name__ == "__main__":
    app.run(host='localhost', port=8080, debug=True, reloader=True)
