local scene = Scenario.new("environment_old_forest_grove_barbed_vines")

-- Model the Barbed Vines action in the Old Forest Grove.
scene:campaign{
  name = "Environment Old Forest Grove Barbed Vines",
  system = "DAGGERHEART",
  gm_mode = "HUMAN",
  theme = "environment"
}

scene:pc("Frodo")

-- The grove lashes out with restraining vines.
scene:start_session("Barbed Vines")

-- Example: Agility reaction or take damage and become Restrained.
-- Missing DSL: apply damage, Restrained condition, and escape checks.
scene:reaction_roll{ actor = "Frodo", trait = "agility", difficulty = 11, outcome = "fear" }
scene:apply_condition{ target = "Frodo", add = { "RESTRAINED" }, source = "barbed_vines" }

scene:end_session()

return scene
