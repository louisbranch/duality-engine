local scene = Scenario.new("ranged_steady_aim")

-- Model the Ranger of the North's Steady Aim advantage spend.
scene:campaign{
  name = "Ranged Steady Aim",
  system = "DAGGERHEART",
  gm_mode = "HUMAN",
  theme = "adversary"
}

scene:pc("Frodo")
scene:adversary("Ranger of the North")

-- The hunter marks Stress to gain advantage on their next attack.
scene:start_session("Steady Aim")

-- Missing DSL: apply stress spend and advantage die.
scene:adversary_attack{ actor = "Ranger of the North", target = "Frodo", difficulty = 0, damage_type = "physical" }

scene:end_session()

return scene
