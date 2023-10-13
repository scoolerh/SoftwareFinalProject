import flask
import json
import random
from gamelogic import *
import requests

#run on port 9000 
app = flask.Flask(__name__)

#0 and 25 are the homes of black and white respectively, not really a part of the board. Plese keep this in mind!
initialState = {'0': '', '1': 'ww', '2': '', '3': '', '4': '', '5': '', '6': 'bbbbbb', '7': '', '8': 'bbb', '9': '', '10': '', '11': '', '12': 'wwwwww',
                 '13': 'bbbbb', '14': '', '15': '', '16': '', '17': 'www', '18': '', '19': 'wwwww', '20': '', '21': '', '22': '', '23': '', '24': 'bb', '25': ''}
p1 = 0
p2 = 0
games = []

#for testing purposes, remove before merging
game = Game(0, p1, p2, initialState)
games.append(game)

'''Print how to use for the user.'''
@app.route("/")
def home():
    return "Welcome to BACKGAMMON. <br>To start a new game, use the /newgame/ route. This will display the starting board and your game id.<br>To see how many wins or losses you have, you can use the /scoreboard/(player)/ route, where player is your username. For Matt, your username is Matt.<br>Enjoy!!!"


'''todo: Allow a user to log in or sign up if they do not have a username'''
@app.route("/login/<username>/")
def findUser(username):
    return "Welcome " + str(username) 

'''Starts a new game for the user and displays the initial board'''
@app.route("/newgame/")
def startGame():
    p1 = 0
    p2 = 0
    gameid = len(games)
    game = Game(gameid, p1, p2, initialState)
    games.append(game)
    currentState = json.dumps(game.board.currentState)
    return requests.get(f"http://frontend:5000/displayboard/{currentState}/").text + "<br>Your game ID is " + str(gameid)

'''todo: Check whose turn it is, if the game is won, have the player make a move'''
@app.route('/play/<gameid>/')
def play(gameid):
    game = games[int(gameid)]
    #plays 10 move
    currentState = json.dumps(game.board.currentState)
    return requests.get(f"http://frontend:5000/displayboard/{currentState}/").text

'''todo: if someone has won, update the database with wins/losses for each player. Print final board.'''
@app.route('/gamewon/')
def won():
    return "Hannah won!"

'''todo: set up SQL database, check if the user is an actual user in the db, then return their win/loss ratio.'''    
@app.route("/scoreboard/<player>/")
def scoreboard(player):
    return requests.get(f"http://db:5000/getprofile/{player}/").text

if __name__ == '__main__':
    app.run(host="0.0.0.0", port=9000, debug=True)