extends Node

# --- SIGNALS ---
# Emitted when the board updates (Numbers, Cards, Players)
signal game_state_received(game_state) 
# Emitted when the server replies to a specific move (Success/Fail)
signal move_reply_received(success, message)
# Emitted on connection loss
signal disconnected()

# --- CONFIGURATION ---
var socket = WebSocketPeer.new()
var server_url = "ws://100.107.2.27:1729/ws" # Match port in utils.go
var _is_connected = false
var current_game_state: GameData.GameState = GameData.GameState.new()

func _ready():
	connect_to_server()
	pass

func connect_to_server():
	print("Attempting to connect to Matetra Server...")
	socket.connect_to_url(server_url)

func _process(_delta):
	socket.poll()
	var state = socket.get_ready_state()
	
	if state == WebSocketPeer.STATE_OPEN:
		if not _is_connected:
			_is_connected = true
			print("Connected!")
			
		while socket.get_available_packet_count():
			var packet_str = socket.get_packet().get_string_from_utf8()
			_handle_incoming_message(packet_str)
			
	elif state == WebSocketPeer.STATE_CLOSED:
		if _is_connected:
			_is_connected = false
			print("Disconnected from server.")
			emit_signal("disconnected")

# --- SENDING COMMANDS (Client -> Server) ---

# 1. Join the Game
func join_game(player_name: String, password: String):
	var payload = GameData.req_register(player_name, password)
	# CHANGED: "REGISTER" -> "ADD_PLAYER" (Matches api.go line 109)
	_send_json("ADD_PLAYER", payload)

# 2. Play a Card (or Preview it)
func play_card(card_index: int, inputs: Array, is_permanent: bool = true):
	var payload = GameData.req_play_card(card_index, inputs, is_permanent)
	# MATCHES: "PLAY_CARD" is correct (Matches api.go line 130)
	_send_json("PLAY_CARD", payload)

# 3. End Turn
func finish_turn():
	var payload = GameData.req_next_turn()
	# CHANGED: "NEXT_TURN" -> "PROCESS_NEXT_TURN" (Matches api.go line 132)
	_send_json("PROCESS_NEXT_TURN", payload)

# --- RECEIVING MESSAGES (Server -> Client) ---

func _handle_incoming_message(json_str: String):
	var json = JSON.new()
	var error = json.parse(json_str)
	if error != OK:
		print("JSON Parse Error: ", json.get_error_message())
		return
	var msg = json.get_data()
	# Keys are lowercase based on 'json:"type"' in api.go
	var type = msg.get("type", "")
	var payload = msg.get("payload", {})

	match type:
		"STATE_UPDATE":
			# Full game sync
			current_game_state.from_dictionary(payload)
			print(str(current_game_state))
			emit_signal("game_state_received", current_game_state)
			
		"PLAY_CARD_REPLY":
			# Direct response to an action (Success/Fail)
			var success = payload.get("success", false)
			var message = payload.get("message", "")
			# "newGameState" is optional in the Go struct (omitempty), so handle null
			var new_state_dict = payload.get("newGameState", null)
			
			print("Move Reply: ", message)
			emit_signal("move_reply_received", success, message)
			
			# If the server sent a new state (it does for previews and valid moves)
			if new_state_dict:
				current_game_state.from_dictionary(new_state_dict)
				print(str(current_game_state))
				emit_signal("game_state_received", current_game_state)
		
		"PLAYER_ADDED":
			# The server sends this specific confirmation (line 126 in api.go)
			# We can just print it, as the STATE_UPDATE usually follows immediately.
			print("Player Added Confirmation: ", payload.get("name"))

		"ERROR":
			print("Server Error: ", payload)

# --- HELPER ---
func _send_json(type: String, payload_data):
	if socket.get_ready_state() != WebSocketPeer.STATE_OPEN:
		print("Error: Socket not open")
		return
		
	var msg = {
		"type": type,
		"payload": payload_data
	}
	socket.send_text(JSON.stringify(msg))
