# NOTE: Processes in the Erlang VM are lightweight and run across 
# all CPUs. While they may seem like native threads, they’re 
# simpler and it’s not uncommon to have thousands of concurrent 
# processes in an Elixir application.
defmodule Example do
  def listen do
    receive do
      {:ok, "hello"} -> IO.puts("World")
    end

    listen()
  end
end

pid = spawn(Example, :listen, [])
# PID<0.108.0>

send(pid, {:ok, "hello"})
World
{:ok, "hello"}

send(pid, :ok)
:ok

# NOTE: Process Linking 
defmodule Example do
  def explode, do: exit(:kaboom)

  def run do
    Process.flag(:trap_exit, true)
    spawn_link(Example, :explode, [])

    receive do
      {:EXIT, _from_pid, reason} -> IO.puts("Exit reason: #{reason}")
    end
  end
end

# NOTE: Process Monitoring
defmodule Example do
  def explode, do: exit(:kaboom)

  def run do
    spawn_monitor(Example, :explode, [])

    receive do
      {:DOWN, _ref, :process, _from_pid, reason} -> IO.puts("Exit reason: #{reason}")
    end
  end
end

# NOTE: Tasks, how to handle expensive ops in the background without blocking
defmodule Example do
  def double(x) do
    :timer.sleep(2000)
    x * 2
  end
end

task = Task.async(Example, :double, [2000])
# %Task{
#   owner: #PID<0.105.0>,
#   pid: #PID<0.114.0>,
#   ref: #Reference<0.2418076177.4129030147.64217>
# }
