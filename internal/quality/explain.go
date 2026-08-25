package quality

import (
	"fmt"
	"strings"

	"github.com/jb843051627/fjord-resonance/internal/model"
)

type Explanation struct {
	Decision model.Decision
	Headline string
	Details  []string
}

func Explain(result model.QualityResult) Explanation {
	headline := "quality accepted"
	if result.Decision == model.DecisionReview {
		headline = "quality needs review"
	}
	if result.Decision == model.DecisionReject {
		headline = "quality rejected"
	}
	details := model.CloneReasons(result.Reasons)
	if len(details) == 0 {
		details = []string{"all measured indicators are within the protocol envelope"}
	}
	return Explanation{Decision: result.Decision, Headline: headline, Details: details}
}

func (e Explanation) String() string {
	return fmt.Sprintf("%s: %s", e.Headline, strings.Join(e.Details, "; "))
}
