local scene = Scenario.new("environment_bree_market_crowd_closes_in")

-- Capture the crowd reaction that splits a PC from the party.
scene:campaign{
  name = "Environment Bree Market Crowd Closes In",
  system = "DAGGERHEART",
  gm_mode = "HUMAN",
  theme = "environment"
}

scene:pc("Frodo")

-- The crowd shifts when a PC splits off.
scene:start_session("Crowd Closes In")

-- Missing DSL: separate the PC from the group and apply positioning.
scene:gm_spend_fear(0):spotlight("Crowd")

scene:end_session()

return scene
