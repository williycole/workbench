defmodule NameBadge do
  # def print(nil, name, nil),
  #   do: "#{name} - OWNER"
  #
  # def print(id, name, nil),
  #   do: "[#{id}] - #{name} - OWNER"
  #
  # def print(nil, name, department),
  #   do: "#{name} - #{String.upcase(department)}"
  #
  # def print(id, name, department),
  #   do: "[#{id}] - #{name} - #{String.upcase(department)}"
  def print(id, name, department) do
    prefix = if id, do: "[#{id}] - ", else: ""
    dept = if department, do: String.upcase(department), else: "OWNER"
    "#{prefix}#{name} - #{dept}"
  end
end
