local scene = Scenario.new("environment_misty_ascent_pitons")

-- Capture the pitons rule that trades stress for a failed tick.
scene:campaign{
  name = "Environment Misty Ascent Pitons",
  system = "DAGGERHEART",
  gm_mode = "HUMAN",
  theme = "environment"
}

scene:pc("Frodo")

-- Pitons let a climber avoid a countdown setback by marking Stress.
scene:start_session("Pitons")

-- Example: on a failed climb, mark Stress instead of ticking up.
-- Missing DSL: intercept failure and apply stress in place of countdown penalty.
scene:action_roll{ actor = "Frodo", trait = "agility", difficulty = 12, outcome = "fear" }

scene:end_session()

return scene
