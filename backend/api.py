from gevent import monkey; monkey.patch_all()
from bottle import *
app = Bottle()
def run():
    app.run(host='localhost', port=8080)

@app.route('/hello')
def hello():
    yield "Hello World!"


if __name__ == "__main__":
    
    app.run(host='localhost', port=8080, debug=True, reloader=True)
    