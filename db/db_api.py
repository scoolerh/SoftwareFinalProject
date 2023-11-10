import psycopg2

conn = psycopg2.connect(database="userstats",
                        user="statmaker",
                        password="secure")

cur = conn.cursor()

cur.execute()