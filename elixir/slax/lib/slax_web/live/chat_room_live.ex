# lib/slax_web/live/chat_room_live.ex
defmodule SlaxWeb.ChatRoomLive do
  use SlaxWeb, :live_view

  def render(assigns) do
    ~H"""
    <Layouts.app flash={@flash}>
      <div>Welcome to the chat!</div>
    </Layouts.app>
    """
  end
end
