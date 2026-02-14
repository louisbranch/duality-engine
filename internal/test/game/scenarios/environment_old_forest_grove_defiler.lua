local scene = Scenario.new("environment_old_forest_grove_defiler")

-- Capture the Defiler fear action summoning a chaos elemental.
scene:campaign{
  name = "Environment Old Forest Grove Defiler",
  system = "DAGGERHEART",
  gm_mode = "HUMAN",
  theme = "environment"
}

scene:pc("Frodo")
scene:adversary("Shadow Wraith")

-- The grove draws in a chaotic threat.
scene:start_session("Defiler")
scene:gm_fear(1)

-- Example: spend Fear to summon an elemental that immediately takes spotlight.
-- Missing DSL: spawn the elemental near a chosen PC and shift spotlight.
scene:gm_spend_fear(1):spotlight("Shadow Wraith")

scene:end_session()

return scene
