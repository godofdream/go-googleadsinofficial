package AdGroupService

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
// The reasons for the adgroup service error.
//
type AdGroupServiceErrorReason string

const (

	//
	// AdGroup with the same name already exists for the campaign.
	//
	AdGroupServiceErrorReasonDUPLICATE_ADGROUP_NAME AdGroupServiceErrorReason = "DUPLICATE_ADGROUP_NAME"

	//
	// AdGroup name is not valid.
	//
	AdGroupServiceErrorReasonINVALID_ADGROUP_NAME AdGroupServiceErrorReason = "INVALID_ADGROUP_NAME"

	//
	// Cannot remove an adgroup, adgroup status can be marked removed
	// using set operator.
	//
	AdGroupServiceErrorReasonUSE_SET_OPERATOR_AND_MARK_STATUS_TO_REMOVED AdGroupServiceErrorReason = "USE_SET_OPERATOR_AND_MARK_STATUS_TO_REMOVED"

	//
	// Advertiser is not allowed to target sites or set site bids that are
	// not on the Google Search Network.
	//
	AdGroupServiceErrorReasonADVERTISER_NOT_ON_CONTENT_NETWORK AdGroupServiceErrorReason = "ADVERTISER_NOT_ON_CONTENT_NETWORK"

	//
	// Bid amount is too big.
	//
	AdGroupServiceErrorReasonBID_TOO_BIG AdGroupServiceErrorReason = "BID_TOO_BIG"

	//
	// AdGroup bid does not match the campaign's bidding strategy.
	//
	AdGroupServiceErrorReasonBID_TYPE_AND_BIDDING_STRATEGY_MISMATCH AdGroupServiceErrorReason = "BID_TYPE_AND_BIDDING_STRATEGY_MISMATCH"

	//
	// AdGroup name is required for Add.
	//
	AdGroupServiceErrorReasonMISSING_ADGROUP_NAME AdGroupServiceErrorReason = "MISSING_ADGROUP_NAME"

	//
	// No link found between the ad group and the label.
	//
	AdGroupServiceErrorReasonADGROUP_LABEL_DOES_NOT_EXIST AdGroupServiceErrorReason = "ADGROUP_LABEL_DOES_NOT_EXIST"

	//
	// The label has already been attached to the ad group.
	//
	AdGroupServiceErrorReasonADGROUP_LABEL_ALREADY_EXISTS AdGroupServiceErrorReason = "ADGROUP_LABEL_ALREADY_EXISTS"

	//
	// The CriterionTypeGroup is not supported for the content bid dimension.
	//
	AdGroupServiceErrorReasonINVALID_CONTENT_BID_CRITERION_TYPE_GROUP AdGroupServiceErrorReason = "INVALID_CONTENT_BID_CRITERION_TYPE_GROUP"

	//
	// The ad group type is not compatible with the campaign channel type.
	//
	AdGroupServiceErrorReasonAD_GROUP_TYPE_NOT_VALID_FOR_ADVERTISING_CHANNEL_TYPE AdGroupServiceErrorReason = "AD_GROUP_TYPE_NOT_VALID_FOR_ADVERTISING_CHANNEL_TYPE"

	//
	// The ad group type is not supported in the country of sale of the campaign.
	//
	AdGroupServiceErrorReasonADGROUP_TYPE_NOT_SUPPORTED_FOR_CAMPAIGN_SALES_COUNTRY AdGroupServiceErrorReason = "ADGROUP_TYPE_NOT_SUPPORTED_FOR_CAMPAIGN_SALES_COUNTRY"

	//
	// Ad groups of AdGroupType.SEARCH_DYNAMIC_ADS can only be added to campaigns that have
	// DynamicSearchAdsSetting attached.
	//
	AdGroupServiceErrorReasonCANNOT_ADD_ADGROUP_OF_TYPE_DSA_TO_CAMPAIGN_WITHOUT_DSA_SETTING AdGroupServiceErrorReason = "CANNOT_ADD_ADGROUP_OF_TYPE_DSA_TO_CAMPAIGN_WITHOUT_DSA_SETTING"
)

//
// Status of this ad group.
//
type AdGroupStatus string

const (

	//
	// <span class="constraint Rejected">Used for return value only. An enumeration could not be processed, typically due to incompatibility with your WSDL version.</span>
	//
	AdGroupStatusUNKNOWN AdGroupStatus = "UNKNOWN"

	//
	// Active.
	//
	AdGroupStatusENABLED AdGroupStatus = "ENABLED"

	//
	// Paused.
	//
	AdGroupStatusPAUSED AdGroupStatus = "PAUSED"

	//
	// Removed.
	//
	AdGroupStatusREMOVED AdGroupStatus = "REMOVED"
)

//
// Defines types of an ad group, specific to a particular campaign channel type.
// This type drives validations that restrict which entities can be added to the ad
// group.
//
type AdGroupType string

const (

	//
	// <span class="constraint Rejected">Used for return value only. An enumeration could not be processed, typically due to incompatibility with your WSDL version.</span>
	//
	AdGroupTypeUNKNOWN AdGroupType = "UNKNOWN"

	//
	// Default AdGroup type for Search Campaigns
	//
	AdGroupTypeSEARCH_STANDARD AdGroupType = "SEARCH_STANDARD"

	//
	// AdGroup type for Dynamic Search Ads campaigns.
	//
	AdGroupTypeSEARCH_DYNAMIC_ADS AdGroupType = "SEARCH_DYNAMIC_ADS"

	//
	// Default AdGroup type for Display Campaigns
	//
	AdGroupTypeDISPLAY_STANDARD AdGroupType = "DISPLAY_STANDARD"

	//
	// Default AdGroup type for Shopping Campaigns serving standard products ads.
	//
	AdGroupTypeSHOPPING_PRODUCT_ADS AdGroupType = "SHOPPING_PRODUCT_ADS"

	//
	// AdGroups limited to serving Showcase/Merchant ads in shopping results.
	//
	AdGroupTypeSHOPPING_SHOWCASE_ADS AdGroupType = "SHOPPING_SHOWCASE_ADS"

	//
	// Ad group type for Universal Shopping Campaigns.
	//
	AdGroupTypeSHOPPING_UNIVERSAL_ADS AdGroupType = "SHOPPING_UNIVERSAL_ADS"
)

//
// Indicates the ad rotation mode selected for the  creatives in the ad group.
//
type AdRotationMode string

const (

	//
	// Invalid type. Should not be used except for detecting values that are incorrect,
	// or values that are not yet known to the user.
	// <span class="constraint Rejected">Used for return value only. An enumeration could not be processed, typically due to incompatibility with your WSDL version.</span>
	//
	AdRotationModeUNKNOWN AdRotationMode = "UNKNOWN"

	//
	// Adwords will optimize your adgroups based on clicks or conversions based on your
	// bidding strategy
	//
	AdRotationModeOPTIMIZE AdRotationMode = "OPTIMIZE"

	//
	// Ads in the adgroups should rotate evenly forever
	//
	AdRotationModeROTATE_FOREVER AdRotationMode = "ROTATE_FOREVER"
)

//
// The reasons for the AdX error.
//
type AdxErrorReason string

const (

	//
	// Attempt to use non-AdX feature by AdX customer.
	//
	AdxErrorReasonUNSUPPORTED_FEATURE AdxErrorReason = "UNSUPPORTED_FEATURE"
)

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

//
// Indicate where a criterion's bid came from: criterion or the adgroup it
// belongs to.
//
type BidSource string

