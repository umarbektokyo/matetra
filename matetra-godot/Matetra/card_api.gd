extends Node

# Reference to the NetworkManager
@onready var network = $"/root/NetworkManager" 

# Mapping CSV 'InputsReq' characters to descriptions for the UI
const INPUT_DESCRIPTIONS = {
	"d": "Roll Dice (1-6)",
	"p": "Select Player Index",
	"n": "Select Number Slot (0-4)",
	"A": "Select Attacker (Player ID)",
	"U": "Select User (Player ID)",
	"c": "Select Card ID" 
}

# --- MAIN FUNCTION TO CALL FROM UI ---
# card_index: The index of the card in the player's hand/queue
# req_string: The 'InputsReq' string from the CSV (e.g., "AnUn")
# user_selections: An array of data gathered from UI interactions
func execute_card_action(card_index: int, req_string: String, user_selections: Array):
	
	var formatted_inputs: Array[int] = []
	
	# Validate that the UI gave us enough data
	if user_selections.size() != req_string.length():
		push_error("Mismatch: Card expects %d inputs, got %d" % [req_string.length(), user_selections.size()])
		return

	# Iterate through the requirement string to format data for Go
	for i in range(req_string.length()):
		var req_char = req_string[i]
		var raw_val = user_selections[i]
		
		match req_char:
			"d": 
				# Dice roll. Ensure it is an integer 1-6.
				formatted_inputs.append(int(raw_val))
			
			"p", "A", "U":
				# Player ID.
				# Note: In functions.go, A (Attacker) and U (User) are treated as Player Indices.
				formatted_inputs.append(int(raw_val))
			
			"n":
				# Number Slot Index (usually 0-4).
				formatted_inputs.append(int(raw_val))
			
			_:
				# Default fallback for simple integers
				formatted_inputs.append(int(raw_val))

	# Construct the final payload for engine.go -> RecordMove
	var payload = {
		"CardIndex": card_index,
		"Inputs": formatted_inputs
	}
	
	print("Sending Move: ", payload)
	network.send_move(payload)

# --- HELPER FOR UI GENERATION ---
# Call this to know what UI elements to show the player
func get_requirements_list(req_string: String) -> Array:
	var requirements = []
	for char in req_string:
		requirements.append(INPUT_DESCRIPTIONS.get(char, "Unknown Input"))
	return requirements
