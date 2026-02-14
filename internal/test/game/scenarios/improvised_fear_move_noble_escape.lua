local scene = Scenario.new("improvised_fear_move_noble_escape")

-- Model the improvised fear move that lets a villain escape.
scene:campaign{
  name = "Improvised Fear Move Noble Escape",
  system = "DAGGERHEART",
  gm_mode = "HUMAN",
  theme = "gm_fear"
}

scene:pc("Sam")
scene:adversary("Corrupt Steward")

-- The GM spends Fear to remove a near-certain victory.
scene:start_session("Noble Escape")
scene:gm_fear(1)

-- Example: the noble reveals a surprise escape to avoid defeat.
-- Missing DSL: encode the improvised fear move effect.
scene:gm_spend_fear(1):spotlight("Corrupt Steward")

-- Close the session after the escape move.
scene:end_session()

return scene
