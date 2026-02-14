local scene = Scenario.new("airship_group_roll")

-- Capture the airship crisis that calls for a group roll.
scene:campaign{
  name = "Airship Group Roll",
  system = "DAGGERHEART",
  gm_mode = "HUMAN",
  theme = "teamwork"
}

scene:pc("Sam")
scene:pc("Frodo")
scene:pc("Gandalf")

-- The spellrider breaks the enchantment keeping the airship aloft.
scene:start_session("Airship Crisis")

-- Example: the GM calls for a group roll to keep the airship flying.
-- Missing DSL: map individual outcomes to a shared consequence.
scene:group_action{
  leader = "Sam",
  leader_trait = "presence",
  difficulty = 14,
  supporters = {
    { name = "Frodo", trait = "agility" },
    { name = "Gandalf", trait = "instinct" }
  }
}

-- Close the session after the group roll.
scene:end_session()

return scene
