local scene = Scenario.new("minion_group_attack_rats")

-- Model the Minion group attack action from the Giant Rat example.
scene:campaign{
  name = "Minion Group Attack Rats",
  system = "DAGGERHEART",
  gm_mode = "HUMAN",
  theme = "minions"
}

scene:pc("Frodo")
scene:adversary("Moria Rats")

-- The GM spends Fear to trigger a group attack.
scene:start_session("Rat Swarm")
scene:gm_fear(1)

-- Example: shared attack roll, 1 damage each, combined.
-- Missing DSL: resolve group attack damage aggregation.
scene:gm_spend_fear(1):spotlight("Moria Rats")
scene:adversary_attack{ actor = "Moria Rats", target = "Frodo", difficulty = 0, damage_type = "physical" }

scene:end_session()

return scene