const (

	//
	// Effective Bid is Adgroup level bid
	//
	BidSourceADGROUP BidSource = "ADGROUP"

	//
	// Effective Bid is Keyword level bid
	//
	BidSourceCRITERION BidSource = "CRITERION"

	//
	// Effective Bid is inherited from Adgroup Bidding Strategy
	//
	BidSourceADGROUP_BIDDING_STRATEGY BidSource = "ADGROUP_BIDDING_STRATEGY"

	//
	// Effective Bid is inherited from Campaign Bidding Strategy
	//
	BidSourceCAMPAIGN_BIDDING_STRATEGY BidSource = "CAMPAIGN_BIDDING_STRATEGY"
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
// Indicates where bidding strategy came from: campaign, adgroup or criterion.
//
type BiddingStrategySource string

const (

	//
	// Bidding strategy is defined on campaign level.
	//
	BiddingStrategySourceCAMPAIGN BiddingStrategySource = "CAMPAIGN"

	//
	// Bidding strategy is defined on adgroup level.
	//
	BiddingStrategySourceADGROUP BiddingStrategySource = "ADGROUP"

	//
	// Bidding strategy is defined on criterion level.
	//
	BiddingStrategySourceCRITERION BiddingStrategySource = "CRITERION"
)

//
// The bidding strategy type. See {@linkplain BiddingStrategyConfiguration}
// for additional information.
//
type BiddingStrategyType string

const (

	//
	// Manual click based bidding where user pays per click. See
	// {@linkplain ManualCpcBiddingScheme} for more details.
	//
	BiddingStrategyTypeMANUAL_CPC BiddingStrategyType = "MANUAL_CPC"

	//
	// Manual impression based bidding where user pays per thousand
	// impressions. See {@linkplain ManualCpmBiddingScheme} for more
	// details.
	//
	BiddingStrategyTypeMANUAL_CPM BiddingStrategyType = "MANUAL_CPM"

	//
	// Page-One Promoted is an automated bid strategy that sets max CPC bids
	// to target impressions on page one or page one promoted slots on
	// google.com. See {@linkplain PageOnePromotedBiddingScheme} for
	// more details.
	//
	BiddingStrategyTypePAGE_ONE_PROMOTED BiddingStrategyType = "PAGE_ONE_PROMOTED"

	//
	// Target Spend (Maximize Clicks) is an automated bid strategy that sets
	// your bids to help get as many clicks as possible within your budget.
	// See {@linkplain TargetSpendBiddingScheme} for more details.
	//
	BiddingStrategyTypeTARGET_SPEND BiddingStrategyType = "TARGET_SPEND"

	//
	// Enhanced CPC is a bidding strategy that raises your bids for clicks
	// that seem more likely to lead to a conversion and lowers them for clicks
	// where they seem less likely. See {@linkplain EnhancedCpcBiddingScheme}
	// for more details.
	//
	BiddingStrategyTypeENHANCED_CPC BiddingStrategyType = "ENHANCED_CPC"

	//
	// Target CPA is an automated bid strategy that sets bids to help get
	// as many conversions as possible at the target cost per acquisition
	// (CPA) you set. See {@linkplain TargetCpaBiddingScheme}
	// for more details.
	//
	BiddingStrategyTypeTARGET_CPA BiddingStrategyType = "TARGET_CPA"

	//
	// Target ROAS is an automated bidding strategy that helps you maximize
	// revenue while averaging a specific target return on average spend (ROAS).
	// See {@linkplain TargetRoasBiddingScheme} for more details.
	//
	BiddingStrategyTypeTARGET_ROAS BiddingStrategyType = "TARGET_ROAS"

	//
	// Maximize conversions is an automated bidding strategy that automatically sets bids to help
	// get the most conversions for your campaign while spending your budget.
	// See {@linkplain MaximizeConversionsBiddingScheme} for more details.
	//
	BiddingStrategyTypeMAXIMIZE_CONVERSIONS BiddingStrategyType = "MAXIMIZE_CONVERSIONS"

	//
	// Maximize conversion value is an automated bidding strategy that automatically sets bids to
	// maximize revenue while spending your budget.
	// See {@linkplain MaximizeConversionValueBiddingScheme} for more details.
	//
	BiddingStrategyTypeMAXIMIZE_CONVERSION_VALUE BiddingStrategyType = "MAXIMIZE_CONVERSION_VALUE"

	//
	// Target Outrank Share is an automated bidding strategy that sets bids
	// based on the target fraction of auctions where the advertiser should
	// outrank a specific competitor. See {@linkplain TargetOutrankShareBiddingScheme}
	// for more details.
	//
	BiddingStrategyTypeTARGET_OUTRANK_SHARE BiddingStrategyType = "TARGET_OUTRANK_SHARE"

	//
	// Special bidding strategy type used to reset the bidding strategy at AdGroup and
	// AdGroupCriterion.
	//
	BiddingStrategyTypeNONE BiddingStrategyType = "NONE"

	//
	// <span class="constraint Rejected">Used for return value only. An enumeration could not be processed, typically due to incompatibility with your WSDL version.</span>
	//
	BiddingStrategyTypeUNKNOWN BiddingStrategyType = "UNKNOWN"
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
// The list of groupings of criteria types.
//
type CriterionTypeGroup string

const (

	//
	// Criteria for targeting keywords. e.g. 'mars cruise'
	// KEYWORD may be used as a content bid dimension. Keywords are always a targeting dimension,
	// so may not be set as a target "ALL" dimension with {@link TargetRestriction}.
	//
	CriterionTypeGroupKEYWORD CriterionTypeGroup = "KEYWORD"

	//
	// Criteria for targeting lists of users.  Lists may represent users with particular
	// interests, or they may represent users who have interacted with an advertiser's site in
	// particular ways.
	//
	CriterionTypeGroupUSER_INTEREST_AND_LIST CriterionTypeGroup = "USER_INTEREST_AND_LIST"

	//
	// Criteria for targeting similar categories of placements, e.g. 'category::Animals>Pets'
	// Used only for content network targeting.
	//
	CriterionTypeGroupVERTICAL CriterionTypeGroup = "VERTICAL"

	//
	// Criteria for targeting gender.
	//
	CriterionTypeGroupGENDER CriterionTypeGroup = "GENDER"

	//
	// Criteria for targeting age ranges.
	//
	CriterionTypeGroupAGE_RANGE CriterionTypeGroup = "AGE_RANGE"

	//
	// Criteria for targeting placements. aka Website. e.g. 'www.flowers4sale.com'
	// This group also includes mobile applications and mobile app categories.
	//
	CriterionTypeGroupPLACEMENT CriterionTypeGroup = "PLACEMENT"

	//
	// Criteria for parental status targeting.
	//
	CriterionTypeGroupPARENT CriterionTypeGroup = "PARENT"

	//
	// Criteria for income range targeting.
	//
	CriterionTypeGroupINCOME_RANGE CriterionTypeGroup = "INCOME_RANGE"

	//
	// Special criteria type group used to reset the existing value of AdGroup's
	// contentBidCriterionTypeGroup.
	//
	CriterionTypeGroupNONE CriterionTypeGroup = "NONE"

	//
	// <span class="constraint Rejected">Used for return value only. An enumeration could not be processed, typically due to incompatibility with your WSDL version.</span>
	//
	CriterionTypeGroupUNKNOWN CriterionTypeGroup = "UNKNOWN"
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
// The reasons for the validation error.
//
type DistinctErrorReason string

const (
	DistinctErrorReasonDUPLICATE_ELEMENT DistinctErrorReason = "DUPLICATE_ELEMENT"

	DistinctErrorReasonDUPLICATE_TYPE DistinctErrorReason = "DUPLICATE_TYPE"
)

type EntityAccessDeniedReason string

const (

	//
	// User did not have read access.
	//
	EntityAccessDeniedReasonREAD_ACCESS_DENIED EntityAccessDeniedReason = "READ_ACCESS_DENIED"

	//
	// User did not have write access.
	//
	EntityAccessDeniedReasonWRITE_ACCESS_DENIED EntityAccessDeniedReason = "WRITE_ACCESS_DENIED"
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
// The reason for the error.
//
type ForwardCompatibilityErrorReason string

const (

	//
	// Invalid value specified for a key in the forward compatibility map.
	//
	ForwardCompatibilityErrorReasonINVALID_FORWARD_COMPATIBILITY_MAP_VALUE ForwardCompatibilityErrorReason = "INVALID_FORWARD_COMPATIBILITY_MAP_VALUE"

	ForwardCompatibilityErrorReasonUNKNOWN ForwardCompatibilityErrorReason = "UNKNOWN"
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

type LabelStatus string

const (

	//
	// The label is enabled.
	//
	LabelStatusENABLED LabelStatus = "ENABLED"

	//
	// The label has been removed.
	//
	LabelStatusREMOVED LabelStatus = "REMOVED"

	//
	// <span class="constraint Rejected">Used for return value only. An enumeration could not be processed, typically due to incompatibility with your WSDL version.</span>
	//
	LabelStatusUNKNOWN LabelStatus = "UNKNOWN"
)

//
// Reason for bidding error.
//
type MultiplierErrorReason string

const (

	//
	// Multiplier value is too high
	//
	MultiplierErrorReasonMULTIPLIER_TOO_HIGH MultiplierErrorReason = "MULTIPLIER_TOO_HIGH"

	//
	// Multiplier value is too low
	//
	MultiplierErrorReasonMULTIPLIER_TOO_LOW MultiplierErrorReason = "MULTIPLIER_TOO_LOW"

	//
	// Too many fractional digits
	//
	MultiplierErrorReasonTOO_MANY_FRACTIONAL_DIGITS MultiplierErrorReason = "TOO_MANY_FRACTIONAL_DIGITS"

	//
	// A multiplier cannot be set for this bidding strategy
	//
	MultiplierErrorReasonMULTIPLIER_NOT_ALLOWED_FOR_BIDDING_STRATEGY MultiplierErrorReason = "MULTIPLIER_NOT_ALLOWED_FOR_BIDDING_STRATEGY"

	//
	// A multiplier cannot be set when there is no base bid (e.g., content max cpc)
	//
	MultiplierErrorReasonMULTIPLIER_NOT_ALLOWED_WHEN_BASE_BID_IS_MISSING MultiplierErrorReason = "MULTIPLIER_NOT_ALLOWED_WHEN_BASE_BID_IS_MISSING"

	//
	// A bid multiplier must be specified
	//
	MultiplierErrorReasonNO_MULTIPLIER_SPECIFIED MultiplierErrorReason = "NO_MULTIPLIER_SPECIFIED"

	//
	// Multiplier causes bid to exceed daily budget
	//
	MultiplierErrorReasonMULTIPLIER_CAUSES_BID_TO_EXCEED_DAILY_BUDGET MultiplierErrorReason = "MULTIPLIER_CAUSES_BID_TO_EXCEED_DAILY_BUDGET"

	//
	// Multiplier causes bid to exceed monthly budget
	//
	MultiplierErrorReasonMULTIPLIER_CAUSES_BID_TO_EXCEED_MONTHLY_BUDGET MultiplierErrorReason = "MULTIPLIER_CAUSES_BID_TO_EXCEED_MONTHLY_BUDGET"

	//
	// Multiplier causes bid to exceed custom budget
	//
	MultiplierErrorReasonMULTIPLIER_CAUSES_BID_TO_EXCEED_CUSTOM_BUDGET MultiplierErrorReason = "MULTIPLIER_CAUSES_BID_TO_EXCEED_CUSTOM_BUDGET"

	//
	// Multiplier causes bid to exceed maximum allowed bid
	//
	MultiplierErrorReasonMULTIPLIER_CAUSES_BID_TO_EXCEED_MAX_ALLOWED_BID MultiplierErrorReason = "MULTIPLIER_CAUSES_BID_TO_EXCEED_MAX_ALLOWED_BID"

	//
	// Multiplier causes bid to become less than the minimum bid allowed
	//
	MultiplierErrorReasonBID_LESS_THAN_MIN_ALLOWED_BID_WITH_MULTIPLIER MultiplierErrorReason = "BID_LESS_THAN_MIN_ALLOWED_BID_WITH_MULTIPLIER"

	//
	// Multiplier type (cpc vs. cpm) needs to match campaign's bidding strategy
	//
	MultiplierErrorReasonMULTIPLIER_AND_BIDDING_STRATEGY_TYPE_MISMATCH MultiplierErrorReason = "MULTIPLIER_AND_BIDDING_STRATEGY_TYPE_MISMATCH"

	MultiplierErrorReasonMULTIPLIER_ERROR MultiplierErrorReason = "MULTIPLIER_ERROR"
)

type NewEntityCreationErrorReason string

const (

	//
	// Do not set the id field while creating new entities.
	//
	NewEntityCreationErrorReasonCANNOT_SET_ID_FOR_ADD NewEntityCreationErrorReason = "CANNOT_SET_ID_FOR_ADD"

	//
	// Creating more than one entity with the same temp ID is not allowed.
	//
	NewEntityCreationErrorReasonDUPLICATE_TEMP_IDS NewEntityCreationErrorReason = "DUPLICATE_TEMP_IDS"

	//
	// Parent object with specified temp id failed validation, so no deep
	// validation will be done for this child entity.
	//
	NewEntityCreationErrorReasonTEMP_ID_ENTITY_HAD_ERRORS NewEntityCreationErrorReason = "TEMP_ID_ENTITY_HAD_ERRORS"
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

type PageOnePromotedBiddingSchemeStrategyGoal string

const (

	//
	// First page on google.com.
	//
	PageOnePromotedBiddingSchemeStrategyGoalPAGE_ONE PageOnePromotedBiddingSchemeStrategyGoal = "PAGE_ONE"

	//
	// Top slots of the first page on google.com.
	//
	PageOnePromotedBiddingSchemeStrategyGoalPAGE_ONE_PROMOTED PageOnePromotedBiddingSchemeStrategyGoal = "PAGE_ONE_PROMOTED"
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
// The reason for the query error.
//
type QueryErrorReason string

const (

	//
	// Exception that happens when trying to parse a query that doesn't match the AWQL grammar.
	//
	QueryErrorReasonPARSING_FAILED QueryErrorReason = "PARSING_FAILED"

	//
	// The provided query is an empty string.
	//
	QueryErrorReasonMISSING_QUERY QueryErrorReason = "MISSING_QUERY"

	//
	// The query does not contain the required SELECT clause or it is not in the
	// correct location.
	//
	QueryErrorReasonMISSING_SELECT_CLAUSE QueryErrorReason = "MISSING_SELECT_CLAUSE"

	//
	// The query does not contain the required FROM clause or it is not in the correct location.
	//
	QueryErrorReasonMISSING_FROM_CLAUSE QueryErrorReason = "MISSING_FROM_CLAUSE"

	//
	// The SELECT clause could not be parsed.
	//
	QueryErrorReasonINVALID_SELECT_CLAUSE QueryErrorReason = "INVALID_SELECT_CLAUSE"

	//
	// The FROM clause could not be parsed.
	//
	QueryErrorReasonINVALID_FROM_CLAUSE QueryErrorReason = "INVALID_FROM_CLAUSE"

	//
	// The WHERE clause could not be parsed.
	//
	QueryErrorReasonINVALID_WHERE_CLAUSE QueryErrorReason = "INVALID_WHERE_CLAUSE"

	//
	// The ORDER BY clause could not be parsed.
	//
	QueryErrorReasonINVALID_ORDER_BY_CLAUSE QueryErrorReason = "INVALID_ORDER_BY_CLAUSE"

	//
	// The LIMIT clause could not be parsed.
	//
	QueryErrorReasonINVALID_LIMIT_CLAUSE QueryErrorReason = "INVALID_LIMIT_CLAUSE"

	//
	// The startIndex in the LIMIT clause does not contain a valid integer.
	//
	QueryErrorReasonINVALID_START_INDEX_IN_LIMIT_CLAUSE QueryErrorReason = "INVALID_START_INDEX_IN_LIMIT_CLAUSE"

	//
	// The pageSize in the LIMIT clause does not contain a valid integer.
	//
	QueryErrorReasonINVALID_PAGE_SIZE_IN_LIMIT_CLAUSE QueryErrorReason = "INVALID_PAGE_SIZE_IN_LIMIT_CLAUSE"

	//
	// The DURING clause could not be parsed.
	//
	QueryErrorReasonINVALID_DURING_CLAUSE QueryErrorReason = "INVALID_DURING_CLAUSE"

	//
	// The minimum date in the DURING clause is not a valid date in YYYYMMDD format.
	//
	QueryErrorReasonINVALID_MIN_DATE_IN_DURING_CLAUSE QueryErrorReason = "INVALID_MIN_DATE_IN_DURING_CLAUSE"

	//
	// The maximum date in the DURING clause is not a valid date in YYYYMMDD format.
	//
	QueryErrorReasonINVALID_MAX_DATE_IN_DURING_CLAUSE QueryErrorReason = "INVALID_MAX_DATE_IN_DURING_CLAUSE"

	//
	// The minimum date in the DURING is after the maximum date.
	//
	QueryErrorReasonMAX_LESS_THAN_MIN_IN_DURING_CLAUSE QueryErrorReason = "MAX_LESS_THAN_MIN_IN_DURING_CLAUSE"

	//
	// The query matched the grammar, but is invalid in some way such as using a service that
	// isn't supported.
	//
	QueryErrorReasonVALIDATION_FAILED QueryErrorReason = "VALIDATION_FAILED"
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
// The reasons for the setting error.
//
type SettingErrorReason string

const (

	//
	// The campaign already has a setting of the type that is being added.
	//
	SettingErrorReasonDUPLICATE_SETTING_TYPE SettingErrorReason = "DUPLICATE_SETTING_TYPE"

	//
	// The campaign setting is not available for this AdWords account.
	//
	SettingErrorReasonSETTING_TYPE_IS_NOT_AVAILABLE SettingErrorReason = "SETTING_TYPE_IS_NOT_AVAILABLE"

	//
	// The setting is not compatible with the campaign.
	//
	SettingErrorReasonSETTING_TYPE_IS_NOT_COMPATIBLE_WITH_CAMPAIGN SettingErrorReason = "SETTING_TYPE_IS_NOT_COMPATIBLE_WITH_CAMPAIGN"

	//
	// The supplied TargetingSetting contains an invalid CriterionTypeGroup. See
	// {@link CriterionTypeGroup} documentation for CriterionTypeGroups allowed in Campaign or
	// AdGroup TargetingSettings.
	//
	SettingErrorReasonTARGETING_SETTING_CONTAINS_INVALID_CRITERION_TYPE_GROUP SettingErrorReason = "TARGETING_SETTING_CONTAINS_INVALID_CRITERION_TYPE_GROUP"

	//
	// Starting with AdWords API v201802, TargetingSetting must not explicitly set any of the
	// Demographic CriterionTypeGroups (AGE_RANGE, GENDER, PARENT, INCOME_RANGE) to false (it's
	// okay to not set them at all, in which case the system will set them to true automatically).
	//
	SettingErrorReasonTARGETING_SETTING_DEMOGRAPHIC_CRITERION_TYPE_GROUPS_MUST_BE_SET_TO_TARGET_ALL SettingErrorReason = "TARGETING_SETTING_DEMOGRAPHIC_CRITERION_TYPE_GROUPS_MUST_BE_SET_TO_TARGET_ALL"

	//
	// Starting with AdWords API v201802,TargetingSetting cannot change any of the Demographic
	// CriterionTypeGroups (AGE_RANGE, GENDER, PARENT, INCOME_RANGE) from true to false.
	//
	SettingErrorReasonTARGETING_SETTING_CANNOT_CHANGE_TARGET_ALL_TO_FALSE_FOR_DEMOGRAPHIC_CRITERION_TYPE_GROUP SettingErrorReason = "TARGETING_SETTING_CANNOT_CHANGE_TARGET_ALL_TO_FALSE_FOR_DEMOGRAPHIC_CRITERION_TYPE_GROUP"

	//
	// At least one feed id should be present.
	//
	SettingErrorReasonDYNAMIC_SEARCH_ADS_SETTING_AT_LEAST_ONE_FEED_ID_MUST_BE_PRESENT SettingErrorReason = "DYNAMIC_SEARCH_ADS_SETTING_AT_LEAST_ONE_FEED_ID_MUST_BE_PRESENT"

	//
	// The supplied DynamicSearchAdsSetting contains an invalid domain name.
	//
	SettingErrorReasonDYNAMIC_SEARCH_ADS_SETTING_CONTAINS_INVALID_DOMAIN_NAME SettingErrorReason = "DYNAMIC_SEARCH_ADS_SETTING_CONTAINS_INVALID_DOMAIN_NAME"

	//
	// The supplied DynamicSearchAdsSetting contains a subdomain name.
	//
	SettingErrorReasonDYNAMIC_SEARCH_ADS_SETTING_CONTAINS_SUBDOMAIN_NAME SettingErrorReason = "DYNAMIC_SEARCH_ADS_SETTING_CONTAINS_SUBDOMAIN_NAME"

	//
	// The supplied DynamicSearchAdsSetting contains an invalid language code.
	//
	SettingErrorReasonDYNAMIC_SEARCH_ADS_SETTING_CONTAINS_INVALID_LANGUAGE_CODE SettingErrorReason = "DYNAMIC_SEARCH_ADS_SETTING_CONTAINS_INVALID_LANGUAGE_CODE"

	//
	// TargetingSettings in search campaigns should not have CriterionTypeGroup.PLACEMENT
	// set to targetAll.
	//
	SettingErrorReasonTARGET_ALL_IS_NOT_ALLOWED_FOR_PLACEMENT_IN_SEARCH_CAMPAIGN SettingErrorReason = "TARGET_ALL_IS_NOT_ALLOWED_FOR_PLACEMENT_IN_SEARCH_CAMPAIGN"

	//
	// Duplicate description in universal app setting description field.
	//
	SettingErrorReasonUNIVERSAL_APP_CAMPAIGN_SETTING_DUPLICATE_DESCRIPTION SettingErrorReason = "UNIVERSAL_APP_CAMPAIGN_SETTING_DUPLICATE_DESCRIPTION"

	//
	// Description line width is too long in universal app setting description field.
	//
	SettingErrorReasonUNIVERSAL_APP_CAMPAIGN_SETTING_DESCRIPTION_LINE_WIDTH_TOO_LONG SettingErrorReason = "UNIVERSAL_APP_CAMPAIGN_SETTING_DESCRIPTION_LINE_WIDTH_TOO_LONG"

	//
	// Universal app setting appId field cannot be modified for COMPLETE campaigns.
	//
	SettingErrorReasonUNIVERSAL_APP_CAMPAIGN_SETTING_APP_ID_CANNOT_BE_MODIFIED SettingErrorReason = "UNIVERSAL_APP_CAMPAIGN_SETTING_APP_ID_CANNOT_BE_MODIFIED"

	//
	// YoutubeVideoMediaIds in universal app setting cannot exceed size limit.
	//
	SettingErrorReasonTOO_MANY_YOUTUBE_MEDIA_IDS_IN_UNIVERSAL_APP_CAMPAIGN SettingErrorReason = "TOO_MANY_YOUTUBE_MEDIA_IDS_IN_UNIVERSAL_APP_CAMPAIGN"

	//
	// ImageMediaIds in universal app setting cannot exceed size limit.
	//
	SettingErrorReasonTOO_MANY_IMAGE_MEDIA_IDS_IN_UNIVERSAL_APP_CAMPAIGN SettingErrorReason = "TOO_MANY_IMAGE_MEDIA_IDS_IN_UNIVERSAL_APP_CAMPAIGN"

	//
	// Media is incompatible for universal app campaign.
	//
	SettingErrorReasonMEDIA_INCOMPATIBLE_FOR_UNIVERSAL_APP_CAMPAIGN SettingErrorReason = "MEDIA_INCOMPATIBLE_FOR_UNIVERSAL_APP_CAMPAIGN"

	//
	// Too many exclamation marks in universal app campaign ad text ideas.
	//
	SettingErrorReasonTOO_MANY_EXCLAMATION_MARKS SettingErrorReason = "TOO_MANY_EXCLAMATION_MARKS"

	//
	// Unspecified campaign setting error.
	// <span class="constraint Rejected">Used for return value only. An enumeration could not be processed, typically due to incompatibility with your WSDL version.</span>
	//
	SettingErrorReasonUNKNOWN SettingErrorReason = "UNKNOWN"
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
// The reasons for errors when querying for stats.
//
type StatsQueryErrorReason string

const (

	//
	// Date is outside of allowed range.
	//
	StatsQueryErrorReasonDATE_NOT_IN_VALID_RANGE StatsQueryErrorReason = "DATE_NOT_IN_VALID_RANGE"
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
// The reasons for the url error.
//
type UrlErrorReason string

const (

	//
	// The tracking url template is invalid.
	//
	UrlErrorReasonINVALID_TRACKING_URL_TEMPLATE UrlErrorReason = "INVALID_TRACKING_URL_TEMPLATE"

	//
	// The tracking url template contains invalid tag.
	//
	UrlErrorReasonINVALID_TAG_IN_TRACKING_URL_TEMPLATE UrlErrorReason = "INVALID_TAG_IN_TRACKING_URL_TEMPLATE"

	//
	// The tracking url template must contain at least one tag (e.g. {lpurl}),
	// This applies only to tracking url template associated with website ads or product ads.
	//
	UrlErrorReasonMISSING_TRACKING_URL_TEMPLATE_TAG UrlErrorReason = "MISSING_TRACKING_URL_TEMPLATE_TAG"

	//
	// The tracking url template must start with a valid protocol (or lpurl tag).
	//
	UrlErrorReasonMISSING_PROTOCOL_IN_TRACKING_URL_TEMPLATE UrlErrorReason = "MISSING_PROTOCOL_IN_TRACKING_URL_TEMPLATE"

	//
	// The tracking url template starts with an invalid protocol.
	//
	UrlErrorReasonINVALID_PROTOCOL_IN_TRACKING_URL_TEMPLATE UrlErrorReason = "INVALID_PROTOCOL_IN_TRACKING_URL_TEMPLATE"

	//
	// The tracking url template  contains illegal characters.
	//
	UrlErrorReasonMALFORMED_TRACKING_URL_TEMPLATE UrlErrorReason = "MALFORMED_TRACKING_URL_TEMPLATE"

	//
	// The tracking url template must contain a host name (or lpurl tag).
	//
	UrlErrorReasonMISSING_HOST_IN_TRACKING_URL_TEMPLATE UrlErrorReason = "MISSING_HOST_IN_TRACKING_URL_TEMPLATE"

	//
	// The tracking url template has an invalid or missing top level domain extension.
	//
	UrlErrorReasonINVALID_TLD_IN_TRACKING_URL_TEMPLATE UrlErrorReason = "INVALID_TLD_IN_TRACKING_URL_TEMPLATE"

	//
	// The tracking url template contains nested occurrences of the same conditional tag
	// (i.e. {ifmobile:{ifmobile:x}}).
	//
	UrlErrorReasonREDUNDANT_NESTED_TRACKING_URL_TEMPLATE_TAG UrlErrorReason = "REDUNDANT_NESTED_TRACKING_URL_TEMPLATE_TAG"

	//
	// The final url is invalid.
	//
	UrlErrorReasonINVALID_FINAL_URL UrlErrorReason = "INVALID_FINAL_URL"

	//
	// The final url contains invalid tag.
	//
	UrlErrorReasonINVALID_TAG_IN_FINAL_URL UrlErrorReason = "INVALID_TAG_IN_FINAL_URL"

	//
	// The final url contains nested occurrences of the same conditional tag
	// (i.e. {ifmobile:{ifmobile:x}}).
	//
	UrlErrorReasonREDUNDANT_NESTED_FINAL_URL_TAG UrlErrorReason = "REDUNDANT_NESTED_FINAL_URL_TAG"

	//
	// The final url must start with a valid protocol.
	//
	UrlErrorReasonMISSING_PROTOCOL_IN_FINAL_URL UrlErrorReason = "MISSING_PROTOCOL_IN_FINAL_URL"

	//
	// The final url starts with an invalid protocol.
	//
	UrlErrorReasonINVALID_PROTOCOL_IN_FINAL_URL UrlErrorReason = "INVALID_PROTOCOL_IN_FINAL_URL"

	//
	// The final url  contains illegal characters.
	//
	UrlErrorReasonMALFORMED_FINAL_URL UrlErrorReason = "MALFORMED_FINAL_URL"

	//
	// The final url must contain a host name.
	//
	UrlErrorReasonMISSING_HOST_IN_FINAL_URL UrlErrorReason = "MISSING_HOST_IN_FINAL_URL"

	//
	// The tracking url template has an invalid or missing top level domain extension.
	//
	UrlErrorReasonINVALID_TLD_IN_FINAL_URL UrlErrorReason = "INVALID_TLD_IN_FINAL_URL"

	//
	// The final mobile url is invalid.
	//
	UrlErrorReasonINVALID_FINAL_MOBILE_URL UrlErrorReason = "INVALID_FINAL_MOBILE_URL"

	//
	// The final mobile url contains invalid tag.
	//
	UrlErrorReasonINVALID_TAG_IN_FINAL_MOBILE_URL UrlErrorReason = "INVALID_TAG_IN_FINAL_MOBILE_URL"

	//
	// The final mobile url contains nested occurrences of the same conditional tag
	// (i.e. {ifmobile:{ifmobile:x}}).
	//
	UrlErrorReasonREDUNDANT_NESTED_FINAL_MOBILE_URL_TAG UrlErrorReason = "REDUNDANT_NESTED_FINAL_MOBILE_URL_TAG"

	//
	// The final mobile url must start with a valid protocol.
	//
	UrlErrorReasonMISSING_PROTOCOL_IN_FINAL_MOBILE_URL UrlErrorReason = "MISSING_PROTOCOL_IN_FINAL_MOBILE_URL"

	//
	// The final mobile url starts with an invalid protocol.
	//
	UrlErrorReasonINVALID_PROTOCOL_IN_FINAL_MOBILE_URL UrlErrorReason = "INVALID_PROTOCOL_IN_FINAL_MOBILE_URL"

	//
	// The final mobile url  contains illegal characters.
	//
	UrlErrorReasonMALFORMED_FINAL_MOBILE_URL UrlErrorReason = "MALFORMED_FINAL_MOBILE_URL"

	//
	// The final mobile url must contain a host name.
	//
	UrlErrorReasonMISSING_HOST_IN_FINAL_MOBILE_URL UrlErrorReason = "MISSING_HOST_IN_FINAL_MOBILE_URL"

	//
	// The tracking url template has an invalid or missing top level domain extension.
	//
	UrlErrorReasonINVALID_TLD_IN_FINAL_MOBILE_URL UrlErrorReason = "INVALID_TLD_IN_FINAL_MOBILE_URL"

	//
	// The final app url is invalid.
	//
	UrlErrorReasonINVALID_FINAL_APP_URL UrlErrorReason = "INVALID_FINAL_APP_URL"

	//
	// The final app url contains invalid tag.
	//
	UrlErrorReasonINVALID_TAG_IN_FINAL_APP_URL UrlErrorReason = "INVALID_TAG_IN_FINAL_APP_URL"

	//
	// The final app url contains nested occurrences of the same conditional tag
	// (i.e. {ifmobile:{ifmobile:x}}).
	//
	UrlErrorReasonREDUNDANT_NESTED_FINAL_APP_URL_TAG UrlErrorReason = "REDUNDANT_NESTED_FINAL_APP_URL_TAG"

	//
	// More than one app url found for the same OS type.
	//
	UrlErrorReasonMULTIPLE_APP_URLS_FOR_OSTYPE UrlErrorReason = "MULTIPLE_APP_URLS_FOR_OSTYPE"

	//
	// The OS type given for an app url is not valid.
	//
	UrlErrorReasonINVALID_OSTYPE UrlErrorReason = "INVALID_OSTYPE"

	//
	// The protocol given for an app url is not valid. (E.g. "android-app://")
	//
	UrlErrorReasonINVALID_PROTOCOL_FOR_APP_URL UrlErrorReason = "INVALID_PROTOCOL_FOR_APP_URL"

	//
	// The package id (app id) given for an app url is not valid.
	//
	UrlErrorReasonINVALID_PACKAGE_ID_FOR_APP_URL UrlErrorReason = "INVALID_PACKAGE_ID_FOR_APP_URL"

	//
	// The number of url custom parameters for an entity exceeds the maximum limit allowed.
	//
	UrlErrorReasonURL_CUSTOM_PARAMETERS_COUNT_EXCEEDS_LIMIT UrlErrorReason = "URL_CUSTOM_PARAMETERS_COUNT_EXCEEDS_LIMIT"

	//
	// The parameter has isRemove set to true but a value that is non-null.
	//
	UrlErrorReasonURL_CUSTOM_PARAMETER_REMOVAL_WITH_NON_NULL_VALUE UrlErrorReason = "URL_CUSTOM_PARAMETER_REMOVAL_WITH_NON_NULL_VALUE"

	//
	// For add operations, there will not be any existing parameters to delete.
	//
	UrlErrorReasonCANNOT_REMOVE_URL_CUSTOM_PARAMETER_IN_ADD_OPERATION UrlErrorReason = "CANNOT_REMOVE_URL_CUSTOM_PARAMETER_IN_ADD_OPERATION"

	//
	// When the doReplace flag is set to true, individual parameters cannot be deleted.
	//
	UrlErrorReasonCANNOT_REMOVE_URL_CUSTOM_PARAMETER_DURING_FULL_REPLACEMENT UrlErrorReason = "CANNOT_REMOVE_URL_CUSTOM_PARAMETER_DURING_FULL_REPLACEMENT"

	//
	// For ADD operations and when the doReplace flag is set to true, custom parameter values
	// cannot be null.
	//
	UrlErrorReasonNULL_CUSTOM_PARAMETER_VALUE_DURING_ADD_OR_FULL_REPLACEMENT UrlErrorReason = "NULL_CUSTOM_PARAMETER_VALUE_DURING_ADD_OR_FULL_REPLACEMENT"

	//
	// An invalid character appears in the parameter key.
	//
	UrlErrorReasonINVALID_CHARACTERS_IN_URL_CUSTOM_PARAMETER_KEY UrlErrorReason = "INVALID_CHARACTERS_IN_URL_CUSTOM_PARAMETER_KEY"

	//
	// An invalid character appears in the parameter value.
	//
	UrlErrorReasonINVALID_CHARACTERS_IN_URL_CUSTOM_PARAMETER_VALUE UrlErrorReason = "INVALID_CHARACTERS_IN_URL_CUSTOM_PARAMETER_VALUE"

	//
	// The url custom parameter value fails url tag validation.
	//
	UrlErrorReasonINVALID_TAG_IN_URL_CUSTOM_PARAMETER_VALUE UrlErrorReason = "INVALID_TAG_IN_URL_CUSTOM_PARAMETER_VALUE"

	//
	// The custom parameter contains nested occurrences of the same conditional tag
	// (i.e. {ifmobile:{ifmobile:x}}).
	//
	UrlErrorReasonREDUNDANT_NESTED_URL_CUSTOM_PARAMETER_TAG UrlErrorReason = "REDUNDANT_NESTED_URL_CUSTOM_PARAMETER_TAG"

	//
	// The protocol (http:// or https://) is missing.
	//
	UrlErrorReasonMISSING_PROTOCOL UrlErrorReason = "MISSING_PROTOCOL"

	//
	// The url is invalid.
	//
	UrlErrorReasonINVALID_URL UrlErrorReason = "INVALID_URL"

	//
	// Destination Url is deprecated.
	//
	UrlErrorReasonDESTINATION_URL_DEPRECATED UrlErrorReason = "DESTINATION_URL_DEPRECATED"

	//
	// The url contains invalid tag.
	//
	UrlErrorReasonINVALID_TAG_IN_URL UrlErrorReason = "INVALID_TAG_IN_URL"

	//
	// The url must contain at least one tag (e.g. {lpurl}),
	// This applies only to urls associated with website ads or product ads.
	//
	UrlErrorReasonMISSING_URL_TAG UrlErrorReason = "MISSING_URL_TAG"

	UrlErrorReasonDUPLICATE_URL_ID UrlErrorReason = "DUPLICATE_URL_ID"

	UrlErrorReasonINVALID_URL_ID UrlErrorReason = "INVALID_URL_ID"

	UrlErrorReasonURL_ERROR UrlErrorReason = "URL_ERROR"
)

type Get struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 get"`

	//
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	//
	ServiceSelector *Selector `xml:"serviceSelector,omitempty"`
}

type GetResponse struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 getResponse"`

	Rval *AdGroupPage `xml:"rval,omitempty"`
}

type Mutate struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 mutate"`

	//
	// <span class="constraint ContentsNotNull">This field must not contain {@code null} elements.</span>
	// <span class="constraint DistinctIds">Elements in this field must have distinct IDs for following {@link Operator}s : SET, REMOVE.</span>
	// <span class="constraint NotEmpty">This field must contain at least one element.</span>
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	// <span class="constraint SupportedOperators">The following {@link Operator}s are supported: ADD, SET.</span>
	//
	Operations []*AdGroupOperation `xml:"operations,omitempty"`
}

type MutateResponse struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 mutateResponse"`

	Rval *AdGroupReturnValue `xml:"rval,omitempty"`
}

type MutateLabel struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 mutateLabel"`

	//
	// <span class="constraint ContentsNotNull">This field must not contain {@code null} elements.</span>
	// <span class="constraint DistinctIds">Elements in this field must have distinct IDs for following {@link Operator}s : ADD, REMOVE.</span>
	// <span class="constraint NotEmpty">This field must contain at least one element.</span>
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	// <span class="constraint SupportedOperators">The following {@link Operator}s are supported: ADD, REMOVE.</span>
	//
	Operations []*AdGroupLabelOperation `xml:"operations,omitempty"`
}

type MutateLabelResponse struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 mutateLabelResponse"`

	Rval *AdGroupLabelReturnValue `xml:"rval,omitempty"`
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

	Rval *AdGroupPage `xml:"rval,omitempty"`
}

type AdGroup struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 AdGroup"`

	//
	// ID of this ad group.
	// <span class="constraint Selectable">This field can be selected using the value "Id".</span><span class="constraint Filterable">This field can be filtered on.</span>
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API for the following {@link Operator}s: ADD.</span>
	// <span class="constraint Required">This field is required and should not be {@code null} when it is contained within {@link Operator}s : SET.</span>
	//
	Id int64 `xml:"id,omitempty"`

	//
	// ID of the campaign with which this ad group is associated.
	// <span class="constraint Selectable">This field can be selected using the value "CampaignId".</span><span class="constraint Filterable">This field can be filtered on.</span>
	// <span class="constraint Required">This field is required and should not be {@code null} when it is contained within {@link Operator}s : ADD.</span>
	//
	CampaignId int64 `xml:"campaignId,omitempty"`

	//
	// Name of the campaign with which this ad group is associated.
	// <span class="constraint Selectable">This field can be selected using the value "CampaignName".</span><span class="constraint Filterable">This field can be filtered on.</span>
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	CampaignName string `xml:"campaignName,omitempty"`

	//
	// Name of this ad group (at most 255 UTF-8 full-width characters).
	// This field is required and should not be {@code null} for ADD operations from v201309 version.
	// <span class="constraint Selectable">This field can be selected using the value "Name".</span><span class="constraint Filterable">This field can be filtered on.</span>
	// <span class="constraint MatchesRegex">AdGroup names must not contain any null (code point 0x0), NL line feed (code point 0xA) or carriage return (code point 0xD) characters. This is checked by the regular expression '[^\x00\x0A\x0D]*'.</span>
	//
	Name string `xml:"name,omitempty"`

	//
	// Status of this ad group.
	// <span class="constraint Selectable">This field can be selected using the value "Status".</span><span class="constraint Filterable">This field can be filtered on.</span>
	//
	Status *AdGroupStatus `xml:"status,omitempty"`

	//
	// List of settings for the AdGroup.
	// <span class="constraint Selectable">This field can be selected using the value "Settings".</span>
	//
	Settings []*Setting `xml:"settings,omitempty"`

	//
	// Labels that are attached to the {@link AdGroup}. To associate an existing {@link Label} to an
	// existing {@link AdGroup}, use {@link AdGroupService#mutateLabel} with ADD
	// operator. To remove an associated {@link Label} from the {@link AdGroup}, use
	// {@link AdGroupService#mutateLabel} with REMOVE operator. To filter on {@link Label}s,
	// use one of {@link Predicate.Operator#CONTAINS_ALL}, {@link Predicate.Operator#CONTAINS_ANY},
	// {@link Predicate.Operator#CONTAINS_NONE} operators with a list of {@link Label} ids.
	// <span class="constraint Selectable">This field can be selected using the value "Labels".</span><span class="constraint Filterable">This field can be filtered on.</span>
	// <span class="constraint CampaignType">This field may not be set for campaign channel subtype UNIVERSAL_APP_CAMPAIGN.</span>
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API for the following {@link Operator}s: REMOVE and SET.</span>
	//
	Labels []*Label `xml:"labels,omitempty"`

	//
	// This Map provides a place to put new features and settings in older versions
	// of the AdWords API in the rare instance we need to introduce a new feature in
	// an older version.
	//
	// It is presently unused.  Do not set a value.
	//
	ForwardCompatibilityMap []*String_StringMapEntry `xml:"forwardCompatibilityMap,omitempty"`

	//
	// Bidding configuration for this ad group. To set the bids on the ad groups use
	// {@link BiddingStrategyConfiguration#bids}. Multiple bids can be set on ad group at the same
	// time. Only the bids that apply to the effective bidding strategy will be used. Effective
	// bidding strategy is considered to be the directly attached strategy or inherited campaign level
	// strategy when there?s no directly attached strategy.
	//
	BiddingStrategyConfiguration *BiddingStrategyConfiguration `xml:"biddingStrategyConfiguration,omitempty"`

	//
	// Allows advertisers to specify a criteria dimension on which to place absolute bids.
	// This is only applicable for campaigns that target only the content network and not search.
	// <span class="constraint Selectable">This field can be selected using the value "ContentBidCriterionTypeGroup".</span><span class="constraint Filterable">This field can be filtered on.</span>
	//
	ContentBidCriterionTypeGroup *CriterionTypeGroup `xml:"contentBidCriterionTypeGroup,omitempty"`

	//
	// ID of the base campaign from which this draft/trial adgroup was created. This
	// field is only returned on get requests.
	// <span class="constraint Selectable">This field can be selected using the value "BaseCampaignId".</span>
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	BaseCampaignId int64 `xml:"baseCampaignId,omitempty"`

	//
	// ID of the base adgroup from which this draft/trial adgroup was created. For
	// base adgroups this is equal to the adgroup ID.  If the adgroup was created
	// in the draft or trial and has no corresponding base adgroup, this field is null.
	// This field is readonly and will be ignored when sent to the API.
	// <span class="constraint Selectable">This field can be selected using the value "BaseAdGroupId".</span>
	//
	BaseAdGroupId int64 `xml:"baseAdGroupId,omitempty"`

	//
	// URL template for constructing a tracking URL.
	//
	// <p>On update, empty string ("") indicates to clear the field.
	// <span class="constraint Selectable">This field can be selected using the value "TrackingUrlTemplate".</span><span class="constraint Filterable">This field can be filtered on.</span>
	//
	TrackingUrlTemplate string `xml:"trackingUrlTemplate,omitempty"`

	//
	// URL template for appending params to Final URL.
	//
	// <p>On update, empty string ("") indicates to clear the field.
	// <p>This field is supported only in test accounts.
	// <span class="constraint Selectable">This field can be selected using the value "FinalUrlSuffix".</span><span class="constraint Filterable">This field can be filtered on.</span>
	//
	FinalUrlSuffix string `xml:"finalUrlSuffix,omitempty"`

	//
	// A list of mappings to be used for substituting URL custom parameter tags in the
	// trackingUrlTemplate, finalUrls, and/or finalMobileUrls.
	// <span class="constraint Selectable">This field can be selected using the value "UrlCustomParameters".</span>
	//
	UrlCustomParameters *CustomParameters `xml:"urlCustomParameters,omitempty"`

	//
	// <span class="constraint Selectable">This field can be selected using the value "AdGroupType".</span><span class="constraint Filterable">This field can be filtered on.</span>
	// <span class="constraint CampaignType">This field may only be set to SHOPPING_UNIVERSAL_ADS for campaign channel type SHOPPING with campaign channel subtype SHOPPING_UNIVERSAL_ADS.</span>
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API for the following {@link Operator}s: SET.</span>
	//
	AdGroupType *AdGroupType `xml:"adGroupType,omitempty"`

	//
	// The ad rotation mode to specify how often the ads in the ad group will be served relative to
	// one another.
	//
	AdGroupAdRotationMode *AdGroupAdRotationMode `xml:"adGroupAdRotationMode,omitempty"`
}

type AdGroupAdRotationMode struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 AdGroupAdRotationMode"`

	//
	// <span class="constraint Selectable">This field can be selected using the value "AdRotationMode".</span><span class="constraint Filterable">This field can be filtered on.</span>
	//
	AdRotationMode *AdRotationMode `xml:"adRotationMode,omitempty"`
}

type AdGroupLabel struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 AdGroupLabel"`

	//
	// The id of the adGroup that the label is applied to.
	// <span class="constraint Required">This field is required and should not be {@code null} when it is contained within {@link Operator}s : ADD, REMOVE.</span>
	//
	AdGroupId int64 `xml:"adGroupId,omitempty"`

	//
	// The id of an existing label to be applied to the ad group.
	// <span class="constraint Required">This field is required and should not be {@code null} when it is contained within {@link Operator}s : ADD, REMOVE.</span>
	//
	LabelId int64 `xml:"labelId,omitempty"`
}

type AdGroupLabelOperation struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 AdGroupLabelOperation"`

	*Operation

	//
	// AdGroupLabel to operate on.
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	//
	Operand *AdGroupLabel `xml:"operand,omitempty"`
}

