from io import TextIOWrapper
import json
from pathlib import Path
from typing import Optional
import blog
from bottle import Request, Response, static_file
from fileinterface import myopen
import auth
import deploy
import cozy

FRONTEND_DIR = Path(__file__).parent / 'resources'

def create_post(username, request: Request, response: Response, *args):
    title = request.forms.title
    color = request.forms.color
    content = request.forms.content
    
    if color == "#ffffff":
        auth.update_coins(username, 10)
    else:
        if auth.get_coins(username) >= 50:
            auth.update_coins(username, -50)
            response.status = 201
            response.add_header('HX-Trigger', "postlist")
            postid = blog.post(title, username, content, color)
            response.body = f"Post created with ID {postid}. You spent 50 mana!"
            return response
        else:
            response.status = 299
            response.body = "Not enough mana to enchant the post."
            return response
    response.status = 201
    response.add_header('HX-Trigger', "postlist")
    postid = blog.post(title, username, content, color)
    response.body = f"Post created with ID {postid}. You earned 10 mana!"
    return response


def delete_post(username, request: Request, response: Response, *args):
    postid = args[0]
    success = blog.delete(postid, username)
    if success:
        return f"Post {postid} deleted."
    else:
        response.status = 403
        return "You do not have permission to delete this post."

def get_coins(username, request: Request, response: Response, *args):
    coins = auth.get_coins(username)
    html = f'You have {coins} mana.'
    response.content_type = 'text/html'
    return html

def get_dashboard(username, request: Request, response: Response, *args):
    with myopen(FRONTEND_DIR / 'dashboard-sect.html', 'r') as f:
        html = f.read().format(username=username)
    response.content_type = 'text/html'
    return html

def get_settings(username, request: Request, response: Response, *args):
    settings = auth.get_prefs(username)
    with myopen(FRONTEND_DIR / 'settings.html', 'r') as f:
        html = f.read().format(
            social=settings.get('social', ''),
        )
    response.content_type = 'text/html'
    return html

def save_settings(username, request: Request, response: Response, *args):
    formdict = {}
    formdict["social"] = request.forms.social
    auth.save_prefs(username, formdict)
    response.add_header('HX-Trigger', "postlist")
    return 'Settings saved!'

routes = {
    "createpost": create_post,
    "deletepost": delete_post,
    "coins": get_coins,
    "dashboard": get_dashboard,
    "getsettings": get_settings,
    "settings": save_settings,
    "cozy": cozy.main,
}