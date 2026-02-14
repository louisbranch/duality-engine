local scene = Scenario.new("environment_bree_market_unexpected_find")

-- Model the marketplace action that reveals a needed item.
scene:campaign{
  name = "Environment Bree Market Unexpected Find",
  system = "DAGGERHEART",
  gm_mode = "HUMAN",
  theme = "environment"
}

scene:pc("Frodo")
scene:npc("Bilbo")

-- A merchant reveals a rare or desired item.
scene:start_session("Unexpected Find")

-- Missing DSL: introduce a quest item and its non-gold cost.
scene:gm_spend_fear(0):spotlight("Bilbo")

scene:end_session()

return scene