type AdGroupLabelReturnValue struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 AdGroupLabelReturnValue"`

	*ListReturnValue

	Value []*AdGroupLabel `xml:"value,omitempty"`

	PartialFailureErrors []*ApiError `xml:"partialFailureErrors,omitempty"`
}

type AdGroupOperation struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 AdGroupOperation"`

	*Operation

	//
	// AdGroup to operate on
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	//
	Operand *AdGroup `xml:"operand,omitempty"`
}

type AdGroupPage struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 AdGroupPage"`

	*Page

	//
	// The result entries in this page.
	//
	Entries []*AdGroup `xml:"entries,omitempty"`
}

type AdGroupReturnValue struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 AdGroupReturnValue"`

	*ListReturnValue

	//
	// List of AdGroups
	//
	Value []*AdGroup `xml:"value,omitempty"`

	//
	// List of partial failure errors.
	//
	PartialFailureErrors []*ApiError `xml:"partialFailureErrors,omitempty"`
}

type AdGroupServiceError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 AdGroupServiceError"`

	*ApiError

	//
	// The error reason represented by an enum.
	//
	Reason *AdGroupServiceErrorReason `xml:"reason,omitempty"`
}

type AdxError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 AdxError"`

	*ApiError

	//
	// The error reason represented by an enum.
	//
	Reason *AdxErrorReason `xml:"reason,omitempty"`
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

