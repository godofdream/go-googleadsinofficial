package TrialService

import (
	"bytes"
	"crypto/tls"
	"encoding/xml"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"time"
)

// against "unused imports"
var _ time.Time
var _ xml.Name

//
// The single reason for the authentication failure.
//
type AuthenticationErrorReason string

const (

	//
	// Authentication of the request failed.
	//
	AuthenticationErrorReasonAUTHENTICATION_FAILED AuthenticationErrorReason = "AUTHENTICATION_FAILED"

	//
	// Client Customer Id is required if CustomerIdMode is set to CLIENT_EXTERNAL_CUSTOMER_ID.
	// Starting version V201409 ClientCustomerId will be required for all requests except
	// for {@link CustomerService#get}
	//
	AuthenticationErrorReasonCLIENT_CUSTOMER_ID_IS_REQUIRED AuthenticationErrorReason = "CLIENT_CUSTOMER_ID_IS_REQUIRED"

	//
	// Client Email is required if CustomerIdMode is set to CLIENT_EXTERNAL_EMAIL_FIELD.
	//
	AuthenticationErrorReasonCLIENT_EMAIL_REQUIRED AuthenticationErrorReason = "CLIENT_EMAIL_REQUIRED"

	//
	// Client customer Id is not a number.
	//
	AuthenticationErrorReasonCLIENT_CUSTOMER_ID_INVALID AuthenticationErrorReason = "CLIENT_CUSTOMER_ID_INVALID"

	//
	// Client customer Id is not a number.
	//
	AuthenticationErrorReasonCLIENT_EMAIL_INVALID AuthenticationErrorReason = "CLIENT_EMAIL_INVALID"

	//
	// Client email is not a valid customer email.
	//
	AuthenticationErrorReasonCLIENT_EMAIL_FAILED_TO_AUTHENTICATE AuthenticationErrorReason = "CLIENT_EMAIL_FAILED_TO_AUTHENTICATE"

	//
	// No customer found for the customer id provided in the header.
	//
	AuthenticationErrorReasonCUSTOMER_NOT_FOUND AuthenticationErrorReason = "CUSTOMER_NOT_FOUND"

	//
	// Client's Google Account is deleted.
	//
	AuthenticationErrorReasonGOOGLE_ACCOUNT_DELETED AuthenticationErrorReason = "GOOGLE_ACCOUNT_DELETED"

	//
	// Google account login token in the cookie is invalid.
	//
	AuthenticationErrorReasonGOOGLE_ACCOUNT_COOKIE_INVALID AuthenticationErrorReason = "GOOGLE_ACCOUNT_COOKIE_INVALID"

	//
	// problem occurred during Google account authentication.
	//
	AuthenticationErrorReasonFAILED_TO_AUTHENTICATE_GOOGLE_ACCOUNT AuthenticationErrorReason = "FAILED_TO_AUTHENTICATE_GOOGLE_ACCOUNT"

	//
	// The user in the google account login token does not match the UserId in the cookie.
	//
	AuthenticationErrorReasonGOOGLE_ACCOUNT_USER_AND_ADS_USER_MISMATCH AuthenticationErrorReason = "GOOGLE_ACCOUNT_USER_AND_ADS_USER_MISMATCH"

	//
	// Login cookie is required for authentication.
	//
	AuthenticationErrorReasonLOGIN_COOKIE_REQUIRED AuthenticationErrorReason = "LOGIN_COOKIE_REQUIRED"

	//
	// User in the cookie is not a valid Ads user.
	//
	AuthenticationErrorReasonNOT_ADS_USER AuthenticationErrorReason = "NOT_ADS_USER"

	//
	// Oauth token in the header is not valid.
	//
	AuthenticationErrorReasonOAUTH_TOKEN_INVALID AuthenticationErrorReason = "OAUTH_TOKEN_INVALID"

	//
	// Oauth token in the header has expired.
	//
	AuthenticationErrorReasonOAUTH_TOKEN_EXPIRED AuthenticationErrorReason = "OAUTH_TOKEN_EXPIRED"

	//
	// Oauth token in the header has been disabled.
	//
	AuthenticationErrorReasonOAUTH_TOKEN_DISABLED AuthenticationErrorReason = "OAUTH_TOKEN_DISABLED"

	//
	// Oauth token in the header has been revoked.
	//
	AuthenticationErrorReasonOAUTH_TOKEN_REVOKED AuthenticationErrorReason = "OAUTH_TOKEN_REVOKED"

	//
	// Oauth token HTTP header is malformed.
	//
	AuthenticationErrorReasonOAUTH_TOKEN_HEADER_INVALID AuthenticationErrorReason = "OAUTH_TOKEN_HEADER_INVALID"

	//
	// Login cookie is not valid.
	//
	AuthenticationErrorReasonLOGIN_COOKIE_INVALID AuthenticationErrorReason = "LOGIN_COOKIE_INVALID"

	//
	// Failed to decrypt the login cookie.
	//
	AuthenticationErrorReasonFAILED_TO_RETRIEVE_LOGIN_COOKIE AuthenticationErrorReason = "FAILED_TO_RETRIEVE_LOGIN_COOKIE"

	//
	// User Id in the header is not a valid id.
	//
	AuthenticationErrorReasonUSER_ID_INVALID AuthenticationErrorReason = "USER_ID_INVALID"
)

//
// The reasons for the authorization error.
//
type AuthorizationErrorReason string

const (

	//
	// Could not complete authorization due to an internal problem.
	//
	AuthorizationErrorReasonUNABLE_TO_AUTHORIZE AuthorizationErrorReason = "UNABLE_TO_AUTHORIZE"

	//
	// Customer has no AdWords account.
	//
	AuthorizationErrorReasonNO_ADWORDS_ACCOUNT_FOR_CUSTOMER AuthorizationErrorReason = "NO_ADWORDS_ACCOUNT_FOR_CUSTOMER"

	//
	// User doesn't have permission to access customer.
	//
	AuthorizationErrorReasonUSER_PERMISSION_DENIED AuthorizationErrorReason = "USER_PERMISSION_DENIED"

	//
	// Effective user doesn't have permission to access this customer.
	//
	AuthorizationErrorReasonEFFECTIVE_USER_PERMISSION_DENIED AuthorizationErrorReason = "EFFECTIVE_USER_PERMISSION_DENIED"

	//
	// Access denied since the customer is not active.
	//
	AuthorizationErrorReasonCUSTOMER_NOT_ACTIVE AuthorizationErrorReason = "CUSTOMER_NOT_ACTIVE"

	//
	// User has read-only permission cannot mutate.
	//
	AuthorizationErrorReasonUSER_HAS_READONLY_PERMISSION AuthorizationErrorReason = "USER_HAS_READONLY_PERMISSION"

	//
	// No customer found.
	//
	AuthorizationErrorReasonNO_CUSTOMER_FOUND AuthorizationErrorReason = "NO_CUSTOMER_FOUND"

	//
	// Developer doesn't have permission to access service.
	//
	AuthorizationErrorReasonSERVICE_ACCESS_DENIED AuthorizationErrorReason = "SERVICE_ACCESS_DENIED"
)

type BiddingErrorsReason string

