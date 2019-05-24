package main

func checkForMapExt(diff *OldDifficultyJSON) bool {
	for _, note := range diff.Notes {
		if note.LineLayer >= 1000 || note.LineLayer <= -1000 {
			return true
		} else if note.LineLayer >= 1000 || note.LineIndex <= -1000 {
			return true
		} else if note.CutDirection >= 1000 || note.CutDirection <= -1000 {
			return true
		}
	}

	return false
}
