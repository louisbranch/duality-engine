local scene = Scenario.new("environment_isengard_ritual_desecrated_ground")

-- Model the Hope die reduction from desecrated ground.
scene:campaign{
  name = "Environment Isengard Ritual Desecrated Ground",
  system = "DAGGERHEART",
  gm_mode = "HUMAN",
  theme = "environment"
}

scene:pc("Frodo")

-- The environment suppresses Hope while the ritual site remains tainted.
scene:start_session("Desecrated Ground")

-- Example: reduce Hope Die to d10 until a progress countdown clears it.
-- Missing DSL: apply Hope die size change and countdown clearance.
scene:countdown_create{ name = "Cleanse Desecration", kind = "progress", current = 0, max = 6, direction = "increase" }

scene:end_session()

return scene
