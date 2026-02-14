local scene = Scenario.new("environment_gondor_court_all_roads")

-- Model disadvantage on Presence rolls that resist imperial norms.
scene:campaign{
  name = "Environment Gondor Court All Roads",
  system = "DAGGERHEART",
  gm_mode = "HUMAN",
  theme = "environment"
}

scene:pc("Frodo")

-- Court etiquette hampers dissenting actions.
scene:start_session("All Roads Lead Here")

-- Missing DSL: apply disadvantage to nonconforming Presence rolls.
scene:action_roll{ actor = "Frodo", trait = "presence", difficulty = 20, outcome = "fear" }

scene:end_session()

return scene
