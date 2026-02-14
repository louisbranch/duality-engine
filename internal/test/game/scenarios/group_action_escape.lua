local scene = Scenario.new("group_action_escape")

-- Recreate the collapsing stronghold group action example.
scene:campaign{
  name = "Group Action Escape",
  system = "DAGGERHEART",
  gm_mode = "HUMAN",
  theme = "teamwork"
}

scene:pc("Sam")
scene:pc("Frodo")
scene:pc("Gandalf")

-- The party bolts for the exit as the stronghold comes down.
scene:start_session("Escape")

-- Sam leads the group action to remember the way out.
-- Missing DSL: assert each participant outcome and any fear/hope changes.
scene:group_action{
  leader = "Sam",
  leader_trait = "instinct",
  difficulty = 12,
  supporters = {
    { name = "Frodo", trait = "presence" },
    { name = "Gandalf", trait = "agility" }
  }
}

-- Close the session after the escape attempt.
scene:end_session()

return scene
