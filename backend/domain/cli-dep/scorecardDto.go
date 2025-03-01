package cli_dep

type ScoreCard struct {
	OverallScore float32 `json:"overallScore"`
}
type ScoreCardDto struct {
	ScoreCard `json:"scorecard"`
}
