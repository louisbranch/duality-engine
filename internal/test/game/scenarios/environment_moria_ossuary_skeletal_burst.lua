local scene = Scenario.new("environment_moria_ossuary_skeletal_burst")

-- Model the skeletal burst shrapnel attack.
scene:campaign{
  name = "Environment Ossuary Skeletal Burst",
  system = "DAGGERHEART",
  gm_mode = "HUMAN",
  theme = "environment"
}

scene:pc("Frodo")

-- The ossuary detonates around the party.
scene:start_session("Skeletal Burst")

-- Missing DSL: apply reaction roll and 4d8+8 damage on failure.
scene:reaction_roll{ actor = "Frodo", trait = "agility", difficulty = 19, outcome = "fear" }

scene:end_session()

return scene