type LabelAttribute struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 LabelAttribute"`

	//
	// Indicates that this instance is a subtype of LabelAttribute.
	// Although this field is returned in the response, it is ignored on input
	// and cannot be selected. Specify xsi:type instead.
	//
	LabelAttributeType string `xml:"LabelAttribute.Type,omitempty"`
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

type BiddingScheme struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 BiddingScheme"`

	//
	// Indicates that this instance is a subtype of BiddingScheme.
	// Although this field is returned in the response, it is ignored on input
	// and cannot be selected. Specify xsi:type instead.
	//
	BiddingSchemeType string `xml:"BiddingScheme.Type,omitempty"`
}

type BiddingStrategyConfiguration struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 BiddingStrategyConfiguration"`

	//
	// Id of the bidding strategy to be associated with the campaign, ad group or ad group criteria. A
	// bidding strategy is created using the BiddingStrategyService ADD operation and is assigned a
	// BiddingStrategyId. The BiddingStrategyId can be shared across campaigns, ad groups and ad group
	// criteria.
	//
	// <p>Starting with v201705, this field cannot be set at the ad group or ad group criterion level.
	// <span class="constraint Selectable">This field can be selected using the value "BiddingStrategyId".</span><span class="constraint Filterable">This field can be filtered on.</span>
	// <span class="constraint AdxEnabled">This is disabled for AdX.</span>
	//
	BiddingStrategyId int64 `xml:"biddingStrategyId,omitempty"`

	//
	// Name of the bidding strategy. This is applicable only for flexible bidding strategies.
	// <span class="constraint Selectable">This field can be selected using the value "BiddingStrategyName".</span><span class="constraint Filterable">This field can be filtered on.</span>
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	BiddingStrategyName string `xml:"biddingStrategyName,omitempty"`

	//
	// The type of the bidding strategy to be attached.
	//
	// <p>For details on portfolio vs. standard availability, see the <a
	// href="https://developers.google.com/adwords/api/docs/guides/bidding">bidding guide</a>.
	//
	// <p>Starting with v201705, this field cannot be set at the ad group or ad group criterion level
	// to any value other than {@code BiddingStrategyType.NONE}.
	// <span class="constraint Selectable">This field can be selected using the value "BiddingStrategyType".</span><span class="constraint Filterable">This field can be filtered on.</span>
	// <span class="constraint CampaignType">This field may only be set to the values: MANUAL_CPC, PAGE_ONE_PROMOTED, TARGET_SPEND, ENHANCED_CPC, TARGET_CPA, TARGET_ROAS, MAXIMIZE_CONVERSIONS, MAXIMIZE_CONVERSION_VALUE, TARGET_OUTRANK_SHARE, NONE, MANUAL_CPM for campaign channel type SEARCH.</span>
	// <span class="constraint CampaignType">This field may only be set to the values: MANUAL_CPC, ENHANCED_CPC, NONE for campaign channel type SHOPPING with ad group type SHOPPING_SHOWCASE_ADS.</span>
	// <span class="constraint CampaignType">This field may only be set to the values: MANUAL_CPC, ENHANCED_CPC, TARGET_ROAS, TARGET_SPEND, NONE for campaign channel type SHOPPING.</span>
	// <span class="constraint CampaignType">This field may only be set to NONE for campaign channel type SHOPPING with campaign channel subtype SHOPPING_UNIVERSAL_ADS.</span>
	// <span class="constraint CampaignType">This field may only be set to the values: MANUAL_CPC, ENHANCED_CPC, TARGET_ROAS, TARGET_SPEND, TARGET_CPA, NONE, MAXIMIZE_CONVERSIONS for campaign channel type DISPLAY.</span>
	// <span class="constraint CampaignType">This field may only be set to NONE for campaign channel type DISPLAY with campaign channel subtype DISPLAY_MOBILE_APP.</span>
	// <span class="constraint CampaignType">This field may only be set to the values: MANUAL_CPC, MAXIMIZE_CONVERSIONS, NONE, PAGE_ONE_PROMOTED, TARGET_CPA, TARGET_OUTRANK_SHARE, TARGET_ROAS, TARGET_SPEND for campaign channel subtype SEARCH_MOBILE_APP.</span>
	// <span class="constraint CampaignType">This field may only be set to NONE for campaign channel subtype UNIVERSAL_APP_CAMPAIGN.</span>
	// <span class="constraint CampaignType">This field may only be set to the values: MANUAL_CPC, ENHANCED_CPC, TARGET_CPA, TARGET_ROAS, NONE for campaign channel type DISPLAY with campaign channel subtype DISPLAY_GMAIL_AD.</span>
	//
	BiddingStrategyType *BiddingStrategyType `xml:"biddingStrategyType,omitempty"`

	//
	// Indicates where the bidding strategy is associated i.e. campaign, ad group or
	// ad group criterion.
	// <span class="constraint Selectable">This field can be selected using the value "BiddingStrategySource".</span><span class="constraint Filterable">This field can be filtered on.</span>
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	BiddingStrategySource *BiddingStrategySource `xml:"biddingStrategySource,omitempty"`

	//
	// The bidding strategy metadata. Bidding strategy can be associated using the {@linkplain
	// BiddingStrategyConfiguration#biddingStrategyType} or the bidding scheme.
	//
	// <p>For details on portfolio vs. standard availability, see the <a
	// href="https://developers.google.com/adwords/api/docs/guides/bidding">bidding guide</a>.
	//
	// <p>Starting with v201705, this field cannot be set at the ad group or ad group criterion level.
	//
	BiddingScheme *BiddingScheme `xml:"biddingScheme,omitempty"`

	//
	// Specifies the bids. Bids can be set only on ad groups and ad group criteria.
	// Bids cannot be set on campaign.
	//
	// Default CPC and CPM bid values will be set if they are not provided during {@linkplain AdGroup}
	// creation. Default CPC and CPM values are minimal billable amounts in local currencies.
	// For example, for US Dollars CPC and CPM default values are $0.01 and $0.01, respectively.
	//
	Bids []*Bids `xml:"bids,omitempty"`

	//
	// The target return on average spend (ROAS). This target can only be set on ad groups. If this
	// ad group's effective bidding strategy is a standard {@code TARGET_ROAS} strategy attached to
	// the campaign, then the target overrides the target roas specified in the campaign's bidding
	// strategy. Otherwise, this value is ignored.
	// <span class="constraint Selectable">This field can be selected using the value "TargetRoasOverride".</span><span class="constraint Filterable">This field can be filtered on.</span>
	// <span class="constraint InRange">This field must be between 0.01 and 1000.0, inclusive.</span>
	//
	TargetRoasOverride float64 `xml:"targetRoasOverride,omitempty"`
}

