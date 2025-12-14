extends Node

func _ready():
	# 1. Listen for Server Responses
	MatetraAPI.game_state_received.connect(_on_game_state)
	MatetraAPI.move_reply_received.connect(_on_move_reply)
	
	print("--- API TESTER READY ---")
	print("Press 'J' to JOIN")
	print("Press 'T' to TEST a move (Preview)")
	print("Press 'Enter' to COMMIT a move (Permanent)")
	print("Press 'Space' to END TURN")

func _input(event):
	if not event.is_pressed():
		return

	if event.is_action_pressed("JoinGame"):
		print(">> Sending JOIN request...")
		MatetraAPI.join_game("TestPlayer", "secret_password")

	elif event.is_action_pressed("TestCardPreview"):
		print(">> Sending PREVIEW...")
		MatetraAPI.play_card(0, [1, 0, 0, 0], false)

	elif event.is_action_pressed("TestCardPermanent"):
		print(">> Sending COMMIT...")
		MatetraAPI.play_card(0, [1, 0, 0, 0], true)

	elif event.is_action_pressed("TestEndTurn"):
		print(">> Sending END TURN request...")
		MatetraAPI.finish_turn()

# --- CALLBACKS ---

func _on_game_state(state: GameData.GameState):
	print("[Client] Board Updated! Current Turn: %d" % state.turn)
	
	# Verify we have cards
	if state.cards.size() > 0:
		var my_first_card = state.cards[0]
		# Only print if I own it (Owner 0 or 1, not Deck -1)
		if my_first_card.owner > -1:
			print("[Client] I have a card: %s (Req: %s)" % [my_first_card.name, my_first_card.inputs_req])

func _on_move_reply(success, message):
	if success:
		print("[Client] Success: %s" % message)
	else:
		push_error("[Client] Failed: %s" % message)
