local scene = Scenario.new("environment_helms_deep_siege_collateral_damage")

-- Model collateral damage from siege weapons after an adversary falls.
scene:campaign{
  name = "Environment Helms Deep Siege Collateral Damage",
  system = "DAGGERHEART",
  gm_mode = "HUMAN",
  theme = "environment"
}

scene:pc("Frodo")

-- A stray attack lands where the fight rages.
scene:start_session("Collateral Damage")
scene:gm_fear(1)

-- Missing DSL: apply reaction roll and damage/stress outcomes.
scene:gm_spend_fear(1):spotlight("Helms Deep Siege")
scene:reaction_roll{ actor = "Frodo", trait = "agility", difficulty = 17, outcome = "fear" }

scene:end_session()

return scene
