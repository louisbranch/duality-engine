local scene = Scenario.new("spellcast_hope_cost")

-- Capture the spellcast roll that costs Hope to cast.
scene:campaign{
  name = "Spellcast Hope Cost",
  system = "DAGGERHEART",
  gm_mode = "HUMAN",
  theme = "spellcast"
}

scene:pc("Gandalf")

-- Gandalf casts a warding door with a Hope cost.
scene:start_session("Arcane Door")

-- Example: Spellcast roll Difficulty 13, success with Fear after spending Hope.
-- Missing DSL: spend Hope to cast and apply the Fear gain to the GM.
scene:action_roll{ actor = "Gandalf", trait = "spellcast", difficulty = 13, outcome = "fear" }

-- Close the session after the spellcast.
scene:end_session()

return scene
