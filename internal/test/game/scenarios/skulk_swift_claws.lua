local scene = Scenario.new("skulk_swift_claws")

-- Model the Fell Beast's Swift Claws leap-and-strike action.
scene:campaign{
  name = "Skulk Swift Claws",
  system = "DAGGERHEART",
  gm_mode = "HUMAN",
  theme = "adversary"
}

scene:pc("Frodo")
scene:adversary("Fell Beast")

-- The wyrm marks Stress to dash in and strike.
scene:start_session("Swift Claws")

-- Example: on hit, deal 2d10+5 and force a Strength reaction to avoid knockback.
-- Missing DSL: apply movement, stress spend, and knockback.
scene:adversary_attack{ actor = "Fell Beast", target = "Frodo", difficulty = 0, damage_type = "physical" }

scene:end_session()

return scene
