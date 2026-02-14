local scene = Scenario.new("ranged_eruption_hazard")

-- Model the Saruman's eruption hazard action.
scene:campaign{
  name = "Ranged Eruption Hazard",
  system = "DAGGERHEART",
  gm_mode = "HUMAN",
  theme = "adversary"
}

scene:pc("Frodo")
scene:adversary("Saruman")

-- The wizard spends Fear to erupt terrain and force reaction rolls.
scene:start_session("Eruption")
scene:gm_fear(1)

-- Example: targets roll Agility 14 or take 2d10 damage and are moved.
-- Missing DSL: apply area hazard, reaction roll, and forced movement.
scene:gm_spend_fear(1):spotlight("Saruman")

scene:end_session()

return scene