const (

	//
	// Cannot transition to new bidding strategy.
	//
	BiddingErrorsReasonBIDDING_STRATEGY_TRANSITION_NOT_ALLOWED BiddingErrorsReason = "BIDDING_STRATEGY_TRANSITION_NOT_ALLOWED"

	//
	// Bidding strategy cannot be overridden by current ad group bidding strategy.
	//
	BiddingErrorsReasonBIDDING_STRATEGY_NOT_COMPATIBLE_WITH_ADGROUP_OVERRIDES BiddingErrorsReason = "BIDDING_STRATEGY_NOT_COMPATIBLE_WITH_ADGROUP_OVERRIDES"

	//
	// Bidding strategy cannot be overridden by current criteria bidding strategy.
	//
	BiddingErrorsReasonBIDDING_STRATEGY_NOT_COMPATIBLE_WITH_ADGROUP_CRITERIA_OVERRIDES BiddingErrorsReason = "BIDDING_STRATEGY_NOT_COMPATIBLE_WITH_ADGROUP_CRITERIA_OVERRIDES"

	//
	// Cannot override campaign bidding strategy.
	//
	BiddingErrorsReasonCAMPAIGN_BIDDING_STRATEGY_CANNOT_BE_OVERRIDDEN BiddingErrorsReason = "CAMPAIGN_BIDDING_STRATEGY_CANNOT_BE_OVERRIDDEN"

	//
	// Cannot override ad group bidding strategy.
	//
	BiddingErrorsReasonADGROUP_BIDDING_STRATEGY_CANNOT_BE_OVERRIDDEN BiddingErrorsReason = "ADGROUP_BIDDING_STRATEGY_CANNOT_BE_OVERRIDDEN"

	//
	// Cannot attach bidding strategy to campaign.
	//
	BiddingErrorsReasonCANNOT_ATTACH_BIDDING_STRATEGY_TO_CAMPAIGN BiddingErrorsReason = "CANNOT_ATTACH_BIDDING_STRATEGY_TO_CAMPAIGN"

	//
	// Cannot attach bidding strategy to ad group.
	//
	BiddingErrorsReasonCANNOT_ATTACH_BIDDING_STRATEGY_TO_ADGROUP BiddingErrorsReason = "CANNOT_ATTACH_BIDDING_STRATEGY_TO_ADGROUP"

	//
	// Cannot attach bidding strategy to criteria.
	//
	BiddingErrorsReasonCANNOT_ATTACH_BIDDING_STRATEGY_TO_ADGROUP_CRITERIA BiddingErrorsReason = "CANNOT_ATTACH_BIDDING_STRATEGY_TO_ADGROUP_CRITERIA"

	//
	// Bidding strategy is not supported or cannot be used as anonymous.
	//
	BiddingErrorsReasonINVALID_ANONYMOUS_BIDDING_STRATEGY_TYPE BiddingErrorsReason = "INVALID_ANONYMOUS_BIDDING_STRATEGY_TYPE"

	//
	// No bids may be set. The bid list must be null or empty.
	//
	BiddingErrorsReasonBIDS_NOT_ALLLOWED BiddingErrorsReason = "BIDS_NOT_ALLLOWED"

	//
	// The bid list contains two or more bids of the same type.
	//
	BiddingErrorsReasonDUPLICATE_BIDS BiddingErrorsReason = "DUPLICATE_BIDS"

	//
	// The bidding scheme does not match the bidding strategy type.
	//
	BiddingErrorsReasonINVALID_BIDDING_SCHEME BiddingErrorsReason = "INVALID_BIDDING_SCHEME"

	//
	// The type does not match the named strategy's type.
	//
	BiddingErrorsReasonINVALID_BIDDING_STRATEGY_TYPE BiddingErrorsReason = "INVALID_BIDDING_STRATEGY_TYPE"

	//
	// The bidding strategy type is missing.
	//
	BiddingErrorsReasonMISSING_BIDDING_STRATEGY_TYPE BiddingErrorsReason = "MISSING_BIDDING_STRATEGY_TYPE"

	//
	// The bid list contains a null entry.
	//
	BiddingErrorsReasonNULL_BID BiddingErrorsReason = "NULL_BID"

	//
	// The bid is invalid.
	//
	BiddingErrorsReasonINVALID_BID BiddingErrorsReason = "INVALID_BID"

	//
	// Bidding strategy is not available for the account type.
	//
	BiddingErrorsReasonBIDDING_STRATEGY_NOT_AVAILABLE_FOR_ACCOUNT_TYPE BiddingErrorsReason = "BIDDING_STRATEGY_NOT_AVAILABLE_FOR_ACCOUNT_TYPE"

	//
	// Conversion tracking is not enabled for the campaign for VBB transition.
	//
	BiddingErrorsReasonCONVERSION_TRACKING_NOT_ENABLED BiddingErrorsReason = "CONVERSION_TRACKING_NOT_ENABLED"

	//
	// Not enough conversions tracked for VBB transitions.
	//
	BiddingErrorsReasonNOT_ENOUGH_CONVERSIONS BiddingErrorsReason = "NOT_ENOUGH_CONVERSIONS"

	//
	// Campaign can not be created with given bidding strategy. It can be transitioned to the
	// strategy, once eligible.
	//
	BiddingErrorsReasonCANNOT_CREATE_CAMPAIGN_WITH_BIDDING_STRATEGY BiddingErrorsReason = "CANNOT_CREATE_CAMPAIGN_WITH_BIDDING_STRATEGY"

	//
	// Cannot target content network only as ad group uses Page One Promoted bidding strategy.
	//
	BiddingErrorsReasonCANNOT_TARGET_CONTENT_NETWORK_ONLY_WITH_AD_GROUP_LEVEL_POP_BIDDING_STRATEGY BiddingErrorsReason = "CANNOT_TARGET_CONTENT_NETWORK_ONLY_WITH_AD_GROUP_LEVEL_POP_BIDDING_STRATEGY"

	//
	// Cannot target content network only as campaign uses Page One Promoted bidding strategy.
	//
	BiddingErrorsReasonCANNOT_TARGET_CONTENT_NETWORK_ONLY_WITH_CAMPAIGN_LEVEL_POP_BIDDING_STRATEGY BiddingErrorsReason = "CANNOT_TARGET_CONTENT_NETWORK_ONLY_WITH_CAMPAIGN_LEVEL_POP_BIDDING_STRATEGY"

	//
	// Budget Optimizer and Target Spend bidding strategies are not supported for campaigns with
	// AdSchedule targeting.
	//
	BiddingErrorsReasonBIDDING_STRATEGY_NOT_SUPPORTED_WITH_AD_SCHEDULE BiddingErrorsReason = "BIDDING_STRATEGY_NOT_SUPPORTED_WITH_AD_SCHEDULE"

	//
	// Pay per conversion is not available to all the customer, only few whitelisted customers
	// can use this.
	//
	BiddingErrorsReasonPAY_PER_CONVERSION_NOT_AVAILABLE_FOR_CUSTOMER BiddingErrorsReason = "PAY_PER_CONVERSION_NOT_AVAILABLE_FOR_CUSTOMER"

	//
	// Pay per conversion is not allowed with Target CPA.
	//
	BiddingErrorsReasonPAY_PER_CONVERSION_NOT_ALLOWED_WITH_TARGET_CPA BiddingErrorsReason = "PAY_PER_CONVERSION_NOT_ALLOWED_WITH_TARGET_CPA"

	//
	// Cannot set bidding strategy to Manual CPM for search network only campaigns.
	//
	BiddingErrorsReasonBIDDING_STRATEGY_NOT_ALLOWED_FOR_SEARCH_ONLY_CAMPAIGNS BiddingErrorsReason = "BIDDING_STRATEGY_NOT_ALLOWED_FOR_SEARCH_ONLY_CAMPAIGNS"

	//
	// The bidding strategy is not supported for use in drafts or experiments.
	//
	BiddingErrorsReasonBIDDING_STRATEGY_NOT_SUPPORTED_IN_DRAFTS_OR_EXPERIMENTS BiddingErrorsReason = "BIDDING_STRATEGY_NOT_SUPPORTED_IN_DRAFTS_OR_EXPERIMENTS"

	//
	// Bidding strategy type does not support product type ad group criterion.
	//
	BiddingErrorsReasonBIDDING_STRATEGY_TYPE_DOES_NOT_SUPPORT_PRODUCT_TYPE_ADGROUP_CRITERION BiddingErrorsReason = "BIDDING_STRATEGY_TYPE_DOES_NOT_SUPPORT_PRODUCT_TYPE_ADGROUP_CRITERION"

	//
	// Bid amount is too small.
	//
	BiddingErrorsReasonBID_TOO_SMALL BiddingErrorsReason = "BID_TOO_SMALL"

	//
	// Bid amount is too big.
	//
	BiddingErrorsReasonBID_TOO_BIG BiddingErrorsReason = "BID_TOO_BIG"

	//
	// Bid has too many fractional digit precision.
	//
	BiddingErrorsReasonBID_TOO_MANY_FRACTIONAL_DIGITS BiddingErrorsReason = "BID_TOO_MANY_FRACTIONAL_DIGITS"

	//
	// EnhancedCpcEnabled cannot be set on <em>portfolio</em> Target Spend strategies.
	//
	BiddingErrorsReasonENHANCED_CPC_ENABLED_NOT_SUPPORTED_ON_PORTFOLIO_TARGET_SPEND_STRATEGY BiddingErrorsReason = "ENHANCED_CPC_ENABLED_NOT_SUPPORTED_ON_PORTFOLIO_TARGET_SPEND_STRATEGY"

	BiddingErrorsReasonUNKNOWN BiddingErrorsReason = "UNKNOWN"
)

//
// The reasons for the budget error.
//
type BudgetErrorReason string

const (

	//
	// The requested budget no longer exists.
	//
	BudgetErrorReasonBUDGET_REMOVED BudgetErrorReason = "BUDGET_REMOVED"

	//
	// Default budget error.
	//
	BudgetErrorReasonBUDGET_ERROR BudgetErrorReason = "BUDGET_ERROR"

	//
	// The budget is associated with at least one campaign, and so the budget cannot be removed.
	//
	BudgetErrorReasonBUDGET_IN_USE BudgetErrorReason = "BUDGET_IN_USE"

	//
	// Customer is not whitelisted for this budget period.
	//
	BudgetErrorReasonBUDGET_PERIOD_NOT_AVAILABLE BudgetErrorReason = "BUDGET_PERIOD_NOT_AVAILABLE"

	//
	// Customer cannot use CampaignService to edit a shared budget.
	//
	BudgetErrorReasonCANNOT_EDIT_SHARED_BUDGET BudgetErrorReason = "CANNOT_EDIT_SHARED_BUDGET"

	//
	// This field is not mutable on implicitly shared budgets
	//
	BudgetErrorReasonCANNOT_MODIFY_FIELD_OF_IMPLICITLY_SHARED_BUDGET BudgetErrorReason = "CANNOT_MODIFY_FIELD_OF_IMPLICITLY_SHARED_BUDGET"

	//
	// Cannot change explicitly shared budgets back to implicitly shared ones.
	//
	BudgetErrorReasonCANNOT_UPDATE_BUDGET_TO_IMPLICITLY_SHARED BudgetErrorReason = "CANNOT_UPDATE_BUDGET_TO_IMPLICITLY_SHARED"

	//
	// An implicit budget without a name cannot be changed to explicitly shared budget.
	//
	BudgetErrorReasonCANNOT_UPDATE_BUDGET_TO_EXPLICITLY_SHARED_WITHOUT_NAME BudgetErrorReason = "CANNOT_UPDATE_BUDGET_TO_EXPLICITLY_SHARED_WITHOUT_NAME"

	//
	// Cannot change an implicitly shared budget to an explicitly shared one.
	//
	BudgetErrorReasonCANNOT_UPDATE_BUDGET_TO_EXPLICITLY_SHARED BudgetErrorReason = "CANNOT_UPDATE_BUDGET_TO_EXPLICITLY_SHARED"

	//
	// Only explicitly shared budgets can be used with multiple campaigns.
	//
	BudgetErrorReasonCANNOT_USE_IMPLICITLY_SHARED_BUDGET_WITH_MULTIPLE_CAMPAIGNS BudgetErrorReason = "CANNOT_USE_IMPLICITLY_SHARED_BUDGET_WITH_MULTIPLE_CAMPAIGNS"

	//
	// A budget with this name already exists.
	//
	BudgetErrorReasonDUPLICATE_NAME BudgetErrorReason = "DUPLICATE_NAME"

	//
	// A money amount was not in the expected currency.
	//
	BudgetErrorReasonMONEY_AMOUNT_IN_WRONG_CURRENCY BudgetErrorReason = "MONEY_AMOUNT_IN_WRONG_CURRENCY"

	//
	// A money amount was less than the minimum CPC for currency.
	//
	BudgetErrorReasonMONEY_AMOUNT_LESS_THAN_CURRENCY_MINIMUM_CPC BudgetErrorReason = "MONEY_AMOUNT_LESS_THAN_CURRENCY_MINIMUM_CPC"

	//
	// A money amount was greater than the maximum allowed.
	//
	BudgetErrorReasonMONEY_AMOUNT_TOO_LARGE BudgetErrorReason = "MONEY_AMOUNT_TOO_LARGE"

	//
	// A money amount was negative.
	//
	BudgetErrorReasonNEGATIVE_MONEY_AMOUNT BudgetErrorReason = "NEGATIVE_MONEY_AMOUNT"

	//
	// A money amount was not a multiple of a minimum unit.
	//
	BudgetErrorReasonNON_MULTIPLE_OF_MINIMUM_CURRENCY_UNIT BudgetErrorReason = "NON_MULTIPLE_OF_MINIMUM_CURRENCY_UNIT"
)

//
// The reasons for the target error.
//
type CampaignErrorReason string

