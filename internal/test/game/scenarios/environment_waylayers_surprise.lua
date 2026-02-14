local scene = Scenario.new("environment_waylayers_surprise")

-- Model the ambushers' surprise action for a sudden strike.
scene:campaign{
  name = "Environment Orc Waylayers Surprise",
  system = "DAGGERHEART",
  gm_mode = "HUMAN",
  theme = "environment"
}

scene:pc("Frodo")
scene:adversary("Orc Waylayers")

-- The ambush begins, shifting the spotlight and adding Fear.
scene:start_session("Ambush")
scene:gm_fear(2)

-- Example: Surprise grants 2 Fear and spotlights an ambusher.
-- Missing DSL: award Fear to the GM and immediate spotlight shift.
scene:gm_spend_fear(2):spotlight("Orc Waylayers")

scene:end_session()

return scene
