local scene = Scenario.new("environment_gondor_court_imperial_decree")

-- Capture the imperial decree ticking a long-term countdown.
scene:campaign{
  name = "Environment Gondor Court Imperial Decree",
  system = "DAGGERHEART",
  gm_mode = "HUMAN",
  theme = "environment"
}

scene:pc("Frodo")

-- The empire accelerates its agenda.
scene:start_session("Imperial Decree")
scene:gm_fear(1)

-- Missing DSL: tick a long-term countdown by 1d4.
scene:gm_spend_fear(1):spotlight("Gondor Court")
scene:countdown_create{ name = "Imperial Agenda", kind = "long_term", current = 0, max = 8, direction = "increase" }

scene:end_session()

return scene