const (

	//
	// A complete campaign cannot go back to being incomplete
	//
	CampaignErrorReasonCANNOT_GO_BACK_TO_INCOMPLETE CampaignErrorReason = "CANNOT_GO_BACK_TO_INCOMPLETE"

	//
	// Cannot target content network.
	//
	CampaignErrorReasonCANNOT_TARGET_CONTENT_NETWORK CampaignErrorReason = "CANNOT_TARGET_CONTENT_NETWORK"

	//
	// Cannot target search network.
	//
	CampaignErrorReasonCANNOT_TARGET_SEARCH_NETWORK CampaignErrorReason = "CANNOT_TARGET_SEARCH_NETWORK"

	//
	// Cannot cover search network without google search network.
	//
	CampaignErrorReasonCANNOT_TARGET_SEARCH_NETWORK_WITHOUT_GOOGLE_SEARCH CampaignErrorReason = "CANNOT_TARGET_SEARCH_NETWORK_WITHOUT_GOOGLE_SEARCH"

	//
	// Cannot target Google Search network for a CPM campaign.
	//
	CampaignErrorReasonCANNOT_TARGET_GOOGLE_SEARCH_FOR_CPM_CAMPAIGN CampaignErrorReason = "CANNOT_TARGET_GOOGLE_SEARCH_FOR_CPM_CAMPAIGN"

	//
	// Must target at least one network.
	//
	CampaignErrorReasonCAMPAIGN_MUST_TARGET_AT_LEAST_ONE_NETWORK CampaignErrorReason = "CAMPAIGN_MUST_TARGET_AT_LEAST_ONE_NETWORK"

	//
	// Only some Google partners are allowed to target partner search network.
	//
	CampaignErrorReasonCANNOT_TARGET_PARTNER_SEARCH_NETWORK CampaignErrorReason = "CANNOT_TARGET_PARTNER_SEARCH_NETWORK"

	//
	// Cannot target content network only as campaign has criteria-level bidding strategy.
	//
	CampaignErrorReasonCANNOT_TARGET_CONTENT_NETWORK_ONLY_WITH_CRITERIA_LEVEL_BIDDING_STRATEGY CampaignErrorReason = "CANNOT_TARGET_CONTENT_NETWORK_ONLY_WITH_CRITERIA_LEVEL_BIDDING_STRATEGY"

	//
	// Cannot modify the start or end date such that the campaign duration would not contain the
	// durations of all runnable trials.
	//
	CampaignErrorReasonCAMPAIGN_DURATION_MUST_CONTAIN_ALL_RUNNABLE_TRIALS CampaignErrorReason = "CAMPAIGN_DURATION_MUST_CONTAIN_ALL_RUNNABLE_TRIALS"

	//
	// Cannot modify dates, budget or campaign name of a trial campaign.
	//
	CampaignErrorReasonCANNOT_MODIFY_FOR_TRIAL_CAMPAIGN CampaignErrorReason = "CANNOT_MODIFY_FOR_TRIAL_CAMPAIGN"

	//
	// Trying to modify the name of an active or paused campaign, where the name is already
	// assigned to another active or paused campaign.
	//
	CampaignErrorReasonDUPLICATE_CAMPAIGN_NAME CampaignErrorReason = "DUPLICATE_CAMPAIGN_NAME"

	//
	// Two fields are in conflicting modes.
	//
	CampaignErrorReasonINCOMPATIBLE_CAMPAIGN_FIELD CampaignErrorReason = "INCOMPATIBLE_CAMPAIGN_FIELD"

	//
	// Campaign name cannot be used.
	//
	CampaignErrorReasonINVALID_CAMPAIGN_NAME CampaignErrorReason = "INVALID_CAMPAIGN_NAME"

	//
	// Given status is invalid.
	//
	CampaignErrorReasonINVALID_AD_SERVING_OPTIMIZATION_STATUS CampaignErrorReason = "INVALID_AD_SERVING_OPTIMIZATION_STATUS"

	//
	// Error in the campaign level tracking url.
	//
	CampaignErrorReasonINVALID_TRACKING_URL CampaignErrorReason = "INVALID_TRACKING_URL"

	//
	// Cannot set both tracking url template and tracking setting. An user has to clear legacy
	// tracking setting in order to add tracking url template.
	//
	CampaignErrorReasonCANNOT_SET_BOTH_TRACKING_URL_TEMPLATE_AND_TRACKING_SETTING CampaignErrorReason = "CANNOT_SET_BOTH_TRACKING_URL_TEMPLATE_AND_TRACKING_SETTING"

	//
	// The maximum number of impressions for Frequency Cap should be an integer greater than 0.
	//
	CampaignErrorReasonMAX_IMPRESSIONS_NOT_IN_RANGE CampaignErrorReason = "MAX_IMPRESSIONS_NOT_IN_RANGE"

	//
	// Only the Day, Week and Month time units are supported.
	//
	CampaignErrorReasonTIME_UNIT_NOT_SUPPORTED CampaignErrorReason = "TIME_UNIT_NOT_SUPPORTED"

	//
	// Operation not allowed on a campaign whose serving status has ended
	//
	CampaignErrorReasonINVALID_OPERATION_IF_SERVING_STATUS_HAS_ENDED CampaignErrorReason = "INVALID_OPERATION_IF_SERVING_STATUS_HAS_ENDED"

	//
	// This budget is exclusively linked to a Campaign that is using @link{Experiment}s
	// so it cannot be shared.
	//
	CampaignErrorReasonBUDGET_CANNOT_BE_SHARED CampaignErrorReason = "BUDGET_CANNOT_BE_SHARED"

	//
	// Campaigns using @link{Experiment}s cannot use a shared budget.
	//
	CampaignErrorReasonCAMPAIGN_CANNOT_USE_SHARED_BUDGET CampaignErrorReason = "CAMPAIGN_CANNOT_USE_SHARED_BUDGET"

	//
	// A different budget cannot be assigned to a campaign when there are running or scheduled
	// trials.
	//
	CampaignErrorReasonCANNOT_CHANGE_BUDGET_ON_CAMPAIGN_WITH_TRIALS CampaignErrorReason = "CANNOT_CHANGE_BUDGET_ON_CAMPAIGN_WITH_TRIALS"

	//
	// No link found between the campaign and the label.
	//
	CampaignErrorReasonCAMPAIGN_LABEL_DOES_NOT_EXIST CampaignErrorReason = "CAMPAIGN_LABEL_DOES_NOT_EXIST"

	//
	// The label has already been attached to the campaign.
	//
	CampaignErrorReasonCAMPAIGN_LABEL_ALREADY_EXISTS CampaignErrorReason = "CAMPAIGN_LABEL_ALREADY_EXISTS"

	//
	// A ShoppingSetting was not found when creating a shopping campaign.
	//
	CampaignErrorReasonMISSING_SHOPPING_SETTING CampaignErrorReason = "MISSING_SHOPPING_SETTING"

	//
	// The country in shopping setting is not an allowed country.
	//
	CampaignErrorReasonINVALID_SHOPPING_SALES_COUNTRY CampaignErrorReason = "INVALID_SHOPPING_SALES_COUNTRY"

	//
	// Shopping merchant is not enabled for Purchases on Google.
	//
	CampaignErrorReasonSHOPPING_MERCHANT_NOT_ALLOWED_FOR_PURCHASES_ON_GOOGLE CampaignErrorReason = "SHOPPING_MERCHANT_NOT_ALLOWED_FOR_PURCHASES_ON_GOOGLE"

	//
	// Purchases on Google not enabled for the shopping campaign's sales country.
	//
	CampaignErrorReasonPURCHASES_ON_GOOGLE_NOT_SUPPORTED_FOR_SHOPPING_SALES_COUNTRY CampaignErrorReason = "PURCHASES_ON_GOOGLE_NOT_SUPPORTED_FOR_SHOPPING_SALES_COUNTRY"

	//
	// A Campaign with channel sub type UNIVERSAL_APP_CAMPAIGN must have a
	// UniversalAppCampaignSetting specified.
	//
	CampaignErrorReasonMISSING_UNIVERSAL_APP_CAMPAIGN_SETTING CampaignErrorReason = "MISSING_UNIVERSAL_APP_CAMPAIGN_SETTING"

	//
	// The requested channel type is not available according to the customer's account setting.
	//
	CampaignErrorReasonADVERTISING_CHANNEL_TYPE_NOT_AVAILABLE_FOR_ACCOUNT_TYPE CampaignErrorReason = "ADVERTISING_CHANNEL_TYPE_NOT_AVAILABLE_FOR_ACCOUNT_TYPE"

	//
	// The AdvertisingChannelSubType is not a valid subtype of the primary channel type.
	//
	CampaignErrorReasonINVALID_ADVERTISING_CHANNEL_SUB_TYPE CampaignErrorReason = "INVALID_ADVERTISING_CHANNEL_SUB_TYPE"

	//
	// At least one conversion must be selected.
	//
	CampaignErrorReasonAT_LEAST_ONE_CONVERSION_MUST_BE_SELECTED CampaignErrorReason = "AT_LEAST_ONE_CONVERSION_MUST_BE_SELECTED"

	//
	// Setting ad rotation mode for a campaign is not allowed.
	// Ad rotation mode at campaign is deprecated.
	//
	CampaignErrorReasonCANNOT_SET_AD_ROTATION_MODE CampaignErrorReason = "CANNOT_SET_AD_ROTATION_MODE"

	//
	// Default error
	// <span class="constraint Rejected">Used for return value only. An enumeration could not be processed, typically due to incompatibility with your WSDL version.</span>
	//
	CampaignErrorReasonUNKNOWN CampaignErrorReason = "UNKNOWN"
)

//
// Enums for the various reasons an error can be thrown as a result of
// ClientTerms violation.
//
type ClientTermsErrorReason string

const (

	//
	// Customer has not agreed to the latest AdWords Terms & Conditions
	//
	ClientTermsErrorReasonINCOMPLETE_SIGNUP_CURRENT_ADWORDS_TNC_NOT_AGREED ClientTermsErrorReason = "INCOMPLETE_SIGNUP_CURRENT_ADWORDS_TNC_NOT_AGREED"
)

//
// The reasons for the database error.
//
type DatabaseErrorReason string

const (

	//
	// A concurrency problem occurred as two threads were attempting to modify same object.
	//
	DatabaseErrorReasonCONCURRENT_MODIFICATION DatabaseErrorReason = "CONCURRENT_MODIFICATION"

	//
	// The permission was denied to access an object.
	//
	DatabaseErrorReasonPERMISSION_DENIED DatabaseErrorReason = "PERMISSION_DENIED"

	//
	// The user's access to an object has been prohibited.
	//
	DatabaseErrorReasonACCESS_PROHIBITED DatabaseErrorReason = "ACCESS_PROHIBITED"

	//
	// Requested campaign belongs to a product that is not supported by the api.
	//
	DatabaseErrorReasonCAMPAIGN_PRODUCT_NOT_SUPPORTED DatabaseErrorReason = "CAMPAIGN_PRODUCT_NOT_SUPPORTED"

	//
	// a duplicate key was detected upon insertion
	//
	DatabaseErrorReasonDUPLICATE_KEY DatabaseErrorReason = "DUPLICATE_KEY"

	//
	// a database error has occurred
	//
	DatabaseErrorReasonDATABASE_ERROR DatabaseErrorReason = "DATABASE_ERROR"

	//
	// an unknown error has occurred
	//
	DatabaseErrorReasonUNKNOWN DatabaseErrorReason = "UNKNOWN"
)

//
// The reasons for the target error.
//
type DateErrorReason string

const (

	//
	// Given field values do not correspond to a valid date.
	//
	DateErrorReasonINVALID_FIELD_VALUES_IN_DATE DateErrorReason = "INVALID_FIELD_VALUES_IN_DATE"

	//
	// Given field values do not correspond to a valid date time.
	//
	DateErrorReasonINVALID_FIELD_VALUES_IN_DATE_TIME DateErrorReason = "INVALID_FIELD_VALUES_IN_DATE_TIME"

	//
	// The string date's format should be yyyymmdd.
	//
	DateErrorReasonINVALID_STRING_DATE DateErrorReason = "INVALID_STRING_DATE"

	//
	// The string date range's format should be yyyymmdd yyyymmdd.
	//
	DateErrorReasonINVALID_STRING_DATE_RANGE DateErrorReason = "INVALID_STRING_DATE_RANGE"

	//
	// The string date time's format should be yyyymmdd hhmmss [tz].
	//
	DateErrorReasonINVALID_STRING_DATE_TIME DateErrorReason = "INVALID_STRING_DATE_TIME"

	//
	// Date is before allowed minimum.
	//
	DateErrorReasonEARLIER_THAN_MINIMUM_DATE DateErrorReason = "EARLIER_THAN_MINIMUM_DATE"

	//
	// Date is after allowed maximum.
	//
	DateErrorReasonLATER_THAN_MAXIMUM_DATE DateErrorReason = "LATER_THAN_MAXIMUM_DATE"

	//
	// Date range bounds are not in order.
	//
	DateErrorReasonDATE_RANGE_MINIMUM_DATE_LATER_THAN_MAXIMUM_DATE DateErrorReason = "DATE_RANGE_MINIMUM_DATE_LATER_THAN_MAXIMUM_DATE"

	//
	// Both dates in range are null.
	//
	DateErrorReasonDATE_RANGE_MINIMUM_AND_MAXIMUM_DATES_BOTH_NULL DateErrorReason = "DATE_RANGE_MINIMUM_AND_MAXIMUM_DATES_BOTH_NULL"
)

//
// The reasons for the date range error.
//
type DateRangeErrorReason string

