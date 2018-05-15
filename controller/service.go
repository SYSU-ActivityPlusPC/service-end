package controller

const secret = "sysu_activity_2018_activity_sysu"

// GetPoster judge if the poster and returns accurate one with given type
func GetPoster(raw string, actType int) string {
	if len(raw) == 0 {
		switch actType {
		// physics
		case 0:
			return "b6f487c6d08921463a6ebc0612d9fe1f.gif"
		// volunteer
		case 1:
			return "ccc55f553829fabb7c15227d79450dae.gif"
		// match
		case 2:
			return "2bee829b10b0a84002cf5cb5c4a3c8f3.gif"
		// show
		case 3:
			return "68dac067d05a98995a353ad8265b1f09.png"
		// speech
		case 4:
			return "a90dc26fbd5299e4053a3bbc39b5afc8.gif"
		// outdoor
		case 5:
			return "e8ae3078dfa14c62ff1e71104ec0b11f.png"
		// relax
		case 6:
			return "b2b71f5f39d3a4389d34ce1b248e9fee.png"
		}
	}
	return raw
}
