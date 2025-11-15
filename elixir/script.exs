IO.inspect("Hello World!")
# String
name = "Peter"
# Integer
age = 32
# Integer in Hex notation
age_hex = 0x20
# Float
height = 190.47
# Float in scientific notation
height_sci = 1.9047e2
# Boolean
adult? = true
# Atom
status = :active
# The 'None/null/Nil' value
address = nil
IO.inspect([name, age, age_hex, height, height_sci, adult?, status, address])
age = 32
age = 21
IO.inspect(age)
_ = 30
_result = 30 + 32
# Interpolate a string with values
name = "Peter"
age = 32
IO.inspect("My name is #{name} and I am #{age} years old")
# Concatenate two strings
IO.inspect("But my real name is " <> "Batman 🦇")