type Bids struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 Bids"`

	//
	// Indicates that this instance is a subtype of Bids.
	// Although this field is returned in the response, it is ignored on input
	// and cannot be selected. Specify xsi:type instead.
	//
	BidsType string `xml:"Bids.Type,omitempty"`
}

type TextLabel struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 TextLabel"`

	*Label
}

type DisplayAttribute struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 DisplayAttribute"`

	*LabelAttribute

	//
	// Background color of the label in RGB format.
	// <span class="constraint MatchesRegex">A background color string must begin with a '#' character followed by either 6 or 3 hexadecimal characters (24 vs. 12 bits). This is checked by the regular expression '^\#([a-fA-F0-9]{6}|[a-fA-F0-9]{3})$'.</span>
	//
	BackgroundColor string `xml:"backgroundColor,omitempty"`

	//
	// A short description of the label.
	// <span class="constraint StringLength">The length of this string should be between 0 and 200, inclusive.</span>
	//
	Description string `xml:"description,omitempty"`
}

type ClientTermsError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 ClientTermsError"`

	*ApiError

	Reason *ClientTermsErrorReason `xml:"reason,omitempty"`
}

type ComparableValue struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 ComparableValue"`

	//
	// Indicates that this instance is a subtype of ComparableValue.
	// Although this field is returned in the response, it is ignored on input
	// and cannot be selected. Specify xsi:type instead.
	//
	ComparableValueType string `xml:"ComparableValue.Type,omitempty"`
}

