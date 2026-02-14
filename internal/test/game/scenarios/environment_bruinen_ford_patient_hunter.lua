local scene = Scenario.new("environment_bruinen_ford_patient_hunter")

-- Capture the river's Patient Hunter fear action summoning a predator.
scene:campaign{
  name = "Environment Bruinen Ford Patient Hunter",
  system = "DAGGERHEART",
  gm_mode = "HUMAN",
  theme = "environment"
}

scene:pc("Frodo")
scene:adversary("Warg")

-- The GM spends Fear to summon a Warg.
scene:start_session("Patient Hunter")
scene:gm_fear(1)

-- Example: summon the Warg within Close range and immediately spotlight it.
-- Missing DSL: place the adversary and trigger its action.
scene:gm_spend_fear(1):spotlight("Warg")

scene:end_session()

return scene
