SELECT 'CREATE DATABASE Backgammon' 
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'Backgammon')\gexec

\c Backgammon;

-- all of this does not run when the docker container starts. I had to do it manually.
-- might be a volumes issue, but I don't know
-- UPDATE: volumes is definitely not working

CREATE TABLE IF NOT EXISTS users(
    username varchar(50) PRIMARY KEY,
    password varchar(50)
);

CREATE TABLE IF NOT EXISTS userstats( 
    username varchar(50),
    FOREIGN KEY (username) REFERENCES Users(username),
    gamesPlayed int,
    wins int,
    losses int
);

CREATE TABLE IF NOT EXISTS games(
    gameId SERIAL PRIMARY KEY,
    status varchar(15), --"new", "open"(?), "paused", "finished"
    white varchar(50),
    black varchar(50),
    FOREIGN KEY (white) REFERENCES Users(username),
    FOREIGN KEY (black) REFERENCES Users(username),
    boardState varchar(1024),
    turn varchar(4),
    winner varchar(4)
);


CREATE ROLE readaccess;
GRANT CONNECT ON DATABASE Backgammon to readaccess;
GRANT USAGE ON SCHEMA PUBLIC TO readaccess;
GRANT SELECT ON userstats TO readaccess;

CREATE USER statmaker WITH PASSWORD 'secure';
GRANT readaccess to statmaker;