const (
	DateRangeErrorReasonDATE_RANGE_ERROR DateRangeErrorReason = "DATE_RANGE_ERROR"

	//
	// Invalid date.
	//
	DateRangeErrorReasonINVALID_DATE DateRangeErrorReason = "INVALID_DATE"

	//
	// The start date was after the end date.
	//
	DateRangeErrorReasonSTART_DATE_AFTER_END_DATE DateRangeErrorReason = "START_DATE_AFTER_END_DATE"

	//
	// Cannot set date to past time
	//
	DateRangeErrorReasonCANNOT_SET_DATE_TO_PAST DateRangeErrorReason = "CANNOT_SET_DATE_TO_PAST"

	//
	// A date was used that is past the system "last" date.
	//
	DateRangeErrorReasonAFTER_MAXIMUM_ALLOWABLE_DATE DateRangeErrorReason = "AFTER_MAXIMUM_ALLOWABLE_DATE"

	//
	// Trying to change start date on a campaign that has started.
	//
	DateRangeErrorReasonCANNOT_MODIFY_START_DATE_IF_ALREADY_STARTED DateRangeErrorReason = "CANNOT_MODIFY_START_DATE_IF_ALREADY_STARTED"
)

//
// The reasons for the validation error.
//
type DistinctErrorReason string

const (
	DistinctErrorReasonDUPLICATE_ELEMENT DistinctErrorReason = "DUPLICATE_ELEMENT"

	DistinctErrorReasonDUPLICATE_TYPE DistinctErrorReason = "DUPLICATE_TYPE"
)

//
// Limits at various levels of the account.
//
type EntityCountLimitExceededReason string

const (

	//
	// Indicates that this request would exceed the number of allowed entities for the AdWords
	// account. The exact entity type and limit being checked can be inferred from
	// {@link #accountLimitType}.
	//
	EntityCountLimitExceededReasonACCOUNT_LIMIT EntityCountLimitExceededReason = "ACCOUNT_LIMIT"

	//
	// Indicates that this request would exceed the number of allowed entities in a Campaign.
	// The exact entity type and limit being checked can be inferred from
	// {@link #accountLimitType}, and the numeric id of the Campaign involved is given by
	// {@link #enclosingId}.
	//
	EntityCountLimitExceededReasonCAMPAIGN_LIMIT EntityCountLimitExceededReason = "CAMPAIGN_LIMIT"

	//
	// Indicates that this request would exceed the number of allowed entities in
	// an ad group.  The exact entity type and limit being checked can be
	// inferred from {@link #accountLimitType}, and the numeric id of the ad group
	// involved is given by {@link #enclosingId}.
	//
	EntityCountLimitExceededReasonADGROUP_LIMIT EntityCountLimitExceededReason = "ADGROUP_LIMIT"

	//
	// Indicates that this request would exceed the number of allowed entities in an ad group ad.
	// The exact entity type and limit being checked can be inferred from {@link #accountLimitType},
	// and the {@link #enclosingId} contains the ad group id followed by the ad id, separated by a
	// single comma (,).
	//
	EntityCountLimitExceededReasonAD_GROUP_AD_LIMIT EntityCountLimitExceededReason = "AD_GROUP_AD_LIMIT"

	//
	// Indicates that this request would exceed the number of allowed entities in an ad group
	// criterion.  The exact entity type and limit being checked can be inferred from
	// {@link #accountLimitType}, and the {@link #enclosingId} contains the ad group id followed by
	// the criterion id, separated by a single comma (,).
	//
	EntityCountLimitExceededReasonAD_GROUP_CRITERION_LIMIT EntityCountLimitExceededReason = "AD_GROUP_CRITERION_LIMIT"

	//
	// Indicates that this request would exceed the number of allowed entities in
	// this shared set.  The exact entity type and limit being checked can be
	// inferred from {@link #accountLimitType}, and the numeric id of the shared
	// set involved is given by {@link #enclosingId}.
	//
	EntityCountLimitExceededReasonSHARED_SET_LIMIT EntityCountLimitExceededReason = "SHARED_SET_LIMIT"

	//
	// Exceeds a limit related to a matching function.
	//
	EntityCountLimitExceededReasonMATCHING_FUNCTION_LIMIT EntityCountLimitExceededReason = "MATCHING_FUNCTION_LIMIT"

	//
	// Specific limit that has been exceeded is unknown (the client may be of an
	// older version than the server).
	//
	EntityCountLimitExceededReasonUNKNOWN EntityCountLimitExceededReason = "UNKNOWN"
)

type EntityNotFoundReason string

const (

	//
	// The specified id refered to an entity which either doesn't exist or is not accessible to the
	// customer. e.g. campaign belongs to another customer.
	//
	EntityNotFoundReasonINVALID_ID EntityNotFoundReason = "INVALID_ID"
)

//
// The reasons for the target error.
//
type IdErrorReason string

const (

	//
	// Id not found
	//
	IdErrorReasonNOT_FOUND IdErrorReason = "NOT_FOUND"
)

//
// The single reason for the internal API error.
//
type InternalApiErrorReason string

const (

	//
	// API encountered an unexpected internal error.
	//
	InternalApiErrorReasonUNEXPECTED_INTERNAL_API_ERROR InternalApiErrorReason = "UNEXPECTED_INTERNAL_API_ERROR"

	//
	// A temporary error occurred during the request. Please retry.
	//
	InternalApiErrorReasonTRANSIENT_ERROR InternalApiErrorReason = "TRANSIENT_ERROR"

	//
	// The cause of the error is not known or only defined in newer versions.
	//
	InternalApiErrorReasonUNKNOWN InternalApiErrorReason = "UNKNOWN"

	//
	// The API is currently unavailable for a planned downtime.
	//
	InternalApiErrorReasonDOWNTIME InternalApiErrorReason = "DOWNTIME"

	//
	// Mutate succeeded but server was unable to build response. Client should not retry mutate.
	//
	InternalApiErrorReasonERROR_GENERATING_RESPONSE InternalApiErrorReason = "ERROR_GENERATING_RESPONSE"
)

//
// The reasons for the validation error.
//
type NotEmptyErrorReason string

const (
	NotEmptyErrorReasonEMPTY_LIST NotEmptyErrorReason = "EMPTY_LIST"
)

//
// The reasons for the validation error.
//
type NullErrorReason string

const (

	//
	// Specified list/container must not contain any null elements
	//
	NullErrorReasonNULL_CONTENT NullErrorReason = "NULL_CONTENT"
)

//
// The reasons for the operation access error.
//
type OperationAccessDeniedReason string

const (

	//
	// Unauthorized invocation of a service's method (get, mutate, etc.)
	//
	OperationAccessDeniedReasonACTION_NOT_PERMITTED OperationAccessDeniedReason = "ACTION_NOT_PERMITTED"

	//
	// Unauthorized ADD operation in invoking a service's mutate method.
	//
	OperationAccessDeniedReasonADD_OPERATION_NOT_PERMITTED OperationAccessDeniedReason = "ADD_OPERATION_NOT_PERMITTED"

	//
	// Unauthorized REMOVE operation in invoking a service's mutate method.
	//
	OperationAccessDeniedReasonREMOVE_OPERATION_NOT_PERMITTED OperationAccessDeniedReason = "REMOVE_OPERATION_NOT_PERMITTED"

	//
	// Unauthorized SET operation in invoking a service's mutate method.
	//
	OperationAccessDeniedReasonSET_OPERATION_NOT_PERMITTED OperationAccessDeniedReason = "SET_OPERATION_NOT_PERMITTED"

	//
	// A mutate action is not allowed on this campaign, from this client.
	//
	OperationAccessDeniedReasonMUTATE_ACTION_NOT_PERMITTED_FOR_CLIENT OperationAccessDeniedReason = "MUTATE_ACTION_NOT_PERMITTED_FOR_CLIENT"

	//
	// This operation is not permitted on this campaign type
	//
	OperationAccessDeniedReasonOPERATION_NOT_PERMITTED_FOR_CAMPAIGN_TYPE OperationAccessDeniedReason = "OPERATION_NOT_PERMITTED_FOR_CAMPAIGN_TYPE"

	//
	// An ADD operation may not set status to REMOVED.
	//
	OperationAccessDeniedReasonADD_AS_REMOVED_NOT_PERMITTED OperationAccessDeniedReason = "ADD_AS_REMOVED_NOT_PERMITTED"

	//
	// This operation is not allowed because the campaign or adgroup is removed.
	//
	OperationAccessDeniedReasonOPERATION_NOT_PERMITTED_FOR_REMOVED_ENTITY OperationAccessDeniedReason = "OPERATION_NOT_PERMITTED_FOR_REMOVED_ENTITY"

	//
	// This operation is not permitted on this ad group type.
	//
	OperationAccessDeniedReasonOPERATION_NOT_PERMITTED_FOR_AD_GROUP_TYPE OperationAccessDeniedReason = "OPERATION_NOT_PERMITTED_FOR_AD_GROUP_TYPE"

	//
	// The reason the invoked method or operation is prohibited is not known
	// (the client may be of an older version than the server).
	//
	OperationAccessDeniedReasonUNKNOWN OperationAccessDeniedReason = "UNKNOWN"
)

//
// This represents an operator that may be presented to an adsapi service.
//
type Operator string

const (

	//
	// The ADD operator.
	//
	OperatorADD Operator = "ADD"

	//
	// The REMOVE operator.
	//
	OperatorREMOVE Operator = "REMOVE"

	//
	// The SET operator (used for updates).
	//
	OperatorSET Operator = "SET"
)

//
// The reasons for the validation error.
//
type OperatorErrorReason string

const (
	OperatorErrorReasonOPERATOR_NOT_SUPPORTED OperatorErrorReason = "OPERATOR_NOT_SUPPORTED"
)

//
// Defines the valid set of operators.
//
type PredicateOperator string

const (

	//
	// Checks if the field is equal to the given value.
	//
	// <p>This operator is used with integers, dates, booleans, strings,
	// enums, and sets.
	//
	PredicateOperatorEQUALS PredicateOperator = "EQUALS"

	//
	// Checks if the field does not equal the given value.
	//
	// <p>This operator is used with integers, booleans, strings, enums,
	// and sets.
	//
	PredicateOperatorNOT_EQUALS PredicateOperator = "NOT_EQUALS"

	//
	// Checks if the field is equal to one of the given values.
	//
	// <p>This operator accepts multiple operands and is used with
	// integers, booleans, strings, and enums.
	//
	PredicateOperatorIN PredicateOperator = "IN"

	//
	// Checks if the field does not equal any of the given values.
	//
	// <p>This operator accepts multiple operands and is used with
	// integers, booleans, strings, and enums.
	//
	PredicateOperatorNOT_IN PredicateOperator = "NOT_IN"

	//
	// Checks if the field is greater than the given value.
	//
	// <p>This operator is used with numbers and dates.
	//
	PredicateOperatorGREATER_THAN PredicateOperator = "GREATER_THAN"

	//
	// Checks if the field is greater or equal to the given value.
	//
	// <p>This operator is used with numbers and dates.
	//
	PredicateOperatorGREATER_THAN_EQUALS PredicateOperator = "GREATER_THAN_EQUALS"

	//
	// Checks if the field is less than the given value.
	//
	// <p>This operator is used with numbers and dates.
	//
	PredicateOperatorLESS_THAN PredicateOperator = "LESS_THAN"

	//
	// Checks if the field is less or equal to than the given value.
	//
	// <p>This operator is used with numbers and dates.
	//
	PredicateOperatorLESS_THAN_EQUALS PredicateOperator = "LESS_THAN_EQUALS"

	//
	// Checks if the field starts with the given value.
	//
	// <p>This operator is used with strings.
	//
	PredicateOperatorSTARTS_WITH PredicateOperator = "STARTS_WITH"

	//
	// Checks if the field starts with the given value, ignoring case.
	//
	// <p>This operator is used with strings.
	//
	PredicateOperatorSTARTS_WITH_IGNORE_CASE PredicateOperator = "STARTS_WITH_IGNORE_CASE"

	//
	// Checks if the field contains the given value as a substring.
	//
	// <p>This operator is used with strings.
	//
	PredicateOperatorCONTAINS PredicateOperator = "CONTAINS"

	//
	// Checks if the field contains the given value as a substring, ignoring
	// case.
	//
	// <p>This operator is used with strings.
	//
	PredicateOperatorCONTAINS_IGNORE_CASE PredicateOperator = "CONTAINS_IGNORE_CASE"

	//
	// Checks if the field does not contain the given value as a substring.
	//
	// <p>This operator is used with strings.
	//
	PredicateOperatorDOES_NOT_CONTAIN PredicateOperator = "DOES_NOT_CONTAIN"

	//
	// Checks if the field does not contain the given value as a substring,
	// ignoring case.
	//
	// <p>This operator is used with strings.
	//
	PredicateOperatorDOES_NOT_CONTAIN_IGNORE_CASE PredicateOperator = "DOES_NOT_CONTAIN_IGNORE_CASE"

	//
	// Checks if the field contains <em>any</em> of the given values.
	//
	// <p>This operator accepts multiple values and is used on sets of numbers
	// or strings.
	//
	PredicateOperatorCONTAINS_ANY PredicateOperator = "CONTAINS_ANY"

	//
	// Checks if the field contains <em>all</em> of the given values.
	//
	// <p>This operator accepts multiple values and is used on sets of numbers
	// or strings.
	//
	PredicateOperatorCONTAINS_ALL PredicateOperator = "CONTAINS_ALL"

	//
	// Checks if the field contains <em>none</em> of the given values.
	//
	// <p>This operator accepts multiple values and is used on sets of numbers
	// or strings.
	//
	PredicateOperatorCONTAINS_NONE PredicateOperator = "CONTAINS_NONE"

	//
	// <span class="constraint Rejected">Used for return value only. An enumeration could not be processed, typically due to incompatibility with your WSDL version.</span>
	//
	PredicateOperatorUNKNOWN PredicateOperator = "UNKNOWN"
)

