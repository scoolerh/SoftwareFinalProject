import psycopg2
from flask import Flask, jsonify, render_template

def sql_conversion():
    app = Flask(__name__)

    #Establishes connection to a SQL database using the Psycopg2
    conn = psycopg2.connect(database="Backgammon", 
                            user="statmaker", 
                            host = 'localhost', 
                            port = 5432, 
                            password="secure")

    #Creates a cursor 
    cur = conn.cursor()
    cur.execute('SELECT * from userstats')
    table = cur.fetchall()

    #'Player': [Username, "Games Played", "Wins", "Losses"]
    postDict = {}

    for row in table:

        player = row["username"]
        gamesPlayed = row["gamesPlayed"]
        gamesWon = row["wins"]
        gamesLost = row["losses"]
        
        postDict.update(
            {
                player: [player, gamesPlayed, gamesWon, gamesLost]
            }
        ) 

    print(postDict)
    jsonify(postDict)

    conn.close()

    return  render_template('scoreboard.html', table_data = postDict)


if __name__ == "__main__":
    sql_conversion()



