defmodule Example.User do
  defstruct name: "Cole", roles: ["cowboy", "nerd", "fighter"]
end

anti_cole = %{cole | name: "AntiCole"}
