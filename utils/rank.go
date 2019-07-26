package utils

// GetDifficultyRank get difficulty rank number from string
func GetDifficultyRank(difficulty string) int {
	switch difficulty {
	case "Easy":
		return 1
	case "Normal":
		return 3
	case "Hard":
		return 5
	case "Expert":
		return 7
	case "ExpertPlus":
		return 9

	default:
		return 3
	}
}

// GetDifficultyString get difficulty string from rank number
func GetDifficultyString(rank int) string {
	switch rank {
	case 1:
		return "Easy"
	case 3:
		return "Normal"
	case 5:
		return "Hard"
	case 7:
		return "Expert"
	case 9:
		return "ExpertPlus"

	default:
		return "Easy"
	}
}

// GetOldDifficultyRank get old difficulty rank number from new rank
func GetOldDifficultyRank(rank int) int {
	switch rank {
	case 1:
		return 0
	case 3:
		return 1
	case 5:
		return 2
	case 7:
		return 3
	case 9:
		return 4

	default:
		return 1
	}
}
