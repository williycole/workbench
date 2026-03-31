IO.inspect(~c/2 + 7 = #{2 + 7}/)
IO.inspect(~C/2 + 7 = #{2 + 7}/)
re = ~r/elixir/
IO.inspect("Elixir" =~ re)

ree = ~r/elixir/i
IO.inspect("Elixir" =~ ree)
