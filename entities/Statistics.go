package entities

type Statistic struct {
	NumberOfSuccesses int     `json:"number_of_successes"`
	NumberOfFailures  int     `json:"number_of_failures"`
	NumberOfTries     int     `json:"number_of_tries"`
	SuccessPercentage float64 `json:"success_percentage"`
}
