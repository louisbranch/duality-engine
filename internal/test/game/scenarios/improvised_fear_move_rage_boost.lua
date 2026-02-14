local scene = Scenario.new("improvised_fear_move_rage_boost")

-- Model the improvised fear move that boosts a solo adversary's damage.
scene:campaign{
  name = "Improvised Fear Move Rage Boost",
  system = "DAGGERHEART",
  gm_mode = "HUMAN",
  theme = "gm_fear"
}

scene:pc("Frodo")
scene:adversary("Uruk-hai Brute")

-- The GM spends Fear to increase a solo adversary's damage output.
scene:start_session("Rage Boost")
scene:gm_fear(2)

-- Example: the adversary flies into a rage for the remainder of the scene.
-- Missing DSL: apply a temporary damage bonus feature.
scene:gm_spend_fear(1):spotlight("Uruk-hai Brute")

-- Close the session after the fear move.
scene:end_session()

return scene
