local scene = Scenario.new("environment_dark_tower_usurpation_ritual_nexus")

-- Capture the ritual backlash on failures with Fear.
scene:campaign{
  name = "Environment Dark Tower Usurpation Ritual Nexus",
  system = "DAGGERHEART",
  gm_mode = "HUMAN",
  theme = "environment"
}

scene:pc("Frodo")
scene:adversary("Saruman")

-- A Fear failure triggers magical backlash.
scene:start_session("Ritual Nexus")

-- Missing DSL: roll 1d4 stress on failure with Fear.
scene:action_roll{ actor = "Frodo", trait = "presence", difficulty = 20, outcome = "fear" }

scene:end_session()

return scene
