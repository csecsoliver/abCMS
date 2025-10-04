
from bottle import Request, Response
from fileinterface import myopen
def get_deployments(username, request, response, *args):
    content = myopen(f'deployments/{username}/deployments.txt', 'r').read()
    html = '<h1>Deployments</h1><ul>'
    
    for deployment in content.splitlines():
        
        toadd = """
        //js
        <li>
            {line} 
            <button hx-get="/user/deployment/{line}" hx-target="#deployment-details" hx-replace="innerHTML">View details</button>
        </li>"""
        
        toadd = toadd.format(line=deployment)
        html += toadd

    html += '</ul>'
    response.content_type = 'text/html'
    return html


def deployment_details(username, request: Request, response: Response, *args):
    deployment_name = args[0]
    html = f'<h1>Details for {deployment_name}</h1><pre>{content}</pre>'
    response.content_type = 'text/html'
    return html





routes = {
    "get_deployments": get_deployments,
    "deployment": deployment_details,
}