type CpaBid struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 CpaBid"`

	*Bids

	//
	// Target cost per acquisition (CPA). This is applicable only at the ad group level.
	//
	// <p>If an ad group-level target is not set and the strategy type is TARGET_CPA,
	// the strategy level target will be used. To set the strategy-level target,
	// set the {@linkplain TargetCpaBiddingScheme#targetCpa} on the strategy's
	// {@linkplain BiddingStrategyConfiguration#biddingScheme}.
	// <span class="constraint Selectable">This field can be selected using the value "TargetCpaBid".</span><span class="constraint Filterable">This field can be filtered on.</span>
	//
	Bid *Money `xml:"bid,omitempty"`

	//
	// The level (ad group, ad group strategy, or campaign strategy) at which the bid was set.
	// This is applicable only at the ad group level.
	// <span class="constraint Selectable">This field can be selected using the value "TargetCpaBidSource".</span><span class="constraint Filterable">This field can be filtered on.</span>
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	BidSource *BidSource `xml:"bidSource,omitempty"`
}

type CpcBid struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 CpcBid"`

	*Bids

	//
	// Max CPC (cost per click) bid.
	// At the ad group level, this represents the default bid applicable for
	// <ul><li>keyword targeting on search network.</li>
	// <li>keywords & placements for content targeting.</li></ul>
	// At the ad group criteria level, this is the max cpc bid.
	// <span class="constraint Selectable">This field can be selected using the value "CpcBid".</span><span class="constraint Filterable">This field can be filtered on.</span>
	//
	Bid *Money `xml:"bid,omitempty"`

	//
	// The level (ad group or criterion) at which the bid was set. This is applicable
	// only at the criteria level.
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	CpcBidSource *BidSource `xml:"cpcBidSource,omitempty"`
}

