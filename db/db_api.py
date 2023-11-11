import psycopg2
import flask

conn = psycopg2.connect(database="userstats",
                        user="statmaker",
                        host = 'localhost'
                        port = 5432
                        password="secure")

cur = conn.cursor()

cur.execute('SELECT * from userstats')
table = cur.fetchall()
username = table[0][0]
print(username)

username = 


