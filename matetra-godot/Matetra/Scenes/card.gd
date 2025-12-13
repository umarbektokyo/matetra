extends Sprite2D

@export var hover_offset: float = -100   
@export var hover_speed: float = 8.0

var original_position: Vector2
var target_position: Vector2
@onready var area = $Area2D

func _ready():
	original_position = position
	target_position = position

	area.mouse_entered.connect(Callable(self, "_on_mouse_entered"))
	area.mouse_exited.connect(Callable(self, "_on_mouse_exited"))
	area.area_entered.connect(Callable(self, "_on_area_entered"))
	area.area_exited.connect(Callable(self, "_on_area_exited"))

func _process(delta):
	position = position.lerp(target_position, hover_speed * delta)

func _on_mouse_entered():
	Global.current_cards_being_held.push_back(self)
	Global.check_cards_count()
	target_position = original_position + Vector2(0, hover_offset)
	$HoverSFX.play()

func _on_mouse_exited():
	target_position = original_position
	
