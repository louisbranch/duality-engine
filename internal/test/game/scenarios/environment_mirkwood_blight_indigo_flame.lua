local scene = Scenario.new("environment_mirkwood_blight_indigo_flame")

-- Capture the knowledge roll about the indigo flame corruption.
scene:campaign{
  name = "Environment Mirkwood Blight Indigo Flame",
  system = "DAGGERHEART",
  gm_mode = "HUMAN",
  theme = "environment"
}

scene:pc("Gandalf")

-- The party studies the corrupted tree.
scene:start_session("Indigo Flame")

-- Missing DSL: map outcome to number of details and stress for extra clue.
scene:action_roll{ actor = "Gandalf", trait = "knowledge", difficulty = 16, outcome = "hope" }

scene:end_session()

return scene