//
// Enums for all the reasons an error can be thrown to the user during
// billing quota checks.
//
type QuotaCheckErrorReason string

const (

	//
	// Customer passed in an invalid token in the header.
	//
	QuotaCheckErrorReasonINVALID_TOKEN_HEADER QuotaCheckErrorReason = "INVALID_TOKEN_HEADER"

	//
	// Customer is marked delinquent.
	//
	QuotaCheckErrorReasonACCOUNT_DELINQUENT QuotaCheckErrorReason = "ACCOUNT_DELINQUENT"

	//
	// Customer is a fraudulent.
	//
	QuotaCheckErrorReasonACCOUNT_INACCESSIBLE QuotaCheckErrorReason = "ACCOUNT_INACCESSIBLE"

	//
	// Inactive Account.
	//
	QuotaCheckErrorReasonACCOUNT_INACTIVE QuotaCheckErrorReason = "ACCOUNT_INACTIVE"

	//
	// Signup not complete
	//
	QuotaCheckErrorReasonINCOMPLETE_SIGNUP QuotaCheckErrorReason = "INCOMPLETE_SIGNUP"

	//
	// Developer token is not approved for production access, and the customer
	// is attempting to access a production account.
	//
	QuotaCheckErrorReasonDEVELOPER_TOKEN_NOT_APPROVED QuotaCheckErrorReason = "DEVELOPER_TOKEN_NOT_APPROVED"

	//
	// Terms and conditions are not signed.
	//
	QuotaCheckErrorReasonTERMS_AND_CONDITIONS_NOT_SIGNED QuotaCheckErrorReason = "TERMS_AND_CONDITIONS_NOT_SIGNED"

	//
	// Monthly budget quota reached.
	//
	QuotaCheckErrorReasonMONTHLY_BUDGET_REACHED QuotaCheckErrorReason = "MONTHLY_BUDGET_REACHED"

	//
	// Monthly budget quota exceeded.
	//
	QuotaCheckErrorReasonQUOTA_EXCEEDED QuotaCheckErrorReason = "QUOTA_EXCEEDED"
)

//
// The reasons for the target error.
//
type RangeErrorReason string

const (
	RangeErrorReasonTOO_LOW RangeErrorReason = "TOO_LOW"

	RangeErrorReasonTOO_HIGH RangeErrorReason = "TOO_HIGH"
)

//
// The reason for the rate exceeded error.
//
type RateExceededErrorReason string

const (

	//
	// Rate exceeded.
	//
	RateExceededErrorReasonRATE_EXCEEDED RateExceededErrorReason = "RATE_EXCEEDED"
)

//
// The reasons for the target error.
//
type ReadOnlyErrorReason string

const (
	ReadOnlyErrorReasonREAD_ONLY ReadOnlyErrorReason = "READ_ONLY"
)

//
// The reasons for the target error.
//
type RejectedErrorReason string

const (

	//
	// Unknown value encountered
	//
	RejectedErrorReasonUNKNOWN_VALUE RejectedErrorReason = "UNKNOWN_VALUE"
)

type RequestErrorReason string

const (

	//
	// Error reason is unknown.
	//
	RequestErrorReasonUNKNOWN RequestErrorReason = "UNKNOWN"

	//
	// Invalid input.
	//
	RequestErrorReasonINVALID_INPUT RequestErrorReason = "INVALID_INPUT"

	//
	// The api version in the request has been discontinued. Please update
	// to the new AdWords API version.
	//
	RequestErrorReasonUNSUPPORTED_VERSION RequestErrorReason = "UNSUPPORTED_VERSION"
)

//
// The reasons for the target error.
//
type RequiredErrorReason string

const (

	//
	// Missing required field.
	//
	RequiredErrorReasonREQUIRED RequiredErrorReason = "REQUIRED"
)

//
// The reasons for the target error.
//
type SelectorErrorReason string

const (

	//
	// The field name is not valid.
	//
	SelectorErrorReasonINVALID_FIELD_NAME SelectorErrorReason = "INVALID_FIELD_NAME"

	//
	// The list of fields is null or empty.
	//
	SelectorErrorReasonMISSING_FIELDS SelectorErrorReason = "MISSING_FIELDS"

	//
	// The list of predicates is null or empty.
	//
	SelectorErrorReasonMISSING_PREDICATES SelectorErrorReason = "MISSING_PREDICATES"

	//
	// Predicate operator does not support multiple values. Multiple values are
	// supported only for {@link Predicate.Operator#IN} and {@link Predicate.Operator#NOT_IN}.
	//
	SelectorErrorReasonOPERATOR_DOES_NOT_SUPPORT_MULTIPLE_VALUES SelectorErrorReason = "OPERATOR_DOES_NOT_SUPPORT_MULTIPLE_VALUES"

	//
	// The predicate enum value is not valid.
	//
	SelectorErrorReasonINVALID_PREDICATE_ENUM_VALUE SelectorErrorReason = "INVALID_PREDICATE_ENUM_VALUE"

	//
	// The predicate operator is empty.
	//
	SelectorErrorReasonMISSING_PREDICATE_OPERATOR SelectorErrorReason = "MISSING_PREDICATE_OPERATOR"

	//
	// The predicate values are empty.
	//
	SelectorErrorReasonMISSING_PREDICATE_VALUES SelectorErrorReason = "MISSING_PREDICATE_VALUES"

	//
	// The predicate field name is not valid.
	//
	SelectorErrorReasonINVALID_PREDICATE_FIELD_NAME SelectorErrorReason = "INVALID_PREDICATE_FIELD_NAME"

	//
	// The predicate operator is not valid.
	//
	SelectorErrorReasonINVALID_PREDICATE_OPERATOR SelectorErrorReason = "INVALID_PREDICATE_OPERATOR"

	//
	// Invalid selection of fields.
	//
	SelectorErrorReasonINVALID_FIELD_SELECTION SelectorErrorReason = "INVALID_FIELD_SELECTION"

	//
	// The predicate value is not valid.
	//
	SelectorErrorReasonINVALID_PREDICATE_VALUE SelectorErrorReason = "INVALID_PREDICATE_VALUE"

	//
	// The sort field name is not valid or the field is not sortable.
	//
	SelectorErrorReasonINVALID_SORT_FIELD_NAME SelectorErrorReason = "INVALID_SORT_FIELD_NAME"

	//
	// Standard error.
	//
	SelectorErrorReasonSELECTOR_ERROR SelectorErrorReason = "SELECTOR_ERROR"

	//
	// Filtering by date range is not supported.
	//
	SelectorErrorReasonFILTER_BY_DATE_RANGE_NOT_SUPPORTED SelectorErrorReason = "FILTER_BY_DATE_RANGE_NOT_SUPPORTED"

	//
	// Selector paging start index is too high.
	//
	SelectorErrorReasonSTART_INDEX_IS_TOO_HIGH SelectorErrorReason = "START_INDEX_IS_TOO_HIGH"

	//
	// The values list in a predicate was too long.
	//
	SelectorErrorReasonTOO_MANY_PREDICATE_VALUES SelectorErrorReason = "TOO_MANY_PREDICATE_VALUES"

	//
	// <span class="constraint Rejected">Used for return value only. An enumeration could not be processed, typically due to incompatibility with your WSDL version.</span>
	//
	SelectorErrorReasonUNKNOWN_ERROR SelectorErrorReason = "UNKNOWN_ERROR"
)

//
// The reasons for Ad Scheduling errors.
//
type SizeLimitErrorReason string

const (

	//
	// The number of entries in the request exceeds the system limit.
	//
	SizeLimitErrorReasonREQUEST_SIZE_LIMIT_EXCEEDED SizeLimitErrorReason = "REQUEST_SIZE_LIMIT_EXCEEDED"

	//
	// The number of entries in the response exceeds the system limit.
	//
	SizeLimitErrorReasonRESPONSE_SIZE_LIMIT_EXCEEDED SizeLimitErrorReason = "RESPONSE_SIZE_LIMIT_EXCEEDED"

	//
	// The account is too large to load.
	//
	SizeLimitErrorReasonINTERNAL_STORAGE_ERROR SizeLimitErrorReason = "INTERNAL_STORAGE_ERROR"

	//
	// <span class="constraint Rejected">Used for return value only. An enumeration could not be processed, typically due to incompatibility with your WSDL version.</span>
	//
	SizeLimitErrorReasonUNKNOWN SizeLimitErrorReason = "UNKNOWN"
)

//
// Possible orders of sorting.
//
type SortOrder string

const (
	SortOrderASCENDING SortOrder = "ASCENDING"

	SortOrderDESCENDING SortOrder = "DESCENDING"
)

//
// The reasons for the target error.
//
type StringFormatErrorReason string

const (
	StringFormatErrorReasonUNKNOWN StringFormatErrorReason = "UNKNOWN"

	//
	// The input string value contains disallowed characters.
	//
	StringFormatErrorReasonILLEGAL_CHARS StringFormatErrorReason = "ILLEGAL_CHARS"

	//
	// The input string value is invalid for the associated field.
	//
	StringFormatErrorReasonINVALID_FORMAT StringFormatErrorReason = "INVALID_FORMAT"
)

//
// The reasons for the target error.
//
type StringLengthErrorReason string

const (
	StringLengthErrorReasonTOO_SHORT StringLengthErrorReason = "TOO_SHORT"

	StringLengthErrorReasonTOO_LONG StringLengthErrorReason = "TOO_LONG"
)

//
// Error codes defined by {@link TrialError}.
//
type TrialErrorReason string

