# e-Kasa Cloud Client for Go

This is a Go client for the [e-Kasa cloud solution](https://github.com/ninedigit/ekasa-cloud) provided by [Nine Digit, s.r.o.](https://ekasa.ninedigit.sk/).

## Installation

```bash
go get ekasa
```

## Usage

### Initializing the Client

The client uses the functional options pattern for configuration.

```go
import (
    "ekasa"
    "time"
)

func main() {
    client := ekasa.NewClient(
        "your_public_key",
        "your_private_key",
        ekasa.WithTimeout(10 * time.Second),
        ekasa.WithTenantID("your_tenant_id"),
        ekasa.WithBaseURL(ekasa.PlaygroundEnvironment), // Use Playground for testing
    )
}
```

### Registering a Receipt

```go
ctx := context.Background()

receiptRequest := ekasa.CreateReceiptRegistrationDto{
    Printer: ekasa.ReceiptPrinterDto{
        Name: "pos",
        Options: ekasa.PosReceiptPrinterOptions{
            OpenDrawer: pointer(true),
        },
    },
    Request: ekasa.CreateReceiptRegistrationRequestDto{
        ReceiptType:      ekasa.ReceiptTypeCashRegister,
        CashRegisterCode: "88812345678900001",
        ExternalID:       "unique-request-id",
        Items: []ekasa.ReceiptRegistrationItemDto{
            {
                Type:      ekasa.ReceiptItemTypePositive,
                Name:      "Coca Cola 0.25l",
                UnitPrice: 1.29,
                VatRate:   20.0,
                Quantity:  ekasa.QuantityDto{Amount: 2, Unit: "ks"},
                Price:     2.58,
            },
        },
        Payments: []ekasa.ReceiptRegistrationPaymentDto{
            {
                Name:   "Hotovosť",
                Amount: 2.58,
            },
        },
    },
    ValidityTimeSpan: 10000,
}

registration, err := client.RegisterReceipt(ctx, &receiptRequest)
if err != nil {
    log.Fatal(err)
}

if registration.State == ekasa.RegistrationStateProcessed {
    fmt.Println("Receipt registered successfully!")
}
```

### Getting Customers

```go
filter := &ekasa.CustomerFilterDto{
    IsActive: pointer(true),
}

result, err := client.GetCustomers(ctx, filter)
if err != nil {
    log.Fatal(err)
}

for _, customer := range result.Items {
    fmt.Printf("Customer: %s %s\n", customer.FirstName, customer.LastName)
}
```

## Configuration Options

- `WithBaseURL(url string)`: Sets the base URL for the API.
- `WithHttpClient(c *http.Client)`: Sets a custom HTTP client.
- `WithTimeout(d time.Duration)`: Sets the request timeout.
- `WithTenantID(id string)`: Sets the tenant ID for all requests.
- `WithProxy(url string)`: Sets the proxy URL.

## Environments

- `ekasa.ProductionEnvironment`: `https://ekasa-cloud.ninedigit.sk/api`
- `ekasa.PlaygroundEnvironment`: `https://ekasa-cloud-int.ninedigit.sk/api`

---

*Helper function for pointers:*
```go
func pointer[T any](v T) *T {
    return &v
}
```
