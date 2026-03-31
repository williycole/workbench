defmodule MixExample.MixProject do
  use Mix.Project

  def project do
    [
      app: :mix_example,
      version: "0.1.0",
      elixir: "~> 1.19",
      start_permanent: Mix.env() == :prod,
      deps: deps()
    ]
  end

  # Run "mix help compile.app" to learn about applications.
  def application do
    [
      extra_applications: [:logger]
    ]
  end

  # Run "mix help deps" to learn about dependencies.
  defp deps do
    [
      {:phoenix, "~> 1.8"},
      {:phoenix_html, "~> 4.3"},
      {:cowboy, "~> 2.12", only: [:dev, :test]},
      {:slime, "~> 1.3"}
    ]
  end
end
