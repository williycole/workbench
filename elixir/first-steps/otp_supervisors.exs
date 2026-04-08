# NOTE: Supervisors are specialized processes with one purpose: monitoring other processes. These supervisors enable us to create fault-tolerant applications by automatically restarting child processes when they fail.
#  The magic of Supervisors is in the Supervisor.start_link/2 function. In addition to starting our supervisor and children, it allows us to define the strategy our supervisor uses for managing child processes.
defmodule SimpleQueue.Application do
  use Application

  def start(_type, _args) do
    children = [
      {SimpleQueue, [1, 2, 3]}
    ]

    opts = [strategy: :one_for_one, name: SimpleQueue.Supervisor]
    Supervisor.start_link(children, opts)
  end
end

# Strategies
# There are currently three different supervision strategies available to supervisors:
#     :one_for_one - Only restart the failed child process.
#     :one_for_all - Restart all child processes in the event of a failure.
#     :rest_for_one - Restart the failed process and any process started after it.
