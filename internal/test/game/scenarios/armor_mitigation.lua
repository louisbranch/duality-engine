local scene = Scenario.new("armor_mitigation")

scene:campaign{
  name = "Armor Mitigation",
  system = "DAGGERHEART",
  gm_mode = "HUMAN",
  theme = "armor"
}

scene:pc("Frodo", { armor = 1 })
scene:adversary("Nazgul")

-- Nazgul pressures Frodo while the GM holds fear to power the assault.
scene:start_session("Armor")
scene:gm_fear(2)

-- Nazgul lands a hit; Frodo is expected to mitigate with armor.
-- Missing DSL: specify damage roll totals and assert armor slot spend/HP loss.
scene:adversary_attack{
  actor = "Nazgul",
  target = "Frodo",
  difficulty = 0,
  damage_type = "physical"
}

scene:end_session()

return scene
