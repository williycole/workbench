defmodule Slax.Chat.Room do
  use Ecto.Schema
  import Ecto.Changeset

  schema "rooms" do
    field :name, :string
    field :topic, :string

    timestamps(type: :utc_datetime)
    # above is shorthand for below
    # field :inserted_at, :utc_datetime
    # field :updated_at, :utc_datetime
  end

  @doc false
  def changeset(room, attrs) do
    room
    |> cast(attrs, [:name, :topic])
    |> validate_required([:name, :topic])
  end
end