type CpmBid struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 CpmBid"`

	*Bids

	//
	// Max CPM (cost per thousand impressions) bid.
	// <span class="constraint Selectable">This field can be selected using the value "CpmBid".</span><span class="constraint Filterable">This field can be filtered on.</span>
	//
	Bid *Money `xml:"bid,omitempty"`

	//
	// The level (ad group or criterion) at which the bid was set. This is applicable
	// only at the criteria level.
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	CpmBidSource *BidSource `xml:"cpmBidSource,omitempty"`
}

type CustomParameter struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 CustomParameter"`

	//
	// The parameter key to be mapped.
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	// <span class="constraint StringLength">The length of this string should be between 1 and 16, inclusive, in UTF-8 bytes, (trimmed).</span>
	//
	Key string `xml:"key,omitempty"`

	//
	// The value this parameter should be mapped to. It should be null if isRemove is true.
	// <span class="constraint StringLength">The length of this string should be between 0 and 200, inclusive, in UTF-8 bytes, (trimmed).</span>
	//
	Value string `xml:"value,omitempty"`

	//
	// On SET operation, indicates that the parameter should be removed from the existing parameters.
	// If set to true, the value field must be null.
	//
	IsRemove bool `xml:"isRemove,omitempty"`
}

type CustomParameters struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 CustomParameters"`

	//
	// The list of custom parameters.
	//
	// <p>On update, all parameters can be cleared by providing an empty or null list and setting
	// doReplace to true.
	//
	Parameters []*CustomParameter `xml:"parameters,omitempty"`

	//
	// On SET operation, indicates that the current parameters should be cleared and replaced
	// with these parameters.
	//
	DoReplace bool `xml:"doReplace,omitempty"`
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

type DistinctError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 DistinctError"`

	*ApiError

	//
	// The error reason represented by an enum.
	//
	Reason *DistinctErrorReason `xml:"reason,omitempty"`
}

type DoubleValue struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 DoubleValue"`

	*NumberValue

	//
	// the underlying double value.
	//
	Number float64 `xml:"number,omitempty"`
}

type EnhancedCpcBiddingScheme struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 EnhancedCpcBiddingScheme"`

	*BiddingScheme
}

type EntityAccessDenied struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 EntityAccessDenied"`

	*ApiError

	//
	// Reason for this error.
	//
	Reason *EntityAccessDeniedReason `xml:"reason,omitempty"`
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

type ExplorerAutoOptimizerSetting struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 ExplorerAutoOptimizerSetting"`

	*Setting

	//
	// <span class="constraint CampaignType">This field may only be set to true for campaign channel type SHOPPING with campaign channel subtype SHOPPING_UNIVERSAL_ADS.</span>
	//
	OptIn bool `xml:"optIn,omitempty"`
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

type ForwardCompatibilityError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 ForwardCompatibilityError"`

	*ApiError

	//
	// The error reason represented by an enum.
	//
	Reason *ForwardCompatibilityErrorReason `xml:"reason,omitempty"`
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

type Label struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 Label"`

	//
	// Id of label.
	// <span class="constraint Required">This field is required and should not be {@code null} when it is contained within {@link Operator}s : SET, REMOVE.</span>
	//
	Id int64 `xml:"id,omitempty"`

	//
	// Name of label.
	// <span class="constraint StringLength">The length of this string should be between 1 and 80, inclusive.</span>
	//
	Name string `xml:"name,omitempty"`

	//
	// Status of the label.
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	Status *LabelStatus `xml:"status,omitempty"`

	//
	// Attributes of the label.
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API for the following {@link Operator}s: REMOVE.</span>
	//
	Attribute *LabelAttribute `xml:"attribute,omitempty"`

	//
	// Indicates that this instance is a subtype of Label.
	// Although this field is returned in the response, it is ignored on input
	// and cannot be selected. Specify xsi:type instead.
	//
	LabelType string `xml:"Label.Type,omitempty"`
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

type LongValue struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 LongValue"`

	*NumberValue

	//
	// the underlying long value.
	//
	Number int64 `xml:"number,omitempty"`
}

type ManualCpcBiddingScheme struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 ManualCpcBiddingScheme"`

	*BiddingScheme

	//
	// The enhanced CPC bidding option for the campaign, which enables
	// bids to be enhanced based on conversion optimizer data. For more
	// information about enhanced CPC, see the
	// <a href="//support.google.com/adwords/answer/2464964"
	// >AdWords Help Center</a>.
	// <span class="constraint Selectable">This field can be selected using the value "EnhancedCpcEnabled".</span><span class="constraint Filterable">This field can be filtered on.</span>
	//
	EnhancedCpcEnabled bool `xml:"enhancedCpcEnabled,omitempty"`
}

type ManualCpmBiddingScheme struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 ManualCpmBiddingScheme"`

	*BiddingScheme

	//
	// This read-only field denotes whether Viewable CPM is enabled, and is computed based on the
	// advertising channel type and subtype. Null unless the bidding strategy type is CPM. Only
	// selectable in CampaignService, using the value ViewableCpmEnabled.
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	ViewableCpmEnabled bool `xml:"viewableCpmEnabled,omitempty"`
}

type MaximizeConversionValueBiddingScheme struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 MaximizeConversionValueBiddingScheme"`

	*BiddingScheme

	//
	// The target return on ad spend (ROAS). This is optional. If set, the bid strategy will
	// maximize revenue while averaging the target return on ad spend. If the target ROAS is high,
	// the bid strategy may not be able to spend the full budget. If the target ROAS is not set, the
	// bid strategy will aim to achieve the highest possible ROAS for the budget.
	// <span class="constraint InRange">This field must be between 0.0 and 1.7976931348623157E308, inclusive.</span>
	//
	TargetRoas float64 `xml:"targetRoas,omitempty"`
}

type MaximizeConversionsBiddingScheme struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 MaximizeConversionsBiddingScheme"`

	*BiddingScheme
}

type Money struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 Money"`

	*ComparableValue

	//
	// Amount in micros. One million is equivalent to one unit.
	//
	MicroAmount int64 `xml:"microAmount,omitempty"`
}

type MultiplierError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 MultiplierError"`

	*ApiError

	Reason *MultiplierErrorReason `xml:"reason,omitempty"`
}

type NewEntityCreationError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 NewEntityCreationError"`

	*ApiError

	//
	// The error reason represented by an enum.
	//
	Reason *NewEntityCreationErrorReason `xml:"reason,omitempty"`
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

type NumberValue struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 NumberValue"`

	*ComparableValue
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

