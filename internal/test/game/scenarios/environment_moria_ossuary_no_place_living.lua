local scene = Scenario.new("environment_moria_ossuary_no_place_living")

-- Model the added Hope cost to clear HP in the ossuary.
scene:campaign{
  name = "Environment Ossuary No Place for the Living",
  system = "DAGGERHEART",
  gm_mode = "HUMAN",
  theme = "environment"
}

scene:pc("Frodo")

-- Healing actions cost extra Hope here.
scene:start_session("No Place for the Living")

-- Missing DSL: apply extra Hope cost to healing or rest effects.
scene:rest{ type = "short", party_size = 1 }

scene:end_session()

return scene
