local scene = Scenario.new("evasion_tie_hit")

-- Highlight that attack totals equal to Evasion still hit.
scene:campaign{
  name = "Evasion Tie Hit",
  system = "DAGGERHEART",
  gm_mode = "HUMAN",
  theme = "battle"
}

scene:pc("Frodo")
scene:adversary("Orc Archer")

-- The archer's attack total ties Frodo's Evasion.
scene:start_session("Tie to Hit")

-- Example: attack total 10 vs Evasion 10 is a hit.
-- Missing DSL: force the adversary attack roll to equal Evasion.
scene:adversary_attack_roll{ actor = "Orc Archer", attack_modifier = 0, advantage = 0, seed = 1 }
scene:apply_adversary_attack_outcome{ targets = { "Frodo" }, difficulty = 10 }

-- Close the session after the tie-hit example.
scene:end_session()

return scene
