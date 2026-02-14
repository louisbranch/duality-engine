local scene = Scenario.new("action_roll_critical_success")

-- Capture the critical success benefits from the example action roll.
scene:campaign{
  name = "Action Roll Critical Success",
  system = "DAGGERHEART",
  gm_mode = "HUMAN",
  theme = "core"
}

scene:pc("Sam")
scene:adversary("Nazgul")

-- A critical success grants Hope, clears Stress, and offers an extra benefit.
scene:start_session("Critical Success")

-- Missing DSL: apply Hope gain, Stress clear, and choose a bonus effect.
scene:action_roll{ actor = "Sam", trait = "agility", difficulty = 14, outcome = "critical" }

scene:end_session()

return scene
