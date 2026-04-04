# NOTE: strings are just a sequence of bytes
string = <<104, 101, 108, 108, 111>>
IO.puts(string)
bytes = string <> <<0>>
# NOTE: puts automatically formats binaries as readable strings
IO.puts(bytes)
IO.inspect(bytes)
# NOTE: charlists
# IO.inspect('hello') # the old way for charlists
# below is the new way
IO.puts(~c"hello")
IO.inspect("hello" <> <<0>>)
# get a characters code point using ? 
# codepoints are basic unicode characters
# that represent one or more bytes depending on UTF-8 encoding
IO.inspect(?Z)
# for anscii
stringa = "\u0061\u0301"
IO.inspect(stringa)
IO.inspect(String.codepoints(stringa))
IO.inspect(String.graphemes(stringa))
IO.inspect(String.length("Hello"))
IO.inspect(String.replace("Hello", "e", "a"))
IO.inspect(String.split("Hello", "e"))
IO.inspect(String.trim(" Hello "))
IO.inspect(String.duplicate("no yes ", 3))
