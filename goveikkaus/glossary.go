package goveikkaus

import (
	"fmt"
)

type GlossaryService service

type GameInfo struct {
	AlsoKnownAs string
	Description string
}

type GameGlossary map[string]GameInfo

func (glossary GameGlossary) Print() {
	for key, gameInfo := range glossary {
		fmt.Printf("############# %s - %s #############\n", gameInfo.AlsoKnownAs, key)
		fmt.Printf("API Term:\t\t\t%s\n", key)
		fmt.Printf("Also known as:\t\t\t%s\n", gameInfo.AlsoKnownAs)
		fmt.Println(gameInfo.Description)
		fmt.Println()
	}
}

func (g *GlossaryService) Get() GameGlossary {
	gameGlossary := GameGlossary{
		"FIXEDODDS": GameInfo{
			AlsoKnownAs: "Pitkäveto",
			Description: `
In fixed odds betting (Pitkäveto), you predict winners or outcomes for 1–20 matches. Stakes vary based on match count, sport, or time. Popular sports include soccer and ice hockey.

You can bet individually or use system betting for multiple combinations. Different bet types for the same match can't be combined, except for Build-a-bet. Odds may change before closing, and the recorded odds on the betting slip are final.
			`,
		},
		"MULTISCORE": GameInfo{
			AlsoKnownAs: "Moniveto",
			Description: `
Multi score (Moniveto) is a variable odds betting game based on the number of goals scored in 2–6 matches or other performance outcomes. Minimum bet ranges from 0.05 to 0.20 euros, with a maximum of 100 euros. Final odds can differ significantly from initial ones due to total bet sums influencing them post-game.
			`,
		},
		"SCORE": GameInfo{
			AlsoKnownAs: "Tulosveto",
			Description: `
Result (or score) betting (Tulosveto) involves betting on the number of goals scored by both teams in the target match or other correct outcomes based on performance. Popular sports for result betting include soccer, ice hockey, basketball, and floorball. The bet ranges from 1.00 to 100.00 euros. Result betting is a variable odds game, where the odds are calculated after the game based on the total sum of bets placed on each outcome. Final odds may significantly differ from the initial ones.
			`,
		},
		"SPORT": GameInfo{
			AlsoKnownAs: "Vakio",
			Description: `
In Vakio, you predict the winners of 6–18 matches in regular game time (1=home win, X=draw, 2=away win), or the outcome of a competition between two or three competitors. Winners are selected for each match.

Vakio 1 always has 13 match options, while other Vakios have 6–18 options. Prize categories vary based on the number of matches, and winnings depend on the number of correct predictions. The price per Vakio line ranges from 0.10 to 0.25 euros per match.

Sports in Vakio mainly include soccer, ice hockey, Formula 1, and individual sports.
			`,
		},
		"WINNER": GameInfo{
			AlsoKnownAs: "Voittajaveto",
			Description: `
Win bet (Voittajaveto) involves betting on winners of events, specific correct result combinations, or outcome options. Popular sports include soccer, ice hockey, winter sports, and Formula 1. Bet ranges from 0.20 to 100.00 euros.

It's a variable odds game where the final odds may differ significantly from the purchase odds. Other forms include Perfecta, Trifecta, Daily Double, and Daily Triple.
			`,
		},
		"PICKTWO": GameInfo{
			AlsoKnownAs: "Päivän pari",
			Description: `
In Daily Double aka "Pick two" (Päivän pari), the subject of the bet is the winners of two different competitions or specific defined result combinations or options.
			`,
		},
		"PICKTHREE": GameInfo{
			AlsoKnownAs: "Päivän trio",
			Description: `
In Daily Triple aka "Pick three" (Päivän trio), the subject of the bet is the winners of three different competitions or specific defined result combinations or options.
			`,
		},
		"PERFECTA": GameInfo{
			AlsoKnownAs: "Superkaksari",
			Description: `
In Perfecta (Superkaksari), the subject of the bet is the winner of the competition and the competitor who finishes second in order of superiority.
			`,
		},
		"TRIFECTA": GameInfo{
			AlsoKnownAs: "Supertripla",
			Description: `
In Trifecta (Supertripla), the subject of the bet is the winner, the second-place finisher, and the third-place finisher in order of superiority.
			`,
		},
	}

	return gameGlossary
}
