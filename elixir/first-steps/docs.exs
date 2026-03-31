# docs.exs
# Run with: elixir docs.exs

defmodule MyApp.Math do
  @moduledoc """
  Demo module showing Elixir doc conventions.
  """

  @doc """
  Adds two numbers.

  ## Examples

      iex> MyApp.Math.add(1, 2)
      3

  """
  @spec add(number(), number()) :: number()
  def add(a, b), do: a + b

  @doc """
  Starts a connector worker for the given tenant.

  Returns `{:ok, pid}` or `{:error, reason}`.
  """
  @spec start_connector(String.t()) :: {:ok, pid()} | {:error, term()}
  def start_connector(_tenant_id) do
    {:ok, self()}
  end

  @doc false
  def internal_helper(x), do: x * 2
end

IO.puts("=== Sanity checks ===")
IO.inspect(MyApp.Math.add(1, 2))
IO.inspect(MyApp.Math.start_connector("tenant_123"))
IO.inspect(MyApp.Math.internal_helper(5))
