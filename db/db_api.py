import psycopg2
from flask import Flask, jsonify, render_template

def sql_conversion():
    app = Flask(__name__)

    #Establishes connection to a SQL database using the Psycopg2
    conn = psycopg2.connect(dbname="backgammon", 
                            user="postgres", 
                            host = 'db', 
                            port = 5432, 
                            password="collective")

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
    #output = jsonify(postDict)
    output = postDict

    conn.close()
    cur.close()

    return output


if __name__ == "__main__":
    sql_conversion()



