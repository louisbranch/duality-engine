local scene = Scenario.new("gm_move_examples")

-- Showcase a few GM move examples tied to roll outcomes.
scene:campaign{
  name = "GM Move Examples",
  system = "DAGGERHEART",
  gm_mode = "HUMAN",
  theme = "gm_move"
}

scene:pc("Gandalf")

-- The GM responds to fear and failure with narrative moves.
scene:start_session("GM Moves")
scene:gm_fear(2)

-- Example: roll with Fear triggers a move showing how the world reacts.
scene:action_roll{ actor = "Gandalf", trait = "presence", difficulty = 12, outcome = "fear" }
scene:gm_spend_fear(1):spotlight("Gandalf")

-- Example: a hard move foreshadows danger even when the door opens.
-- Missing DSL: encode the specific GM move type and consequence.
scene:gm_spend_fear(1):spotlight("Gandalf")

-- Close the session after the GM move sequence.
scene:end_session()

return scene
