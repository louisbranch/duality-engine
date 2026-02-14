local scene = Scenario.new("gm_fear_spend_chain")

scene:campaign{
  name = "GM Fear Spend Chain",
  system = "DAGGERHEART",
  gm_mode = "HUMAN",
  theme = "gm_fear"
}

scene:pc("Frodo")
scene:adversary("Nazgul")

-- The GM starts with fear and spends it in quick succession.
scene:start_session("GM Fear")
scene:gm_fear(5)

-- Two spotlight spends show how fear accelerates the GM's cadence.
scene:gm_spend_fear(1):spotlight("Nazgul")
scene:gm_spend_fear(2):spotlight("Nazgul")

scene:end_session()

return scene
