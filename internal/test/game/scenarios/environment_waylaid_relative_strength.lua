local scene = Scenario.new("environment_waylaid_relative_strength")

-- Model Waylaid using the highest adversary Difficulty.
scene:campaign{
  name = "Environment Waylaid Relative Strength",
  system = "DAGGERHEART",
  gm_mode = "HUMAN",
  theme = "environment"
}

scene:pc("Frodo")
scene:adversary("Orc Sniper")
scene:adversary("Orc Lackey")

-- The ambushed difficulty matches the toughest adversary.
scene:start_session("Relative Strength")

-- Missing DSL: derive environment Difficulty from highest adversary.
scene:action_roll{ actor = "Frodo", trait = "instinct", difficulty = 0, outcome = "fear" }

scene:end_session()

return scene