const (

	//
	// <span class="constraint Rejected">Used for return value only. An enumeration could not be processed, typically due to incompatibility with your WSDL version.</span>
	//
	TrialErrorReasonUNKNOWN TrialErrorReason = "UNKNOWN"

	//
	// Trial status cannot be updated from the current status to the requested target status.
	//
	TrialErrorReasonINVALID_STATUS_TRANSITION TrialErrorReason = "INVALID_STATUS_TRANSITION"

	//
	// Cannot create a trial from a campaign using an explicitly shared budget.
	//
	TrialErrorReasonCANNOT_USE_TRIAL_WITH_SHARED_BUDGET TrialErrorReason = "CANNOT_USE_TRIAL_WITH_SHARED_BUDGET"

	//
	// Cannot create a trial as long as the campaign has a running or scheduled Advertiser Campaign
	// Experiment.
	//
	TrialErrorReasonCANNOT_CREATE_TRIAL_WHEN_CAMPAIGN_HAS_ACTIVE_EXPERIMENTS TrialErrorReason = "CANNOT_CREATE_TRIAL_WHEN_CAMPAIGN_HAS_ACTIVE_EXPERIMENTS"

	//
	// Cannot create a trial for a base campaign, which is deleted.
	//
	TrialErrorReasonCANNOT_CREATE_TRIAL_FOR_DELETED_BASE_CAMPAIGN TrialErrorReason = "CANNOT_CREATE_TRIAL_FOR_DELETED_BASE_CAMPAIGN"

	//
	// Cannot create a trial from a draft, which has a status other than proposed.
	//
	TrialErrorReasonCANNOT_CREATE_TRIAL_FOR_NON_PROPOSED_DRAFT TrialErrorReason = "CANNOT_CREATE_TRIAL_FOR_NON_PROPOSED_DRAFT"

	//
	// This customer is not allowed to create a trial.
	//
	TrialErrorReasonCUSTOMER_CANNOT_CREATE_TRIAL TrialErrorReason = "CUSTOMER_CANNOT_CREATE_TRIAL"

	//
	// This campaign is not allowed to create a trial.
	//
	TrialErrorReasonCAMPAIGN_CANNOT_CREATE_TRIAL TrialErrorReason = "CAMPAIGN_CANNOT_CREATE_TRIAL"

	//
	// Trying to use a trial name which is already assigned to another campaign or trial.
	//
	TrialErrorReasonNAME_ALREADY_IN_USE TrialErrorReason = "NAME_ALREADY_IN_USE"

	//
	// Trying to set a trial duration which overlaps with another trial.
	//
	TrialErrorReasonTRIAL_DURATIONS_MUST_NOT_OVERLAP TrialErrorReason = "TRIAL_DURATIONS_MUST_NOT_OVERLAP"

	//
	// All non-archived trials must start and end within their campaign's duration.
	//
	TrialErrorReasonTRIAL_DURATION_MUST_BE_WITHIN_CAMPAIGN_DURATION TrialErrorReason = "TRIAL_DURATION_MUST_BE_WITHIN_CAMPAIGN_DURATION"
)

//
// Status of a trial.
//
type TrialStatus string

const (

	//
	// Invalid status. Should not be used except for detecting values that are
	// incorrect, or values that are not yet known to the user.
	// <span class="constraint Rejected">Used for return value only. An enumeration could not be processed, typically due to incompatibility with your WSDL version.</span>
	//
	TrialStatusUNKNOWN TrialStatus = "UNKNOWN"

	//
	// The trial campaign is being created.
	//
	TrialStatusCREATING TrialStatus = "CREATING"

	//
	// The trial campaign is fully created. The trial is currently running, scheduled
	// to run in the future or has ended based on its end date.The advertiser cannot
	// set this status directly. A trial with the status CREATING will be updated to
	// ACTIVE when it is fully created.
	//
	TrialStatusACTIVE TrialStatus = "ACTIVE"

	//
	// The advertiser requested to merge changes in the trial back into the original
	// campaigns. The update to the original campaign will be kicked off asynchronously
	// and the status will be updated to PROMOTED or PROMOTE_FAILED upon completion.
	//
	TrialStatusPROMOTING TrialStatus = "PROMOTING"

	//
	// The process to merge changes in the trial back to the original campaign has
	// completedly successfully. The advertiser cannot set this status directly. To
	// move the trial to this status, set the trial to status PROMOTING and the status
	// will be updated to PROMOTED when the changes are applied to the original
	// campaign.
	//
	TrialStatusPROMOTED TrialStatus = "PROMOTED"

	//
	// The advertiser archived the campaign trial.
	//
	TrialStatusARCHIVED TrialStatus = "ARCHIVED"

	//
	// The trial campaign failed to create. More details about the errors are available
	// through getErrors in the TrialService API.The advertiser cannot set this status
	// directly.
	//
	TrialStatusCREATION_FAILED TrialStatus = "CREATION_FAILED"

	//
	// The promotion failed after it was partially applied. Promote cannot be attempted
	// again safely, so the issue must be corrected in the original campaign. More
	// details about the errors are available through getErrors in the TrialService
	// API.The advertiser cannot set this status directly. To promote the trial, set
	// the trial in state PROMOTING and the status will be updated to PROMOTE_FAILED if
	// errors are encountered while applying changes to the original campaign.
	//
	TrialStatusPROMOTE_FAILED TrialStatus = "PROMOTE_FAILED"

	//
	// The advertiser has graduated the trial campaign to a standalone campaign,
	// existing independently of the trial.
	//
	TrialStatusGRADUATED TrialStatus = "GRADUATED"

	//
	// The advertiser has halted the trial.
	//
	TrialStatusHALTED TrialStatus = "HALTED"
)

type Get struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 get"`

	//
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	//
	Selector *Selector `xml:"selector,omitempty"`
}

type GetResponse struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 getResponse"`

	Rval *TrialPage `xml:"rval,omitempty"`
}

type Mutate struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 mutate"`

	//
	// <span class="constraint ContentsNotNull">This field must not contain {@code null} elements.</span>
	// <span class="constraint DistinctIds">Elements in this field must have distinct IDs for following {@link Operator}s : SET, REMOVE.</span>
	// <span class="constraint NotEmpty">This field must contain at least one element.</span>
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	//
	Operations []*TrialOperation `xml:"operations,omitempty"`
}

type MutateResponse struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 mutateResponse"`

	Rval *TrialReturnValue `xml:"rval,omitempty"`
}

type Query struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 query"`

	//
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	//
	Query string `xml:"query,omitempty"`
}

type QueryResponse struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 queryResponse"`

	Rval *TrialPage `xml:"rval,omitempty"`
}

type ApiError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 ApiError"`

	//
	// The OGNL field path to identify cause of error.
	//
	FieldPath string `xml:"fieldPath,omitempty"`

	//
	// A parsed copy of the field path. For example, the field path "operations[1].operand"
	// corresponds to this list: {FieldPathElement(field = "operations", index = 1),
	// FieldPathElement(field = "operand", index = null)}.
	//
	FieldPathElements []*FieldPathElement `xml:"fieldPathElements,omitempty"`

	//
	// The data that caused the error.
	//
	Trigger string `xml:"trigger,omitempty"`

	//
	// A simple string representation of the error and reason.
	//
	ErrorString string `xml:"errorString,omitempty"`

	//
	// Indicates that this instance is a subtype of ApiError.
	// Although this field is returned in the response, it is ignored on input
	// and cannot be selected. Specify xsi:type instead.
	//
	ApiErrorType string `xml:"ApiError.Type,omitempty"`
}

type ApiException struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 ApiException"`

	*ApplicationException

	//
	// List of errors.
	//
	Errors []*ApiError `xml:"errors,omitempty"`
}

type ApplicationException struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 ApplicationException"`

	//
	// Error message.
	//
	Message string `xml:"message,omitempty"`

	//
	// Indicates that this instance is a subtype of ApplicationException.
	// Although this field is returned in the response, it is ignored on input
	// and cannot be selected. Specify xsi:type instead.
	//
	ApplicationExceptionType string `xml:"ApplicationException.Type,omitempty"`
}

type AuthenticationError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 AuthenticationError"`

	*ApiError

	//
	// The error reason represented by an enum.
	//
	Reason *AuthenticationErrorReason `xml:"reason,omitempty"`
}

type AuthorizationError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 AuthorizationError"`

	*ApiError

	//
	// The error reason represented by an enum.
	//
	Reason *AuthorizationErrorReason `xml:"reason,omitempty"`
}

type BiddingErrors struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 BiddingErrors"`

	*ApiError

	//
	// The error reason represented by an enum.
	//
	Reason *BiddingErrorsReason `xml:"reason,omitempty"`
}

type BudgetError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 BudgetError"`

	*ApiError

	//
	// The error reason represented by an enum.
	//
	Reason *BudgetErrorReason `xml:"reason,omitempty"`
}

type CampaignError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 CampaignError"`

	*ApiError

	//
	// The error reason represented by an enum.
	//
	Reason *CampaignErrorReason `xml:"reason,omitempty"`
}

type ClientTermsError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 ClientTermsError"`

	*ApiError

	Reason *ClientTermsErrorReason `xml:"reason,omitempty"`
}

type DatabaseError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 DatabaseError"`

	*ApiError

	//
	// The error reason represented by an enum.
	//
	Reason *DatabaseErrorReason `xml:"reason,omitempty"`
}

type DateError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 DateError"`

	*ApiError

	//
	// The error reason represented by an enum.
	//
	Reason *DateErrorReason `xml:"reason,omitempty"`
}

type DateRange struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 DateRange"`

	//
	// the lower bound of this date range, inclusive.
	//
	Min string `xml:"min,omitempty"`

	//
	// the upper bound of this date range, inclusive.
	//
	Max string `xml:"max,omitempty"`
}

type DateRangeError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 DateRangeError"`

	*ApiError

	//
	// The error reason represented by an enum.
	//
	Reason *DateRangeErrorReason `xml:"reason,omitempty"`
}

type DistinctError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 DistinctError"`

	*ApiError

	//
	// The error reason represented by an enum.
	//
	Reason *DistinctErrorReason `xml:"reason,omitempty"`
}

type EntityCountLimitExceeded struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 EntityCountLimitExceeded"`

	*ApiError

	//
	// Specifies which level's limit was exceeded.
	//
	Reason *EntityCountLimitExceededReason `xml:"reason,omitempty"`

	//
	// Id of the entity whose limit was exceeded.
	//
	EnclosingId string `xml:"enclosingId,omitempty"`

	//
	// The limit which was exceeded.
	//
	Limit int32 `xml:"limit,omitempty"`

	//
	// The account limit type which was exceeded.
	//
	AccountLimitType string `xml:"accountLimitType,omitempty"`

	//
	// The count of existing entities.
	//
	ExistingCount int32 `xml:"existingCount,omitempty"`
}

type EntityNotFound struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 EntityNotFound"`

	*ApiError

	//
	// Reason for this error.
	//
	Reason *EntityNotFoundReason `xml:"reason,omitempty"`
}

type FieldPathElement struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 FieldPathElement"`

	//
	// The name of a field in lower camelcase. (e.g. "biddingStrategy")
	//
	Field string `xml:"field,omitempty"`

	//
	// For list fields, this is a 0-indexed position in the list. Null for non-list fields.
	//
	Index int32 `xml:"index,omitempty"`
}

type IdError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 IdError"`

	*ApiError

	//
	// The error reason represented by an enum.
	//
	Reason *IdErrorReason `xml:"reason,omitempty"`
}

type InternalApiError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 InternalApiError"`

	*ApiError

	//
	// The error reason represented by an enum.
	//
	Reason *InternalApiErrorReason `xml:"reason,omitempty"`
}

type ListReturnValue struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 ListReturnValue"`

	//
	// Indicates that this instance is a subtype of ListReturnValue.
	// Although this field is returned in the response, it is ignored on input
	// and cannot be selected. Specify xsi:type instead.
	//
	ListReturnValueType string `xml:"ListReturnValue.Type,omitempty"`
}

