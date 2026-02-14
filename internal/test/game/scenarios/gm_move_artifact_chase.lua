local scene = Scenario.new("gm_move_artifact_chase")

-- Model a GM move that steals an artifact and launches a chase.
scene:campaign{
  name = "GM Move Artifact Chase",
  system = "DAGGERHEART",
  gm_mode = "HUMAN",
  theme = "gm_move"
}

scene:pc("Gandalf")
scene:adversary("Golum")

-- The GM introduces a chase after a sudden theft.
scene:start_session("Artifact Theft")
scene:gm_fear(1)

-- Example: GM move steals an artifact and forces a chase.
-- Missing DSL: represent item theft and chase trigger.
scene:gm_spend_fear(1):spotlight("Golum")
scene:countdown_create{ name = "Thief Escape", kind = "consequence", current = 0, max = 4, direction = "increase" }

scene:end_session()

return scene
