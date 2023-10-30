import psycopg2
import psqlConfig_web

class BackgammonDB:

    def __init__(self):
        try:
            self.conn = psycopg2.connect(database = psqlConfig_web.database, user = psqlConfig_web.username, password = psqlConfig_web.password, host = "localhost")
            print("established database connection")
        except(Exception, psycopg2.Error) as error:
            print("Error connecting to PostgreSQL,", error)

    def cursor_init(self):
        try:
            curr = self.conn.cursor()
            return curr
        except(Exception, psycopg2.Error) as error:
            print("Error connecting to PostgreSQL", error)
            return(error)
        
    def register(self, username, password):
        curr = self.cursor_init()
        try:
            curr.execute("INSERT INTO Users (username, password) VALUES (%s, %s);", (username, password))
            self.conn.commit()
            count = curr.rowcount
            print(count, "new user successfully added to table")
            return("200")
        except(Exception, psycopg2.Error) as error:
            print("Error inserting new user", error)
            return(error)
        
    def newGame(self, black, white, boardState):
        curr = self.cursor_init()
        try:
            curr.execute("INSERT INTO Games (white, black, status, boardstate) VALUES (%s, %s, %s, %s) RETURNING gameId;", (white, black, 'new', boardState))
            game_id = curr.fetchone()[0]
            self.conn.commit()
            print(f"Game with ID {game_id} successfully added to the table.")
            return str(game_id)
        except(Exception, psycopg2.Error) as error:
            print("Error inserting new game", error)
            return str(error)
        
    def userstats(self, username):
        curr = self.cursor_init()
        try:
            #make sure this is working
            curr.execute("SELECT username=%s FROM TABLE Users", username)
        except(Exception, psycopg2.Error) as error:
            print("Error getting user stats", error)
            return str(error)

    def highscores():
        #get all usernames, sorted by #wins
        pass


if __name__ == "__main__":
    bestsellers = BackgammonDB()
    bestsellers.conn.close()