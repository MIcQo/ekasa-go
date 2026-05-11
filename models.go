package ekasa

import (
	"time"
)

// --- Common Models ---

type ProblemDetails struct {
	Type     string `json:"type,omitempty"`
	Title    string `json:"title,omitempty"`
	Status   int    `json:"status,omitempty"`
	Detail   string `json:"detail,omitempty"`
	Instance string `json:"instance,omitempty"`
	Code     int    `json:"code,omitempty"`
	TraceID  string `json:"traceId,omitempty"`
}

type ValidationProblemDetails struct {
	ProblemDetails
	Errors map[string][]string `json:"errors,omitempty"`
}

type QuantityDto struct {
	Amount float64 `json:"amount"`
	Unit   string  `json:"unit"`
}

type CustomerDto struct {
	ID                   string                                `json:"id"`
	Type                 string                                `json:"type"` // CustomerIdType
	CreationTime         time.Time                             `json:"creationTime,omitempty"`
	CreatorID            string                                `json:"creatorId,omitempty"`
	LastModificationTime *time.Time                            `json:"lastModificationTime,omitempty"`
	LastModifierID       string                                `json:"lastModifierId,omitempty"`
	ConcurrencyStamp     string                                `json:"concurrencyStamp,omitempty"`
	TenantID             string                                `json:"tenantId,omitempty"`
	IsActive             bool                                  `json:"isActive"`
	Status               string                                `json:"status,omitempty"` // CustomerState
	ExternalID           string                                `json:"externalId,omitempty"`
	ActivationTime       *time.Time                            `json:"activationTime,omitempty"`
	ExpirationTime       *time.Time                            `json:"expirationTime,omitempty"`
	FirstName            string                                `json:"firstName,omitempty"`
	LastName             string                                `json:"lastName,omitempty"`
	Gender               string                                `json:"gender,omitempty"`
	BirthDate            string                                `json:"birthDate,omitempty"`
	Email                string                                `json:"email,omitempty"`
	Phone                string                                `json:"phone,omitempty"`
	IsCompany            bool                                  `json:"isCompany"`
	Company              *CustomerCompanyDto                   `json:"company,omitempty"`
	Address              *CustomerAddressDto                   `json:"address,omitempty"`
	CreditBalance        *CreditDto                            `json:"creditBalance,omitempty"`
	Transactions         []CustomerCreditBalanceTransactionDto `json:"creditBalanceTransactions,omitempty"`
	Cards                []CustomerCardDto                     `json:"cards,omitempty"`
	CreditRate           float64                               `json:"creditRate,omitempty"`
	DiscountRate         float64                               `json:"discountRate,omitempty"`
	Note                 string                                `json:"note,omitempty"`
	Meta                 map[string]interface{}                `json:"meta,omitempty"`
	UserMeta             []UserRawMetaDto                      `json:"userMeta,omitempty"`
}

type CustomerCompanyDto struct {
	Name  string `json:"name,omitempty"`
	CRN   string `json:"crn,omitempty"`
	VatID string `json:"vatId,omitempty"`
	TaxID string `json:"taxId,omitempty"`
}

type CustomerAddressDto struct {
	Street      string             `json:"street,omitempty"`
	City        string             `json:"city,omitempty"`
	PostalCode  string             `json:"postalCode,omitempty"`
	Country     string             `json:"country,omitempty"`
	Coordinates *GeoCoordinatesDto `json:"coordinates,omitempty"`
	Note        string             `json:"note,omitempty"`
}

type GeoCoordinatesDto struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}

type CreditDto struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency,omitempty"`
}

type CustomerCreditBalanceTransactionDto struct {
	ID                  string                 `json:"id"`
	CreationTime        time.Time              `json:"creationTime"`
	CreatorID           string                 `json:"creatorId,omitempty"`
	SequenceNumber      int                    `json:"sequenceNumber"`
	ExternalID          string                 `json:"externalId,omitempty"`
	Amount              CreditDto              `json:"amount"`
	Type                string                 `json:"type"` // CustomerCreditBalanceTransactionType
	EndingCreditBalance CreditDto              `json:"endingCreditBalance"`
	Note                string                 `json:"note,omitempty"`
	Meta                map[string]interface{} `json:"meta,omitempty"`
	UserMeta            []UserRawMetaDto       `json:"userMeta,omitempty"`
	CustomerID          string                 `json:"customerId"`
}

