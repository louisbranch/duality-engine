local scene = Scenario.new("skulk_icicle_barb")

-- Capture the Fell Beast's Icicle Barb group attack.
scene:campaign{
  name = "Skulk Icicle Barb",
  system = "DAGGERHEART",
  gm_mode = "HUMAN",
  theme = "adversary"
}

scene:pc("Frodo")
scene:pc("Sam")
scene:adversary("Fell Beast")

-- The wyrm spends Fear to pin a group with barbs.
scene:start_session("Icicle Barb")
scene:gm_fear(1)

-- Example: targets take 2d4 damage and become Restrained until they break free.
-- Missing DSL: apply group attack resolution and Restrained condition.
scene:gm_spend_fear(1):spotlight("Fell Beast")
scene:apply_condition{ target = "Frodo", add = { "RESTRAINED" }, source = "icicle_barb" }
scene:apply_condition{ target = "Sam", add = { "RESTRAINED" }, source = "icicle_barb" }

scene:end_session()

return scene
