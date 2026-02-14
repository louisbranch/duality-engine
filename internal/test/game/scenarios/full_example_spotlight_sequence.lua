local scene = Scenario.new("full_example_spotlight_sequence")

-- Follow the example-of-play spotlight order across multiple adversaries.
scene:campaign{
  name = "Full Example Spotlight Sequence",
  system = "DAGGERHEART",
  gm_mode = "HUMAN",
  theme = "spotlight"
}

scene:pc("Sam")
scene:pc("Frodo")
scene:pc("Gandalf")
scene:pc("Aragorn")
scene:adversary("Orc Archer One")
scene:adversary("Orc Archer Two")
scene:adversary("Nazgul")
scene:adversary("Orc Raiders")

-- The GM chains spotlights as threats activate in sequence.
scene:start_session("Spotlight Sequence")
scene:gm_fear(4)

-- Example: archers fire, dredges swarm, then the knight takes center stage.
scene:gm_spend_fear(1):spotlight("Orc Archer One")
scene:gm_spend_fear(1):spotlight("Orc Archer Two")
scene:gm_spend_fear(1):spotlight("Orc Raiders")
scene:gm_spend_fear(1):spotlight("Nazgul")

-- Close the session after the spotlight chain resolves.
scene:end_session()

return scene
