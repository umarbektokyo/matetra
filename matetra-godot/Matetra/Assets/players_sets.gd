extends Control

@onready var Player1Label = $PanelContainer/VBoxContainer/Label
@onready var Player2Label = $PanelContainer/VBoxContainer/Label2
@onready var Player3Label = $PanelContainer/VBoxContainer/Label3
@onready var Player4Label = $PanelContainer/VBoxContainer/Label4

func _process(delta):
	if Input.is_action_pressed("operation"):
		# Send to API
		CardAPI.execute_card_action(Global.card_index, Global.inputs_req, Global.ui_selection_data)
