defmodule Rules do
  # Please implement the eat_ghost?/2 function
  def eat_ghost?(power_pellet_active?, touching_ghost?),
    do: power_pellet_active? and touching_ghost?

  # Please implement the score?/2 function
  def score?(touching_power_pellet?, touching_dot?),
    do: touching_power_pellet? or touching_dot?

  # Please implement the lose?/2 function
  def lose?(power_pellet_active?, touching_ghost?),
    do: !power_pellet_active? and touching_ghost?

  # Please implement the win?/3 function
  def win?(has_eaten_all_dots?, power_pellet_active?, touching_ghost?),
    do: !lose?(power_pellet_active?, touching_ghost?) and has_eaten_all_dots?
end
