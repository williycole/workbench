a_list = [1, 2, 3]
b_list = ["four", 5.00, 6]
c_list = a_list ++ b_list
d_list = ["foo", :bar, 42] -- [42, "bar"]
IO.inspect(hd(a_list))
IO.inspect(tl(a_list))
# first | last item in a list
[head | tail] = [3.14, :pie, "Apple"]
IO.inspect(head)
IO.inspect(tail)
# Tupeles are commoanly used as a mechanism to return additional info from functions
a_tuple = {3.14, :pie, "Apple"}
IO.inspect(elem(a_tuple, 0))
# Keyword lists are often used to pass options to a function
# ordered by the dev, keys must be atoms, keys don't have to be unique
# share preformance with lists
key_word_list = [{:foo, "bar"}, {:hello, "world"}]
IO.inspect(hd(key_word_list))
# maps are the go to key value store
a_map = %{:foo => "bar", "hello" => :world}
IO.inspect(a_map[:foo])
IO.inspect(a_map["hello"])
IO.inspect(Map.get(a_map, :foo))
IO.inspect(a_map.foo)
# update a map, but its a new map bc immutability
# this updarte only works where keys already exist
b_map = %{a_map | foo: "baz"}
IO.inspect(b_map)
c_map = Map.put(a_map, :hello, "world")
IO.inspect(c_map)
