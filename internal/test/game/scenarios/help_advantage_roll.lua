local scene = Scenario.new("help_advantage_roll")

-- Mirror the help action that grants an advantage die.
scene:campaign{
  name = "Help Advantage Roll",
  system = "DAGGERHEART",
  gm_mode = "HUMAN",
  theme = "help"
}

scene:pc("Aragorn")
scene:pc("Gandalf")

-- Gandalf spends Hope to help Aragorn's Instinct roll.
scene:start_session("Help an Ally")

-- Example: help adds an advantage die and a +3 bonus from assistance.
-- Missing DSL: spend Hope on help and apply the advantage die.
scene:action_roll{ actor = "Aragorn", trait = "instinct", difficulty = 10, outcome = "fear" }

-- Close the session after the assisted roll.
scene:end_session()

return scene
