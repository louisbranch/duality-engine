local scene = Scenario.new("environment_mirkwood_blight_grasping_vines")

-- Model the grasping vines restrain + vulnerable action.
scene:campaign{
  name = "Environment Mirkwood Blight Grasping Vines",
  system = "DAGGERHEART",
  gm_mode = "HUMAN",
  theme = "environment"
}

scene:pc("Gandalf")

-- Vines whip out and bind a target.
scene:start_session("Grasping Vines")

-- Missing DSL: apply Restrained + Vulnerable, escape roll damage, and Hope loss.
scene:reaction_roll{ actor = "Gandalf", trait = "agility", difficulty = 16, outcome = "fear" }
scene:apply_condition{ target = "Gandalf", add = { "RESTRAINED", "VULNERABLE" }, source = "grasping_vines" }

scene:end_session()

return scene
