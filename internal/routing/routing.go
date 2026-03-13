package routing

import "strings"

type Decision struct {
	ExecutorAgent string // Satya | Elon | Pixel | Vector | ...
	Reason        string
}

func DecideExecutor(immediate bool, title, desc string) Decision {
	if immediate {
		return Decision{ExecutorAgent: "Satya", Reason: "immediate=true => Satya executes"}
	}

	// v1: if engineering-like => Elon (manager), else Satya
	text := strings.ToLower(title + " " + desc)
	if strings.Contains(text, "api") || strings.Contains(text, "backend") || strings.Contains(text, "frontend") ||
		strings.Contains(text, "deploy") || strings.Contains(text, "infra") || strings.Contains(text, "security") ||
		strings.Contains(text, "bug") || strings.Contains(text, "feature") {
		return Decision{ExecutorAgent: "Elon", Reason: "engineering keywords => delegate to Elon"}
	}
	return Decision{ExecutorAgent: "Satya", Reason: "non-engineering or unknown => Satya handles"}
}
