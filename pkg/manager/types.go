package manager

// see https://github.com/SUSE-Enceladus/csp-billing-adapter#csp-config
type CSPConfig struct {
	Timestamp           string                  `json:"timestamp"`
        BillingAPIAccessOK  bool                    `json:"billing_api_access_ok"`
        Expire		    string                  `json:"expire"`
	LastBilled	    string                  `json:"last_billed,omitempty"`
	Errors              []string                `json:"errors,omitempty"`
	BaseProduct         string                  `json:"base_product,omitempty"`
	// we don't care about the rest of the fields so ignoring them for now
	Ignore              map[string]interface{}  `json:"-"`
}
