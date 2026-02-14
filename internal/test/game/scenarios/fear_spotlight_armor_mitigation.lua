local scene = Scenario.new("fear_spotlight_armor_mitigation")

-- Recreate a fear-triggered spotlight shift with armor mitigation.
scene:campaign{
  name = "Fear Spotlight Armor Mitigation",
  system = "DAGGERHEART",
  gm_mode = "HUMAN",
  theme = "gm_fear"
}

scene:pc("Frodo", { armor = 1 })
scene:adversary("Uruk-hai Brute")

-- Frodo strikes, the roll lands on Fear, and the GM takes over.
scene:start_session("Spotlight Shift")
scene:gm_fear(6)

scene:attack{
  actor = "Frodo",
  target = "Uruk-hai Brute",
  trait = "instinct",
  difficulty = 0,
  outcome = "fear",
  damage_type = "physical"
}

-- The GM spotlights the adversary breaking free from Vulnerable.
scene:apply_condition{ target = "Uruk-hai Brute", add = { "VULNERABLE" } }
scene:gm_spend_fear(1):spotlight("Uruk-hai Brute")
scene:apply_condition{ target = "Uruk-hai Brute", remove = { "VULNERABLE" }, source = "break_free" }

-- The adversary counterattacks for 9 damage; armor reduces Major to Minor.
-- Missing DSL: set the adversary hit, damage total, and armor slot spend.
scene:adversary_attack{
  actor = "Uruk-hai Brute",
  target = "Frodo",
  difficulty = 0,
  damage_type = "physical"
}

-- Close the session after the spotlight exchange.
scene:end_session()

return scene
