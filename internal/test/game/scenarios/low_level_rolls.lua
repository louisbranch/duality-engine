local scene = Scenario.new("low_level_rolls")

-- Set a low-stakes skirmish to exercise roll plumbing.
scene:campaign{ name = "Low Level Rolls", system = "DAGGERHEART", gm_mode = "HUMAN" }
scene:pc("Sam")
scene:adversary("Goblin")

-- Open a session to run through low-level rolls.
scene:start_session("Low Level")

-- Resolve Sam's action roll and apply the attack outcome.
scene:action_roll{ actor = "Sam", trait = "instinct", difficulty = 8, outcome = "hope" }
scene:apply_roll_outcome{}
scene:apply_attack_outcome{ targets = { "Goblin" } }

-- Roll damage for Sam's attack.
scene:damage_roll{ actor = "Sam", damage_dice = { { sides = 6, count = 2 } }, modifier = 1, seed = 21 }

-- Resolve the goblin's counterattack and apply the outcome.
scene:adversary_attack_roll{ actor = "Goblin", attack_modifier = 1, advantage = 1, seed = 33 }
scene:apply_adversary_attack_outcome{ targets = { "Sam" }, difficulty = 10 }

-- Resolve Sam's reaction roll and apply the outcome.
scene:reaction_roll{ actor = "Sam", trait = "agility", difficulty = 8, outcome = "hope" }
scene:apply_reaction_outcome{}

-- Close the session once the roll chain completes.
scene:end_session()
return scene
