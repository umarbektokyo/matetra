extends Button

func _on_pressed() -> void:
	pass # Replace with function body.
	var serverURLinput = LineEdit.new()
	serverURLinput.connect("text_changed", Callable(self, "user_input"))
	get_parent().add_child(serverURLinput)

func user_input(server_url_entered: String) -> void:
	MatetraAPI.server_url = server_url_entered
	MatetraAPI.connect_to_server()
	MatetraAPI.join_game("TestPlayer", "secret_password")
