from gevent import monkey
monkey.patch_all()
from bottle import app, route, post, request, response, Bottle
from bottle_cors_plugin import cors_plugin

import auth
import userroutes as ur # these are all the routes needing authorization

import auth
app = Bottle()

def run():
    auth.getsecret("cs")
    auth.getsecret("ss")
    app.run(host='0.0.0.0', port=25556)


@app.route('/hello')
def hello():
    yield "Hello World!"

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

app.install(cors_plugin('*'))
if __name__ == "__main__":
    app.run(host='localhost', port=8080, debug=True, reloader=True)
