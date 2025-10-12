from io import TextIOWrapper
import json
from pathlib import Path
from typing import Optional
import blog
from bottle import Request, Response, static_file
from fileinterface import myopen
import auth
import deploy

FRONTEND_DIR = Path(__file__).parent / 'resources'


# def get_deployments(username, request, response, *args):
#     with myopen(f'deployments/{username}/deployments.txt', 'r') as f:
#         content = f.read()
#     html = '<ul>'
#     for deployment in content.splitlines():
        
#         toadd = """
#         <li>
#             {line} 
#             <button hx-get="/user/deployment/{line}" hx-target="#deployment-details" hx-replace="innerHTML">View details</button>
#         </li>"""
        
#         toadd = toadd.format(line=deployment)
#         html += toadd

#     html += '</ul>'
#     response.content_type = 'text/html'
#     return html


# def deployment(username, request: Request, response: Response, *args):
#     match args[0]:
#         case "details":
#             return deployment_details(username, request, response, *args[1:])
#         case "create":
#             return create_deployment(username, request, response, *args[1:])
#         case _:
#             response.status = 404
#             return "Not Found"

# def deployment_details(username: str, request: Request, response: Response, *args):
    
#     deployment_name: str = args[0]
#     content: str = myopen(f'deployments/{username}/{deployment_name}.po', 'r').read()
    
#     if len(content) == 0:
#         response.status = 404
#         return "Not Found"
    
#     content: dict = json.loads(content)
    
#     html = f'<h1>Details for {deployment_name}</h1><ul>'
    
#     for key, value in content.items():
#         html += f'<li><strong>{key}:</strong> {value}</li>'
        
#     html += '</ul>'
#     response.content_type = 'text/html'
#     return html


# def create_deployment(username, request: Request, response: Response, *args):
#     depdict = {}
#     depdict["author"] = username
#     depdict["git"] = request.forms.git # input url
#     depdict["bash"] = request.forms.bash # textarea
#     depdict["port"] = request.forms.port # input number
#     depdict["subdomain"] = request.forms.subdomain
#     depdict["name"] = request.forms.name
#     if Path(f'deployments/{username}/{depdict["name"]}.json').exists():
#         response.status = 400
#         return "Deployment with that name already exists."
#     with myopen(f'deployments/{username}/{depdict["name"]}.json', 'w') as f:
#         json.dump(depdict, f)
#         with myopen(f'deployments/{username}/deployments.txt', 'a') as f:
#             f.write(f'{depdict["name"]}\n')
#         deploy
            
    
    
#     html = "Successfully created deployment!"
#     response.content_type = 'text/html'
#     return html

def create_post(username, request: Request, response: Response, *args):
    title = request.forms.title
    color = request.forms.color
    content = request.forms.content
    
    postid = blog.post(title, username, content, color)
    if color == "#ffffff":
        auth.update_coins(username, 10)
    else:
        if auth.get_coins(username) >= 50:
            auth.update_coins(username, -50)
            response.status = 201
            response.add_header('HX-Trigger', "postlist")
            return f"Post created with ID {postid}. You spent 50 coins!"
        else:
            response.status = 400
            return "Not enough coins to enchant the post."
    response.status = 201
    response.add_header('HX-Trigger', "postlist")
    return f"Post created with ID {postid}. You earned 10 coins!"


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
    html = f'You have {coins} coins.'
    response.content_type = 'text/html'
    return html

def get_dashboard(username, request: Request, response: Response, *args):
    with myopen(FRONTEND_DIR / 'dashboard-sect.html', 'r') as f:
        html = f.read().format(username=username)
    response.content_type = 'text/html'
    return html

routes = {
    # "deployments": get_deployments,
    # "deployment": deployment,
    # "createdeployment": create_deployment,
    "createpost": create_post,
    "deletepost": delete_post,
    "coins": get_coins,
    "dashboard": get_dashboard
}