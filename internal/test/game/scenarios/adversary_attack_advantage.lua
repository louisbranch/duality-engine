local scene = Scenario.new("adversary_attack_advantage")

scene:campaign{
  name = "Adversary Attack Advantage",
  system = "DAGGERHEART",
  gm_mode = "HUMAN",
  theme = "ambush"
}

scene:pc("Frodo", { armor = 1 })
scene:adversary("Saruman")

-- Saruman ambushes Frodo with a clear edge.
scene:start_session("Ambush")

-- Advantage and an attack modifier tilt the roll in Saruman's favor.
-- Missing DSL: specify damage rolls and armor/HP consequences.
scene:adversary_attack{
  actor = "Saruman",
  target = "Frodo",
  difficulty = 0,
  attack_modifier = 2,
  advantage = 1,
  damage_type = "physical"
}

scene:end_session()

return scene
