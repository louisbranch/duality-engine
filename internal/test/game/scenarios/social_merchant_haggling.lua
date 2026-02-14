local scene = Scenario.new("social_merchant_haggling")

-- Model social mechanics for haggling with a merchant.
scene:campaign{
  name = "Social Bilbo Haggling",
  system = "DAGGERHEART",
  gm_mode = "HUMAN",
  theme = "social"
}

scene:pc("Frodo")
scene:adversary("Bree Merchant")

-- The merchant rewards success and penalizes poor rolls.
scene:start_session("Haggling")

-- Example: success grants discounts, failure adds stress and disadvantage.
-- Missing DSL: apply Preferential Treatment and The Runaround effects.
scene:action_roll{ actor = "Frodo", trait = "presence", difficulty = 12, outcome = "fear" }

scene:end_session()

return scene
