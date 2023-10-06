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

class Board():
    '''inputs: initial state'''
    def __init__(self, initialState):
        self.currentState = initialState

    def doMove(self, playerColor, move): 
        #move looks like (slot to move from, # of places to move)
        currState = self.currentState
        #might want to do two moves!! Figure out how to handle this
        originalSpace = move[0]
        newSpace = originalSpace + int(move[1])
        originalSpaceState = currState[str(originalSpace)]
        currState[str(originalSpace)] = originalSpaceState[0:-1]
        newSpaceState = currState[str(newSpace)]
        currState[str(newSpace)] = newSpaceState + playerColor
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
        return [die1, die2]
    
    def getPossibleMoves(self, player, diceRoll, board):
        currState = board.currentState
        possibleMoves = []
        die1 = diceRoll[0]
        die2 = diceRoll[1]
        i = 0
        if player.color == "w":
            print("white to move")
            for i in range(24):
                i+=1
                if "w" in currState[str(i)]:
                    print("found slot with white at", i)
                    if 24-i>die1:
                        goalPlace = currState[str(i+die1)]
                        if not("b" in goalPlace and len(goalPlace) >= 2):
                            possibleMoves.append((i, die1))
                            print("adding a possible move: ", (i, die1))
                    if 24-i>die2:
                        goalPlace = currState[str(i+die2)]
                        if not("b" in goalPlace and len(goalPlace) >= 2):
                            possibleMoves.append((i, die2))
                            print("adding a possible move: ", (i, die2))
        elif player.color == "b":
            print("black to move")
            for i in range(24):
                i+=1
                if "b" in currState[str(i)]:
                    print("found slot with black at", i)
                    if i>= die1:
                        goalPlace = currState[str(i-die1)]
                        if not("w" in goalPlace and len(goalPlace) >= 2):
                            possibleMoves.append((i, -die1))
                            print("adding a possible move: ", (i, -die1))
                    if i>die2:
                        goalPlace = currState[str(i-die2)]
                        if not("w" in goalPlace and len(goalPlace) >= 2):
                            possibleMoves.append((i, -die2))
                            print("adding a possible move: ", (i, -die2))
        return possibleMoves
    
    def getMove(self, player):
        #gets possible moves, gets move from player
        diceRoll = self.rollDice() #diceroll might have to be an endpoint instead of being called from getMove!
        possibleMoves = self.getPossibleMoves(player, diceRoll, self.board)
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
        #more complicated logic should come here
        
