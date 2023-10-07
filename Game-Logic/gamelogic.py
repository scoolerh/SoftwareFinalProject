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
        die = int(move[1])
        #will want to do two moves!! Figure out how to handle this
        originalSpace = move[0]
        newSpace = originalSpace + die
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
        if player2 == 0:
            self.player2 = AIPlayer("b")   
        self.board = Board(initialState)
        self.gid = id
    
    def rollDice(self):
        die1 = random.randint(1,6)
        die2 = random.randint(1,6)
        return [die1, die2]
    
    def getPossibleMoves(self, player, dice, board):
        currState = board.currentState
        possibleMoves = []
        i = 0
        if player.color == "w":
            print("white to move")
            for i in range(24):
                i+=1
                if "w" in currState[str(i)]:
                    for die in dice:
                        if 24-i>die:
                            goalPlace = currState[str(i+die)]
                            if not("b" in goalPlace and len(goalPlace) >= 2):
                                possibleMoves.append((i, die))


                        # if 24-i>die2:
                        #     goalPlace = currState[str(i+die2)]
                        #     if not("b" in goalPlace and len(goalPlace) >= 2):
                        #         possibleMoves.append((i, die2))
                        #         print("adding a possible move: ", (i, die2))
        elif player.color == "b":
            print("black to move")
            for i in range(24):
                i+=1
                if "b" in currState[str(i)]:
                    for die in dice:
                        if i>= die:
                            goalPlace = currState[str(i-die)]
                            if not("w" in goalPlace and len(goalPlace) >= 2):
                                possibleMoves.append((i, -die))
                        # if i>die2:
                        #     goalPlace = currState[str(i-die2)]
                        #     if not("w" in goalPlace and len(goalPlace) >= 2):
                        #         possibleMoves.append((i, -die2))
                        #         print("adding a possible move: ", (i, -die2))
        return possibleMoves
    
    def move(self, player):
        #gets possible moves, gets move from player
        #moves = []
        dice = self.rollDice() #diceroll might have to be an endpoint instead of being called from move!
        while len(dice) != 0:
            print("dice: ", dice)
            possibleMoves = self.getPossibleMoves(player, dice, self.board) #figure out how to pass on the turn when no moves are possible
            print("possible moves: ", possibleMoves)
            if len(possibleMoves) == 0:
                print("no possible move. Passing on the turn to ", player.color)
                return currState
            move = player.getMove(possibleMoves)
            print("move: ", move)
            if player.color == "b":
                die = -move[1]
            elif player.color == "w":
                die = move[1]
            currState = self.board.doMove(player.color, move)
            try: 
                dice.remove(die)
            except:
                print("tried to remove {} from {}".format(die, dice))
            #call doMove here, so player can see move before doing the second move
        #return game state rather than the actual move
        #fix the endpoint to reflect changes
        return currState
    
    

class Player():
    '''send move to game class'''
    def __init__(self):
        #what needs to be common for the two?
        #some move function, but they need to be implemented differently.
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

    def getMove(self, possibleMoves):
        if self.color == "b":
            return (possibleMoves[-1]) #minus because black moves backwards, white moves forwards
        elif self.color == "w":
            return (possibleMoves[0])
        else:
            return("Error. Not a valid color. Must be b or w")
        #more complicated logic should come here
        
