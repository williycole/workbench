defmodule BirdCount do
  def today([]), do: nil
  def today([head | _tail]), do: head

  def increment_day_count([]), do: [1]
  def increment_day_count([head | tail]), do: [head + 1 | tail]

  def has_day_without_birds?([]), do: false
  def has_day_without_birds?([0 | _tail]), do: true
  def has_day_without_birds?([_head | tail]), do: has_day_without_birds?(tail)

  def total([]), do: 0
  def total([head | tail]), do: head + total(tail)

  def busy_days([]), do: 0
  def busy_days([head | tail]) when head >= 5, do: 1 + busy_days(tail)
  def busy_days([_head | tail]), do: busy_days(tail)
end

# NOTE: for refercne without recursion
#
# defmodule BirdCount do
#   def today([]), do: nil
#   def today(list), do: list |> List.first()
#
#   def increment_day_count([]), do: [1]
#
#   # short hand
#   def increment_day_count(list), do: List.update_at(list, 0, &(&1 + 1))
#   # traditional way
#   # def increment_day_count(list), do: List.update_at(list, 0, fn x -> x + 1 end)
#
#   # short hand
#   def has_day_without_birds?(list), do: Enum.any?(list, &(&1 === 0))
#   # traditional way
#   # def has_day_without_birds?(list), do: Enum.any?(list, fn x -> x === 0 end)
#
#   def total(list), do: Enum.sum(list)
#
#   def busy_days(list), do: Enum.count(list, &(&1 >= 5))
# end
