import random
import flask

app = flask.Flask(__name__)

app.route('/')
def home():
    return "Welcome to BACKGAMMON"

if __name__ == '__main__':
    app.run(host="0.0.0.0", port=8095, debug=True)