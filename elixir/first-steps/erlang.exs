defmodule Example do
  def timed(fun, args) do
    {time, result} = :timer.tc(fun, args)
    IO.puts("Time: #{time} μs")
    IO.puts("Result: #{result}")
  end
end

Example.timed(fn n -> n * n * n end, [100])

# Erlang cang be covered by mix too
# if an erlang library ins't in mix 
# and is on git you can ref it like 
# below
def deps do
  [{:png, github: "yuce/png"}]
end
