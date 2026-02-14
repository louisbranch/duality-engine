local scene = Scenario.new("environment_waylayers_relative_strength")

-- Model Orc Waylayers using the highest adversary Difficulty.
scene:campaign{
  name = "Environment Orc Waylayers Relative Strength",
  system = "DAGGERHEART",
  gm_mode = "HUMAN",
  theme = "environment"
}

scene:pc("Frodo")
scene:adversary("Orc Sniper")
scene:adversary("Orc Lackey")

-- The ambush difficulty matches the toughest adversary.
scene:start_session("Relative Strength")

-- Missing DSL: derive environment Difficulty from highest adversary.
scene:action_roll{ actor = "Frodo", trait = "instinct", difficulty = 0, outcome = "hope" }

scene:end_session()

return scene
