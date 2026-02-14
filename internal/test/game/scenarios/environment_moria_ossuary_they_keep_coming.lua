local scene = Scenario.new("environment_moria_ossuary_they_keep_coming")

-- Model the undead reinforcements action in the ossuary.
scene:campaign{
  name = "Environment Ossuary They Keep Coming",
  system = "DAGGERHEART",
  gm_mode = "HUMAN",
  theme = "environment"
}

scene:pc("Frodo")
scene:adversary("Orc Rabble")
scene:adversary("Uruk-hai")

-- The necromancer calls in more undead.
scene:start_session("They Just Keep Coming")
scene:gm_fear(1)

-- Missing DSL: summon 1d6 rotted zombies, two perfected, or a legion.
scene:gm_spend_fear(1):spotlight("Orc Rabble")

scene:end_session()

return scene
