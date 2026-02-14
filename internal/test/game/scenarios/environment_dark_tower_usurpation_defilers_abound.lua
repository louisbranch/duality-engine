local scene = Scenario.new("environment_dark_tower_usurpation_defilers_abound")

-- Capture summoning Orc Shock Troops and their group attack.
scene:campaign{
  name = "Environment Dark Tower Usurpation Defilers Abound",
  system = "DAGGERHEART",
  gm_mode = "HUMAN",
  theme = "environment"
}

scene:pc("Frodo")
scene:adversary("Orc Shock Troops")

-- The usurper calls in shock troops.
scene:start_session("Defilers Abound")
scene:gm_fear(2)

-- Missing DSL: summon 1d4+2 troops and trigger their group attack.
scene:gm_spend_fear(2):spotlight("Orc Shock Troops")

scene:end_session()

return scene
