defmodule Slax.Chat do
  alias Slax.Accounts.Scope
  alias Slax.Chat.Message
  alias Slax.Chat.Room
  alias Slax.Repo

  import Ecto.Query

  @pubsub Slax.PubSub

  def subscribe_to_room(room) do
    Phoenix.PubSub.subscribe(@pubsub, topic(room.id))
  end

  def unsubscribe_from_room(room) do
    Phoenix.PubSub.unsubscribe(@pubsub, topic(room.id))
  end

  defp topic(room_id), do: "chat_room:#{room_id}"

  def change_message(message, attrs \\ %{}, scope) do
    Message.changeset(message, attrs, scope)
  end

  def create_message(room, attrs, scope) do
    with {:ok, message} <-
           %Message{room: room}
           |> Message.changeset(attrs, scope)
           |> Repo.insert() do
      message = Repo.preload(message, :user)
      Phoenix.PubSub.broadcast!(@pubsub, topic(room.id), {:new_message, message})
      {:ok, message}
    end
  end

  def delete_message_by_id(id, %Scope{user: user}) do
    message = Repo.get_by!(Message, id: id, user_id: user.id)

    Repo.delete(message)
    Phoenix.PubSub.broadcast!(@pubsub, topic(message.room_id), {:message_deleted, message})
  end

  def list_messages_in_room(%Room{id: room_id}) do
    Message
    |> where([m], m.room_id == ^room_id)
    |> order_by([m], asc: :inserted_at, asc: :id)
    |> preload(:user)
    |> Repo.all()
  end

  def list_rooms do
    Repo.all(from Room, order_by: [asc: :name])
  end

  def get_room!(id) do
    Repo.get!(Room, id)
  end

  def change_room(room, attrs \\ %{}) do
    Room.changeset(room, attrs)
  end

  def create_room(attrs) do
    %Room{}
    |> Room.changeset(attrs)
    |> Repo.insert()
  end

  def update_room(%Room{} = room, attrs) do
    room
    |> Room.changeset(attrs)
    |> Repo.update()
  end
end
