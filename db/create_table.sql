CREATE TABLE IF NOT EXISTS Users(
    username varchar(50) PRIMARY KEY,
    password varchar(50),
    gamesPlayed int,
    wins int,
    losses int
);

CREATE TABLE IF NOT EXISTS Games(
    gameId SERIAL PRIMARY KEY,
    status varchar(15), --"new", "open"(?), "paused", "finished"
    white varchar(50),
    black varchar(50),
    boardState varchar(1024),
    turn varchar(4),
    winner varchar(4)
);