package ekasa

import (
	"encoding/json"
	"testing"
)

func TestModels_JSONNullability(t *testing.T) {
	// Test that nil pointers are omitted and values are present
	trueVal := true
	header := "Test Header"

	req := ReceiptRegistrationRequestDto{
		ExternalID: "ext-1",
		HeaderText: &header,
		Paragon:    false, // bool (not pointer) will be present
	}

	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if _, ok := m["invoiceNumber"]; ok {
		t.Error("invoiceNumber should be omitted when nil")
	}

	if val, ok := m["headerText"]; !ok || val != header {
		t.Errorf("headerText should be %s, got %v", header, val)
	}

	if val, ok := m["paragon"]; !ok || val != false {
		t.Errorf("paragon should be false, got %v", val)
	}

	// Test printer options
	printer := ReceiptPrinterDto{
		Name: "pos",
		Options: PosReceiptPrinterOptions{
			OpenDrawer: &trueVal,
		},
	}

	data, err = json.Marshal(printer)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	json.Unmarshal(data, &m)
	opts := m["options"].(map[string]interface{})
	if val, ok := opts["openDrawer"]; !ok || val != true {
		t.Errorf("openDrawer should be true, got %v", val)
	}
}
