local scene = Scenario.new("experience_spend_modifier")

-- Model spending Hope to apply an Experience modifier.
scene:campaign{
  name = "Experience Spend Modifier",
  system = "DAGGERHEART",
  gm_mode = "HUMAN",
  theme = "experience"
}

scene:pc("Frodo")

-- Frodo uses a relevant Experience by spending Hope for a modifier.
scene:start_session("Experience Modifier")

-- Missing DSL: spend Hope and apply an Experience bonus to the roll.
scene:action_roll{ actor = "Frodo", trait = "presence", difficulty = 12, outcome = "hope" }

scene:end_session()

return scene
