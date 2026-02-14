local scene = Scenario.new("environment_bruinen_ford_dangerous_crossing")

-- Model the dangerous crossing progress countdown.
scene:campaign{
  name = "Environment Bruinen Ford Dangerous Crossing",
  system = "DAGGERHEART",
  gm_mode = "HUMAN",
  theme = "environment"
}

scene:pc("Frodo")

-- Crossing requires a progress countdown and can trigger undertow.
scene:start_session("Dangerous Crossing")

-- Example: Progress Countdown (4) with failure + Fear triggering Undertow.
-- Missing DSL: tie failure with Fear to immediate undertow action.
scene:countdown_create{ name = "Bruinen Ford Crossing", kind = "progress", current = 0, max = 4, direction = "increase" }
scene:action_roll{ actor = "Frodo", trait = "agility", difficulty = 10, outcome = "fear" }

scene:end_session()

return scene
