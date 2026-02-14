local scene = Scenario.new("environment_gondor_court_rival_vassals")

-- Capture the rival vassals social pressure in the court.
scene:campaign{
  name = "Environment Gondor Court Gondor Vassals",
  system = "DAGGERHEART",
  gm_mode = "HUMAN",
  theme = "environment"
}

scene:pc("Frodo")
scene:npc("Gondor Vassals")

-- Courtiers compete for favor and feed intrigue.
scene:start_session("Gondor Vassals")

-- Missing DSL: model ongoing social pressure and favor exchanges.
scene:gm_spend_fear(0):spotlight("Gondor Vassals")

scene:end_session()

return scene
