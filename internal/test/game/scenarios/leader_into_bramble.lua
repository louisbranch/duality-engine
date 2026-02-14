local scene = Scenario.new("leader_into_bramble")

-- Model the Mirkwood Warden's Into the Bramble fear action.
scene:campaign{
  name = "Leader Into the Bramble",
  system = "DAGGERHEART",
  gm_mode = "HUMAN",
  theme = "adversary"
}

scene:pc("Frodo")
scene:adversary("Mirkwood Warden")
scene:adversary("Mirkwood Archer")

-- The leader spends Fear to reposition allies and hide them.
scene:start_session("Bramble Ambush")
scene:gm_fear(1)

-- Example: spotlight up to 1d4 allies and grant Hidden.
-- Missing DSL: move allies to cover and apply Hidden until they attack.
scene:gm_spend_fear(1):spotlight("Mirkwood Warden")
scene:apply_condition{ target = "Mirkwood Archer", add = { "HIDDEN" }, source = "into_bramble" }

scene:end_session()

return scene
