package expectid

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"modulus/kyc/common"
)

var _ = Describe("Response", func() {
	Describe("toResult", func() {
		Context("when using positive response", func() {
			var response = &Response{
				IDNumber: 2073386242,
				SummaryResult: SummaryResult{
					Key:     "id.success",
					Message: "PASS",
				},
				Results: Results{
					Key:     "result.match",
					Message: "ID Located",
				},
			}

			Context("when summary result is enabled", func() {
				It("should approve", func() {
					result, err := response.toResult(true)

					Expect(err).NotTo(HaveOccurred())
					Expect(result.Details).To(BeNil())
					Expect(result.Status).To(Equal(common.Approved))
				})

			})

			Context("when summary result is disabled", func() {
				It("should approve", func() {
					result, err := response.toResult(false)

					Expect(err).NotTo(HaveOccurred())
					Expect(result.Details).To(BeNil())
					Expect(result.Status).To(Equal(common.Approved))
				})
			})

			Context("when using sandbox response with ID Notes", func() {
				It("should approve with notes", func() {
					var response = &Response{
						IDNumber: 2073386292,
						SummaryResult: SummaryResult{
							Key:     "id.success",
							Message: "PASS",
						},
						Results: Results{
							Key:     "result.match",
							Message: "ID Located",
						},
						Qualifiers: &Qualifiers{
							Qualifiers: []Qualifier{
								Qualifier{
									Key:     "resultcode.address.does.not.match",
									Message: "Address Does Not Match",
								},
								Qualifier{
									Key:     "resultcode.street.number.does.not.match",
									Message: "Street Number Does Not Match",
								},
								Qualifier{
									Key:     "resultcode.street.name.does.not.match",
									Message: "Street Name Does Not Match",
								},
							},
						},
					}

					result, err := response.toResult(false)

					Expect(result.Status).To(Equal(common.Approved))
					Expect(result.Details).NotTo(BeNil())
					Expect(result.Details.Finality).To(Equal(common.Unknown))
					Expect(result.Details.Reasons).To(HaveLen(3))
					Expect(result.Details.Reasons[0]).To(Equal("Address Does Not Match"))
					Expect(result.Details.Reasons[1]).To(Equal("Street Number Does Not Match"))
					Expect(result.Details.Reasons[2]).To(Equal("Street Name Does Not Match"))
					Expect(err).NotTo(HaveOccurred())
				})
			})
		})

		Context("when using negative response", func() {
			It("should deny", func() {
				var response = &Response{
					IDNumber: 2073386264,
					SummaryResult: SummaryResult{
						Key:     "id.failure",
						Message: "FAIL",
					},
					Results: Results{
						Key:     "result.match.restricted",
						Message: "result.match.restricted",
					},
					Qualifiers: &Qualifiers{
						Qualifiers: []Qualifier{
							Qualifier{
								Key:     "resultcode.coppa.alert",
								Message: "COPPA Alert",
							},
						},
					},
				}

				result, err := response.toResult(false)

				Expect(result.Status).To(Equal(common.Denied))
				Expect(result.Details).NotTo(BeNil())
				Expect(result.Details.Finality).To(Equal(common.Unknown))
				Expect(result.Details.Reasons).To(HaveLen(1))
				Expect(result.Details.Reasons[0]).To(Equal("COPPA Alert"))
				Expect(err).NotTo(HaveOccurred())
			})

			It("should deny", func() {
				var response = &Response{
					IDNumber: 2073457900,
					SummaryResult: SummaryResult{
						Key:     "id.failure",
						Message: "FAIL",
					},
					Results: Results{
						Key:     "result.no.match",
						Message: "ID Not Located",
					},
				}

				result, err := response.toResult(false)

				Expect(result.Status).To(Equal(common.Denied))
				Expect(result.Details).To(BeNil())
				Expect(err).NotTo(HaveOccurred())
			})
		})

	})
})
