defmodule HighSchoolSweetheart do
  def first_letter(name), do: name |> String.trim() |> String.first()

  def initial(name), do: first_letter(String.upcase(name)) <> "."

  def initials(full_name) do
    name = String.split(full_name, " ")
    first_letter = Enum.at(name, 0) |> initial()
    second_letter = Enum.at(name, 1) |> initial()
    first_letter <> " " <> second_letter
  end

  def pair(full_name1, full_name2) do
    # ❤-------------------❤
    # |  X. X.  +  X. X.  |
    # ❤-------------------❤
    #
    divider = "❤-------------------❤"
    i1 = initials(full_name1)
    i2 = initials(full_name2)

    """
    #{divider}
    |  #{i1}  +  #{i2}  |
    #{divider}
    """
  end
end
