list = [1, 2, 3, 4, 5]
IO.inspect(list)
res = for x <- list, do: x * 2
IO.inspect(res)
keyword_list = [a: 1, b: 2, c: 3]
IO.inspect(keyword_list)
vals = for {_key, val} <- keyword_list, do: val * 2
IO.inspect(vals)

res2 =
  for n <- list, times <- 1..n do
    String.duplicate("*", times)
  end

IO.inspect(res2)

into_map = for {k, v} <- [a: "a", b: "b", c: "c"], into: %{}, do: {k, v}
IO.inspect(into_map)

# NOTE: how you can use some of this for helpers and the like
defmodule Helpers do
  def kv_into_map(list) do
    for {k, v} <- list, into: %{}, do: {k, v}
  end
end

IO.inspect(Helpers.kv_into_map(one: 1, two: 2, three: 3))
