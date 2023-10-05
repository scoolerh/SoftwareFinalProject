import flask
import json

app = flask.Flask(__name__)

profiles = [{'Id': 1, 'Name': 'Sophie', 'Wins': 1, 'Losses': 10},{'Id': 2, 'Name': 'Matt', 'Wins': 10, 'Losses': 1},{'Id': 3, 'Name': 'Hannah', 'Wins': 500, 'Losses': 0},{'Id': 4, 'Name': 'Selma', 'Wins': 5, 'Losses': 5}, {'Id': 5, 'Name': 'Bryan', 'Wins': 0, 'Losses': 0}]
games = [{'Id': 0, 'Turn': '', 'Roll': '', 'Move': ''}]

@app.route('/getprofile/<username>/')
def getProfile(username):
    for profile in profiles: 
        if (profile['Name'] == username): 
            return json.dumps(profile)
    return "Who even are you"

@app.route('/getgame/<gameId>/')
def getGame(gameId):
    for game in games: 
        if (game['Id'] == gameId):
            return json.dumps(game)
    return "Did you even play a game here"

if __name__ == '__main__':
    app.run(host="0.0.0.0", port=7788, debug=True)