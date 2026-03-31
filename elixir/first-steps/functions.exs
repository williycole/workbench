# handle_result = fn
#   {:ok, result} -> IO.puts("Handling result...")
#   {:ok, _} -> IO.puts("This would be never run as previous will be matched beforehand.")
#   {:error} -> IO.puts("An error has occurred!")
# end

# NOTE: Pattern matching and recursion
defmodule Length do
  def of([]), do: 0
  def of([_ | tail]), do: 1 + of(tail)
end

IO.puts(Length.of([1, 2, 3]))
IO.puts(Length.of([]))

defmodule Greeter do
  def hello() do
    hello("world")
  end

  def hello(name1, name2) when is_binary(name1) and is_binary(name2) do
    hello([name1, name2])
  end

  def hello(names) when is_list(names) do
    names = Enum.join(names, ", ")
    hello(names)
  end

  def hello(name) when is_binary(name) do
    phrase() <> name
  end

  defp phrase, do: "Hello, "
end

# "Hello, world"
IO.puts(Greeter.hello())
# "Hello, John"
IO.puts(Greeter.hello("John"))
# "Hello, John, Mary"   
IO.puts(Greeter.hello("John", "Mary"))

defmodule Greeter2 do
  def hello(%{name: person_name} = person) do
    IO.puts("Hello, " <> person_name)
    IO.inspect(person)
  end
end

person = %{name: "Fred", age: "95", favorite_color: "Taupe"}
Greeter2.hello(person)

defmodule Greeter3 do
  # NOTE: Elixir doens't like defualt args with multiple matching functions so 
  # we use a function head
  def hello(names, language_code \\ "en")

  def hello(names, language_code) when is_list(names) do
    comma_separated_names = Enum.join(names, ", ")

    hello(comma_separated_names, language_code)
  end

  def hello(name, language_code) when is_binary(name) do
    phrase(language_code) <> name
  end

  defp phrase("en"), do: "Hello, "
  defp phrase("es"), do: "Hola, "
end

IO.puts(Greeter3.hello("John"))
