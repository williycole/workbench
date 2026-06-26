defmodule BirdCount do
  def today([]), do: nil

  def today(list), do: list |> List.first()

  def increment_day_count([]), do: [1]
  def increment_day_count(list), do: List.update_at(list, 0, fn x -> x + 1 end)

  def has_day_without_birds?(list), do: Enum.any?(list, fn x -> x === 0 end)

  def total(list), do: Enum.sum(list)

  def busy_days(list), do: Enum.count(list, &(&1 >= 5))
end
