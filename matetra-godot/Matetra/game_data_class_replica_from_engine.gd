extends Node
class_name GameData

# --- ENUMS (CONSTANTS) ---
# Mirrors internal logic for readable code
const INPUT_DICE = "d"      # Requires a dice roll (1-6)
const INPUT_NUMBER = "n"    # Requires a Number Slot Index (0-4)
const INPUT_PLAYER = "p"    # Requires a Player ID (Generic)
const INPUT_ATTACKER = "A"  # Requires a Player ID (Target/Opponent)
const INPUT_USER = "U"      # Requires a Player ID (Self)
const INPUT_CARD = "c"      # Requires a Card ID (Not fully implemented yet)

# ==========================================
# 1. DATA MODELS (Mirrors model.go)
# ==========================================

# Represents a single number slot on the board
class NumberSlot:
	var value: int = 0
	var mark: String = "n" # 'n' = normal, 'F' = Fibonacci
	
	func _init(data = {}):
		if data.is_empty(): return
		# Go sends "Value" (Capitalized)
		value = int(data.get("Value", 0)) 
		mark = data.get("Mark", "n")

# Represents a Card
class Card:
	var name: String = ""
	var description: String = ""
	var type: String = ""
	var method: String = ""
	var owner: int = -1
	var inputs: Array = []
	var inputs_req: String = ""
	
	func _init(data = {}):
		if data.is_empty(): return
		# Go sends Capitalized keys for GameState structs
		name = data.get("Name", "Unknown")
		description = data.get("Description", "")
		type = data.get("Type", "Function")
		method = data.get("Method", "")
		owner = int(data.get("Owner", -1))
		inputs = data.get("Inputs", [])
		inputs_req = data.get("InputsReq", "")

# Represents a Player
class Player:
	var name: String = ""
	var hash_id: String = "" # The password hash/ID
	
	func _init(data = {}):
		if data.is_empty(): return
		name = data.get("Name", "Unknown")
		hash_id = data.get("Hash", "")

# Represents the Full Game State
class GameState:
	var game_id: String = ""
	var turn: int = 0
	var done_status: Array = [false, false] # [Player1, Player2]
	var queue
	var players: Array[Player] = []
	var cards: Array[Card] = []
	var numbers: Array = [] # Array of Arrays of NumberSlot
	
	# This function acts as the "Decoder" from the raw JSON payload
	func from_dictionary(payload: Dictionary):
		game_id = payload.get("GameID", "")
		turn = int(payload.get("Turn", 0))
		done_status = payload.get("Done", [false, false])
		queue = payload.get("Queue", [])
		
		# Parse Players
		players.clear()
		for p_data in payload.get("Players", []):
			players.append(Player.new(p_data))
			
		# Parse Cards
		cards.clear()
		for c_data in payload.get("Cards", []):
			cards.append(Card.new(c_data))
			
		# Parse Numbers (2D Array)
		numbers.clear()
		var raw_nums = payload.get("Numbers", [])
		for player_row in raw_nums: # Iterate over players
			var row_slots: Array[NumberSlot] = []
			for n_data in player_row: # Iterate over slots
				row_slots.append(NumberSlot.new(n_data))
			numbers.append(row_slots)

# ==========================================
# 2. REQUEST BUILDERS (Mirrors api.go Payloads)
# ==========================================
# Use these static functions to create requests. 
# This ensures you never misspell a key like "card_index".

static func req_register(username: String, password: String) -> Dictionary:
	# Matches PlayerPayload in api.go
	return {
		"name": username,
		"hash": password
	}

static func req_play_card(index: int, inputs: Array, permanent: bool) -> Dictionary:
	# Matches CardPlayPayload in api.go
	return {
		"card_index": index,
		"inputs": inputs,
		"permanent": permanent
	}

static func req_next_turn() -> Dictionary:
	
	return {
		"done": true
	}
