local scene = Scenario.new("environment_mirkwood_blight_choking_ash")

-- Model the looping choking ash countdown.
scene:campaign{
  name = "Environment Mirkwood Blight Choking Ash",
  system = "DAGGERHEART",
  gm_mode = "HUMAN",
  theme = "environment"
}

scene:pc("Gandalf")

-- The ash periodically forces reaction rolls.
scene:start_session("Choking Ash")

-- Missing DSL: loop countdown and apply direct damage with half on success.
scene:countdown_create{ name = "Choking Ash", kind = "loop", current = 0, max = 4, direction = "increase" }
scene:reaction_roll{ actor = "Gandalf", trait = "strength", difficulty = 16, outcome = "fear" }

scene:end_session()

return scene