type NotEmptyError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 NotEmptyError"`

	*ApiError

	//
	// The error reason represented by an enum.
	//
	Reason *NotEmptyErrorReason `xml:"reason,omitempty"`
}

type NullError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 NullError"`

	*ApiError

	//
	// The error reason represented by an enum.
	//
	Reason *NullErrorReason `xml:"reason,omitempty"`
}

type Operation struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 Operation"`

	//
	// Operator.
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	//
	Operator *Operator `xml:"operator,omitempty"`

	//
	// Indicates that this instance is a subtype of Operation.
	// Although this field is returned in the response, it is ignored on input
	// and cannot be selected. Specify xsi:type instead.
	//
	OperationType string `xml:"Operation.Type,omitempty"`
}

type OperationAccessDenied struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 OperationAccessDenied"`

	*ApiError

	Reason *OperationAccessDeniedReason `xml:"reason,omitempty"`
}

type OperatorError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 OperatorError"`

	*ApiError

	//
	// The error reason represented by an enum.
	//
	Reason *OperatorErrorReason `xml:"reason,omitempty"`
}

type OrderBy struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 OrderBy"`

	//
	// The field to sort the results on.
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	//
	Field string `xml:"field,omitempty"`

	//
	// The order to sort the results on. The default sort order is {@link SortOrder#ASCENDING}.
	//
	SortOrder *SortOrder `xml:"sortOrder,omitempty"`
}

type Page struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 Page"`

	//
	// Total number of entries in the result that this page is a part of.
	//
	TotalNumEntries int32 `xml:"totalNumEntries,omitempty"`

	//
	// Indicates that this instance is a subtype of Page.
	// Although this field is returned in the response, it is ignored on input
	// and cannot be selected. Specify xsi:type instead.
	//
	PageType string `xml:"Page.Type,omitempty"`
}

type Paging struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 Paging"`

	//
	// Index of the first result to return in this page.
	// <span class="constraint InRange">This field must be greater than or equal to 0.</span>
	//
	StartIndex int32 `xml:"startIndex,omitempty"`

	//
	// Maximum number of results to return in this page. Set this to a reasonable value to limit
	// the number of results returned per page.
	// <span class="constraint InRange">This field must be greater than or equal to 0.</span>
	//
	NumberResults int32 `xml:"numberResults,omitempty"`
}

type Predicate struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 Predicate"`

	//
	// The field by which to filter the returned data. Possible values are marked Filterable on
	// the entity's reference page. For example, for predicates for the
	// CampaignService {@link Selector selector}, refer to the filterable fields from the
	// {@link Campaign} reference page.
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	//
	Field string `xml:"field,omitempty"`

	//
	// The operator to use for filtering the data returned.
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	//
	Operator *PredicateOperator `xml:"operator,omitempty"`

	//
	// The values by which to filter the field. The {@link Operator#CONTAINS_ALL},
	// {@link Operator#CONTAINS_ANY}, {@link Operator#CONTAINS_NONE}, {@link Operator#IN}
	// and {@link Operator#NOT_IN} take multiple values. All others take a single value.
	// <span class="constraint ContentsNotNull">This field must not contain {@code null} elements.</span>
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	//
	Values []string `xml:"values,omitempty"`
}

type QuotaCheckError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 QuotaCheckError"`

	*ApiError

	Reason *QuotaCheckErrorReason `xml:"reason,omitempty"`
}

type RangeError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 RangeError"`

	*ApiError

	//
	// The error reason represented by an enum.
	//
	Reason *RangeErrorReason `xml:"reason,omitempty"`
}

type RateExceededError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 RateExceededError"`

	*ApiError

	//
	// The error reason represented by an enum.
	//
	Reason *RateExceededErrorReason `xml:"reason,omitempty"`

	//
	// Cause of the rate exceeded error.
	//
	RateName string `xml:"rateName,omitempty"`

	//
	// The scope of the rate (ACCOUNT/DEVELOPER).
	//
	RateScope string `xml:"rateScope,omitempty"`

	//
	// The amount of time (in seconds) the client should wait before retrying the request.
	//
	RetryAfterSeconds int32 `xml:"retryAfterSeconds,omitempty"`
}

type ReadOnlyError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 ReadOnlyError"`

	*ApiError

	//
	// The error reason represented by an enum.
	//
	Reason *ReadOnlyErrorReason `xml:"reason,omitempty"`
}

type RejectedError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 RejectedError"`

	*ApiError

	//
	// The error reason represented by an enum.
	//
	Reason *RejectedErrorReason `xml:"reason,omitempty"`
}

type RequestError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 RequestError"`

	*ApiError

	Reason *RequestErrorReason `xml:"reason,omitempty"`
}

type RequiredError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 RequiredError"`

	*ApiError

	//
	// The error reason represented by an enum.
	//
	Reason *RequiredErrorReason `xml:"reason,omitempty"`
}

type Selector struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 Selector"`

	//
	// List of fields to select.
	// <a href="/adwords/api/docs/appendix/selectorfields">Possible values</a>
	// are marked {@code Selectable} on the entity's reference page.
	// For example, for the {@code CampaignService} selector, refer to the
	// selectable fields from the {@link Campaign} reference page.
	// <span class="constraint ContentsDistinct">This field must contain distinct elements.</span>
	// <span class="constraint ContentsNotNull">This field must not contain {@code null} elements.</span>
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	//
	Fields []string `xml:"fields,omitempty"`

	//
	// Specifies how an entity (eg. adgroup, campaign, criterion, ad) should be filtered.
	// <span class="constraint ContentsNotNull">This field must not contain {@code null} elements.</span>
	//
	Predicates []*Predicate `xml:"predicates,omitempty"`

	//
	// Range of dates for which you want to include data. If this value is omitted,
	// results for all dates are returned.
	// <p class="note"><b>Note:</b> This field is only used by the report download
	// service. For all other services, it is ignored.</p>
	// <span class="constraint DateRangeWithinRange">This range must be contained within the range [19700101, 20380101].</span>
	//
	DateRange *DateRange `xml:"dateRange,omitempty"`

	//
	// The fields on which you want to sort, and the sort order. The order in the list is
	// significant: The first element in the list indicates the primary sort order, the next
	// specifies the secondary sort order and so on.
	//
	Ordering []*OrderBy `xml:"ordering,omitempty"`

	//
	// Pagination information.
	//
	Paging *Paging `xml:"paging,omitempty"`
}

type SelectorError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 SelectorError"`

	*ApiError

	//
	// The error reason represented by enum.
	//
	Reason *SelectorErrorReason `xml:"reason,omitempty"`
}

type SizeLimitError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 SizeLimitError"`

	*ApiError

	//
	// The error reason represented by an enum.
	//
	Reason *SizeLimitErrorReason `xml:"reason,omitempty"`
}

type SoapHeader struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 SoapHeader"`

	//
	// The header identifies the customer id of the client of the AdWords manager, if an AdWords
	// manager is acting on behalf of their client or the customer id of the advertiser managing their
	// own account.
	//
	ClientCustomerId string `xml:"clientCustomerId,omitempty"`

	//
	// Developer token to identify that the person making the call has enough
	// quota.
	//
	DeveloperToken string `xml:"developerToken,omitempty"`

	//
	// UserAgent is used to track distribution of API client programs and
	// application usage. The client is responsible for putting in a meaningful
	// value for tracking purposes. To be clear this is not the same as an HTTP
	// user agent.
	//
	UserAgent string `xml:"userAgent,omitempty"`

	//
	// Used to validate the request without executing it.
	//
	ValidateOnly bool `xml:"validateOnly,omitempty"`

	//
	// If true, API will try to commit as many error free operations as possible and
	// report the other operations' errors.
	//
	// <p>Ignored for non-mutate calls.
	//
	PartialFailure bool `xml:"partialFailure,omitempty"`
}

type SoapResponseHeader struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 SoapResponseHeader"`

	//
	// Unique id that identifies this request. If developers have any support issues, sending us
	// this id will enable us to find their request more easily.
	//
	RequestId string `xml:"requestId,omitempty"`

	//
	// The name of the service being invoked.
	//
	ServiceName string `xml:"serviceName,omitempty"`

	//
	// The name of the method being invoked.
	//
	MethodName string `xml:"methodName,omitempty"`

	//
	// Number of operations performed for this SOAP request.
	//
	Operations int64 `xml:"operations,omitempty"`

	//
	// Elapsed time in milliseconds between the AdWords API receiving the request and sending the
	// response.
	//
	ResponseTime int64 `xml:"responseTime,omitempty"`
}

type StringFormatError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 StringFormatError"`

	*ApiError

	Reason *StringFormatErrorReason `xml:"reason,omitempty"`
}

type StringLengthError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 StringLengthError"`

	*ApiError

	//
	// The error reason represented by an enum.
	//
	Reason *StringLengthErrorReason `xml:"reason,omitempty"`
}

type Trial struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 Trial"`

	//
	// The id of this trial.
	// <span class="constraint Selectable">This field can be selected using the value "Id".</span><span class="constraint Filterable">This field can be filtered on.</span>
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API for the following {@link Operator}s: ADD.</span>
	// <span class="constraint Required">This field is required and should not be {@code null} when it is contained within {@link Operator}s : SET.</span>
	//
	Id int64 `xml:"id,omitempty"`

	//
	// Id of the base campaign, which will be the control arm of this trial.
	// <span class="constraint Selectable">This field can be selected using the value "BaseCampaignId".</span><span class="constraint Filterable">This field can be filtered on.</span>
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API for the following {@link Operator}s: SET.</span>
	// <span class="constraint Required">This field is required and should not be {@code null} when it is contained within {@link Operator}s : ADD.</span>
	//
	BaseCampaignId int64 `xml:"baseCampaignId,omitempty"`

	//
	// Valid id of the draft this trial is based on.
	// <span class="constraint Selectable">This field can be selected using the value "DraftId".</span><span class="constraint Filterable">This field can be filtered on.</span>
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API for the following {@link Operator}s: SET.</span>
	// <span class="constraint Required">This field is required and should not be {@code null} when it is contained within {@link Operator}s : ADD.</span>
	//
	DraftId int64 `xml:"draftId,omitempty"`

	//
	// Id of the new budget to assign to the trial campaign when graduating a trial.
	//
	// <p>Required for {@link Operator#SET SET} operations, when changing the {@link #status} to
	// {@code GRADUATED}, and read-only otherwise.
	//
	// <p>When graduating a trial, the same constraints apply to this field as for a budget id passed
	// to {@code CampaignService} when creating a new campaign.
	//
	// <p>{@code GET} operations always return a {@code null} budget id.
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API for the following {@link Operator}s: ADD.</span>
	//
	BudgetId int64 `xml:"budgetId,omitempty"`

	//
	// The name of this trial. Must not conflict with the name of another trial or campaign.
	// <span class="constraint Selectable">This field can be selected using the value "Name".</span><span class="constraint Filterable">This field can be filtered on.</span>
	// <span class="constraint Required">This field is required and should not be {@code null} when it is contained within {@link Operator}s : ADD.</span>
	// <span class="constraint StringLength">The length of this string should be between 1 and 1024, inclusive, in UTF-8 bytes, (trimmed).</span>
	//
	Name string `xml:"name,omitempty"`

	//
	// Date the trial begins. On add, defaults to the the base campaign's start date or the
	// current day in the parent account's local timezone (whichever is later).
	// <span class="constraint Selectable">This field can be selected using the value "StartDate".</span><span class="constraint Filterable">This field can be filtered on.</span>
	//
	StartDate string `xml:"startDate,omitempty"`

	//
	// Date the campaign ends. On add, defaults to the base campaign's end date.
	// <span class="constraint Selectable">This field can be selected using the value "EndDate".</span><span class="constraint Filterable">This field can be filtered on.</span>
	//
	EndDate string `xml:"endDate,omitempty"`

	//
	// Traffic share to be directed to the trial arm of this trial, i.e. the arm containing the
	// trial changes (in percent). The remainder of the traffic (100 - {@code trafficSplitPercent})
	// will be directed to the base campaign.
	// <span class="constraint Selectable">This field can be selected using the value "TrafficSplitPercent".</span><span class="constraint Filterable">This field can be filtered on.</span>
	// <span class="constraint InRange">This field must be between 1 and 99, inclusive.</span>
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API for the following {@link Operator}s: SET.</span>
	// <span class="constraint Required">This field is required and should not be {@code null} when it is contained within {@link Operator}s : ADD.</span>
	//
	TrafficSplitPercent int32 `xml:"trafficSplitPercent,omitempty"`

	//
	// Status of this trial. Note that a running trial will always be ACTIVE, but not all ACTIVE
	// trials are currently running: they may have ended or been scheduled for the future.
	// <span class="constraint Selectable">This field can be selected using the value "Status".</span><span class="constraint Filterable">This field can be filtered on.</span>
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API for the following {@link Operator}s: ADD.</span>
	//
	Status *TrialStatus `xml:"status,omitempty"`

	//
	// Id of the trial campaign. This will be null if the Trial has
	// status CREATING.
	// <span class="constraint Selectable">This field can be selected using the value "TrialCampaignId".</span><span class="constraint Filterable">This field can be filtered on.</span>
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	TrialCampaignId int64 `xml:"trialCampaignId,omitempty"`
}

