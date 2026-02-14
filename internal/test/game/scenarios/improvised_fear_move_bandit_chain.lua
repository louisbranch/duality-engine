local scene = Scenario.new("improvised_fear_move_bandit_chain")

-- Capture the bandit fear-move chain with multiple spotlights.
scene:campaign{
  name = "Improvised Fear Move Orc Chain",
  system = "DAGGERHEART",
  gm_mode = "HUMAN",
  theme = "gm_fear"
}

scene:pc("Sam", { armor = 0 })
scene:adversary("Orc Captain")
scene:adversary("Orc Raider")
scene:adversary("Orc Minions")

-- The GM spends Fear to escalate the bandit ambush.
scene:start_session("Orc Ambush")
scene:gm_fear(5)

-- Example: spotlight Orc Captain with a sudden ambush move.
scene:gm_spend_fear(1):spotlight("Orc Captain")

-- Example: spotlight Orc Raider and swing with a multi-target action.
-- Missing DSL: use Better Surrounded to hit all targets in range.
scene:gm_spend_fear(1):spotlight("Orc Raider")
scene:adversary_attack{ actor = "Orc Raider", target = "Sam", difficulty = 0, damage_type = "physical" }

-- Example: spotlight minions and spend Fear for a group attack.
-- Missing DSL: apply group attack damage to the target.
scene:gm_spend_fear(1):spotlight("Orc Minions")
scene:adversary_attack{ actor = "Orc Minions", target = "Sam", difficulty = 0, damage_type = "physical" }

-- Close the session after the bandit chain.
scene:end_session()

return scene
