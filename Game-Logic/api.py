import flask
from data import getProfiles

app = flask.Flask(__name__)

@app.route("/")
def home():
    return "Welcome to BACKGAMMON"

@app.route("/login/")
def findUser():
    username = "Sophie"
    userid = 0
    profiles = getProfiles()
    for profile in profiles: 
        if (profile['Name'] == username):
            userid = profile['Id']
    if (userid == 0):
        return "No user found"
    return "Welcome " + profile['Name'] + ", player ID " + profile['Id']

if __name__ == '__main__':
    app.run(host="0.0.0.0", port=8095, debug=True)