local scene = Scenario.new("wild_flame_minion_blast")

-- Capture the Wild Flame multi-target spell against minions.
scene:campaign{
  name = "Wild Flame Minion Blast",
  system = "DAGGERHEART",
  gm_mode = "HUMAN",
  theme = "battle"
}

scene:pc("Gandalf")
scene:adversary("Orc Raider A")
scene:adversary("Orc Raider B")
scene:adversary("Nazgul")

-- Gandalf targets two minions and the knight with a single spell.
scene:start_session("Wild Flame")

-- Example: damage roll applies to all targets, minion thresholds delete extras.
-- Missing DSL: apply Minion (4) overflow and stress marking.
scene:multi_attack{
  actor = "Gandalf",
  targets = { "Orc Raider A", "Orc Raider B", "Nazgul" },
  trait = "spellcast",
  difficulty = 0,
  outcome = "fear",
  damage_type = "magic",
  damage_dice = { { sides = 6, count = 2 }, { sides = 10, count = 1 } }
}

-- Close the session after the blast.
scene:end_session()

return scene
