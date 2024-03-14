package goveikkaus

import (
	"io"
	"log"
	"os"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("GlossaryService", func() {
	Describe("GameGlossary", func() {
		Context("Print", func() {
			It("should print the game glossary", func() {
				gameGlossary := GameGlossary{
					"FIXEDODDS": {
						AlsoKnownAs: "Pitkäveto",
						Description: "In fixed odds betting (Pitkäveto), you predict winners or outcomes for 1–20 matches. Stakes vary based on match count, sport, or time. Popular sports include soccer and ice hockey.",
					},
					// Add more entries as needed
				}

				// Capture printed output
				old := os.Stdout
				r, w, _ := os.Pipe()
				os.Stdout = w

				// Call the Print method
				gameGlossary.Print()

				// Restore stdout
				w.Close()
				os.Stdout = old

				// Read captured output from the pipe
				var output strings.Builder
				if _, err := io.Copy(&output, r); err != nil {
					log.Fatalf("did not expect to fail in this unit-test. ERR: %v", err)
				}

				// Close the read end of the pipe
				r.Close()

				// Get the captured output as a string
				capturedOutput := output.String()

				// Expected output
				expectedOutput := `
############# Pitkäveto - FIXEDODDS #############
API Term:			FIXEDODDS
Also known as:			Pitkäveto
In fixed odds betting (Pitkäveto), you predict winners or outcomes for 1–20 matches. Stakes vary based on match count, sport, or time. Popular sports include soccer and ice hockey.
`

				// Compare the captured output with the expected output
				Ω(strings.TrimSpace(capturedOutput)).Should(Equal(strings.TrimSpace(expectedOutput)))
			})
		})
	})

	Describe("GlossaryService", func() {
		Context("Get", func() {
			It("should return the game glossary", func() {
				service := &GlossaryService{}

				// Call the Get method
				gameGlossary := service.Get()

				// Check if the returned game glossary is not nil
				Ω(gameGlossary).ShouldNot(BeNil())

				// Add more assertions as needed
			})
		})
	})
})
