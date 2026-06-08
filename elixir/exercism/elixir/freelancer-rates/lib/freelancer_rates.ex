defmodule FreelancerRates do
  # Please implement the daily_rate/1 function
  def daily_rate(hourly_rate), do: 8 * hourly_rate / 1

  # Please implement the apply_discount/2 function
  def apply_discount(before_discount, discount) do
    sub = before_discount * discount / 100
    (before_discount - sub) / 1
  end

  # Please implement the monthly_rate/2 function
  def monthly_rate(hourly_rate, discount) do
    monthly = 22 * daily_rate(hourly_rate)
    ceil(apply_discount(monthly, discount))
  end

  # Please implement the days_in_budget/3 function
  def days_in_budget(budget, hourly_rate, discount) do
    day_by_rate = budget / daily_rate(hourly_rate)
    Float.floor(apply_discount(day_by_rate, discount), 1)
  end
end