type TrialError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 TrialError"`

	*ApiError

	Reason *TrialErrorReason `xml:"reason,omitempty"`
}

type TrialOperation struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 TrialOperation"`

	*Operation

	//
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	//
	Operand *Trial `xml:"operand,omitempty"`
}

type TrialPage struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 TrialPage"`

	*Page

	Entries []*Trial `xml:"entries,omitempty"`
}

type TrialReturnValue struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 TrialReturnValue"`

	*ListReturnValue

	Value []*Trial `xml:"value,omitempty"`

	PartialFailureErrors []*ApiError `xml:"partialFailureErrors,omitempty"`
}

type TrialServiceInterface struct {
	client *SOAPClient
}

func NewTrialServiceInterface(url string, tls bool, auth *BasicAuth) *TrialServiceInterface {
	if url == "" {
		url = ""
	}
	client := NewSOAPClient(url, tls, auth)

	return &TrialServiceInterface{
		client: client,
	}
}

func NewTrialServiceInterfaceWithTLSConfig(url string, tlsCfg *tls.Config, auth *BasicAuth) *TrialServiceInterface {
	if url == "" {
		url = ""
	}
	client := NewSOAPClientWithTLSConfig(url, tlsCfg, auth)

	return &TrialServiceInterface{
		client: client,
	}
}

func (service *TrialServiceInterface) AddHeader(header interface{}) {
	service.client.AddHeader(header)
}

// Backwards-compatible function: use AddHeader instead
func (service *TrialServiceInterface) SetHeader(header interface{}) {
	service.client.AddHeader(header)
}

// Error can be either of the following types:
//
//   - ApiException
/*
   Loads a TrialPage containing a list of {@link Trial} objects matching the selector.

   @param selector defines which subset of all available trials to return, the sort order, and
   which fields to include

   @return Returns a page of matching trial objects.
   @throws com.google.ads.api.services.common.error.ApiException if errors occurred while
   retrieving the results.
*/
func (service *TrialServiceInterface) Get(request *Get) (*GetResponse, error) {
	response := new(GetResponse)
	err := service.client.Call("", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// Error can be either of the following types:
//
//   - ApiException
/*
   Creates new trials, updates properties and controls the life cycle of existing trials.
   See {@link TrialService} for details on the trial life cycle.

   @return Returns the list of updated Trials, in the same order as the {@code operations} list.
   @throws com.google.ads.api.services.common.error.ApiException if errors occurred while
   processing the request.
*/
func (service *TrialServiceInterface) Mutate(request *Mutate) (*MutateResponse, error) {
	response := new(MutateResponse)
	err := service.client.Call("", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// Error can be either of the following types:
//
//   - ApiException
/*
   Loads a TrialPage containing a list of {@link Trial} objects matching the query.

   @param query defines which subset of all available trials to return, the sort order, and
   which fields to include

   @return Returns a page of matching trial objects.
   @throws com.google.ads.api.services.common.error.ApiException if errors occurred while
   retrieving the results.
*/
func (service *TrialServiceInterface) Query(request *Query) (*QueryResponse, error) {
	response := new(QueryResponse)
	err := service.client.Call("", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

var timeout = time.Duration(30 * time.Second)

func dialTimeout(network, addr string) (net.Conn, error) {
	return net.DialTimeout(network, addr, timeout)
}

type SOAPEnvelope struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Envelope"`
	Header  *SOAPHeader
	Body    SOAPBody
}

type SOAPHeader struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Header"`

	Items []interface{} `xml:",omitempty"`
}

type SOAPBody struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Body"`

	Fault   *SOAPFault  `xml:",omitempty"`
	Content interface{} `xml:",omitempty"`
}

type SOAPFault struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Fault"`

	Code   string `xml:"faultcode,omitempty"`
	String string `xml:"faultstring,omitempty"`
	Actor  string `xml:"faultactor,omitempty"`
	Detail string `xml:"detail,omitempty"`
}

const (
	// Predefined WSS namespaces to be used in
	WssNsWSSE string = "http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-secext-1.0.xsd"
	WssNsWSU  string = "http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-utility-1.0.xsd"
	WssNsType string = "http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-username-token-profile-1.0#PasswordText"
)

type WSSSecurityHeader struct {
	XMLName   xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ wsse:Security"`
	XmlNSWsse string   `xml:"xmlns:wsse,attr"`

	MustUnderstand string `xml:"mustUnderstand,attr,omitempty"`

	Token *WSSUsernameToken `xml:",omitempty"`
}

type WSSUsernameToken struct {
	XMLName   xml.Name `xml:"wsse:UsernameToken"`
	XmlNSWsu  string   `xml:"xmlns:wsu,attr"`
	XmlNSWsse string   `xml:"xmlns:wsse,attr"`

	Id string `xml:"wsu:Id,attr,omitempty"`

	Username *WSSUsername `xml:",omitempty"`
	Password *WSSPassword `xml:",omitempty"`
}

type WSSUsername struct {
	XMLName   xml.Name `xml:"wsse:Username"`
	XmlNSWsse string   `xml:"xmlns:wsse,attr"`

	Data string `xml:",chardata"`
}

type WSSPassword struct {
	XMLName   xml.Name `xml:"wsse:Password"`
	XmlNSWsse string   `xml:"xmlns:wsse,attr"`
	XmlNSType string   `xml:"Type,attr"`

	Data string `xml:",chardata"`
}

type BasicAuth struct {
	Login    string
	Password string
}

type SOAPClient struct {
	url     string
	tlsCfg  *tls.Config
	auth    *BasicAuth
	headers []interface{}
}

// **********
// Accepted solution from http://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-golang
// Author: Icza - http://stackoverflow.com/users/1705598/icza

const (
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func randStringBytesMaskImprSrc(n int) string {
	src := rand.NewSource(time.Now().UnixNano())
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return string(b)
}

// **********

func NewWSSSecurityHeader(user, pass, mustUnderstand string) *WSSSecurityHeader {
	hdr := &WSSSecurityHeader{XmlNSWsse: WssNsWSSE, MustUnderstand: mustUnderstand}
	hdr.Token = &WSSUsernameToken{XmlNSWsu: WssNsWSU, XmlNSWsse: WssNsWSSE, Id: "UsernameToken-" + randStringBytesMaskImprSrc(9)}
	hdr.Token.Username = &WSSUsername{XmlNSWsse: WssNsWSSE, Data: user}
	hdr.Token.Password = &WSSPassword{XmlNSWsse: WssNsWSSE, XmlNSType: WssNsType, Data: pass}
	return hdr
}

func (b *SOAPBody) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	if b.Content == nil {
		return xml.UnmarshalError("Content must be a pointer to a struct")
	}

	var (
		token    xml.Token
		err      error
		consumed bool
	)

Loop:
	for {
		if token, err = d.Token(); err != nil {
			return err
		}

		if token == nil {
			break
		}

		switch se := token.(type) {
		case xml.StartElement:
			if consumed {
				return xml.UnmarshalError("Found multiple elements inside SOAP body; not wrapped-document/literal WS-I compliant")
			} else if se.Name.Space == "http://schemas.xmlsoap.org/soap/envelope/" && se.Name.Local == "Fault" {
				b.Fault = &SOAPFault{}
				b.Content = nil

				err = d.DecodeElement(b.Fault, &se)
				if err != nil {
					return err
				}

				consumed = true
			} else {
				if err = d.DecodeElement(b.Content, &se); err != nil {
					return err
				}

				consumed = true
			}
		case xml.EndElement:
			break Loop
		}
	}

	return nil
}

func (f *SOAPFault) Error() string {
	return f.String
}

func NewSOAPClient(url string, insecureSkipVerify bool, auth *BasicAuth) *SOAPClient {
	tlsCfg := &tls.Config{
		InsecureSkipVerify: insecureSkipVerify,
	}
	return NewSOAPClientWithTLSConfig(url, tlsCfg, auth)
}

func NewSOAPClientWithTLSConfig(url string, tlsCfg *tls.Config, auth *BasicAuth) *SOAPClient {
	return &SOAPClient{
		url:    url,
		tlsCfg: tlsCfg,
		auth:   auth,
	}
}

func (s *SOAPClient) AddHeader(header interface{}) {
	s.headers = append(s.headers, header)
}

func (s *SOAPClient) Call(soapAction string, request, response interface{}) error {
	envelope := SOAPEnvelope{}

	if s.headers != nil && len(s.headers) > 0 {
		soapHeader := &SOAPHeader{Items: make([]interface{}, len(s.headers))}
		copy(soapHeader.Items, s.headers)
		envelope.Header = soapHeader
	}

	envelope.Body.Content = request
	buffer := new(bytes.Buffer)

	encoder := xml.NewEncoder(buffer)
	//encoder.Indent("  ", "    ")

	if err := encoder.Encode(envelope); err != nil {
		return err
	}

	if err := encoder.Flush(); err != nil {
		return err
	}

	log.Println(buffer.String())

	req, err := http.NewRequest("POST", s.url, buffer)
	if err != nil {
		return err
	}
	if s.auth != nil {
		req.SetBasicAuth(s.auth.Login, s.auth.Password)
	}

	req.Header.Add("Content-Type", "text/xml; charset=\"utf-8\"")
	req.Header.Add("SOAPAction", soapAction)

	req.Header.Set("User-Agent", "gowsdl/0.1")
	req.Close = true

	tr := &http.Transport{
		TLSClientConfig: s.tlsCfg,
		Dial:            dialTimeout,
	}

	client := &http.Client{Transport: tr}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	rawbody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if len(rawbody) == 0 {
		log.Println("empty response")
		return nil
	}

	log.Println(string(rawbody))
	respEnvelope := new(SOAPEnvelope)
	respEnvelope.Body = SOAPBody{Content: response}
	err = xml.Unmarshal(rawbody, respEnvelope)
	if err != nil {
		return err
	}

	fault := respEnvelope.Body.Fault
	if fault != nil {
		return fault
	}

	return nil
}
