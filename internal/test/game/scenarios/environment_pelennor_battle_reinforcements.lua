local scene = Scenario.new("environment_pelennor_battle_reinforcements")

-- Model reinforcements arriving mid-battle.
scene:campaign{
  name = "Environment Pelennor Battle Reinforcements",
  system = "DAGGERHEART",
  gm_mode = "HUMAN",
  theme = "environment"
}

scene:pc("Frodo")
scene:adversary("Gondor Knight")
scene:adversary("Uruk-hai Minions")

-- A fresh force joins the fight.
scene:start_session("Reinforcements")

-- Missing DSL: spawn new adversaries and spotlight the knight.
scene:gm_spend_fear(0):spotlight("Gondor Knight")

scene:end_session()

return scene
