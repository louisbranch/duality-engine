local scene = Scenario.new("environment_moria_ossuary_aura_of_death")

-- Capture undead healing from the aura of death.
scene:campaign{
  name = "Environment Ossuary Aura of Death",
  system = "DAGGERHEART",
  gm_mode = "HUMAN",
  theme = "environment"
}

scene:pc("Frodo")
scene:adversary("Uruk-hai")

-- The aura restores undead HP and Stress.
scene:start_session("Aura of Death")

-- Missing DSL: roll d4 and distribute healing across undead.
scene:gm_spend_fear(0):spotlight("Uruk-hai")

scene:end_session()

return scene