type PageOnePromotedBiddingScheme struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 PageOnePromotedBiddingScheme"`

	*BiddingScheme

	//
	// Specifies the strategy goal: where impressions are desired to be shown on
	// search result pages.
	//
	StrategyGoal *PageOnePromotedBiddingSchemeStrategyGoal `xml:"strategyGoal,omitempty"`

	//
	// Strategy maximum bid limit in advertiser local currency micro units.
	// This upper limit applies to all keywords managed by the strategy.
	// <span class="constraint InRange">This field must be greater than or equal to 0.</span>
	//
	BidCeiling *Money `xml:"bidCeiling,omitempty"`

	//
	// Bid Multiplier to be applied to the relevant bid estimate (depending on the strategyGoal)
	// in determining a keyword's new max cpc bid.
	// <span class="constraint InRange">This field must be greater than or equal to 0.</span>
	//
	BidModifier float64 `xml:"bidModifier,omitempty"`

	//
	// Controls whether the strategy always follows bid estimate changes, or only
	// increases. If false, always set a keyword's new bid to the current bid estimate.
	// If true, only updates a keyword's bid if the current bid estimate is
	// greater than the current bid.
	//
	BidChangesForRaisesOnly bool `xml:"bidChangesForRaisesOnly,omitempty"`

	//
	// Controls whether the strategy is allowed to raise bids when the throttling rate
	// of the budget it is serving out of rises above a threshold.
	//
	RaiseBidWhenBudgetConstrained bool `xml:"raiseBidWhenBudgetConstrained,omitempty"`

	//
	// Controls whether the strategy is allowed to raise bids on keywords with lower-range
	// quality scores.
	//
	RaiseBidWhenLowQualityScore bool `xml:"raiseBidWhenLowQualityScore,omitempty"`
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

type QueryError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 QueryError"`

	*ApiError

	Reason *QueryErrorReason `xml:"reason,omitempty"`

	Message string `xml:"message,omitempty"`
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

type Setting struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 Setting"`

	//
	// Indicates that this instance is a subtype of Setting.
	// Although this field is returned in the response, it is ignored on input
	// and cannot be selected. Specify xsi:type instead.
	//
	SettingType string `xml:"Setting.Type,omitempty"`
}

type SettingError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 SettingError"`

	*ApiError

	//
	// The setting error reason.
	//
	Reason *SettingErrorReason `xml:"reason,omitempty"`
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

type StatsQueryError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 StatsQueryError"`

	*ApiError

	//
	// The error reason represented by an enum.
	//
	Reason *StatsQueryErrorReason `xml:"reason,omitempty"`
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

type String_StringMapEntry struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 String_StringMapEntry"`

	Key string `xml:"key,omitempty"`

	Value string `xml:"value,omitempty"`
}

type TargetCpaBiddingScheme struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 TargetCpaBiddingScheme"`

	*BiddingScheme

	//
	// Average cost per acquisition (CPA) target. This target should be greater than or equal to
	// minimum billable unit based on the currency for the account.
	// <span class="constraint Selectable">This field can be selected using the value "TargetCpa".</span><span class="constraint Filterable">This field can be filtered on.</span>
	//
	TargetCpa *Money `xml:"targetCpa,omitempty"`

	//
	// Maximum cpc bid limit that applies to all keywords managed by the strategy.
	// <span class="constraint InRange">This field must be greater than or equal to 0.</span>
	//
	MaxCpcBidCeiling *Money `xml:"maxCpcBidCeiling,omitempty"`

	//
	// Minimum cpc bid limit that applies to all keywords managed by the strategy.
	// <span class="constraint InRange">This field must be greater than or equal to 0.</span>
	//
	MaxCpcBidFloor *Money `xml:"maxCpcBidFloor,omitempty"`
}

type TargetOutrankShareBiddingScheme struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 TargetOutrankShareBiddingScheme"`

	*BiddingScheme

	//
	// Specifies the target fraction (in micros) of auctions where the advertiser should outrank the
	// competitor. The advertiser outranks the competitor in an auction if either the advertiser
	// appears above the competitor in the search results, or appears in the search results when the
	// competitor does not.
	// <span class="constraint InRange">This field must be between 1 and 1000000, inclusive.</span>
	//
	TargetOutrankShare int32 `xml:"targetOutrankShare,omitempty"`

	//
	// Competitor's visible domain URL.
	//
	CompetitorDomain string `xml:"competitorDomain,omitempty"`

	//
	// Ceiling on max CPC bids.
	//
	MaxCpcBidCeiling *Money `xml:"maxCpcBidCeiling,omitempty"`

	//
	// Controls whether the strategy always follows bid estimate changes, or only increases. If false,
	// always sets a keyword's new bid to the estimate that will meet the target. If true, only
	// updates a keyword's bid if the current bid estimate is greater than the current bid.
	//
	BidChangesForRaisesOnly bool `xml:"bidChangesForRaisesOnly,omitempty"`

	//
	// Controls whether the strategy is allowed to raise bids on keywords with lower-range quality
	// scores.
	//
	RaiseBidWhenLowQualityScore bool `xml:"raiseBidWhenLowQualityScore,omitempty"`
}

type TargetingSettingDetail struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 TargetingSettingDetail"`

	//
	// The criterion type group that these settings apply to.
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	//
	CriterionTypeGroup *CriterionTypeGroup `xml:"criterionTypeGroup,omitempty"`

	//
	// If true, criteria of this type can be used to modify bidding but will not restrict targeting
	// of ads. This is equivalent to "Bid only" in the AdWords user interface.
	// If false, restricts your ads to showing only for the criteria you have selected for this
	// CriterionTypeGroup. This is equivalent to "Target and Bid" in the AdWords user interface.
	// The default setting for a CriterionTypeGroup is false ("Target and Bid").
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	//
	TargetAll bool `xml:"targetAll,omitempty"`
}

type TargetRoasBiddingScheme struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 TargetRoasBiddingScheme"`

	*BiddingScheme

	//
	// The target return on average spend (ROAS).
	// <span class="constraint InRange">This field must be between 0.01 and 1000.0, inclusive.</span>
	//
	TargetRoas float64 `xml:"targetRoas,omitempty"`

	//
	// Maximum bid limit that applies to all keywords managed by the strategy.
	// <span class="constraint InRange">This field must be greater than or equal to 0.</span>
	//
	BidCeiling *Money `xml:"bidCeiling,omitempty"`

	//
	// Minimum bid limit that applies to all keywords managed by the strategy.
	// <span class="constraint InRange">This field must be greater than or equal to 0.</span>
	//
	BidFloor *Money `xml:"bidFloor,omitempty"`
}

type TargetSpendBiddingScheme struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 TargetSpendBiddingScheme"`

	*BiddingScheme

	//
	// The largest max CPC bid that can be set by the TargetSpend bidder.
	// <span class="constraint InRange">This field must be greater than or equal to 0.</span>
	//
	BidCeiling *Money `xml:"bidCeiling,omitempty"`

	//
	// A spend target under which to maximize clicks. The TargetSpend bidder will
	// attempt to spend the smaller of this value or the natural throttling spend
	// amount. If not specified, the budget is used as the spend target.
	// <span class="constraint InRange">This field must be greater than or equal to 0.</span>
	//
	SpendTarget *Money `xml:"spendTarget,omitempty"`
}

type TargetingSetting struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 TargetingSetting"`

	*Setting

	//
	// The list of per-criterion-type-group targeting settings.
	//
	Details []*TargetingSettingDetail `xml:"details,omitempty"`
}

type UrlError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 UrlError"`

	*ApiError

	//
	// The error reason represented by an enum.
	//
	Reason *UrlErrorReason `xml:"reason,omitempty"`
}

type AdGroupServiceInterface struct {
	client *SOAPClient
}

func NewAdGroupServiceInterface(url string, tls bool, auth *BasicAuth) *AdGroupServiceInterface {
	if url == "" {
		url = ""
	}
	client := NewSOAPClient(url, tls, auth)

	return &AdGroupServiceInterface{
		client: client,
	}
}

func NewAdGroupServiceInterfaceWithTLSConfig(url string, tlsCfg *tls.Config, auth *BasicAuth) *AdGroupServiceInterface {
	if url == "" {
		url = ""
	}
	client := NewSOAPClientWithTLSConfig(url, tlsCfg, auth)

	return &AdGroupServiceInterface{
		client: client,
	}
}

func (service *AdGroupServiceInterface) AddHeader(header interface{}) {
	service.client.AddHeader(header)
}

// Backwards-compatible function: use AddHeader instead
func (service *AdGroupServiceInterface) SetHeader(header interface{}) {
	service.client.AddHeader(header)
}

// Error can be either of the following types:
//
//   - ApiException
/*
   Returns a list of all the ad groups specified by the selector
   from the target customer's account.

   @param serviceSelector The selector specifying the {@link AdGroup}s to return.
   @return List of adgroups identified by the selector.
   @throws ApiException when there is at least one error with the request.
*/
func (service *AdGroupServiceInterface) Get(request *Get) (*GetResponse, error) {
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
   Adds, updates, or removes ad groups.
   <p class="note"><b>Note:</b> {@link AdGroupOperation} does not support the
   {@code REMOVE} operator. To remove an ad group, set its
   {@link AdGroup#status status} to {@code REMOVED}.</p>

   @param operations List of unique operations. The same ad group cannot be
   specified in more than one operation.
   @return The updated adgroups.
*/
func (service *AdGroupServiceInterface) Mutate(request *Mutate) (*MutateResponse, error) {
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
   Adds labels to the {@linkplain AdGroup ad group} or removes {@linkplain Label label}s from the
   {@linkplain AdGroup ad group}.
   <p>{@code ADD} -- Apply an existing label to an existing {@linkplain AdGroup ad group}.
   The {@code adGroupId} must reference an existing {@linkplain AdGroup ad group}. The
   {@code labelId} must reference an existing {@linkplain Label label}.
   <p>{@code REMOVE} -- Removes the link between the specified {@linkplain AdGroup ad group}
   and a {@linkplain Label label}.</p>

   @param operations the operations to apply.
   @return a list of {@linkplain AdGroupLabel}s where each entry in the list is the result of
   applying the operation in the input list with the same index. For an
   add operation, the returned AdGroupLabel contains the AdGroupId and the LabelId.
   In the case of a remove operation, the returned AdGroupLabel will only have AdGroupId.
   @throws ApiException when there are one or more errors with the request.
*/
func (service *AdGroupServiceInterface) MutateLabel(request *MutateLabel) (*MutateLabelResponse, error) {
	response := new(MutateLabelResponse)
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
   Returns the list of ad groups that match the query.

   @param query The SQL-like AWQL query string
   @return A list of adgroups
   @throws ApiException
*/
func (service *AdGroupServiceInterface) Query(request *Query) (*QueryResponse, error) {
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
