# Matetra-go
This is a digital implementation of matetra card game. It might be confusing, buggy, or even fully broken, but it will forsure be fun!

Idea: create two services: matetra-api, which will be responsible for handling the logic, storage, and backend of the game in general and matetra-tui which will handle the user interactions in the terminal. We could build a godot game for it as well but those are worries for later.

## matetra-api
.matetra extension when game is paused and  (JSON but named .matetra because it's cool?)
Data stored:
list {
  matetra_version: int (idk, just to ensure no corrupted files in the future)
  cards: list[cardID, display_string, operation_code, owner] (list of all cards, probably loaded from a database at the beginning of each game.)
  turns: playerID (changes to the next increment of the player)
  players: int
  players: list {
    playerID
    order: int
    name: string
    password: string
    set: list[idk what datatype, it must be adaptive to a super-small or super-big numbers. Would be useful to store to 4 significant digits but together with scientific notation?] (every player can only have 5 numbers in their set)
    actions: list[cardID, set(number's index which ranges from 0 to 4)]
  }
}
Idk if I'm using a correct way to store data, it seems really arbitary... idk...

Code-wise, here is how I think the game will go:
Repat:
- If it's your turn (defender):
  - You can roll a dice which will fill out empty number slots (you can't roll once you don't have enough slots)
  - You can use your cards and set your numbers as inputs
- If its the opponent's turn (attacker):
  - You can use your car on the chosen number of the defender. You must use your numbers as inputs for everything but 1 attacked number as shown in the card.
- All the applied cards are stored in an "actions" field, where they will be applied in order at the end of the term. (There are special cards which can cancel the previous card, change places for the numbers and cards, dublicate all cards from one number and so on.)
- Every player has a "done" button. Once they've clicked it they can unclick it; however, if everyone's button is clicked, the turn has ended.
- Apply all the cards' code on the numbers.
- Update to the next player's turn
- Fill everybody's hands with 6 cards

We should make so it can be played with multiple people around the world... idk how to do it. Is it a good idea to make so a person can host an http server and others can conenct to it, making so every player has an access to the contents of the same game? We can protect commands with a password which users have to input that then gets send each time to the api. I'm not worried about cybersecurity, it's just a proof of concept to play around with friends for now. 