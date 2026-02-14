local scene = Scenario.new("environment_bree_outpost_shakedown")

-- Model the crime boss shakedown that the PCs can intervene in.
scene:campaign{
  name = "Environment Bree Outpost Shakedown",
  system = "DAGGERHEART",
  gm_mode = "HUMAN",
  theme = "environment"
}

scene:pc("Frodo")
scene:adversary("Orc Boss")

-- The party witnesses intimidation at a general goods store.
scene:start_session("Shakedown")

-- Example: the environment action introduces a threat without a roll.
-- Missing DSL: represent the narrative prompt and resulting tension.
scene:gm_spend_fear(0):spotlight("Orc Boss")

scene:end_session()

return scene
