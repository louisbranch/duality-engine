local scene = Scenario.new("temporary_armor_bonus")

-- Echo the temporary armor bonus example tied to a rest.
scene:campaign{
  name = "Temporary Armor Bonus",
  system = "DAGGERHEART",
  gm_mode = "HUMAN",
  theme = "armor"
}

scene:pc("Gandalf", { armor = 3 })

-- Gandalf gains temporary armor until the next rest.
scene:start_session("Armor Bonus")

-- Example: Gandalf's Armor Score increases by 2, then resets on rest.
-- Missing DSL: apply temporary armor bonus and clear Armor Slots on rest.
scene:rest{ type = "short", party_size = 1 }

-- Close the session after the rest clears temporary armor.
scene:end_session()

return scene
