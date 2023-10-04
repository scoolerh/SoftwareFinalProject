'''
class Board
    - Which pieces where
    - On board, captured, finished
class Game
    - whose turn is it
    - check for win (later)
    - get move
    - do move
    - know valid moves (later)
    - validate move (could be part of above)
class Player
    - subclasses: Human and AI
    - figure out what move to do
    - send it to game class
    
'''

#go or python 
#docker - what in each container
#one endpoint? many endpoints, connect them together later
#frontend tech: HTTP communication. as easy as python, html, css? yes we can do render_template (ask about react, maybe use view)
# probably HTTP for communication. Hannah can decide view/html etc. 

#idea for how board state should look:
#{1: 'ww', 2: '', }
class Board():
    '''inputs: initial state'''
    def __init__(self, initialState):
        self.currentState = initialState

    def move(player, move):
        #move piece
        #self.currentState = newState
        pass

class Game():
    def __init__(self, player1, player2):
        self.player1 = player1
        self.player2 = player2


class Player():
    '''send move to game class'''
    def __init__(self):
        #get players from login/user interface?
        #will there be a route for choosing players 
        #(ie i want to play w a friend or w an ai)
        #have to figure out how we know who is playing
        pass



class HumanPlayer(Player):
    def __init__(self):
            pass
class AIPlayer(Player):
    def __init__(self):
        pass