local scene = Scenario.new("leader_ferocious_defense")

-- Model Ferocious Defense increasing Difficulty after taking HP damage.
scene:campaign{
  name = "Leader Ferocious Defense",
  system = "DAGGERHEART",
  gm_mode = "HUMAN",
  theme = "adversary"
}

scene:pc("Frodo")
scene:adversary("Mirkwood Warden")

-- The Mirkwood Warden hardens after a damaging hit.
scene:start_session("Ferocious Defense")

-- Example: after marking HP, Difficulty increases by 1 until they mark HP.
-- Missing DSL: apply Difficulty increase after HP loss.
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
