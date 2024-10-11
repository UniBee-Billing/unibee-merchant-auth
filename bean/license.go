package bean

type Plan struct {
	// amount, cent, without tax
	Amount *int64 `json:"amount,omitempty"`
	// binded recurring addon planIds，split with ,
	BindingAddonIds *string `json:"bindingAddonIds,omitempty"`
	// binded onetime addon planIds，split with ,
	BindingOnetimeAddonIds *string `json:"bindingOnetimeAddonIds,omitempty"`
	// whether cancel at subscripiton first trial end，0-false | 1-true, will pass to cancelAtPeriodEnd of subscription
	CancelAtTrialEnd *int32 `json:"cancelAtTrialEnd,omitempty"`
	// create utc time
	CreateTime *int64 `json:"createTime,omitempty"`
	// currency
	Currency *string `json:"currency,omitempty"`
	// description
	Description *string `json:"description,omitempty"`
	// external_user_id
	ExternalPlanId  *string `json:"externalPlanId,omitempty"`
	ExtraMetricData *string `json:"extraMetricData,omitempty"`
	// who pay the gas, merchant|user
	GasPayer *string `json:"gasPayer,omitempty"`
	// home_url
	HomeUrl *string `json:"homeUrl,omitempty"`
	Id      *int64  `json:"id,omitempty"`
	// image_url
	ImageUrl *string `json:"imageUrl,omitempty"`
	// period unit count
	IntervalCount *int32 `json:"intervalCount,omitempty"`
	// period unit,day|month|year|week
	IntervalUnit *string `json:"intervalUnit,omitempty"`
	// merchant id
	MerchantId *int64                  `json:"merchantId,omitempty"`
	Metadata   *map[string]interface{} `json:"metadata,omitempty"`
	// PlanName
	PlanName *string `json:"planName,omitempty"`
	// product id
	ProductId *int64 `json:"productId,omitempty"`
	// 1-UnPublish,2-Publish, Use For Display Plan At UserPortal
	PublishStatus *int32 `json:"publishStatus,omitempty"`
	// status，1-editing，2-active，3-inactive，4-expired
	Status *int32 `json:"status,omitempty"`
	// TaxPercentage 1000 = 10%
	TaxPercentage *int32 `json:"taxPercentage,omitempty"`
	// price of trial period
	TrialAmount *int64  `json:"trialAmount,omitempty"`
	TrialDemand *string `json:"trialDemand,omitempty"`
	// duration of trial
	TrialDurationTime *int64 `json:"trialDurationTime,omitempty"`
	// type，1-main plan，2-addon plan
	Type *int32 `json:"type,omitempty"`
}

type MerchantVersion struct {
	Name      string `json:"name" dc:"Name"`
	IsPaid    bool   `json:"isPaid" dc:"IsPaid"`
	Expired   bool   `json:"expired" dc:"Expired"`
	Plan      *Plan  `json:"plan" dc:"Plan"`
	StartTime int64  `json:"startTime" dc:"StartTime,UTC, The Start Time Of Plan,0 for free"`
	EndTime   int64  `json:"endTime" dc:"EndTime,UTC, The End Time Of Plan,0 for free"`
}

type License struct {
	OwnerEmail string           `json:"ownerEmail" dc:"OwnerEmail"`
	Version    *MerchantVersion `json:"version" dc:"Version Info"`
	License    string           `json:"license" dc:"License, Premium Version will contain License"`
}
