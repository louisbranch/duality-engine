local scene = Scenario.new("progress_countdown_climb")

-- Model the mountain ascent progress countdown from the example of play.
scene:campaign{
  name = "Progress Countdown Climb",
  system = "DAGGERHEART",
  gm_mode = "HUMAN",
  theme = "countdown"
}

scene:pc("Frodo")
scene:pc("Sam")

-- The GM sets a shorter progress countdown due to helpful guidance.
scene:start_session("Whitecrest Ascent")

-- Example: progress countdown starts at 3 instead of 5.
scene:countdown_create{ name = "Whitecrest Ascent", kind = "progress", current = 3, max = 3, direction = "decrease" }

-- Sam succeeds with Fear, advancing the climb despite consequences.
-- Missing DSL: tick countdown by result tier and award Hope/Fear changes.
scene:action_roll{ actor = "Sam", trait = "agility", difficulty = 12, outcome = "fear" }
scene:countdown_update{ name = "Whitecrest Ascent", delta = -1, reason = "climb" }

-- Close the session after the ascent advances.
scene:end_session()

return scene
