from io import TextIOWrapper
import json
from typing import Optional
import blog
from bottle import Request, Response
from fileinterface import myopen
import auth

def get_deployments(username, request, response, *args):
    content = myopen(f'deployments/{username}/deployments.txt', 'r').read()
    html = '<h1>Deployments</h1><ul>'
    
    for deployment in content.splitlines():
        
        toadd = """
        <li>
            {line} 
            <button hx-get="/user/deployment/{line}" hx-target="#deployment-details" hx-replace="innerHTML">View details</button>
        </li>"""
        
        toadd = toadd.format(line=deployment)
        html += toadd

    html += '</ul>'
    response.content_type = 'text/html'
    return html


def deployment(username, request: Request, response: Response, *args):
    match args[0]:
        case "details":
            return deployment_details(username, request, response, *args[1:])
        case "create":
            return create_deployment(username, request, response, *args[1:])
        
        case _:
            response.status = 404
            return "Not Found"

def deployment_details(username: str, request: Request, response: Response, *args):
    
    deployment_name: str = args[0]
    content: str = myopen(f'deployments/{username}/{deployment_name}.po', 'r').read()
    
    if len(content) == 0:
        response.status = 404
        return "Not Found"
    
    content: dict = json.loads(content)
    
    html = f'<h1>Details for {deployment_name}</h1><ul>'
    
    for key, value in content.items():
        html += f'<li><strong>{key}:</strong> {value}</li>'
        
    html += '</ul>'
    response.content_type = 'text/html'
    return html


def create_deployment(username, request: Request, response: Response, *args):
    depdict = {}
    depdict["author"] = username
    depdict["git"] = request.forms.git # input url
    depdict["bash"] = request.forms.bash 
    depdict["port"] = request.forms.port # input number
    depdict["subdomain"] = request.forms.subdomain
    depdict["name"] = request.forms.name

    with myopen(f'deployments/{username}/{depdict["name"]}.json', 'w') as f:
        json.dump(depdict, f)
        with myopen(f'deployments/{username}/deployments.txt', 'a') as f:
            f.write(f'{depdict["name"]}\n')
            
    
    
    html = "Successfully created deployment!"
    response.content_type = 'text/html'
    return html

def create_post(username, request: Request, response: Response, *args):
    title = request.forms.title
    content = request.forms.content
    postid = blog.post(title, username, content)
    # Award coins for posting
    auth.update_coins(username, 10)
    response.status = 201
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

routes = {
    "deployments": get_deployments,
    "deployment": deployment,
    "createdeployment": create_deployment,
    "createpost": create_post,
    "deletepost": delete_post,
    "coins": get_coins,
}