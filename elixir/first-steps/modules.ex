defmodule Example.User do
  @derive {Inspect, only: [:name]}
  defstruct name: nil, roles: []

  def hey_user(name), do: IO.puts("Hey #{name}")
end

defmodule Example.ScriptAlias do
  alias Example.User, as: Hi

  def run do
    cole = %Example.User{name: "Cole", roles: ["cowboy", "nerd", "fighter"]}
    anti_cole = %{cole | name: "AntiCole"}
    IO.inspect(cole)
    IO.inspect(anti_cole)
    # NOTE: Alias as this without the atom as: hi
    # User.hey_user(cole.name)
    # or like this for shorthand
    Hi.hey_user(cole.name)
    # NOTE: no alias 
    # User.hey_user(cole.name)
  end
end

Example.ScriptAlias.run()

# NOTE: you can also import specific functions like
# import List, except: [last: 1] for example.
# You can also import specifics like 
# import List, only: :functions
# import List, only: :macros

# NOTE: require is there too
# defmodule Example do
# require SuperMacros
# SuperMacros.do_stuff
# end

# NOTE: this is that templateing bit that GenServers and such use
# swaps out __using__ for the module name
# defmodule Hello do
#   defmacro __using__(opts) do
#     greeting = Keyword.get(opts, :greeting, "Hi")
#     quote do
#       def hello(name), do: unquote(greeting) <> ", " <> name
#     end
#   end
# end
