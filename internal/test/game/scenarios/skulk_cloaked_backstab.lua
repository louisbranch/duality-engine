local scene = Scenario.new("skulk_cloaked_backstab")

-- Model a Skulk using Cloaked to set up a Backstab.
scene:campaign{
  name = "Skulk Cloaked Backstab",
  system = "DAGGERHEART",
  gm_mode = "HUMAN",
  theme = "adversary"
}

scene:pc("Frodo")
scene:adversary("Orc Stalker")

-- The shadow hides, then strikes with advantage for boosted damage.
scene:start_session("Cloaked Backstab")

-- Example: Cloaked grants Hidden; Backstab replaces damage on advantaged hit.
-- Missing DSL: apply Hidden and swap in 1d6+6 damage on advantage.
scene:apply_condition{ target = "Orc Stalker", add = { "HIDDEN" }, source = "cloaked" }
scene:adversary_attack{ actor = "Orc Stalker", target = "Frodo", difficulty = 0, damage_type = "physical" }

scene:end_session()

return scene
