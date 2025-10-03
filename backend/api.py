import blog
from bottle import app, route, post, request, response
from gevent import monkey
monkey.patch_all()
app = Bottle()


def run():
    app.run(host='localhost', port=8080)


@app.route('/hello')
def hello():
    yield "Hello World!"


@app.post('/admin/auth')
def auth():
    pass


if __name__ == "__main__":

    app.run(host='localhost', port=8080, debug=True, reloader=True)
