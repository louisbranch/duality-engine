local scene = Scenario.new("improvised_fear_move_shadow")

-- Showcase an improvised fear move that shifts the scene.
scene:campaign{
  name = "Improvised Fear Move Shadow",
  system = "DAGGERHEART",
  gm_mode = "HUMAN",
  theme = "gm_fear"
}

scene:pc("Frodo")

-- The GM spends fear after a success-with-fear to escalate the chase.
scene:start_session("Fear Move")
scene:gm_fear(2)

scene:action_roll{ actor = "Frodo", trait = "instinct", difficulty = 12, outcome = "fear" }

-- Example: the GM spends fear to introduce a looming shadow.
-- Missing DSL: encode the narrative fear move effect.
scene:gm_spend_fear(1):spotlight("Frodo")

-- Close the session after the fear move.
scene:end_session()

return scene
