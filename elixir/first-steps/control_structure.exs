# NOTE: in elixir we still do have if else statments, but pattern matching
# is really the thing to lean on. Use if else for like single boolean condtions
greet = fn
  :ok, name ->
    if name == "cole" do
      "hello, #{name}, you are learning elixir"
    else
      "hello, #{name}, you should learn elixir"
    end

  :error, _ ->
    "something went wrong"
end

# NOTE: the proper way to do if checks is for guard clauses
# think negative space programming
def process_alarm(alarm) do
  if is_nil(alarm.resolved_at), do: notify_oncall(alarm)

  log_alarm(alarm)
end

# NOTE: so if you can simplify to pattern matching do so
greet2 = fn
  nil, _ -> "no name given"
  :ok, "cole" -> "hello, cole, you are learning elixir"
  :ok, name -> "hello, #{name}, you should learn elixir"
  :error, _ -> "something went wrong"
end

IO.inspect(greet2.(:ok, "cole"))
IO.inspect(greet2.(:ok, "bob"))
IO.inspect(greet2.(:ok, "bob"))

# NOTE: example of when you may want an if
greet3 = fn
  nil, _ ->
    "no name given"

  :ok, "cole" ->
    if System.get_env("PROD"), do: Logger.info("cole logged in")
    "hello, cole, you are learning elixir"

  :ok, name ->
    "hello, #{name}, you should learn elixir"

  :error, _ ->
    "something went wrong"
end

# NOTE: case example, branching on a single value
def handle_alarm(alarm) do
  case alarm.severity do
    :critical -> page_oncall(alarm)
    :warning -> send_slack(alarm)
    :info -> log_alarm(alarm)
    _ -> :ok
  end
end

# NOTE: cond, multiple unrelated boolean checks
def alarm_priority(alarm) do
  cond do
    alarm.severity == :critical and alarm.unacknowledged -> 1
    alarm.severity == :critical -> 2
    alarm.unacknowledged -> 3
    true -> 4
  end
end

# NOTE: with is a sequence of steps where each step must succeed for the next to run.
def process_alarm(tenant_id, params) do
  # with is like an anonymus holder for our data
  # its like setting a struct in go but we arent actually 
  #
  # doing like {tenant.id = res.tenant.id, alarm.id = res.alarm.id}
  # instead we get what the data looks like, its 'form'
  # and we handle it accordingly, if at any point we can't get that data
  # to completley form our data, we bail out
  with {:ok, tenant} <- Repo.fetch(Tenant, tenant_id),
       {:ok, alarm} <- Alarms.create(tenant, params),
       {:ok, _notif} <- Notifications.send(alarm) do
    # all 3 steps succeeded, return the alarm
    {:ok, alarm}
    # bail out
  else
    {:error, :not_found} -> {:error, "tenant not found"}
    {:error, changeset} -> {:error, changeset}
  end
end
