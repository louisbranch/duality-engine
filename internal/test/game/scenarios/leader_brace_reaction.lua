local scene = Scenario.new("leader_brace_reaction")

-- Capture the Mirkwood Warden Brace reaction reducing HP loss.
scene:campaign{
  name = "Leader Brace Reaction",
  system = "DAGGERHEART",
  gm_mode = "HUMAN",
  theme = "adversary"
}

scene:pc("Frodo")
scene:adversary("Mirkwood Warden")

-- The Mirkwood Warden marks Stress to reduce HP marked.
scene:start_session("Brace")

-- Example: when the Mirkwood Warden marks HP, they can mark Stress to mark 1 fewer.
-- Missing DSL: reduce HP loss and spend Stress on reaction.
scene:attack{
  actor = "Frodo",
  target = "Mirkwood Warden",
  trait = "instinct",
  difficulty = 0,
  outcome = "hope",
  damage_type = "physical"
}

scene:end_session()

return scene
