local scene = Scenario.new("environment_misty_ascent_fall")

-- Model the fall action that escalates damage by countdown state.
scene:campaign{
  name = "Environment Misty Ascent Fall",
  system = "DAGGERHEART",
  gm_mode = "HUMAN",
  theme = "environment"
}

scene:pc("Frodo")

-- A handhold fails, risking a deadly fall.
scene:start_session("Misty Fall")
scene:gm_fear(1)

-- Example: spend Fear, if not saved next action, damage scales by countdown.
-- Missing DSL: defer the damage until a follow-up action fails to save.
scene:gm_spend_fear(1):spotlight("Misty Ascent")

scene:end_session()

return scene
