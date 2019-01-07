package complyadvantage

import (
	"testing"

	"modulus/kyc/common"

	"github.com/stretchr/testify/assert"
)

func TestToResult(t *testing.T) {
	r := Response{
		Content: &Content{
			Data: Data{
				ID:        123,
				TotalHits: 1,
				Hits: []Hit{
					Hit{
						Doc: Doc{
							EntityType: "fake_type",
						},
					},
				},
			},
		},
	}

	res, err := r.toResult()

	assert := assert.New(t)

	assert.Nil(err)
	assert.Equal(common.Denied, res.Status)
	assert.NotNil(res.Details)
	assert.Equal(common.Unknown, res.Details.Finality)
	assert.Len(res.Details.Reasons, 2)
	assert.Equal("Search ID: 123", res.Details.Reasons[0])
	assert.Equal("Possible false positive. Please, inspect case details on the ComplyAdvantage site.", res.Details.Reasons[1])
	assert.Empty(res.ErrorCode)
	assert.Nil(res.StatusCheck)

}
