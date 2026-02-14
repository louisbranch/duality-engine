local scene = Scenario.new("environment_bree_outpost_broken_compass")

-- Capture the adventuring society passive at an outpost town.
scene:campaign{
  name = "Environment Bree Outpost Broken Compass",
  system = "DAGGERHEART",
  gm_mode = "HUMAN",
  theme = "environment"
}

scene:pc("Frodo")
scene:npc("Elrond")

-- The society offers boasts, rumors, and rivalries.
scene:start_session("Broken Compass")

-- Example: a passive feature that sets social tension and leads.
-- Missing DSL: represent the ongoing social pressure from the society.
scene:gm_spend_fear(0):spotlight("Elrond")

scene:end_session()

return scene
