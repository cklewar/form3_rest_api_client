package response

// OrganisationAccountData is ...
type OrganisationAccountData struct {
	Data OrganisationAccount
}

// OrganisationAccountAttributes is ...
type OrganisationAccountAttributes struct {
	AlternativeBankAccountNames []string `json:"alternative_bank_account_names"`
	BankID                      string   `json:"bank_id"`
	BankIDCode                  string   `json:"bank_id_code"`
	BaseCurrency                string   `json:"base_currency"`
	Bic                         string   `json:"bic"`
	Country                     string   `json:"country"`
}

// OrganisationAccount is ...
type OrganisationAccount struct {
	Attributes     OrganisationAccountAttributes
	CreatedOn      string `json:"created_on"`
	ID             string `json:"id"`
	ModifiedOn     string `json:"modified_on "`
	OrganisationID string `json:"organisation_id "`
	Type           string `json:"type"`
	Version        int    `json:"version"`
	Links          Links
}

// Links is ...
type Links struct {
	Self string `json:"self"`
}
