defmodule LanguageList do
  def new(), do: []

  def add(list, language), do: [language | list]

  def remove(list), do: List.delete(list, List.first(list))

  def first(list), do: list |> List.first()

  def count(list), do: list |> Enum.count()

  def functional_list?(list),
    do:
      list
      |> Enum.any?(&(&1 in ["Clojure", "Haskell", "Erlang", "Elixir", "F#"]))
end
