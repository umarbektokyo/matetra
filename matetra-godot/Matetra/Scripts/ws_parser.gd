extends Node

# An attempt to get any sort of feedback V2 (v1 was not committed since it, simply put, sucked)
var socket = MatetraAPI.socket
var ws_url: String = MatetraAPI.server_url

"""func _process(_delta):
	socket.poll()
	var state = socket.get_ready_state()
	
	if state == WebSocketPeer.STATE_OPEN:
		_while_packages()

func _while_packages():
	print("getting")
	var packet_str = socket.get_packet().get_string_from_utf8()
	_retrieve_info(packet_str)

func _retrieve_info(info): 
	print("WS_PARSER: {info}".format([info], "{info}"))"""

# game_state_receive just gives back the node that is receiving the state, e.g.:
# Node(matetra_godot_api.gd)::[signal]game_state_received
