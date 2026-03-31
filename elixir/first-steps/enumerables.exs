IO.puts("🌞 There are over 70 functions for working with enumerables")
# there are over 70 functions for working with enumerables
a_all = Enum.all?([1, 2, 3], fn x -> x > 0 end)
b_all = Enum.all?(["foo", "bar", "baz"], fn s -> String.length(s) > 0 end)
a_b_all = [a_all, b_all]
# print the whole list
IO.inspect(a_b_all)
# print each value on a new line
Enum.each(a_b_all, &IO.puts/1)
# returns true if at least one val is true
c_any = Enum.any?(["foo", "bar", "hello"], fn s -> String.length(s) == 5 end)
IO.inspect(c_any)
# if you need to break up a collection this is your boy
# here we break up by 2
d_chunk = Enum.chunk_every([1, 2, 3, 4, 5, 6], 2)
# print the whole list
IO.inspect(d_chunk)
# print each value on a new line
Enum.each(d_chunk, &IO.inspect/1)
# do  a thing to each value and get a new list
a_map = [1, 2, 3, 4, 5, 6]
b_map = Enum.map(a_map, fn x -> x * 4 end)
IO.inspect(b_map)
IO.inspect(a_map)
some_list = [3, 2, 1, 5, 7, 2]
m = Enum.min(some_list)
IO.inspect(m)
n = Enum.max(some_list)
IO.inspect(some_list)
m_empty = Enum.min([], fn -> :empty end)
IO.inspect(m_empty)
# filter on
IO.puts("☀️ Filters and Reduce")
zero_remainders = Enum.filter(some_list, fn y -> rem(y, 2) == 0 end)
IO.inspect(zero_remainders)
two_list = [2, 2]
sum_of_10_plus_each_value = Enum.reduce(two_list, 10, fn x, acc -> x + acc end)
IO.puts("🙉 Sorts and Uniques")
IO.inspect(sum_of_10_plus_each_value)
basic_sort = Enum.sort([3, 2, 5, 2, 1, 8])
IO.inspect(basic_sort)
custom_sort = Enum.sort([%{:val => 4}, %{:val => 1}], fn x, y -> x[:val] > y[:val] end)
IO.inspect(custom_sort)
IO.inspect(Enum.sort([%{:count => 4}, %{:count => 1}]))
des_sort = Enum.sort([2, 3, 4, 5, 1, 8, 2, 1], :desc)
IO.inspect(des_sort)
unique = Enum.uniq(des_sort)
IO.inspect(unique)
unique_by = Enum.uniq_by([%{x: 1, y: 1}, %{x: 2, y: 1}, %{x: 3, y: 3}], fn coord -> coord.y end)
IO.inspect(unique_by)
IO.puts("💅Using captures and fancy short hand for percision syntax")
capture_a = Enum.map([1, 2, 3], fn num -> num + 3 end)
capture_b = Enum.map([1, 2, 3], &(&1 + 3))
plus_3 = &(&1 + 3)
capture_c = Enum.map([1, 2, 3], plus_3)
IO.inspect(capture_a)
IO.inspect(capture_b)
IO.inspect(capture_c)
