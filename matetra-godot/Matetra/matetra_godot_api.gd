extends Node

# Signal to notify other parts of the game when state updates
signal game_state_received(state_data)
signal connection_established()
signal connection_closed()

var socket = WebSocketPeer.new()
var websocket_url = "ws://localhost:1729/ws" # Default port from utils.go

func _ready():
	connect_to_server()

func connect_to_server():
	print("Connecting to Matetra Server...")
	var err = socket.connect_to_url(websocket_url)
	if err != OK:
		print("Unable to connect")
		set_process(false)

func _process(_delta):
	socket.poll()
	var state = socket.get_ready_state()
	
	if state == WebSocketPeer.STATE_OPEN:
		while socket.get_available_packet_count():
			var packet = socket.get_packet()
			var msg_string = packet.get_string_from_utf8()
			_handle_message(msg_string)
			
	elif state == WebSocketPeer.STATE_CLOSED:
		var code = socket.get_close_code()
		print("WebSocket closed with code: %d" % code)
		emit_signal("connection_closed")
		set_process(false)

func _handle_message(json_str: String):
	var json = JSON.new()
	var parse_result = json.parse(json_str)
	
	if parse_result != OK:
		print("JSON Parse Error")
		return

	var msg = json.get_data()
	
	# Matches "type" from main.go
	match msg.get("Type", ""):
		"STATE_UPDATE":
			emit_signal("game_state_received", msg["Payload"])
		"ERROR":
			print("Server Error: ", msg["Payload"])

# Public function to send moves
func send_move(payload: Dictionary):
	if socket.get_ready_state() != WebSocketPeer.STATE_OPEN:
		print("Socket not ready")
		return
		
	# The engine expects a message wrapper. 
	# Based on typical Go websocket patterns, we wrap the payload.
	var packet = {
		"Type": "MOVE", # Assuming generic "MOVE" type, or "APPLY_CARD"
		"Payload": payload
	}
	socket.send_text(JSON.stringify(packet))
