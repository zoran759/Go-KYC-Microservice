package verification

import (
	"errors"
	"time"

	"modulus/kyc/common"
)

const (
	docStatusValid   = "SUBMITTED|VALID"
	docStatusInvalid = "SUBMITTED|INVALID"

	notReady = "UNVERIFIED"
)

// Response represents the API response on KYC related requests.
type Response struct {
	ID           string             `json:"_id"`
	Documents    []ResponseDocument `json:"documents"`
	Permission   string             `json:"permission"`
	RefreshToken string             `json:"refresh_token"`
}

// ResponseDocument represents document object from the API response.
type ResponseDocument struct {
	ID              string                `json:"id"`
	PermissionScope string                `json:"permission_scope"`
	VirtualDocs     []ResponseSubDocument `json:"virtual_docs"`
	PhysicalDocs    []ResponseSubDocument `json:"physical_docs"`
}

// ResponseSubDocument represents sub-document object from the API response.
type ResponseSubDocument struct {
	ID          string `json:"id"`
	Type        string `json:"document_type"`
	LastUpdated int64  `json:"last_updated"`
	Status      string `json:"status"`
}

// ToKYCResult processes the response and generates the verification result.
func (r Response) ToKYCResult() (result common.KYCResult, err error) {
	if r.Permission != notReady {
		result.Status = common.Approved
		return
	}

	if len(r.Documents) == 0 {
		err = errors.New("documents for verification are missing, please, supply one")
		return
	}

	reasons := []string{}
	denied := false

	for _, doc := range r.Documents {
		if doc.PermissionScope != notReady {
			continue
		}

		reasons = append(reasons, "Docs set permission: "+doc.PermissionScope)

		for _, vdoc := range doc.VirtualDocs {
			if vdoc.Status == docStatusValid {
				continue
			}
			if vdoc.Status == docStatusInvalid {
				denied = true
			}
			reasons = append(reasons, "Virtual doc | type: "+vdoc.Type+" | status: "+vdoc.Status)
		}

		for _, phdoc := range doc.PhysicalDocs {
			if phdoc.Status == docStatusValid {
				continue
			}
			if phdoc.Status == docStatusInvalid {
				denied = true
			}
			reasons = append(reasons, "Physical doc | type: "+phdoc.Type+" | status: "+phdoc.Status)
		}
	}

	if denied {
		result.Status = common.Denied
		if len(reasons) > 0 {
			result.Details = &common.KYCDetails{
				Reasons: reasons,
			}
		}
		return
	}

	result.Status = common.Unclear
	result.StatusCheck = &common.KYCStatusCheck{
		Provider:    common.SynapseFI,
		ReferenceID: r.ID,
		LastCheck:   time.Now(),
	}

	return
}
