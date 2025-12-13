extends Area2D
#
#@onready var collision_shape = $CollisionShape2D
#
#var original_shape_size: Vector2
#var original_shape_position: Vector2
#
## Keep track of currently overlapping areas

#var overlapping_areas: Array = []
#
#func _ready():
	#if collision_shape.shape is RectangleShape2D:
		#original_shape_size = collision_shape.shape.size
		#original_shape_position = collision_shape.position
	#else:
		#push_error("This only works with RectangleShape2D")
#
#func _on_area_entered(body):
	#overlapping_areas.append(body)
	#update_collision_shape()
#
#func _on_area_exited(body):
	#return
#
#func update_collision_shape():
	#var left_limit = -original_shape_size.x / 2
	#var right_limit = original_shape_size.x / 2
#
	#for other_area in overlapping_areas:
		#if not (other_area.collision_shape.shape is RectangleShape2D):
			#continue
#
		## Global positions
		#var my_left = global_position.x + left_limit
		#var my_right = global_position.x + right_limit
		#var other_left = other_area.global_position.x - other_area.collision_shape.shape.size.x / 2
		#var other_right = other_area.global_position.x + other_area.collision_shape.shape.size.x / 2
#
		## Compute overlap
		#var overlap_left = max(my_left, other_left)
		#var overlap_right = min(my_right, other_right)
#
		#if overlap_left < overlap_right:
			#if other_area.global_position.x < global_position.x:
				## Other card is to the left -> shrink left side
				#left_limit = max(left_limit, overlap_right - global_position.x)
			#else:
				## Other card is to the right -> shrink right side
				#right_limit = min(right_limit, overlap_left - global_position.x)
#
	## Apply new collision shape size and position
	#var new_size_x = max(right_limit - left_limit, 1) # prevent zero or negative
	#collision_shape.shape.size.x = new_size_x
	#collision_shape.position.x = (left_limit + right_limit) / 2
