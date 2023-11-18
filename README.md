# CS 347 Final Project: The Collective

## Contributors 

Sophie Quinn, Hannah Scooler, Selma Vangstein, Bryan Yang 

## Description of our project
This an application built using go, postgreSQL, HTML, CSS and javascript.
The app is built with a hope of giving the user a smooth and joyful Backgammon experience. We have the following functions to ensure this:
- On the home page we have the rules of Backgammon listed in case we have a user that is new to the game, or one who wants to fresh up on the rules before playing
- We give the option to register or log in. These forms should make sure that the user chooses a unique username, does not log in as the users that are supposed to be available (the two AIs and 'guest'), and checks for the correct password. If the user wants to play without having their stats stored - no worries! Just start a new game without logging in yet - we will get to this flow shortly
- To increase competitiveness wins and losses, we store wins and losses to each user, giving them the option to view their personal stats and the competitive leaderboard (if the other group did their job).
- The last button on the home page is the most important one - New Game. Click this button when you want to start playing! The most important functionality here before moving on to the actual game is about who actually gets to play. You get the option to play a guest player, play a friend, play the AI (with two different difficulties), or watch the AI play itself. If you want to play yourself, just choose 'guest' for one of the users, and yourself as the others! This is a great way to make your stats look better.
- If you want to play another user (non-guest), make sure they are already registered. This is to ensure that the newgame-experience is smooth to use. Also, when letting this player log in: we do not accept a sloppy login in the newgame-flow. If you fail to successfully log in, you are just assigned as a guest. Keep track of your usernames and passwords, it is an important lesson.
- About the AI: We have two AIs: easy (nicknamed 'steve') and hard (nicknamed 'joe'). We do not claim that any of them are difficult to play, we are only saying that one thinks slightly more than the other. Steve is an impatient soul, and always chooses the first move he sees. Joe, on the other hand, takes into consideration how many towers and blots he ends up with, and how many pips will be left for each player after the move - effectively taking into regard whether a piece is captured. We do not claim that this is an optimal or even a good strategy for playing the game, but it is a strategy.
- All common rules of backgammon should be implemented, with capturing and with bearing off only when pieces are in the home board, and the more complicated rules that follows bearing off. We list the possible moves for the user to choose from, so you can do each move with only a click of a button - this felt more convenient than the less efficient solutions out there.
- We do not have a doubling die as we do not condone gambling. We also dont have a pause/resume button because we always encourage people to finish what they started - this is an important value of our team. 
- When a player is able to bear all their pieces off the board, the game is declared finished, and the players are taken to the "win" page, where the winner is clearly stated, and we present the options to play again or return to the home page
- We have a database running in the background. We have a user table that handles registering and logins, it stores all usernames and passwords. This table is hidden from the group that only works with stats. A second table stores user stats: number of games, wins, and losses for each of the users, including the AIs. The third table store information about each game: the last saved gamestate, the current status of the game (for instance 'finished'), the white and black player, and the winner (if there is one). Currently we only update this when the game is created and when it is finished, however this can be changed in the future if there is a wish for that.

## Project File structure
├── backgammon _(contains code associated with frontend, frontend API, and game logic)_  
|   ├── app  
|   |   ├── html _(directory containing all html files/templates)_  
|   |   ├── api_finctions.go _(contains helper functions for api.go)_  
|   |   └── api.go _(api to interface with a frontend)_  
|   └── game  
|       ├── ai.go _(logic for AI players)_  
|       └── gamelogic _(logic for rules and updating db)_  
├── db  
|   ├── db_setup.sql _(sets up db tables)_  
|   ├── psqlConfig_master.py _(db credentials with read/write access)_  
|   └── psqlConfig_readaccess.py _(db credentials with read only access)_  
├── pgdata _(directory containing the postgresql db)_  
├── compose.yaml _(Compose file to run all Dockerfiles. Sets up envionments/ports)_  
└── README.md  


## How to Use 
Our app should hopefully be easy to use. Follow these steps and you should be in for a joyful experience!
(numbers here will be fixed when I am sure we are done with adding things)
1. Clone the repo and make sure you are on main
2. Download Docker if you don't have it, and make sure you have the docker engine running
3. Do docker compose build and docker compose up
4. Wait until you see "successfully connected to database" before you try to load the page
4. Go to localhost:9000 to get to our home page
5. Here, you can choose to freshen up on the game rules, register a new account, or go straight to new game to play as a guest
6. On this page, you choose who will play the game. You can play another guest (or just keep the computer to yourself and play alone), you can let a friend or enemy log in, you can play an AI, or you can watch the two AI's play each other.
7. After submitting your choices, click start to begin your joyful Backgammon experience!
(. roll to decide who starts (?) . The move-buttons show you what moves are possible with your dice roll. Chose one of them to make your first move! When you have used all the dice, the turn is passed to the opponent. Bear off all your pieces before your opponent to win the game!)

## A note on cheating
Please don't. 
We do not condone such actions. However, if you need a confidence boost, something to brighten your day, there are multiple ways you can do that by manipulating the url of this site, for instance by fabricating moves, or declaring yourself the winner. It is also technically possible to use the back-button to undo a move, but this is also considered cheating, so don't do it. The consequences of doing this is not thoroughly tested.
