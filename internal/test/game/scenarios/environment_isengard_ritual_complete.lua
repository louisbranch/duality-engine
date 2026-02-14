local scene = Scenario.new("environment_isengard_ritual_complete")

-- Model the ritual leader's protection reaction.
scene:campaign{
  name = "Environment Isengard Ritual Complete",
  system = "DAGGERHEART",
  gm_mode = "HUMAN",
  theme = "environment"
}

scene:pc("Frodo")
scene:adversary("Saruman")
scene:adversary("Orc Raider")

-- An ally steps in to take a hit meant for the leader.
scene:start_session("Complete the Ritual")

-- Missing DSL: redirect an attack to the ally by marking Stress.
scene:attack{ actor = "Frodo", target = "Saruman", trait = "instinct", difficulty = 0, outcome = "hope", damage_type = "physical" }

scene:end_session()

return scene
