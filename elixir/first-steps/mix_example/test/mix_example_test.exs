defmodule MixExampleTest do
  use ExUnit.Case
  doctest MixExample

  test "greets the world" do
    assert MixExample.hello() == :world
  end
end
