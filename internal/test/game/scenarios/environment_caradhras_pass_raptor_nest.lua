local scene = Scenario.new("environment_caradhras_pass_raptor_nest")

-- Capture the raptor nest reaction that summons predators.
scene:campaign{
  name = "Environment Caradhras Pass Raptor Nest",
  system = "DAGGERHEART",
  gm_mode = "HUMAN",
  theme = "environment"
}

scene:pc("Frodo")
scene:adversary("Great Eagles")

-- The PCs enter a hunting ground and predators appear.
scene:start_session("Raptor Nest")

-- Missing DSL: spawn two eagles at Very Far range.
scene:gm_spend_fear(0):spotlight("Great Eagles")

scene:end_session()

return scene
