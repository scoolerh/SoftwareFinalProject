import random

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

    def move(self, player, move):
        #move piece
        #self.currentState = newState
        pass

class Game():
    def __init__(self, player1, player2): #id for human, maybe 0 for AI?
        self.player1 = player1
        self.player2 = player2
    
    def rollDice(self):
        die1 = random.randint(1,6)
        die2 = random.randint(1,6)
        return (die1, die2)
    
    def getPossibleMoves(self, player, board): #should this player just be b(lack) or w(hite)?
        #this is simple code that won't work, it just shows the idea!
        possibleMoves = ()
        i = 0
        for entry in board:
            if entry[0] == player:
                possibleMoves.append(i, 2)
            i+=1
        return possibleMoves
    
    def getMove(self, player):
        possibleMoves = self.getPossibleMoves(player, Board.board)
        move = player.move(possibleMoves)
        return


class Player():
    '''send move to game class'''
    def __init__(self):
        #what needs to be common for the two?
        #some move function, but they need to be implemented differently.
        #should players keep track of their own pieces? Or should the game just send it a list of possible moves?
        # # Or just the full game state for now?
        pass



class HumanPlayer(Player):
    #get players from login/user interface?
    #will there be a route for choosing players 
    #(ie i want to play w a friend or w an ai)
    #have to figure out how we know who is playing
        #some ID thing, but the backend game logic does not need to care, as long as we have it stored somewhere to send to the db
    def __init__(self):
            pass
class AIPlayer(Player):
    def __init__(self): #for the future: difficulty level?
        pass
    def move(possibleMoves):
        #remember that "first piece" is the opposite for the two players, since they play the board in opposite directions
            #just think about how to handle that
        #more complicated logic should come here
        #find first piece that belongs to player
        #tell to move it 2 places or so.
        return possibleMoves[0]
        