local scene = Scenario.new("environment_helms_deep_siege_reinforcements")

-- Capture the reinforcements action that brings in new foes.
scene:campaign{
  name = "Environment Helms Deep Siege Reinforcements",
  system = "DAGGERHEART",
  gm_mode = "HUMAN",
  theme = "environment"
}

scene:pc("Frodo")
scene:adversary("Gondor Knight")
scene:adversary("Uruk-hai Minions")

-- New forces arrive within Far range.
scene:start_session("Reinforcements")

-- Missing DSL: spawn adversaries based on party size and spotlight the knight.
scene:gm_spend_fear(0):spotlight("Gondor Knight")

scene:end_session()

return scene
