local scene = Scenario.new("environment_old_forest_grove_overgrown")

-- Capture the Overgrown Battlefield investigation in the Old Forest Grove.
scene:campaign{
  name = "Environment Old Forest Grove Overgrown",
  system = "DAGGERHEART",
  gm_mode = "HUMAN",
  theme = "environment"
}

scene:pc("Frodo")

-- The grove reveals its history on a successful Instinct roll.
scene:start_session("Old Forest Grove")

-- Example: success with Hope learns all details; failure can mark Stress for one.
-- Missing DSL: apply graded information gain and stress option.
scene:action_roll{ actor = "Frodo", trait = "instinct", difficulty = 11, outcome = "hope" }

scene:end_session()

return scene
