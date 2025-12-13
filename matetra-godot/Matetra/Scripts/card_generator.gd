extends Node2D

# Hard-coded title values list:
var card_title = ["Addition", "Subtraction", "Multiplication", "Division", "Absolute Value", "Inverse", 
"Negative", "Positive", "Square Root", "Factorial", "Square", "Cosine", "Base-10 Logarithm",
"Exponential", "Natural Logarithm", "Sine", "Tangent", "Logarithm", "Base-Root", "Power", 
"Summation", "Product", "Second Order Polynomial", "First Order Polynomial", "Identity Element", 
"Closure Element", "Distributive Element", "Commutative Element", "Inverse Element", 
"Identity Property", "Closure Property", "Distributive Property", "Commutative Property",
"Inverse Property", "Pascal's Triangle", "Pythagorean Theorem", 
"Fundamental Theorem of Arithmetic", "Euler's Number", "Negative", "Sheldon's Number", "Googol",
"The Answer", "Phi", "Zero", "Pi", "Lucky Number", "2nd Perfect Number", "1st Perfect Number", 
"Fibonacci Number", "Symmetrical Number", "Tau", "Scientific Notation", "Graham's number"]

# Title line variable to make my life easier:
var title_line = "[font_size=40]
────────────────────
[/font_size]"

func generate_card(title, info):
	
	# If/else statement that deals with font size and keeping it from becoming 
	# too large for the text box:
	if title.length() <= 9:  # 9 is the max number of characters for the title at size 40
		$Title.text = title + title_line
	else:
		# Decreasing the font size proportionally to the length of the text:
		var title_size = ceil((15/float(title.length())) * 40)
		# Putting the text together:
		$Title.text = "[font_size={size}]{title}[/font_size]".format({"size": title_size, "title": title}) + title_line 
	
	# Another if/else statement to dynamically change the font size of the info:
	if info.length() <= 176:  # 176 characters is the max before a scroll bar appears
		$Card_info.text = info
	else:
		# The rest is the same as the title, minus the title_line
		var info_size = ceil((176/float(info.length())) * 40)
		$Card_info.text = "[font_size={size}]{info}[/font_size]".format({"size": info_size, "info": info})

func _ready() -> void:
	# Example generation:
	generate_card("Addition", "Add a two numbers--one from your set, one from another player--together, and add it to your set.")
