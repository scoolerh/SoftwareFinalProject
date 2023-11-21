import psycopg2
from flask import Flask, jsonify, render_template
import json

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

    columns = [desc[0] for desc in cur.description]

    data = []
    for row in table:
        row_data = {}
        for idx, column in enumerate(columns):
            row_data[column] = str(row[idx])
        
        data.append(row_data)

    conn.close()
    cur.close()

    with app.app_context():
        print(data)
        output = json.dumps(data)
        print(output)
        return output


if __name__ == "__main__":
    data = sql_conversion()



