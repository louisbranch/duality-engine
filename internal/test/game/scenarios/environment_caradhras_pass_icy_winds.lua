local scene = Scenario.new("environment_caradhras_pass_icy_winds")

-- Model the looping icy winds countdown.
scene:campaign{
  name = "Environment Caradhras Pass Icy Winds",
  system = "DAGGERHEART",
  gm_mode = "HUMAN",
  theme = "environment"
}

scene:pc("Frodo")

-- The pass inflicts stress at regular intervals.
scene:start_session("Icy Winds")

-- Example: countdown loop 4 triggers Strength reaction or Stress.
-- Missing DSL: implement looping countdown and cold gear advantage.
scene:countdown_create{ name = "Icy Winds", kind = "loop", current = 0, max = 4, direction = "increase" }
scene:reaction_roll{ actor = "Frodo", trait = "strength", difficulty = 15, outcome = "fear" }

scene:end_session()

return scene
