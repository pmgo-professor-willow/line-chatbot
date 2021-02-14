package utils

type Type int

const (
	NONE Type = iota
	BUG
	DARK
	DRAGON
	ELECTRIC
	FAIRY
	FIGHTING
	FIRE
	FLYING
	GHOST
	GRASS
	GROUND
	ICE
	NORMAL
	POISON
	PSYCHIC
	ROCK
	STEEL
	WATER
)

var Types = []Type{
	BUG,
	DARK,
	DRAGON,
	ELECTRIC,
	FAIRY,
	FIGHTING,
	FIRE,
	FLYING,
	GHOST,
	GRASS,
	GROUND,
	ICE,
	NORMAL,
	POISON,
	PSYCHIC,
	ROCK,
	STEEL,
	WATER,
}

type PokemonTypeRateTable map[Type]float32

type PokemonTypeEffectivenessTable map[Type]PokemonTypeRateTable

var typeEffectivenessTable = PokemonTypeEffectivenessTable{
	BUG: PokemonTypeRateTable{
		BUG:      1,
		DARK:     1.6,
		DRAGON:   1,
		ELECTRIC: 1,
		FAIRY:    0.625,
		FIGHTING: 0.625,
		FIRE:     0.625,
		FLYING:   0.625,
		GHOST:    0.625,
		GRASS:    1.6,
		GROUND:   1,
		ICE:      1,
		NORMAL:   1,
		POISON:   0.625,
		PSYCHIC:  1.6,
		ROCK:     1,
		STEEL:    0.625,
		WATER:    1,
	},
	DARK: PokemonTypeRateTable{
		BUG:      1,
		DARK:     0.625,
		DRAGON:   1,
		ELECTRIC: 1,
		FAIRY:    0.625,
		FIGHTING: 0.625,
		FIRE:     1,
		FLYING:   1,
		GHOST:    1.6,
		GRASS:    1,
		GROUND:   1,
		ICE:      1,
		NORMAL:   1,
		POISON:   1,
		PSYCHIC:  1.6,
		ROCK:     1,
		STEEL:    1,
		WATER:    1,
	},
	DRAGON: PokemonTypeRateTable{
		BUG:      1,
		DARK:     1,
		DRAGON:   1.6,
		ELECTRIC: 1,
		FAIRY:    0.390625,
		FIGHTING: 1,
		FIRE:     1,
		FLYING:   1,
		GHOST:    1,
		GRASS:    1,
		GROUND:   1,
		ICE:      1,
		NORMAL:   1,
		POISON:   1,
		PSYCHIC:  1,
		ROCK:     1,
		STEEL:    0.625,
		WATER:    1,
	},
	ELECTRIC: PokemonTypeRateTable{
		BUG:      1,
		DARK:     1,
		DRAGON:   0.625,
		ELECTRIC: 0.625,
		FAIRY:    1,
		FIGHTING: 1,
		FIRE:     1,
		FLYING:   1.6,
		GHOST:    1,
		GRASS:    0.625,
		GROUND:   0.390625,
		ICE:      1,
		NORMAL:   1,
		POISON:   1,
		PSYCHIC:  1,
		ROCK:     1,
		STEEL:    1,
		WATER:    1.6,
	},
	FAIRY: PokemonTypeRateTable{
		BUG:      1,
		DARK:     1.6,
		DRAGON:   1.6,
		ELECTRIC: 1,
		FAIRY:    1,
		FIGHTING: 1.6,
		FIRE:     0.625,
		FLYING:   1,
		GHOST:    1,
		GRASS:    1,
		GROUND:   1,
		ICE:      1,
		NORMAL:   1,
		POISON:   0.625,
		PSYCHIC:  1,
		ROCK:     1,
		STEEL:    0.625,
		WATER:    1,
	},
	FIGHTING: PokemonTypeRateTable{
		BUG:      0.625,
		DARK:     1.6,
		DRAGON:   1,
		ELECTRIC: 1,
		FAIRY:    0.625,
		FIGHTING: 1,
		FIRE:     1,
		FLYING:   0.625,
		GHOST:    0.390625,
		GRASS:    1,
		GROUND:   1,
		ICE:      1.6,
		NORMAL:   1.6,
		POISON:   0.625,
		PSYCHIC:  0.625,
		ROCK:     1.6,
		STEEL:    1.6,
		WATER:    1,
	},
	FIRE: PokemonTypeRateTable{
		BUG:      1.6,
		DARK:     1,
		DRAGON:   0.625,
		ELECTRIC: 1,
		FAIRY:    1,
		FIGHTING: 1,
		FIRE:     0.625,
		FLYING:   1,
		GHOST:    1,
		GRASS:    1.6,
		GROUND:   1,
		ICE:      1.6,
		NORMAL:   1,
		POISON:   1,
		PSYCHIC:  1,
		ROCK:     0.625,
		STEEL:    1.6,
		WATER:    0.625,
	},
	FLYING: PokemonTypeRateTable{
		BUG:      1.6,
		DARK:     1,
		DRAGON:   1,
		ELECTRIC: 0.625,
		FAIRY:    1,
		FIGHTING: 1.6,
		FIRE:     1,
		FLYING:   1,
		GHOST:    1,
		GRASS:    1.6,
		GROUND:   1,
		ICE:      1,
		NORMAL:   1,
		POISON:   1,
		PSYCHIC:  1,
		ROCK:     0.625,
		STEEL:    0.625,
		WATER:    1,
	},
	GHOST: PokemonTypeRateTable{
		BUG:      1,
		DARK:     0.625,
		DRAGON:   1,
		ELECTRIC: 1,
		FAIRY:    1,
		FIGHTING: 1,
		FIRE:     1,
		FLYING:   1,
		GHOST:    1.6,
		GRASS:    1,
		GROUND:   1,
		ICE:      1,
		NORMAL:   0.390625,
		POISON:   1,
		PSYCHIC:  1.6,
		ROCK:     1,
		STEEL:    1,
		WATER:    1,
	},
	GRASS: PokemonTypeRateTable{
		BUG:      0.625,
		DARK:     1,
		DRAGON:   0.625,
		ELECTRIC: 1,
		FAIRY:    1,
		FIGHTING: 1,
		FIRE:     0.625,
		FLYING:   0.625,
		GHOST:    1,
		GRASS:    0.625,
		GROUND:   1.6,
		ICE:      1,
		NORMAL:   1,
		POISON:   0.625,
		PSYCHIC:  1,
		ROCK:     1.6,
		STEEL:    0.625,
		WATER:    1.6,
	},
	GROUND: PokemonTypeRateTable{
		BUG:      0.625,
		DARK:     1,
		DRAGON:   1,
		ELECTRIC: 1.6,
		FAIRY:    1,
		FIGHTING: 1,
		FIRE:     1.6,
		FLYING:   0.390625,
		GHOST:    1,
		GRASS:    0.625,
		GROUND:   1,
		ICE:      1,
		NORMAL:   1,
		POISON:   1.6,
		PSYCHIC:  1,
		ROCK:     1.6,
		STEEL:    1.6,
		WATER:    1,
	},
	ICE: PokemonTypeRateTable{
		BUG:      1,
		DARK:     1,
		DRAGON:   1.6,
		ELECTRIC: 1,
		FAIRY:    1,
		FIGHTING: 1,
		FIRE:     0.625,
		FLYING:   1.6,
		GHOST:    1,
		GRASS:    1.6,
		GROUND:   1.6,
		ICE:      0.625,
		NORMAL:   1,
		POISON:   1,
		PSYCHIC:  1,
		ROCK:     1,
		STEEL:    0.625,
		WATER:    0.625,
	},
	NORMAL: PokemonTypeRateTable{
		BUG:      1,
		DARK:     1,
		DRAGON:   1,
		ELECTRIC: 1,
		FAIRY:    1,
		FIGHTING: 1,
		FIRE:     1,
		FLYING:   1,
		GHOST:    0.390625,
		GRASS:    1,
		GROUND:   1,
		ICE:      1,
		NORMAL:   1,
		POISON:   1,
		PSYCHIC:  1,
		ROCK:     0.625,
		STEEL:    0.625,
		WATER:    1,
	},
	POISON: PokemonTypeRateTable{
		BUG:      1,
		DARK:     1,
		DRAGON:   1,
		ELECTRIC: 1,
		FAIRY:    1.6,
		FIGHTING: 1,
		FIRE:     1,
		FLYING:   1,
		GHOST:    0.625,
		GRASS:    1.6,
		GROUND:   0.625,
		ICE:      1,
		NORMAL:   1,
		POISON:   0.625,
		PSYCHIC:  1,
		ROCK:     0.625,
		STEEL:    0.390625,
		WATER:    1,
	},
	PSYCHIC: PokemonTypeRateTable{
		BUG:      1,
		DARK:     0.390625,
		DRAGON:   1,
		ELECTRIC: 1,
		FAIRY:    1,
		FIGHTING: 1.6,
		FIRE:     1,
		FLYING:   1,
		GHOST:    1,
		GRASS:    1,
		GROUND:   1,
		ICE:      1,
		NORMAL:   1,
		POISON:   1.6,
		PSYCHIC:  0.625,
		ROCK:     1,
		STEEL:    0.625,
		WATER:    1,
	},
	ROCK: PokemonTypeRateTable{
		BUG:      1.6,
		DARK:     1,
		DRAGON:   1,
		ELECTRIC: 1,
		FAIRY:    1,
		FIGHTING: 0.625,
		FIRE:     1.6,
		FLYING:   1.6,
		GHOST:    1,
		GRASS:    1,
		GROUND:   0.625,
		ICE:      1.6,
		NORMAL:   1,
		POISON:   1,
		PSYCHIC:  1,
		ROCK:     1,
		STEEL:    0.625,
		WATER:    1,
	},
	STEEL: PokemonTypeRateTable{
		BUG:      1,
		DARK:     1,
		DRAGON:   1,
		ELECTRIC: 0.625,
		FAIRY:    1.6,
		FIGHTING: 1,
		FIRE:     0.625,
		FLYING:   1,
		GHOST:    1,
		GRASS:    1,
		GROUND:   1,
		ICE:      1.6,
		NORMAL:   1,
		POISON:   1,
		PSYCHIC:  1,
		ROCK:     1.6,
		STEEL:    0.625,
		WATER:    0.625,
	},
	WATER: PokemonTypeRateTable{
		BUG:      1,
		DARK:     1,
		DRAGON:   0.625,
		ELECTRIC: 1,
		FAIRY:    1,
		FIGHTING: 1,
		FIRE:     1.6,
		FLYING:   1,
		GHOST:    1,
		GRASS:    0.625,
		GROUND:   1.6,
		ICE:      1,
		NORMAL:   1,
		POISON:   1,
		PSYCHIC:  1,
		ROCK:     1.6,
		STEEL:    1,
		WATER:    0.625,
	},
}

