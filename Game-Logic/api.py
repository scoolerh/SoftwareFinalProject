import flask
#run on port 7800 
app = flask.Flask(__name__)

@app.route("/")
def home():
    return "Welcome to BACKGAMMON"

@app.route("/login/<username>/")
def findUser(username):
    return "Welcome " + str(username) 

if __name__ == '__main__':
    app.run(host="0.0.0.0", port=7800, debug=True)