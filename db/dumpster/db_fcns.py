# import psycopg2
# import psqlConfig_master

# class BackgammonDB:

#     def __init__(self):
#         try:
#             self.conn = psycopg2.connect(database = psqlConfig_master.database, user = psqlConfig_master.username, password = psqlConfig_master.password, host = "localhost")
#             print("established database connection")
#         except(Exception, psycopg2.Error) as error:
#             print("Error connecting to PostgreSQL,", error)

#     def cursor_init(self):
#         try:
#             curr = self.conn.cursor()
#             return curr
#         except(Exception, psycopg2.Error) as error:
#             print("Error connecting to PostgreSQL", error)
#             return(error)
        
#     def register(self, username, password):
#         curr = self.cursor_init()
#         try:
#             curr.execute("INSERT INTO Users (username, password) VALUES (%s, %s);", (username, password))
#             self.conn.commit()
#             count = curr.rowcount
#             print(count, "new user successfully added to table")
#             curr.execute("INSERT INTO Userstats (username, gamesPlayed, wins, losses) VALUES (%s, %s, %s, %s);", (username, '0', '0', '0'))
#             self.conn.commit()
#             print("row successfully opened in userstats")
#             return("200")
#         except(Exception, psycopg2.Error) as error:
#             print("Error inserting new user", error)
#             return str(error)
        
#     def newGame(self, black, white, boardState):
#         curr = self.cursor_init()
#         try:
#             curr.execute("INSERT INTO Games (white, black, status, boardstate) VALUES (%s, %s, %s, %s) RETURNING gameId;", (white, black, 'new', boardState))
#             game_id = curr.fetchone()[0]
#             self.conn.commit()
#             print(f"Game with ID {game_id} successfully added to the table.")
#             return str(game_id)
#         except(Exception, psycopg2.Error) as error:
#             print("Error inserting new game", error)
#             return str(error)
        
#     def updateStats(self, username, result):
#         curr = self.cursor_init()
#         try:
#             curr.execute("UPDATE Userstats SET gamesPlayed = gamesPlayed + 1 WHERE username = %s; ", (username,))
#             self.conn.commit()
#             print("updated gamesplayed")
#         except(Exception, psycopg2.Error) as error:
#             print("Error updating stats", error)

#         try:
#             if result == "win":
#                 curr.execute("UPDATE Userstats SET wins = wins + 1 WHERE username = %s; ", (username,))
#             else:
#                 curr.execute("UPDATE Userstats SET losses = losses + 1 WHERE username = %s; ", (username,))
#             self.conn.commit()
#             print("updated win/loss")
#             return("200")
#         except(Exception, psycopg2.Error) as error:
#             print("Error updating stats", error)
#             return str(error)

# if __name__ == "__main__":
#     bestsellers = BackgammonDB()
#     bestsellers.conn.close()