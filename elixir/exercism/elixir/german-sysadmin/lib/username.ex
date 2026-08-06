defmodule Username do
  # base case: empty list -> stop recursion
  def sanitize([]), do: []

  # one clause for a non-empty list: peel head off, `case` decides what to do with it
  def sanitize([head | tail]) do
    case head do
      # German chars -> splice ASCII pair in, recurse tail
      ?ä -> ~c"ae" ++ sanitize(tail)
      ?ö -> ~c"oe" ++ sanitize(tail)
      ?ü -> ~c"ue" ++ sanitize(tail)
      ?ß -> ~c"ss" ++ sanitize(tail)
      # keep lowercase a-z or underscore -> prepend head, recurse tail
      h when h in ?a..?z or h == ?_ -> [h | sanitize(tail)]
      # anything else -> drop, recurse tail
      _ -> sanitize(tail)
    end
  end
end

defmodule UsernameTake2 do
  # base case: empty list -> stop recursion, return empty
  def sanitize([]), do: []
  # German chars: swap for ASCII pair, keep walking tail.
  # Read left->right: IF list head is ä, THEN build ~c"ae" concatenated (++) onto the sanitized tail.
  #   [?ä | t]  = list whose head is the codepoint of ä, rest of list is `t`
  #   ~c"ae"    = the charlist [?a, ?e], the ASCII replacement
  #   ++        = glue two lists together
  #   sanitize(t) = run the same function on the leftover chars
  # ä -> ae
  def sanitize([?ä | t]), do: ~c"ae" ++ sanitize(t)
  # ö -> oe
  def sanitize([?ö | t]), do: ~c"oe" ++ sanitize(t)
  # ü -> ue
  def sanitize([?ü | t]), do: ~c"ue" ++ sanitize(t)
  # ß -> ss
  def sanitize([?ß | t]), do: ~c"ss" ++ sanitize(t)
  # keep: lowercase a-z or underscore -> prepend head, recurse tail
  def sanitize([h | t]) when h in ?a..?z or h == ?_, do: [h | sanitize(t)]

  # catch-all: anything else (digits, symbols, uppercase) -> drop, recurse tail
  def sanitize([_ | t]), do: sanitize(t)
end
