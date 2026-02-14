local scene = Scenario.new("orc_archer_opportunist")

-- Highlight damage doubling from the Opportunist feature.
scene:campaign{
  name = "Orc Archer Opportunist",
  system = "DAGGERHEART",
  gm_mode = "HUMAN",
  theme = "damage"
}

scene:pc("Frodo", { armor = 1 })
scene:adversary("Orc Archer")

-- The archer attacks while allies crowd the target.
scene:start_session("Opportunist Shot")

-- Example: 1d8+1 damage is doubled when multiple foes are Very Close.
-- Missing DSL: apply the Opportunist doubling and armor mitigation.
scene:adversary_attack{
  actor = "Orc Archer",
  target = "Frodo",
  difficulty = 0,
  damage_type = "physical"
}

-- Close the session after the opportunist shot.
scene:end_session()

return scene
