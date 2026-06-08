defmodule LogLevel do
  # Please implement the to_label/2 function
  def to_label(level, legacy?) do
    #
    cond do
      # no support for legacy apps
      level == 0 and !legacy? -> :trace
      level == 0 and legacy? -> :unknown
      level == 5 and !legacy? -> :fatal
      level == 5 and legacy? -> :unknown
      level > 5 or level < 0 -> :unknown
      # support for legacy apps
      level == 1 -> :debug
      level == 2 -> :info
      level == 3 -> :warning
      level == 4 -> :error
    end
  end

  # Please implement the alert_recipient/2 function
  def alert_recipient(level, legacy?) do
  end
end
