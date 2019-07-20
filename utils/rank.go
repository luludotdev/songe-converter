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
