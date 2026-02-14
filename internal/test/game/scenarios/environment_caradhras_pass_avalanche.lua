local scene = Scenario.new("environment_caradhras_pass_avalanche")

-- Capture the avalanche action and reaction roll consequences.
scene:campaign{
  name = "Environment Caradhras Pass Avalanche",
  system = "DAGGERHEART",
  gm_mode = "HUMAN",
  theme = "environment"
}

scene:pc("Frodo")

-- The GM triggers an avalanche down the pass.
scene:start_session("Avalanche")
scene:gm_fear(1)

-- Example: reaction roll or take 2d20 damage, knocked to Far, mark Stress.
-- Missing DSL: apply movement and damage severity on failure.
scene:gm_spend_fear(1):spotlight("Caradhras Pass")
scene:reaction_roll{ actor = "Frodo", trait = "agility", difficulty = 15, outcome = "fear" }

scene:end_session()

return scene