func GetWeaknessTypes(selfTypes []string, minRate, maxRate float32) []string {
	// Init
	finalRateTable := make(map[Type]float32)
	for _, typeEnum := range Types {
		finalRateTable[typeEnum] = 1
	}

	for _, selfType := range selfTypes {
		defenderTypeEnum := NONE

		switch selfType {
		case "蟲":
			defenderTypeEnum = BUG
			break
		case "惡":
			defenderTypeEnum = DARK
			break
		case "龍":
			defenderTypeEnum = DRAGON
			break
		case "電":
			defenderTypeEnum = ELECTRIC
			break
		case "妖精":
			defenderTypeEnum = FAIRY
			break
		case "火":
			defenderTypeEnum = FIRE
			break
		case "格鬥":
			defenderTypeEnum = FIGHTING
			break
		case "飛行":
			defenderTypeEnum = FLYING
			break
		case "幽靈":
			defenderTypeEnum = GHOST
			break
		case "草":
			defenderTypeEnum = GRASS
			break
		case "地面":
			defenderTypeEnum = GROUND
			break
		case "冰":
			defenderTypeEnum = ICE
			break
		case "一般":
			defenderTypeEnum = NORMAL
			break
		case "毒":
			defenderTypeEnum = POISON
			break
		case "超能力":
			defenderTypeEnum = PSYCHIC
			break
		case "岩石":
			defenderTypeEnum = ROCK
			break
		case "鋼":
			defenderTypeEnum = STEEL
			break
		case "水":
			defenderTypeEnum = WATER
			break
		}

		for _, attackerTypeEnum := range Types {
			finalRateTable[attackerTypeEnum] *= typeEffectivenessTable[attackerTypeEnum][defenderTypeEnum]
		}
	}

	results := []string{}

	if (finalRateTable[BUG] >= minRate) && (finalRateTable[BUG] <= maxRate) {
		results = append(results, "蟲")
	}
	if (finalRateTable[DARK] >= minRate) && (finalRateTable[DARK] <= maxRate) {
		results = append(results, "惡")
	}
	if (finalRateTable[DRAGON] >= minRate) && (finalRateTable[DRAGON] <= maxRate) {
		results = append(results, "龍")
	}
	if (finalRateTable[ELECTRIC] >= minRate) && (finalRateTable[ELECTRIC] <= maxRate) {
		results = append(results, "電")
	}
	if (finalRateTable[FAIRY] >= minRate) && (finalRateTable[FAIRY] <= maxRate) {
		results = append(results, "妖精")
	}
	if (finalRateTable[FIGHTING] >= minRate) && (finalRateTable[FIGHTING] <= maxRate) {
		results = append(results, "格鬥")
	}
	if (finalRateTable[FIRE] >= minRate) && (finalRateTable[FIRE] <= maxRate) {
		results = append(results, "火")
	}
	if (finalRateTable[FLYING] >= minRate) && (finalRateTable[FLYING] <= maxRate) {
		results = append(results, "飛行")
	}
	if (finalRateTable[GHOST] >= minRate) && (finalRateTable[GHOST] <= maxRate) {
		results = append(results, "幽靈")
	}
	if (finalRateTable[GRASS] >= minRate) && (finalRateTable[GRASS] <= maxRate) {
		results = append(results, "草")
	}
	if (finalRateTable[GROUND] >= minRate) && (finalRateTable[GROUND] <= maxRate) {
		results = append(results, "地面")
	}
	if (finalRateTable[ICE] >= minRate) && (finalRateTable[ICE] <= maxRate) {
		results = append(results, "冰")
	}
	if (finalRateTable[NORMAL] >= minRate) && (finalRateTable[NORMAL] <= maxRate) {
		results = append(results, "一般")
	}
	if (finalRateTable[POISON] >= minRate) && (finalRateTable[POISON] <= maxRate) {
		results = append(results, "毒")
	}
	if (finalRateTable[PSYCHIC] >= minRate) && (finalRateTable[PSYCHIC] <= maxRate) {
		results = append(results, "超能力")
	}
	if (finalRateTable[ROCK] >= minRate) && (finalRateTable[ROCK] <= maxRate) {
		results = append(results, "岩石")
	}
	if (finalRateTable[STEEL] >= minRate) && (finalRateTable[STEEL] <= maxRate) {
		results = append(results, "鋼")
	}
	if (finalRateTable[WATER] >= minRate) && (finalRateTable[WATER] <= maxRate) {
		results = append(results, "水")
	}

	return results
}
