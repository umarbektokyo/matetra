extends Button

var server_pid = -1

func _on_pressed():
	# 1. Locate the Go binary (assuming it's in the same folder as the game)
	var exe_path = OS.get_executable_path().get_base_dir() + "/matetra-server"
	if OS.get_name() == "Windows":
		exe_path += ".exe"
	
	# 2. Define the arguments (just like you typed in the terminal)
	# equivalent to: matetra-server start "MyGodotLobby"
	var args = ["start", "MyGodotLobby"]
	
	# 3. Run the process in the background
	# The -1 argument keeps it independent of the main thread
	server_pid = OS.create_process(exe_path, args)
	
	if server_pid != -1:
		print("Server started successfully! PID: ", server_pid)
		# Now automatically trigger the "Join" logic for the host
		MatetraAPI.server_url = "ws://100.97.115.40:1729/ws"
		MatetraAPI.connect_to_server()
	else:
		printerr("Failed to start server. Is the executable in the right folder?")
		
func _process(delta):
	if Input.is_action_just_pressed("KillServer"):
		if server_pid != -1:
			OS.kill(server_pid)
			print ("Server Killed!")
		else:
			printerr ("Server not killed! Is server running? Check if port is open on Windows with netstat -ano | findstr :1729")
# Clean up: Kill the server when Godot closes
func _exit_tree():
	if server_pid != -1:
		OS.kill(server_pid)
