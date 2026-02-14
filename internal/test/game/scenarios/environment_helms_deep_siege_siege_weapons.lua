local scene = Scenario.new("environment_helms_deep_siege_siege_weapons")

-- Model the siege weapons countdown breaching the walls.
scene:campaign{
  name = "Environment Helms Deep Siege Siege Weapons",
  system = "DAGGERHEART",
  gm_mode = "HUMAN",
  theme = "environment"
}

scene:pc("Frodo")

-- Siege weapons grind down defenses.
scene:start_session("Siege Weapons")

-- Missing DSL: activate consequence countdown and shift to Helms Deep Siege.
scene:countdown_create{ name = "Breach the Walls", kind = "consequence", current = 0, max = 6, direction = "increase" }

scene:end_session()

return scene
