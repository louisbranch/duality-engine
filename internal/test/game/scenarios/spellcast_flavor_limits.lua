local scene = Scenario.new("spellcast_flavor_limits")

-- Capture the example where flavor doesn't grant extra effects.
scene:campaign{
  name = "Spellcast Flavor Limits",
  system = "DAGGERHEART",
  gm_mode = "HUMAN",
  theme = "spellcast"
}

scene:pc("Gandalf")
scene:adversary("Saruman")

-- Flavoring a warding circle doesn't add extra damage.
scene:start_session("Flavor Limits")

-- Missing DSL: enforce that narration doesn't modify damage.
scene:action_roll{ actor = "Gandalf", trait = "spellcast", difficulty = 12, outcome = "hope" }

scene:end_session()

return scene
