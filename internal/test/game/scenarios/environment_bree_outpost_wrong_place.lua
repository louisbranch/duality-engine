local scene = Scenario.new("environment_bree_outpost_wrong_place")

-- Capture the ambush by thieves in a dark alley.
scene:campaign{
  name = "Environment Bree Outpost Wrong Place",
  system = "DAGGERHEART",
  gm_mode = "HUMAN",
  theme = "environment"
}

scene:pc("Frodo")
scene:adversary("Orc Captain")
scene:adversary("Orc Lackeys")
scene:adversary("Orc Lieutenant")

-- Thieves emerge at close range when the party is isolated.
scene:start_session("Wrong Place, Wrong Time")
scene:gm_fear(1)

-- Example: spend Fear to introduce a robber group at Close range.
-- Missing DSL: spawn multiple adversaries based on party size.
scene:gm_spend_fear(1):spotlight("Orc Captain")

scene:end_session()

return scene
