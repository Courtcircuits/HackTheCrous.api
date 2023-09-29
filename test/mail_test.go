package test

import (
	"testing"

	"github.com/Courtcircuits/HackTheCrous.api/util"
)

// RUN THIS NOT OFTEN

func TestSendConfirmationMail(t *testing.T) {
	err := util.SendConfirmationMail("hackthecroustkt@yopmail.com", "12345")
	if err != nil {
		t.Errorf("should not throw : %v\n", err)
	}
}
