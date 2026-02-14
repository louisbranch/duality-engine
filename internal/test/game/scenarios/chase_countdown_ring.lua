local scene = Scenario.new("chase_countdown_ring")

-- Model the ring chase with competing countdowns.
scene:campaign{
  name = "Chase Countdown Ring",
  system = "DAGGERHEART",
  gm_mode = "HUMAN",
  theme = "countdown"
}

scene:pc("Sam")
scene:pc("Frodo")
scene:adversary("Golum")

-- The PCs chase a thief across a market with progress and consequence clocks.
scene:start_session("Market Chase")

-- Example: PC progress countdown uses a d6; thief starts at 3 on a consequence countdown.
-- Missing DSL: represent d6 countdown dice and starting values.
scene:countdown_create{ name = "PC Progress", kind = "progress", current = 0, max = 6, direction = "increase" }
scene:countdown_create{ name = "Thief Escape", kind = "consequence", current = 0, max = 3, direction = "increase" }

-- Sam rolls with help and hope, advancing the PC countdown.
-- Missing DSL: tie action roll outcomes to countdown ticks.
scene:action_roll{ actor = "Sam", trait = "agility", difficulty = 15, outcome = "hope" }
scene:countdown_update{ name = "PC Progress", delta = 2, reason = "gain_ground" }

-- Close the session after the chase advances.
scene:end_session()

return scene
