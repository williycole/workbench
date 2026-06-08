defmodule GuessingGame do
  def compare(secret_number, guess) when secret_number == guess,
    do: "Correct"

  def compare(secret_number, guess)
      when is_nil(guess) or
             guess == :no_guess,
      do: "Make a guess"

  def compare(secret_number) do
    "Make a guess"
  end

  def compare(secret_number, guess)
      when secret_number + 1 == guess or
             secret_number - 1 == guess do
    "So close"
  end

  def compare(secret_number, guess) do
    cond do
      secret_number < guess -> "Too high"
      secret_number > guess -> "Too low"
    end
  end
end