type CustomerCardDto struct {
	ID               string                 `json:"id"`
	CreationTime     time.Time              `json:"creationTime"`
	CreatorID        string                 `json:"creatorId,omitempty"`
	LastModification *time.Time             `json:"lastModificationTime,omitempty"`
	LastModifierID   string                 `json:"lastModifierId,omitempty"`
	TenantID         string                 `json:"tenantId,omitempty"`
	ExternalID       string                 `json:"externalId,omitempty"`
	ConcurrencyStamp string                 `json:"concurrencyStamp,omitempty"`
	IsActive         bool                   `json:"isActive"`
	IsVirtual        bool                   `json:"isVirtual"`
	SerialNumber     string                 `json:"serialNumber"`
	Processor        string                 `json:"processor"` // CustomerCardProcessorNames
	Status           string                 `json:"status"`    // CustomerCardState
	StatusTime       *time.Time             `json:"statusTime,omitempty"`
	StatusReason     string                 `json:"statusReason,omitempty"`
	ActivationTime   time.Time              `json:"activationTime"`
	ExpirationTime   *time.Time             `json:"expirationTime,omitempty"`
	Note             string                 `json:"note,omitempty"`
	Meta             map[string]interface{} `json:"meta,omitempty"`
	UserMeta         []UserRawMetaDto       `json:"userMeta,omitempty"`
	CustomerID       string                 `json:"customerId"`
}

type UserRawMetaDto struct {
	UserID string                 `json:"userId"`
	Meta   map[string]interface{} `json:"meta"`
}

type GetCustomerListResultDto struct {
	Items       []CustomerDto `json:"items"`
	RequestTime time.Time     `json:"requestTime"`
}

type CustomerFilterDto struct {
	IDs               []string   `json:"ids,omitempty"`
	ExternalID        string     `json:"externalId,omitempty"`
	ModifiedAfter     *time.Time `json:"modifiedAfter,omitempty"`
	IsActive          *bool      `json:"isActive,omitempty"`
	CardID            string     `json:"cardId,omitempty"`
	CardSerialNumbers []string   `json:"cardSerialNumbers,omitempty"`
}

// --- Registration Models ---

type CreateReceiptRegistrationDto struct {
	Printer          ReceiptPrinterDto                   `json:"printer"`
	Request          CreateReceiptRegistrationRequestDto `json:"request"`
	ValidityTimeSpan int                                 `json:"validityTimeSpan"`
}

type ReceiptPrinterDto struct {
	Name    string      `json:"name"` // pos, pdf, email
	Options interface{} `json:"options,omitempty"`
}

type PosReceiptPrinterOptions struct {
	OpenDrawer        *bool `json:"openDrawer,omitempty"`
	PrintLogo         *bool `json:"printLogo,omitempty"`
	LogoMemoryAddress *int  `json:"logoMemoryAddress,omitempty"`
}

type EmailReceiptPrinterOptions struct {
	To                   string  `json:"to"`
	RecipientDisplayName *string `json:"recipientDisplayName,omitempty"`
	Subject              *string `json:"subject,omitempty"`
	Body                 *string `json:"body,omitempty"`
}

type PdfReceiptPrinterOptions struct {
}

type CreateReceiptRegistrationRequestDto struct {
	ExternalID       string                          `json:"externalId"`
	CashRegisterCode string                          `json:"cashRegisterCode"`
	ReceiptType      string                          `json:"receiptType"` // ReceiptType
	IssueDate        *time.Time                      `json:"issueDate,omitempty"`
	InvoiceNumber    *string                         `json:"invoiceNumber,omitempty"`
	ParagonNumber    *int                            `json:"paragonNumber,omitempty"`
	Amount           *float64                        `json:"amount,omitempty"`
	HeaderText       *string                         `json:"headerText,omitempty"`
	FooterText       *string                         `json:"footerText,omitempty"`
	Customer         *CustomerDto                    `json:"customer,omitempty"`
	Items            []ReceiptRegistrationItemDto    `json:"items,omitempty"`
	Payments         []ReceiptRegistrationPaymentDto `json:"payments,omitempty"`
}

type ReceiptRegistrationItemDto struct {
	Type               string      `json:"type"` // ReceiptItemType
	Name               string      `json:"name"`
	Price              float64     `json:"price"`
	UnitPrice          float64     `json:"unitPrice"`
	Quantity           QuantityDto `json:"quantity"`
	VatRate            float64     `json:"vatRate"`
	ReferenceReceiptID *string     `json:"referenceReceiptId,omitempty"`
	SpecialRegulation  *string     `json:"specialRegulation,omitempty"`
	VoucherNumber      *string     `json:"voucherNumber,omitempty"`
	Seller             *SellerDto  `json:"seller,omitempty"`
	Description        *string     `json:"description,omitempty"`
}

type ReceiptRegistrationPaymentDto struct {
	Name   string  `json:"name"`
	Amount float64 `json:"amount"`
}

type SellerDto struct {
	ID   string `json:"id"`
	Type string `json:"type"` // SellerIdType
}

type ReceiptRegistrationDto struct {
	ID                 string                          `json:"id"`
	Request            ReceiptRegistrationRequestDto   `json:"request"`
	Printer            ReceiptPrinterDto               `json:"printer"`
	CreationDate       time.Time                       `json:"creationDate"`
	CreatedBy          string                          `json:"createdBy"`
	NotificationDate   *time.Time                      `json:"notificationDate,omitempty"`
	ValidityTimeSpan   int                             `json:"validityTimeSpan"`
	AcceptationDate    *time.Time                      `json:"acceptationDate,omitempty"`
	CompletionTimeSpan *int                            `json:"completionTimeSpan,omitempty"`
	CompletionDate     *time.Time                      `json:"completionDate,omitempty"`
	State              string                          `json:"state"` // RegistrationState
	Error              *RegistrationErrorDto           `json:"error,omitempty"`
	RejectionReason    *RegistrationRejectionReasonDto `json:"rejectionReason,omitempty"`
}

