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
#frontend tech: HTTP communication. as easy as python, html, css? yes we can do render_template (ask about react, maybe use view(vue?))
# probably HTTP for communication. Hannah can decide view/html etc. 

#idea for how board state should look:
#{1: 'ww', 2: '', 3: '', 4: '', 5: '', 6: 'bbbbbb', etc...}
class Board():
    '''inputs: initial state'''
    def __init__(self, initialState):
        self.currentState = initialState

    def move(self, player, move):
        #move piece
        #self.currentState = newState
        pass

    def doMove(self, playerColor, move): 
        #move looks like (slot to move from, # of places to move)
        #player format should be 'b' or 'w', if not this then translate
        currState = self.currentState
        #might want to do two moves
        originalSpace = move[0]
        if playerColor == "w":
            newSpace = originalSpace + int(move[1])
        if playerColor == "b":
            newSpace = originalSpace - int(move[1])
        updatedOriginalSpace = currState[str(originalSpace)]
        currState[str(originalSpace)] = updatedOriginalSpace[0:-1]
        updatedNewSpace = currState[str(newSpace)]
        currState[str(newSpace)] = updatedNewSpace + playerColor
        self.currentState = currState
        return self.currentState

class Game():
    def __init__(self, id, player1, player2, initialState): #id for human, maybe 0 for AI? 
        if player1 == 0:
            self.player1 = AIPlayer("w")
        else:
            print("!")
        if player2 == 0:
            self.player2 = AIPlayer("b")   
        else: 
            print("!")
        self.board = Board(initialState)
        self.gid = id
    
    def rollDice(self):
        die1 = random.randint(1,6)
        die2 = random.randint(1,6)
        return (die1, die2)
    #initialState = {1: 'ww', 2: '', 3: '', 4: '', 5: '', 6: 'bbbbbb', 7: '', 8: 'bbb', 9: '', 10: '', 11: '', 12: 'wwwwww',
    #             13: 'bbbbb', 14: '', 15: '', 16: '', 17: 'www', 18: '', 19: 'wwwww', 20: '', 21: '', 22: '', 23: '', 24: 'bb'}
    def getPossibleMoves(self, player, board):
        #this is simple code that won't work, it just shows the idea!
        currState = board.currentState
        possibleMoves = []
        i = 0
        for i in range(24):
            i+=1
            if player.color in currState[str(i)]:
                possibleMoves.append((i, 2)) #2 is dice roll
        return possibleMoves
    
    def getMove(self, player):
        #gets possible moves, gets move from player
        possibleMoves = self.getPossibleMoves(player, self.board)
        move = player.move(possibleMoves)
        return move
    
    

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
            self.color = "w" #make this an input later
            pass
    
class AIPlayer(Player):
    def __init__(self, color): #for the future: difficulty level?
        self.color = color
        pass

    def move(self, possibleMoves):
        if self.color == "b":
            return (possibleMoves[-1]) #minus because black moves backwards, white moves forwards
        elif self.color == "w":
            return (possibleMoves[0])
        else:
            return("Error. Not a valid color. Must be b or w")
        #remember that "first piece" is the opposite for the two players, since they play the board in opposite directions
            #just think about how to handle that
        #more complicated logic should come here
        
