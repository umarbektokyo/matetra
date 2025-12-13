extends Node

var card_index = 5 # Example index in hand
var inputs_req = "AnUn" # From CSV for "Addition"

var ui_selection_data = [
selected_attacker_player_id, 
selected_attacker_slot, 
selected_user_player_id, 
selected_user_slot
]

# These values would come from the user clicking things in your scene
var selected_attacker_player_id = 2 # Opponent
var selected_attacker_slot = 3
var selected_user_player_id = 0     # Self
var selected_user_slot = 1
