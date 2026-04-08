# NOTE: how to raise an error, in elixir, by convention, we only 
# raise exceptions, invalid user input we don't raise like this
# raise "Oh no!"
# raise ArgumentError, message: "the arg value was invalid"
try do
  opts
  |> Keyword.fetch!(:source_file)
  |> File.read!()
rescue
  # NOTE: in elixir try catch we can rescue on multiple exceptions
  e in KeyError -> IO.puts("missing :source_file option")
  e in File.Error -> IO.puts("unable to read source file")
after
  # NOTE: after is used most commonly with connections and files
  # that should be closed
  IO.puts("done reading source file")
  File.close(file)
end

# NOTE: can handle custom errors
defmodule ExampleError do
  defexception message: "an example error has occurred"
end

# NOTE: in old elixir code you may see throws, or as stopgaps
# when libraries don't provide adequate apis
try do
  for x <- 0..10 do
    if x == 5, do: throw(x)
    IO.puts(x)
  end
catch
  x -> "Caught: #{x}"
end

# NOTE: elixir also has exit, this occurs when a process dies
try do
  exit("oh no!")
catch
  :exit, _ -> "exit blocked"
end
