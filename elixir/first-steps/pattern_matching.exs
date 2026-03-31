# NOTE: Pattern Matching is how we define Success
list = [1, 2, 3]
IO.inspect(list)

# this is ok bc list is set to this
# [2,3,4] = list would give an error
[1, 2, 3] = list
IO.inspect(list)

# Will error
# [] = list
# IO.inspect([] = list)

# NOTE: In go we have multiple ifs instead a func
# but in elixir we push it to function heads
# and pattern match
def calc(:add, nums), do: {:ok, Enum.sum(nums)}
# &-/2 is a capture opperator that says subtract the last two numbers
def calc(:sub, nums), do: {:ok, Enum.reduce(nums, &-/2)}
# similar here with &*/2 for multiply
def calc(:mul, nums), do: {:ok, Enum.reduce(nums, &*/2)}
def calc(:div, nums), do: {:ok, Enum.reduce(nums, &//2)}
def calc(_, _), do: {:error, "Invalid operation"}

# NOTE: Pin Operator 
# Without ^, current_user_id gets rebound to whatever author_id is on the post
# the unauthorized check never fires, every user can edit every post.
# With ^, the match only succeeds if the post's author_id equals the logged-in user's ID.
def handle_in("update_post", %{"post_id" => post_id, "author_id" => author_id}, socket) do
  current_user_id = socket.assigns.user_id

  case Posts.get_post(post_id) do
    %{author_id: ^current_user_id} -> {:ok, Posts.update(post_id)}
    %{author_id: _} -> {:error, :unauthorized}
  end
end

# // Go
# tenant, err := repo.Fetch(tenantID)
# if err != nil { return err }

# alarm, err := alarms.Create(tenant, params)
# if err != nil { return err }
#
## Elixir
# with {:ok, tenant} <- Repo.fetch(Tenant, tenant_id),
# {:ok, alarm}  <- Alarms.create(tenant, params) do
# {:ok, alarm}
# end
