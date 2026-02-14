local scene = Scenario.new("fireball_orc_pack_multi")

-- Capture the fireball example against multiple targets.
scene:campaign{
  name = "Fireball Orc Pack",
  system = "DAGGERHEART",
  gm_mode = "HUMAN",
  theme = "battle"
}

scene:pc("Gandalf")
scene:adversary("Orc Pack A")
scene:adversary("Orc Pack B")

-- Gandalf casts Fireball to catch multiple orc packs at once.
scene:start_session("Fireball")

-- Example: one roll applied to multiple targets.
-- Missing DSL: assert per-target outcomes and damage tiers.
scene:multi_attack{
  actor = "Gandalf",
  targets = { "Orc Pack A", "Orc Pack B" },
  trait = "spellcast",
  difficulty = 0,
  outcome = "hope",
  damage_type = "magic",
  damage_dice = { { sides = 6, count = 2 } }
}

-- Close the session after the multi-target strike.
scene:end_session()

return scene
