local scene = Scenario.new("environment_misty_ascent_progress")

-- Model the Misty Ascent progress countdown and roll outcomes.
scene:campaign{
  name = "Environment Misty Ascent Progress",
  system = "DAGGERHEART",
  gm_mode = "HUMAN",
  theme = "environment"
}

scene:pc("Frodo")

-- The party climbs using a progress countdown.
scene:start_session("Misty Ascent")

-- Example: Progress Countdown (12) ticks based on roll outcomes.
-- Missing DSL: apply tick deltas by outcome tiers.
scene:countdown_create{ name = "Misty Ascent", kind = "progress", current = 0, max = 12, direction = "increase" }
scene:action_roll{ actor = "Frodo", trait = "agility", difficulty = 12, outcome = "hope" }
scene:countdown_update{ name = "Misty Ascent", delta = 2, reason = "success_with_hope" }

scene:end_session()

return scene