type ReceiptRegistrationRequestDto struct {
	ExternalID       string                          `json:"externalId"`
	CashRegisterCode string                          `json:"cashRegisterCode"`
	ReceiptType      string                          `json:"receiptType"`
	IssueDate        *time.Time                      `json:"issueDate,omitempty"`
	OrpCreateDate    *time.Time                      `json:"orpCreateDate,omitempty"`
	ReceiptNumber    *int                            `json:"receiptNumber,omitempty"`
	InvoiceNumber    *string                         `json:"invoiceNumber,omitempty"`
	Paragon          bool                            `json:"paragon"`
	ParagonNumber    *int                            `json:"paragonNumber,omitempty"`
	DIC              *string                         `json:"dic,omitempty"`
	ICDPH            *string                         `json:"icDph,omitempty"`
	ICO              *string                         `json:"ico,omitempty"`
	Amount           *float64                        `json:"amount,omitempty"`
	Customer         *CustomerDto                    `json:"customer,omitempty"`
	BasicVatAmount   *float64                        `json:"basicVatAmount,omitempty"`
	ReducedVatAmount *float64                        `json:"reducedVatAmount,omitempty"`
	TaxFreeAmount    *float64                        `json:"taxFreeAmount,omitempty"`
	TaxBaseBasic     *float64                        `json:"taxBaseBasic,omitempty"`
	TaxBaseReduced   *float64                        `json:"taxBaseReduced,omitempty"`
	Items            []ReceiptRegistrationItemDto    `json:"items,omitempty"`
	Payments         []ReceiptRegistrationPaymentDto `json:"payments,omitempty"`
	OKP              *string                         `json:"okp,omitempty"`
	PKP              *string                         `json:"pkp,omitempty"`
	HeaderText       *string                         `json:"headerText,omitempty"`
	FooterText       *string                         `json:"footerText,omitempty"`
	RequestID        *string                         `json:"requestId,omitempty"`
	RequestDate      *time.Time                      `json:"requestDate,omitempty"`
	SendingCount     *int                            `json:"sendingCount,omitempty"`
	ReceiptID        *string                         `json:"receiptId,omitempty"`
	OrpProcessDate   *time.Time                      `json:"orpProcessDate,omitempty"`
	EKasaError       *EKasaErrorDto                  `json:"eKasaError,omitempty"`
}

type EKasaErrorDto struct {
	Message string `json:"message"`
	Code    *int   `json:"code,omitempty"`
}

type RegistrationErrorDto struct {
	Message string  `json:"message"`
	Code    *string `json:"code,omitempty"`
	Source  *string `json:"source,omitempty"`
	TraceID *string `json:"traceId,omitempty"`
}

type RegistrationRejectionReasonDto struct {
	Message string `json:"message"`
	Code    *int   `json:"code,omitempty"`
}

type ReceiptRegistrationStateChangeResultDto struct {
	IsSuccessful bool                    `json:"isSuccessful"`
	Registration *ReceiptRegistrationDto `json:"registration,omitempty"`
}

// --- Enums & Constants ---

const (
	CustomerIdTypeICO   = "ICO"
	CustomerIdTypeDIC   = "DIC"
	CustomerIdTypeICDPH = "ICDPH"
	CustomerIdTypeOther = "Other"
)

const (
	RegistrationStateCreated          = "Created"
	RegistrationStateNotified         = "Notified"
	RegistrationStateExpired          = "Expired"
	RegistrationStateAccepted         = "Accepted"
	RegistrationStateCanceled         = "Canceled"
	RegistrationStateTimedOut         = "TimedOut"
	RegistrationStateRejected         = "Rejected"
	RegistrationStateProcessed        = "Processed"
	RegistrationStateProcessedOffline = "ProcessedOffline"
	RegistrationStateProcessFailed    = "ProcessFailed"
	RegistrationStateFailed           = "Failed"
)

const (
	ReceiptTypeCashRegister   = "CashRegister"
	ReceiptTypeInvalid        = "Invalid"
	ReceiptTypeParagon        = "Paragon"
	ReceiptTypeInvoice        = "Invoice"
	ReceiptTypeInvoiceParagon = "InvoiceParagon"
	ReceiptTypeDeposit        = "Deposit"
	ReceiptTypeWithdraw       = "Withdraw"
)

const (
	ReceiptItemTypePositive          = "Positive"
	ReceiptItemTypeReturnedContainer = "ReturnedContainer"
	ReceiptItemTypeReturned          = "Returned"
	ReceiptItemTypeCorrection        = "Correction"
	ReceiptItemTypeDiscount          = "Discount"
	ReceiptItemTypeAdvance           = "Advance"
	ReceiptItemTypeVoucher           = "Voucher"
)
