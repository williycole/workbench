IO.puts("Elixir rocks" |> String.upcase() |> String.split())

raw_csv = "sensor-1,CRITICAL,42.5\nsensor-2,warning,18.0\nsensor-3,INFO,99.1"

# NOTE: an example of how to manipulate data with pipes
# this should not be done on large data sets
result =
  raw_csv
  |> String.split("\n")
  |> Enum.map(&String.split(&1, ","))
  |> Enum.map(fn [id, severity, value] ->
    %{id: id, severity: String.downcase(severity), value: String.to_float(value)}
  end)
  |> Enum.filter(&(&1.value > 20.0))
  |> Enum.sort_by(& &1.value, :desc)

IO.inspect(result)
