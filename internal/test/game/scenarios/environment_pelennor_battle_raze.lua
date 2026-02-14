local scene = Scenario.new("environment_pelennor_battle_raze")

-- Model the raze-and-pillage escalation.
scene:campaign{
  name = "Environment Pelennor Battle Raze",
  system = "DAGGERHEART",
  gm_mode = "HUMAN",
  theme = "environment"
}

scene:pc("Frodo")

-- The battle escalates with fire or abduction.
scene:start_session("Raze and Pillage")

-- Missing DSL: apply narrative escalation and objective shifts.
scene:gm_spend_fear(0):spotlight("Battlefield")

scene:end_session()

return scene
