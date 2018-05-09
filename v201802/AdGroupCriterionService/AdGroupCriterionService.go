package AdGroupCriterionService

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
// The reasons for the target error.
//
type AdGroupCriterionErrorReason string

const (

	//
	// No link found between the AdGroupCriterion and the label.
	//
	AdGroupCriterionErrorReasonAD_GROUP_CRITERION_LABEL_DOES_NOT_EXIST AdGroupCriterionErrorReason = "AD_GROUP_CRITERION_LABEL_DOES_NOT_EXIST"

	//
	// The label has already been attached to the AdGroupCriterion.
	//
	AdGroupCriterionErrorReasonAD_GROUP_CRITERION_LABEL_ALREADY_EXISTS AdGroupCriterionErrorReason = "AD_GROUP_CRITERION_LABEL_ALREADY_EXISTS"

	//
	// Negative AdGroupCriterion cannot have labels.
	//
	AdGroupCriterionErrorReasonCANNOT_ADD_LABEL_TO_NEGATIVE_CRITERION AdGroupCriterionErrorReason = "CANNOT_ADD_LABEL_TO_NEGATIVE_CRITERION"

	//
	// Too many operations for a single call.
	//
	AdGroupCriterionErrorReasonTOO_MANY_OPERATIONS AdGroupCriterionErrorReason = "TOO_MANY_OPERATIONS"

	//
	// Negative ad group criteria are not updateable.
	//
	AdGroupCriterionErrorReasonCANT_UPDATE_NEGATIVE AdGroupCriterionErrorReason = "CANT_UPDATE_NEGATIVE"

	//
	// Concrete type of criterion (keyword v.s. placement) is required for
	// ADD and SET operations.
	//
	AdGroupCriterionErrorReasonCONCRETE_TYPE_REQUIRED AdGroupCriterionErrorReason = "CONCRETE_TYPE_REQUIRED"

	//
	// Bid is incompatible with ad group's bidding settings.
	//
	AdGroupCriterionErrorReasonBID_INCOMPATIBLE_WITH_ADGROUP AdGroupCriterionErrorReason = "BID_INCOMPATIBLE_WITH_ADGROUP"

	//
	// Cannot target and exclude the same criterion at once.
	//
	AdGroupCriterionErrorReasonCANNOT_TARGET_AND_EXCLUDE AdGroupCriterionErrorReason = "CANNOT_TARGET_AND_EXCLUDE"

	//
	// The URL of a placement is invalid.
	//
	AdGroupCriterionErrorReasonILLEGAL_URL AdGroupCriterionErrorReason = "ILLEGAL_URL"

	//
	// Keyword text was invalid.
	//
	AdGroupCriterionErrorReasonINVALID_KEYWORD_TEXT AdGroupCriterionErrorReason = "INVALID_KEYWORD_TEXT"

	//
	// Destination URL was invalid.
	//
	AdGroupCriterionErrorReasonINVALID_DESTINATION_URL AdGroupCriterionErrorReason = "INVALID_DESTINATION_URL"

	//
	// The destination url must contain at least one tag (e.g. {lpurl})
	//
	AdGroupCriterionErrorReasonMISSING_DESTINATION_URL_TAG AdGroupCriterionErrorReason = "MISSING_DESTINATION_URL_TAG"

	//
	// Keyword-level cpm bid is not supported
	//
	AdGroupCriterionErrorReasonKEYWORD_LEVEL_BID_NOT_SUPPORTED_FOR_MANUALCPM AdGroupCriterionErrorReason = "KEYWORD_LEVEL_BID_NOT_SUPPORTED_FOR_MANUALCPM"

	//
	// For example, cannot add a biddable ad group criterion that had been removed.
	//
	AdGroupCriterionErrorReasonINVALID_USER_STATUS AdGroupCriterionErrorReason = "INVALID_USER_STATUS"

	//
	// Criteria type cannot be targeted for the ad group. Either the account
	// is restricted to keywords only, the criteria type is incompatible
	// with the campaign's bidding strategy, or the criteria type can only
	// be applied to campaigns.
	//
	AdGroupCriterionErrorReasonCANNOT_ADD_CRITERIA_TYPE AdGroupCriterionErrorReason = "CANNOT_ADD_CRITERIA_TYPE"

	//
	// Criteria type cannot be excluded for the ad group. Refer to the
	// documentation for a specific criterion to check if it is excludable.
	//
	AdGroupCriterionErrorReasonCANNOT_EXCLUDE_CRITERIA_TYPE AdGroupCriterionErrorReason = "CANNOT_EXCLUDE_CRITERIA_TYPE"

	//
	// Ad group is invalid due to the product partitions it contains.
	//
	AdGroupCriterionErrorReasonINVALID_PRODUCT_PARTITION_HIERARCHY AdGroupCriterionErrorReason = "INVALID_PRODUCT_PARTITION_HIERARCHY"

	//
	// Product partition unit cannot have children.
	//
	AdGroupCriterionErrorReasonPRODUCT_PARTITION_UNIT_CANNOT_HAVE_CHILDREN AdGroupCriterionErrorReason = "PRODUCT_PARTITION_UNIT_CANNOT_HAVE_CHILDREN"

	//
	// Subdivided product partitions must have an "others" case.
	//
	AdGroupCriterionErrorReasonPRODUCT_PARTITION_SUBDIVISION_REQUIRES_OTHERS_CASE AdGroupCriterionErrorReason = "PRODUCT_PARTITION_SUBDIVISION_REQUIRES_OTHERS_CASE"

	//
	// Dimension type of product partition must be the same as that of its siblings.
	//
	AdGroupCriterionErrorReasonPRODUCT_PARTITION_REQUIRES_SAME_DIMENSION_TYPE_AS_SIBLINGS AdGroupCriterionErrorReason = "PRODUCT_PARTITION_REQUIRES_SAME_DIMENSION_TYPE_AS_SIBLINGS"

	//
	// Product partition cannot be added to the ad group because it already exists.
	//
	AdGroupCriterionErrorReasonPRODUCT_PARTITION_ALREADY_EXISTS AdGroupCriterionErrorReason = "PRODUCT_PARTITION_ALREADY_EXISTS"

	//
	// Product partition referenced in the operation was not found in the ad group.
	//
	AdGroupCriterionErrorReasonPRODUCT_PARTITION_DOES_NOT_EXIST AdGroupCriterionErrorReason = "PRODUCT_PARTITION_DOES_NOT_EXIST"

	//
	// Recursive removal failed because product partition subdivision is being created or modified
	// in this request.
	//
	AdGroupCriterionErrorReasonPRODUCT_PARTITION_CANNOT_BE_REMOVED AdGroupCriterionErrorReason = "PRODUCT_PARTITION_CANNOT_BE_REMOVED"

	//
	// Product partition type is not allowed for specified AdGroupCriterion type.
	//
	AdGroupCriterionErrorReasonINVALID_PRODUCT_PARTITION_TYPE AdGroupCriterionErrorReason = "INVALID_PRODUCT_PARTITION_TYPE"

	//
	// Product partition in an ADD operation specifies a non temporary CriterionId.
	//
	AdGroupCriterionErrorReasonPRODUCT_PARTITION_ADD_MAY_ONLY_USE_TEMP_ID AdGroupCriterionErrorReason = "PRODUCT_PARTITION_ADD_MAY_ONLY_USE_TEMP_ID"

	//
	// Partial failure is not supported for shopping campaign mutate operations.
	//
	AdGroupCriterionErrorReasonCAMPAIGN_TYPE_NOT_COMPATIBLE_WITH_PARTIAL_FAILURE AdGroupCriterionErrorReason = "CAMPAIGN_TYPE_NOT_COMPATIBLE_WITH_PARTIAL_FAILURE"

	//
	// Operations in the mutate request changes too many shopping ad groups. Please split
	// requests for multiple shopping ad groups across multiple requests.
	//
	AdGroupCriterionErrorReasonOPERATIONS_FOR_TOO_MANY_SHOPPING_ADGROUPS AdGroupCriterionErrorReason = "OPERATIONS_FOR_TOO_MANY_SHOPPING_ADGROUPS"

	//
	// Not allowed to modify url fields of an ad group criterion if there are duplicate elements
	// for that ad group criterion in the request.
	//
	AdGroupCriterionErrorReasonCANNOT_MODIFY_URL_FIELDS_WITH_DUPLICATE_ELEMENTS AdGroupCriterionErrorReason = "CANNOT_MODIFY_URL_FIELDS_WITH_DUPLICATE_ELEMENTS"

	//
	// Cannot set url fields without also setting final urls.
	//
	AdGroupCriterionErrorReasonCANNOT_SET_WITHOUT_FINAL_URLS AdGroupCriterionErrorReason = "CANNOT_SET_WITHOUT_FINAL_URLS"

	//
	// Cannot clear final urls if final mobile urls exist.
	//
	AdGroupCriterionErrorReasonCANNOT_CLEAR_FINAL_URLS_IF_FINAL_MOBILE_URLS_EXIST AdGroupCriterionErrorReason = "CANNOT_CLEAR_FINAL_URLS_IF_FINAL_MOBILE_URLS_EXIST"

	//
	// Cannot clear final urls if final app urls exist.
	//
	AdGroupCriterionErrorReasonCANNOT_CLEAR_FINAL_URLS_IF_FINAL_APP_URLS_EXIST AdGroupCriterionErrorReason = "CANNOT_CLEAR_FINAL_URLS_IF_FINAL_APP_URLS_EXIST"

	//
	// Cannot clear final urls if tracking url template exists.
	//
	AdGroupCriterionErrorReasonCANNOT_CLEAR_FINAL_URLS_IF_TRACKING_URL_TEMPLATE_EXISTS AdGroupCriterionErrorReason = "CANNOT_CLEAR_FINAL_URLS_IF_TRACKING_URL_TEMPLATE_EXISTS"

	//
	// Cannot clear final urls if url custom parameters exist.
	//
	AdGroupCriterionErrorReasonCANNOT_CLEAR_FINAL_URLS_IF_URL_CUSTOM_PARAMETERS_EXIST AdGroupCriterionErrorReason = "CANNOT_CLEAR_FINAL_URLS_IF_URL_CUSTOM_PARAMETERS_EXIST"

	//
	// Cannot set both destination url and final urls.
	//
	AdGroupCriterionErrorReasonCANNOT_SET_BOTH_DESTINATION_URL_AND_FINAL_URLS AdGroupCriterionErrorReason = "CANNOT_SET_BOTH_DESTINATION_URL_AND_FINAL_URLS"

	//
	// Cannot set both destination url and tracking url template.
	//
	AdGroupCriterionErrorReasonCANNOT_SET_BOTH_DESTINATION_URL_AND_TRACKING_URL_TEMPLATE AdGroupCriterionErrorReason = "CANNOT_SET_BOTH_DESTINATION_URL_AND_TRACKING_URL_TEMPLATE"

	//
	// Final urls are not supported for this criterion type.
	//
	AdGroupCriterionErrorReasonFINAL_URLS_NOT_SUPPORTED_FOR_CRITERION_TYPE AdGroupCriterionErrorReason = "FINAL_URLS_NOT_SUPPORTED_FOR_CRITERION_TYPE"

	//
	// Final mobile urls are not supported for this criterion type.
	//
	AdGroupCriterionErrorReasonFINAL_MOBILE_URLS_NOT_SUPPORTED_FOR_CRITERION_TYPE AdGroupCriterionErrorReason = "FINAL_MOBILE_URLS_NOT_SUPPORTED_FOR_CRITERION_TYPE"

	AdGroupCriterionErrorReasonUNKNOWN AdGroupCriterionErrorReason = "UNKNOWN"
)

//
// The entity type that exceeded the limit.
//
type AdGroupCriterionLimitExceededCriteriaLimitType string

const (
	AdGroupCriterionLimitExceededCriteriaLimitTypeADGROUP_KEYWORD AdGroupCriterionLimitExceededCriteriaLimitType = "ADGROUP_KEYWORD"

	AdGroupCriterionLimitExceededCriteriaLimitTypeADGROUP_WEBSITE AdGroupCriterionLimitExceededCriteriaLimitType = "ADGROUP_WEBSITE"

	AdGroupCriterionLimitExceededCriteriaLimitTypeADGROUP_CRITERION AdGroupCriterionLimitExceededCriteriaLimitType = "ADGROUP_CRITERION"

	AdGroupCriterionLimitExceededCriteriaLimitTypeUNKNOWN AdGroupCriterionLimitExceededCriteriaLimitType = "UNKNOWN"
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

type AgeRangeAgeRangeType string

const (
	AgeRangeAgeRangeTypeAGE_RANGE_18_24 AgeRangeAgeRangeType = "AGE_RANGE_18_24"

	AgeRangeAgeRangeTypeAGE_RANGE_25_34 AgeRangeAgeRangeType = "AGE_RANGE_25_34"

	AgeRangeAgeRangeTypeAGE_RANGE_35_44 AgeRangeAgeRangeType = "AGE_RANGE_35_44"

	AgeRangeAgeRangeTypeAGE_RANGE_45_54 AgeRangeAgeRangeType = "AGE_RANGE_45_54"

	AgeRangeAgeRangeTypeAGE_RANGE_55_64 AgeRangeAgeRangeType = "AGE_RANGE_55_64"

	AgeRangeAgeRangeTypeAGE_RANGE_65_UP AgeRangeAgeRangeType = "AGE_RANGE_65_UP"

	AgeRangeAgeRangeTypeAGE_RANGE_UNDETERMINED AgeRangeAgeRangeType = "AGE_RANGE_UNDETERMINED"

	//
	// <span class="constraint Rejected">Used for return value only. An enumeration could not be processed, typically due to incompatibility with your WSDL version.</span>
	//
	AgeRangeAgeRangeTypeUNKNOWN AgeRangeAgeRangeType = "UNKNOWN"
)

//
// The possible types of App Payment Model.
//
type AppPaymentModelAppPaymentModelType string

const (

	//
	// Represents paid-for apps.
	//
	AppPaymentModelAppPaymentModelTypeAPP_PAYMENT_MODEL_PAID AppPaymentModelAppPaymentModelType = "APP_PAYMENT_MODEL_PAID"

	//
	// <span class="constraint Rejected">Used for return value only. An enumeration could not be processed, typically due to incompatibility with your WSDL version.</span>
	//
	AppPaymentModelAppPaymentModelTypeUNKNOWN AppPaymentModelAppPaymentModelType = "UNKNOWN"
)

//
// The possible os types for an AppUrl
//
type AppUrlOsType string

const (

	//
	// The Apple IOS operating system,
	//
	AppUrlOsTypeOS_TYPE_IOS AppUrlOsType = "OS_TYPE_IOS"

	//
	// The Android operating system.
	//
	AppUrlOsTypeOS_TYPE_ANDROID AppUrlOsType = "OS_TYPE_ANDROID"

	AppUrlOsTypeUNKNOWN AppUrlOsType = "UNKNOWN"
)

//
// Approval status for the criterion.
// Note: there are more states involved but we only expose two to users.
//
type ApprovalStatus string

const (

	//
	// Criterion with no reportable policy problems.
	//
	ApprovalStatusAPPROVED ApprovalStatus = "APPROVED"

	//
	// Criterion that is yet to be reviewed.
	//
	ApprovalStatusPENDING_REVIEW ApprovalStatus = "PENDING_REVIEW"

	//
	// Criterion that is under review.
	//
	ApprovalStatusUNDER_REVIEW ApprovalStatus = "UNDER_REVIEW"

	//
	// Criterion disapproved due to policy violation.
	//
	ApprovalStatusDISAPPROVED ApprovalStatus = "DISAPPROVED"
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
// The reasons for the target error.
//
type CollectionSizeErrorReason string

const (
	CollectionSizeErrorReasonTOO_FEW CollectionSizeErrorReason = "TOO_FEW"

	CollectionSizeErrorReasonTOO_MANY CollectionSizeErrorReason = "TOO_MANY"
)

//
// The types of criteria.
//
type CriterionType string

const (

	//
	// Content label for exclusion.
	//
	CriterionTypeCONTENT_LABEL CriterionType = "CONTENT_LABEL"

	//
	// Keyword. e.g. 'mars cruise'
	//
	CriterionTypeKEYWORD CriterionType = "KEYWORD"

	//
	// Placement. aka Website. e.g. 'www.flowers4sale.com'
	//
	CriterionTypePLACEMENT CriterionType = "PLACEMENT"

	//
	// Vertical, e.g. 'category::Animals>Pets'  This is for vertical targeting on the content
	// network.
	//
	CriterionTypeVERTICAL CriterionType = "VERTICAL"

	//
	// User lists, are links to sets of users defined by the advertiser.
	//
	CriterionTypeUSER_LIST CriterionType = "USER_LIST"

	//
	// User interests, categories of interests the user is interested in.
	//
	CriterionTypeUSER_INTEREST CriterionType = "USER_INTEREST"

	//
	// Mobile applications to target.
	//
	CriterionTypeMOBILE_APPLICATION CriterionType = "MOBILE_APPLICATION"

	//
	// Mobile application categories to target.
	//
	CriterionTypeMOBILE_APP_CATEGORY CriterionType = "MOBILE_APP_CATEGORY"

	//
	// Product partition (product group) in a shopping campaign.
	//
	CriterionTypePRODUCT_PARTITION CriterionType = "PRODUCT_PARTITION"

	//
	// IP addresses to exclude.
	//
	CriterionTypeIP_BLOCK CriterionType = "IP_BLOCK"

	//
	// Webpages of an advertiser's website to target.
	//
	CriterionTypeWEBPAGE CriterionType = "WEBPAGE"

	//
	// Languages to target.
	//
	CriterionTypeLANGUAGE CriterionType = "LANGUAGE"

	//
	// Geographic regions to target.
	//
	CriterionTypeLOCATION CriterionType = "LOCATION"

	//
	// Age Range to exclude.
	//
	CriterionTypeAGE_RANGE CriterionType = "AGE_RANGE"

	//
	// Mobile carriers to target.
	//
	CriterionTypeCARRIER CriterionType = "CARRIER"

	//
	// Mobile operating system versions to target.
	//
	CriterionTypeOPERATING_SYSTEM_VERSION CriterionType = "OPERATING_SYSTEM_VERSION"

	//
	// Mobile devices to target.
	//
	CriterionTypeMOBILE_DEVICE CriterionType = "MOBILE_DEVICE"

	//
	// Gender to exclude.
	//
	CriterionTypeGENDER CriterionType = "GENDER"

	//
	// Parent to target and exclude.
	//
	CriterionTypePARENT CriterionType = "PARENT"

	//
	// Proximity (area within a radius) to target.
	//
	CriterionTypePROXIMITY CriterionType = "PROXIMITY"

	//
	// Platforms to target.
	//
	CriterionTypePLATFORM CriterionType = "PLATFORM"

	//
	// Representing preferred content bid modifier.
	//
	CriterionTypePREFERRED_CONTENT CriterionType = "PREFERRED_CONTENT"

	//
	// AdSchedule or specific days and time intervals to target.
	//
	CriterionTypeAD_SCHEDULE CriterionType = "AD_SCHEDULE"

	//
	// Targeting based on location groups.
	//
	CriterionTypeLOCATION_GROUPS CriterionType = "LOCATION_GROUPS"

	//
	// Scope of products. Contains a list of product dimensions, all of which a product has to match
	// to be included in the campaign.
	//
	CriterionTypePRODUCT_SCOPE CriterionType = "PRODUCT_SCOPE"

	//
	// YouTube video to target.
	//
	CriterionTypeYOUTUBE_VIDEO CriterionType = "YOUTUBE_VIDEO"

	//
	// YouTube channel to target.
	//
	CriterionTypeYOUTUBE_CHANNEL CriterionType = "YOUTUBE_CHANNEL"

	//
	// Enables advertisers to target paid apps.
	//
	CriterionTypeAPP_PAYMENT_MODEL CriterionType = "APP_PAYMENT_MODEL"

	//
	// Income range to target and exclude.
	//
	CriterionTypeINCOME_RANGE CriterionType = "INCOME_RANGE"

	//
	// Interaction type to bid modify.
	//
	CriterionTypeINTERACTION_TYPE CriterionType = "INTERACTION_TYPE"

	//
	// <span class="constraint Rejected">Used for return value only. An enumeration could not be processed, typically due to incompatibility with your WSDL version.</span>
	//
	CriterionTypeUNKNOWN CriterionType = "UNKNOWN"
)

type CriterionErrorReason string

const (

	//
	// Concrete type of criterion is required for ADD and SET operations.
	//
	CriterionErrorReasonCONCRETE_TYPE_REQUIRED CriterionErrorReason = "CONCRETE_TYPE_REQUIRED"

	//
	// The category requested for exclusion is invalid.
	//
	CriterionErrorReasonINVALID_EXCLUDED_CATEGORY CriterionErrorReason = "INVALID_EXCLUDED_CATEGORY"

	//
	// Invalid keyword criteria text.
	//
	CriterionErrorReasonINVALID_KEYWORD_TEXT CriterionErrorReason = "INVALID_KEYWORD_TEXT"

	//
	// Keyword text should be less than 80 chars.
	//
	CriterionErrorReasonKEYWORD_TEXT_TOO_LONG CriterionErrorReason = "KEYWORD_TEXT_TOO_LONG"

	//
	// Keyword text has too many words.
	//
	CriterionErrorReasonKEYWORD_HAS_TOO_MANY_WORDS CriterionErrorReason = "KEYWORD_HAS_TOO_MANY_WORDS"

	//
	// Keyword text has invalid characters or symbols.
	//
	CriterionErrorReasonKEYWORD_HAS_INVALID_CHARS CriterionErrorReason = "KEYWORD_HAS_INVALID_CHARS"

	//
	// Invalid placement URL.
	//
	CriterionErrorReasonINVALID_PLACEMENT_URL CriterionErrorReason = "INVALID_PLACEMENT_URL"

	//
	// Invalid user list criterion.
	//
	CriterionErrorReasonINVALID_USER_LIST CriterionErrorReason = "INVALID_USER_LIST"

	//
	// Invalid user interest criterion.
	//
	CriterionErrorReasonINVALID_USER_INTEREST CriterionErrorReason = "INVALID_USER_INTEREST"

	//
	// Placement URL has wrong format.
	//
	CriterionErrorReasonINVALID_FORMAT_FOR_PLACEMENT_URL CriterionErrorReason = "INVALID_FORMAT_FOR_PLACEMENT_URL"

	//
	// Placement URL is too long.
	//
	CriterionErrorReasonPLACEMENT_URL_IS_TOO_LONG CriterionErrorReason = "PLACEMENT_URL_IS_TOO_LONG"

	//
	// Indicates the URL contains an illegal character.
	//
	CriterionErrorReasonPLACEMENT_URL_HAS_ILLEGAL_CHAR CriterionErrorReason = "PLACEMENT_URL_HAS_ILLEGAL_CHAR"

	//
	// Indicates the URL contains multiple comma separated URLs.
	//
	CriterionErrorReasonPLACEMENT_URL_HAS_MULTIPLE_SITES_IN_LINE CriterionErrorReason = "PLACEMENT_URL_HAS_MULTIPLE_SITES_IN_LINE"

	//
	// Indicates the domain is blacklisted.
	//
	CriterionErrorReasonPLACEMENT_IS_NOT_AVAILABLE_FOR_TARGETING_OR_EXCLUSION CriterionErrorReason = "PLACEMENT_IS_NOT_AVAILABLE_FOR_TARGETING_OR_EXCLUSION"

	//
	// Invalid vertical path.
	//
	CriterionErrorReasonINVALID_VERTICAL_PATH CriterionErrorReason = "INVALID_VERTICAL_PATH"

	//
	// Indicates the placement is a YouTube vertical channel, which is no longer supported.
	//
	CriterionErrorReasonYOUTUBE_VERTICAL_CHANNEL_DEPRECATED CriterionErrorReason = "YOUTUBE_VERTICAL_CHANNEL_DEPRECATED"

	//
	// Indicates the placement is a YouTube demographic channel, which is no longer supported.
	//
	CriterionErrorReasonYOUTUBE_DEMOGRAPHIC_CHANNEL_DEPRECATED CriterionErrorReason = "YOUTUBE_DEMOGRAPHIC_CHANNEL_DEPRECATED"

	//
	// YouTube urls are not supported in Placement criterion. Use YouTubeChannel and
	// YouTubeVideo criterion instead.
	//
	CriterionErrorReasonYOUTUBE_URL_UNSUPPORTED CriterionErrorReason = "YOUTUBE_URL_UNSUPPORTED"

	//
	// Criteria type can not be excluded by the customer,
	// like AOL account type cannot target site type criteria.
	//
	CriterionErrorReasonCANNOT_EXCLUDE_CRITERIA_TYPE CriterionErrorReason = "CANNOT_EXCLUDE_CRITERIA_TYPE"

	//
	// Criteria type can not be targeted.
	//
	CriterionErrorReasonCANNOT_ADD_CRITERIA_TYPE CriterionErrorReason = "CANNOT_ADD_CRITERIA_TYPE"

	//
	// Product filter in the product criteria has invalid characters.
	// Operand and the argument in the filter can not have "==" or "&+".
	//
	CriterionErrorReasonINVALID_PRODUCT_FILTER CriterionErrorReason = "INVALID_PRODUCT_FILTER"

	//
	// Product filter in the product criteria is translated to a string as
	// operand1==argument1&+operand2==argument2, maximum allowed length for
	// the string is 255 chars.
	//
	CriterionErrorReasonPRODUCT_FILTER_TOO_LONG CriterionErrorReason = "PRODUCT_FILTER_TOO_LONG"

	//
	// Not allowed to exclude similar user list.
	//
	CriterionErrorReasonCANNOT_EXCLUDE_SIMILAR_USER_LIST CriterionErrorReason = "CANNOT_EXCLUDE_SIMILAR_USER_LIST"

	//
	// Not allowed to target a closed user list.
	//
	CriterionErrorReasonCANNOT_ADD_CLOSED_USER_LIST CriterionErrorReason = "CANNOT_ADD_CLOSED_USER_LIST"

	//
	// Not allowed to add display only UserLists to search only campaigns.
	//
	CriterionErrorReasonCANNOT_ADD_DISPLAY_ONLY_LISTS_TO_SEARCH_ONLY_CAMPAIGNS CriterionErrorReason = "CANNOT_ADD_DISPLAY_ONLY_LISTS_TO_SEARCH_ONLY_CAMPAIGNS"

	//
	// Not allowed to add display only UserLists to search plus campaigns.
	//
	CriterionErrorReasonCANNOT_ADD_DISPLAY_ONLY_LISTS_TO_SEARCH_CAMPAIGNS CriterionErrorReason = "CANNOT_ADD_DISPLAY_ONLY_LISTS_TO_SEARCH_CAMPAIGNS"

	//
	// Not allowed to add display only UserLists to shopping campaigns.
	//
	CriterionErrorReasonCANNOT_ADD_DISPLAY_ONLY_LISTS_TO_SHOPPING_CAMPAIGNS CriterionErrorReason = "CANNOT_ADD_DISPLAY_ONLY_LISTS_TO_SHOPPING_CAMPAIGNS"

	//
	// Not allowed to add User interests to search only campaigns.
	//
	CriterionErrorReasonCANNOT_ADD_USER_INTERESTS_TO_SEARCH_CAMPAIGNS CriterionErrorReason = "CANNOT_ADD_USER_INTERESTS_TO_SEARCH_CAMPAIGNS"

	//
	// Not allowed to set bids for this criterion type in search campaigns
	//
	CriterionErrorReasonCANNOT_SET_BIDS_ON_CRITERION_TYPE_IN_SEARCH_CAMPAIGNS CriterionErrorReason = "CANNOT_SET_BIDS_ON_CRITERION_TYPE_IN_SEARCH_CAMPAIGNS"

	//
	// Final URLs, URL Templates and CustomParameters cannot be set for the criterion
	// types of Gender, AgeRange, UserList, Placement, MobileApp, and MobileAppCategory
	// in search campaigns and shopping campaigns.
	//
	CriterionErrorReasonCANNOT_ADD_URLS_TO_CRITERION_TYPE_FOR_CAMPAIGN_TYPE CriterionErrorReason = "CANNOT_ADD_URLS_TO_CRITERION_TYPE_FOR_CAMPAIGN_TYPE"

	//
	// IP address is not valid.
	//
	CriterionErrorReasonINVALID_IP_ADDRESS CriterionErrorReason = "INVALID_IP_ADDRESS"

	//
	// IP format is not valid.
	//
	CriterionErrorReasonINVALID_IP_FORMAT CriterionErrorReason = "INVALID_IP_FORMAT"

	//
	// Mobile application is not valid.
	//
	CriterionErrorReasonINVALID_MOBILE_APP CriterionErrorReason = "INVALID_MOBILE_APP"

	//
	// Mobile application category is not valid.
	//
	CriterionErrorReasonINVALID_MOBILE_APP_CATEGORY CriterionErrorReason = "INVALID_MOBILE_APP_CATEGORY"

	//
	// The CriterionId does not exist or is of the incorrect type.
	//
	CriterionErrorReasonINVALID_CRITERION_ID CriterionErrorReason = "INVALID_CRITERION_ID"

	//
	// The Criterion is not allowed to be targeted.
	//
	CriterionErrorReasonCANNOT_TARGET_CRITERION CriterionErrorReason = "CANNOT_TARGET_CRITERION"

	//
	// The criterion is not allowed to be targeted as it is deprecated.
	//
	CriterionErrorReasonCANNOT_TARGET_OBSOLETE_CRITERION CriterionErrorReason = "CANNOT_TARGET_OBSOLETE_CRITERION"

	//
	// The CriterionId is not valid for the type.
	//
	CriterionErrorReasonCRITERION_ID_AND_TYPE_MISMATCH CriterionErrorReason = "CRITERION_ID_AND_TYPE_MISMATCH"

	//
	// Distance for the radius for the proximity criterion is invalid.
	//
	CriterionErrorReasonINVALID_PROXIMITY_RADIUS CriterionErrorReason = "INVALID_PROXIMITY_RADIUS"

	//
	// Units for the distance for the radius for the proximity criterion is invalid.
	//
	CriterionErrorReasonINVALID_PROXIMITY_RADIUS_UNITS CriterionErrorReason = "INVALID_PROXIMITY_RADIUS_UNITS"

	//
	// Street address is too short.
	//
	CriterionErrorReasonINVALID_STREETADDRESS_LENGTH CriterionErrorReason = "INVALID_STREETADDRESS_LENGTH"

	//
	// City name in the address is too short.
	//
	CriterionErrorReasonINVALID_CITYNAME_LENGTH CriterionErrorReason = "INVALID_CITYNAME_LENGTH"

	//
	// Region code in the address is too short.
	//
	CriterionErrorReasonINVALID_REGIONCODE_LENGTH CriterionErrorReason = "INVALID_REGIONCODE_LENGTH"

	//
	// Region name in the address is not valid.
	//
	CriterionErrorReasonINVALID_REGIONNAME_LENGTH CriterionErrorReason = "INVALID_REGIONNAME_LENGTH"

	//
	// Postal code in the address is not valid.
	//
	CriterionErrorReasonINVALID_POSTALCODE_LENGTH CriterionErrorReason = "INVALID_POSTALCODE_LENGTH"

	//
	// Country code in the address is not valid.
	//
	CriterionErrorReasonINVALID_COUNTRY_CODE CriterionErrorReason = "INVALID_COUNTRY_CODE"

	//
	// Latitude for the GeoPoint is not valid.
	//
	CriterionErrorReasonINVALID_LATITUDE CriterionErrorReason = "INVALID_LATITUDE"

	//
	// Longitude for the GeoPoint is not valid.
	//
	CriterionErrorReasonINVALID_LONGITUDE CriterionErrorReason = "INVALID_LONGITUDE"

	//
	// The Proximity input is not valid. Both address and geoPoint cannot be null.
	//
	CriterionErrorReasonPROXIMITY_GEOPOINT_AND_ADDRESS_BOTH_CANNOT_BE_NULL CriterionErrorReason = "PROXIMITY_GEOPOINT_AND_ADDRESS_BOTH_CANNOT_BE_NULL"

	//
	// The Proximity address cannot be geocoded to a valid lat/long.
	//
	CriterionErrorReasonINVALID_PROXIMITY_ADDRESS CriterionErrorReason = "INVALID_PROXIMITY_ADDRESS"

	//
	// User domain name is not valid.
	//
	CriterionErrorReasonINVALID_USER_DOMAIN_NAME CriterionErrorReason = "INVALID_USER_DOMAIN_NAME"

	//
	// Length of serialized criterion parameter exceeded size limit.
	//
	CriterionErrorReasonCRITERION_PARAMETER_TOO_LONG CriterionErrorReason = "CRITERION_PARAMETER_TOO_LONG"

	//
	// Time interval in the AdSchedule overlaps with another AdSchedule.
	//
	CriterionErrorReasonAD_SCHEDULE_TIME_INTERVALS_OVERLAP CriterionErrorReason = "AD_SCHEDULE_TIME_INTERVALS_OVERLAP"

	//
	// AdSchedule time interval cannot span multiple days.
	//
	CriterionErrorReasonAD_SCHEDULE_INTERVAL_CANNOT_SPAN_MULTIPLE_DAYS CriterionErrorReason = "AD_SCHEDULE_INTERVAL_CANNOT_SPAN_MULTIPLE_DAYS"

	//
	// AdSchedule time interval specified is invalid,
	// endTime cannot be earlier than startTime.
	//
	CriterionErrorReasonAD_SCHEDULE_INVALID_TIME_INTERVAL CriterionErrorReason = "AD_SCHEDULE_INVALID_TIME_INTERVAL"

	//
	// The number of AdSchedule entries in a day exceeds the limit.
	//
	CriterionErrorReasonAD_SCHEDULE_EXCEEDED_INTERVALS_PER_DAY_LIMIT CriterionErrorReason = "AD_SCHEDULE_EXCEEDED_INTERVALS_PER_DAY_LIMIT"

	//
	// CriteriaId does not match the interval of the AdSchedule specified.
	//
	CriterionErrorReasonAD_SCHEDULE_CRITERION_ID_MISMATCHING_FIELDS CriterionErrorReason = "AD_SCHEDULE_CRITERION_ID_MISMATCHING_FIELDS"

	//
	// Cannot set bid modifier for this criterion type.
	//
	CriterionErrorReasonCANNOT_BID_MODIFY_CRITERION_TYPE CriterionErrorReason = "CANNOT_BID_MODIFY_CRITERION_TYPE"

	//
	// Cannot bid modify criterion, since it is opted out of the campaign.
	//
	CriterionErrorReasonCANNOT_BID_MODIFY_CRITERION_CAMPAIGN_OPTED_OUT CriterionErrorReason = "CANNOT_BID_MODIFY_CRITERION_CAMPAIGN_OPTED_OUT"

	//
	// Cannot set bid modifier for a negative criterion.
	//
	CriterionErrorReasonCANNOT_BID_MODIFY_NEGATIVE_CRITERION CriterionErrorReason = "CANNOT_BID_MODIFY_NEGATIVE_CRITERION"

	//
	// Bid Modifier already exists. Use SET operation to update.
	//
	CriterionErrorReasonBID_MODIFIER_ALREADY_EXISTS CriterionErrorReason = "BID_MODIFIER_ALREADY_EXISTS"

	//
	// Feed Id is not allowed in these Location Groups.
	//
	CriterionErrorReasonFEED_ID_NOT_ALLOWED CriterionErrorReason = "FEED_ID_NOT_ALLOWED"

	//
	// The account may not use the requested criteria type. For example, some
	// accounts are restricted to keywords only.
	//
	CriterionErrorReasonACCOUNT_INELIGIBLE_FOR_CRITERIA_TYPE CriterionErrorReason = "ACCOUNT_INELIGIBLE_FOR_CRITERIA_TYPE"

	//
	// The requested criteria type cannot be used with campaign or ad group bidding strategy.
	//
	CriterionErrorReasonCRITERIA_TYPE_INVALID_FOR_BIDDING_STRATEGY CriterionErrorReason = "CRITERIA_TYPE_INVALID_FOR_BIDDING_STRATEGY"

	//
	// The Criterion is not allowed to be excluded.
	//
	CriterionErrorReasonCANNOT_EXCLUDE_CRITERION CriterionErrorReason = "CANNOT_EXCLUDE_CRITERION"

	//
	// The criterion is not allowed to be removed. For example, we cannot remove any
	// of the platform criterion.
	//
	CriterionErrorReasonCANNOT_REMOVE_CRITERION CriterionErrorReason = "CANNOT_REMOVE_CRITERION"

	//
	// The combined length of product dimension values of the product scope criterion is too long.
	//
	CriterionErrorReasonPRODUCT_SCOPE_TOO_LONG CriterionErrorReason = "PRODUCT_SCOPE_TOO_LONG"

	//
	// Product scope contains too many dimensions.
	//
	CriterionErrorReasonPRODUCT_SCOPE_TOO_MANY_DIMENSIONS CriterionErrorReason = "PRODUCT_SCOPE_TOO_MANY_DIMENSIONS"

	//
	// The combined length of product dimension values of the product partition criterion is too
	// long.
	//
	CriterionErrorReasonPRODUCT_PARTITION_TOO_LONG CriterionErrorReason = "PRODUCT_PARTITION_TOO_LONG"

	//
	// Product partition contains too many dimensions.
	//
	CriterionErrorReasonPRODUCT_PARTITION_TOO_MANY_DIMENSIONS CriterionErrorReason = "PRODUCT_PARTITION_TOO_MANY_DIMENSIONS"

	//
	// The product dimension is invalid (e.g. dimension contains illegal value, dimension type is
	// represented with wrong class, etc). Product dimension value can not contain "==" or "&+".
	//
	CriterionErrorReasonINVALID_PRODUCT_DIMENSION CriterionErrorReason = "INVALID_PRODUCT_DIMENSION"

	//
	// Product dimension type is either invalid for campaigns of this type or cannot be used in the
	// current context. BIDDING_CATEGORY_Lx and PRODUCT_TYPE_Lx product dimensions must be used in
	// ascending order of their levels: L1, L2, L3, L4, L5... The levels must be specified
	// sequentially and start from L1. Furthermore, an "others" product partition cannot be
	// subdivided with a dimension of the same type but of a higher level ("others"
	// BIDDING_CATEGORY_L3 can be subdivided with BRAND but not with BIDDING_CATEGORY_L4).
	//
	CriterionErrorReasonINVALID_PRODUCT_DIMENSION_TYPE CriterionErrorReason = "INVALID_PRODUCT_DIMENSION_TYPE"

	//
	// Bidding categories do not form a valid path in the Shopping bidding category taxonomy.
	//
	CriterionErrorReasonINVALID_PRODUCT_BIDDING_CATEGORY CriterionErrorReason = "INVALID_PRODUCT_BIDDING_CATEGORY"

	//
	// ShoppingSetting must be added to the campaign before ProductScope criteria can be added.
	//
	CriterionErrorReasonMISSING_SHOPPING_SETTING CriterionErrorReason = "MISSING_SHOPPING_SETTING"

	//
	// Matching function is invalid.
	//
	CriterionErrorReasonINVALID_MATCHING_FUNCTION CriterionErrorReason = "INVALID_MATCHING_FUNCTION"

	//
	// Filter parameters not allowed for location groups targeting.
	//
	CriterionErrorReasonLOCATION_FILTER_NOT_ALLOWED CriterionErrorReason = "LOCATION_FILTER_NOT_ALLOWED"

	//
	// Given location filter parameter is invalid for location groups targeting.
	//
	CriterionErrorReasonLOCATION_FILTER_INVALID CriterionErrorReason = "LOCATION_FILTER_INVALID"

	//
	// Criteria type cannot be associated with a campaign and its ad group(s) simultaneously.
	//
	CriterionErrorReasonCANNOT_ATTACH_CRITERIA_AT_CAMPAIGN_AND_ADGROUP CriterionErrorReason = "CANNOT_ATTACH_CRITERIA_AT_CAMPAIGN_AND_ADGROUP"

	CriterionErrorReasonUNKNOWN CriterionErrorReason = "UNKNOWN"
)

//
// The way a criterion is used - biddable or negative.
//
type CriterionUse string

const (

	//
	// Biddable (positive) criterion
	//
	CriterionUseBIDDABLE CriterionUse = "BIDDABLE"

	//
	// Negative criterion
	//
	CriterionUseNEGATIVE CriterionUse = "NEGATIVE"
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

type GenderGenderType string

const (
	GenderGenderTypeGENDER_MALE GenderGenderType = "GENDER_MALE"

	GenderGenderTypeGENDER_FEMALE GenderGenderType = "GENDER_FEMALE"

	GenderGenderTypeGENDER_UNDETERMINED GenderGenderType = "GENDER_UNDETERMINED"
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
// Income percentile ranges.
//
type IncomeRangeIncomeRangeType string

const (

	//
	// Users for which income range is undetermined.
	//
	IncomeRangeIncomeRangeTypeINCOME_RANGE_UNDETERMINED IncomeRangeIncomeRangeType = "INCOME_RANGE_UNDETERMINED"

	//
	// Users in [0%, 50%) income percentile range.
	//
	IncomeRangeIncomeRangeTypeINCOME_RANGE_0_50 IncomeRangeIncomeRangeType = "INCOME_RANGE_0_50"

	//
	// Users in [50%, 60%) income percentile range.
	//
	IncomeRangeIncomeRangeTypeINCOME_RANGE_50_60 IncomeRangeIncomeRangeType = "INCOME_RANGE_50_60"

	//
	// Users in [60%, 70%) income percentile range.
	//
	IncomeRangeIncomeRangeTypeINCOME_RANGE_60_70 IncomeRangeIncomeRangeType = "INCOME_RANGE_60_70"

	//
	// Users in [70%, 80%) income percentile range.
	//
	IncomeRangeIncomeRangeTypeINCOME_RANGE_70_80 IncomeRangeIncomeRangeType = "INCOME_RANGE_70_80"

	//
	// Users in [80%, 90%) income percentile range.
	//
	IncomeRangeIncomeRangeTypeINCOME_RANGE_80_90 IncomeRangeIncomeRangeType = "INCOME_RANGE_80_90"

	//
	// Users in [90%, 100%] income percentile range.
	//
	IncomeRangeIncomeRangeTypeINCOME_RANGE_90_UP IncomeRangeIncomeRangeType = "INCOME_RANGE_90_UP"

	//
	// <span class="constraint Rejected">Used for return value only. An enumeration could not be processed, typically due to incompatibility with your WSDL version.</span>
	//
	IncomeRangeIncomeRangeTypeUNKNOWN IncomeRangeIncomeRangeType = "UNKNOWN"
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
// Match type of a keyword. i.e. the way we match a keyword string with
// search queries.
//
type KeywordMatchType string

const (

	//
	// Exact match
	//
	KeywordMatchTypeEXACT KeywordMatchType = "EXACT"

	//
	// Phrase match
	//
	KeywordMatchTypePHRASE KeywordMatchType = "PHRASE"

	//
	// Broad match
	//
	KeywordMatchTypeBROAD KeywordMatchType = "BROAD"
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
// The reasons for errors when using pagination.
//
type PagingErrorReason string

const (

	//
	// The start index value cannot be a negative number.
	//
	PagingErrorReasonSTART_INDEX_CANNOT_BE_NEGATIVE PagingErrorReason = "START_INDEX_CANNOT_BE_NEGATIVE"

	//
	// The number of results cannot be a negative number.
	//
	PagingErrorReasonNUMBER_OF_RESULTS_CANNOT_BE_NEGATIVE PagingErrorReason = "NUMBER_OF_RESULTS_CANNOT_BE_NEGATIVE"

	//
	// <span class="constraint Rejected">Used for return value only. An enumeration could not be processed, typically due to incompatibility with your WSDL version.</span>
	//
	PagingErrorReasonUNKNOWN PagingErrorReason = "UNKNOWN"
)

//
// The possible types of parents.
//
type ParentParentType string

const (
	ParentParentTypePARENT_PARENT ParentParentType = "PARENT_PARENT"

	ParentParentTypePARENT_NOT_A_PARENT ParentParentType = "PARENT_NOT_A_PARENT"

	ParentParentTypePARENT_UNDETERMINED ParentParentType = "PARENT_UNDETERMINED"

	//
	// <span class="constraint Rejected">Used for return value only. An enumeration could not be processed, typically due to incompatibility with your WSDL version.</span>
	//
	ParentParentTypeUNKNOWN ParentParentType = "UNKNOWN"
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
// A canonical product condition.
//
type ProductCanonicalConditionCondition string

const (
	ProductCanonicalConditionConditionNEW ProductCanonicalConditionCondition = "NEW"

	ProductCanonicalConditionConditionUSED ProductCanonicalConditionCondition = "USED"

	ProductCanonicalConditionConditionREFURBISHED ProductCanonicalConditionCondition = "REFURBISHED"

	//
	// <span class="constraint Rejected">Used for return value only. An enumeration could not be processed, typically due to incompatibility with your WSDL version.</span>
	//
	ProductCanonicalConditionConditionUNKNOWN ProductCanonicalConditionCondition = "UNKNOWN"
)

//
// Type of product dimension.
//
type ProductDimensionType string

const (

	//
	// <span class="constraint Rejected">Used for return value only. An enumeration could not be processed, typically due to incompatibility with your WSDL version.</span>
	//
	ProductDimensionTypeUNKNOWN ProductDimensionType = "UNKNOWN"

	ProductDimensionTypeBIDDING_CATEGORY_L1 ProductDimensionType = "BIDDING_CATEGORY_L1"

	ProductDimensionTypeBIDDING_CATEGORY_L2 ProductDimensionType = "BIDDING_CATEGORY_L2"

	ProductDimensionTypeBIDDING_CATEGORY_L3 ProductDimensionType = "BIDDING_CATEGORY_L3"

	ProductDimensionTypeBIDDING_CATEGORY_L4 ProductDimensionType = "BIDDING_CATEGORY_L4"

	ProductDimensionTypeBIDDING_CATEGORY_L5 ProductDimensionType = "BIDDING_CATEGORY_L5"

	ProductDimensionTypeBRAND ProductDimensionType = "BRAND"

	ProductDimensionTypeCANONICAL_CONDITION ProductDimensionType = "CANONICAL_CONDITION"

	ProductDimensionTypeCUSTOM_ATTRIBUTE_0 ProductDimensionType = "CUSTOM_ATTRIBUTE_0"

	ProductDimensionTypeCUSTOM_ATTRIBUTE_1 ProductDimensionType = "CUSTOM_ATTRIBUTE_1"

	ProductDimensionTypeCUSTOM_ATTRIBUTE_2 ProductDimensionType = "CUSTOM_ATTRIBUTE_2"

	ProductDimensionTypeCUSTOM_ATTRIBUTE_3 ProductDimensionType = "CUSTOM_ATTRIBUTE_3"

	ProductDimensionTypeCUSTOM_ATTRIBUTE_4 ProductDimensionType = "CUSTOM_ATTRIBUTE_4"

	ProductDimensionTypeOFFER_ID ProductDimensionType = "OFFER_ID"

	ProductDimensionTypePRODUCT_TYPE_L1 ProductDimensionType = "PRODUCT_TYPE_L1"

	ProductDimensionTypePRODUCT_TYPE_L2 ProductDimensionType = "PRODUCT_TYPE_L2"

	ProductDimensionTypePRODUCT_TYPE_L3 ProductDimensionType = "PRODUCT_TYPE_L3"

	ProductDimensionTypePRODUCT_TYPE_L4 ProductDimensionType = "PRODUCT_TYPE_L4"

	ProductDimensionTypePRODUCT_TYPE_L5 ProductDimensionType = "PRODUCT_TYPE_L5"

	ProductDimensionTypeCHANNEL ProductDimensionType = "CHANNEL"

	ProductDimensionTypeCHANNEL_EXCLUSIVITY ProductDimensionType = "CHANNEL_EXCLUSIVITY"
)

//
// Type of a product partition in a shopping campaign.
//
type ProductPartitionType string

const (

	//
	// <span class="constraint Rejected">Used for return value only. An enumeration could not be processed, typically due to incompatibility with your WSDL version.</span>
	//
	ProductPartitionTypeUNKNOWN ProductPartitionType = "UNKNOWN"

	//
	// Subdivision of products along some product dimension.
	//
	ProductPartitionTypeSUBDIVISION ProductPartitionType = "SUBDIVISION"

	//
	// Unit which either defines a bid or delegates bidding to other campaigns.
	//
	ProductPartitionTypeUNIT ProductPartitionType = "UNIT"
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
// Channel specifies where the item is sold: online or in local stores.
//
type ShoppingProductChannel string

const (

	//
	// <span class="constraint Rejected">Used for return value only. An enumeration could not be processed, typically due to incompatibility with your WSDL version.</span>
	//
	ShoppingProductChannelUNKNOWN ShoppingProductChannel = "UNKNOWN"

	//
	// The item is sold online.
	//
	ShoppingProductChannelONLINE ShoppingProductChannel = "ONLINE"

	//
	// The item is sold in local stores.
	//
	ShoppingProductChannelLOCAL ShoppingProductChannel = "LOCAL"
)

//
// Channel exclusivity specifies whether an item is sold exclusively in one channel
// or through multiple channels.
//
type ShoppingProductChannelExclusivity string

const (

	//
	// <span class="constraint Rejected">Used for return value only. An enumeration could not be processed, typically due to incompatibility with your WSDL version.</span>
	//
	ShoppingProductChannelExclusivityUNKNOWN ShoppingProductChannelExclusivity = "UNKNOWN"

	//
	// The item is sold through one channel only, either local stores or online as
	// indicated by its ShoppingProductChannel.
	//
	ShoppingProductChannelExclusivitySINGLE_CHANNEL ShoppingProductChannelExclusivity = "SINGLE_CHANNEL"

	//
	// The item is matched to its online or local stores counterpart, indicating it is
	// available for purchase in both ShoppingProductChannels.
	//
	ShoppingProductChannelExclusivityMULTI_CHANNEL ShoppingProductChannelExclusivity = "MULTI_CHANNEL"
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
// Reported by system to reflect the criterion's serving status.
//
type SystemServingStatus string

const (

	//
	// Criterion is eligible to serve.
	//
	SystemServingStatusELIGIBLE SystemServingStatus = "ELIGIBLE"

	//
	// Indicates low search volume.
	// <p>For more information, visit
	// <a href="https://support.google.com/adwords/answer/2616014">Low Search Volume</a>.</p>
	//
	SystemServingStatusRARELY_SERVED SystemServingStatus = "RARELY_SERVED"
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

//
// Membership status of the user list.
//
type CriterionUserListMembershipStatus string

const (

	//
	// Open status - list is accruing members and can be targeted to.
	//
	CriterionUserListMembershipStatusOPEN CriterionUserListMembershipStatus = "OPEN"

	//
	// Closed status - No new members being added. Can not be used for targeting a new campaign.
	// Existing campaigns can still work as long as the list is not removed as a targeting criteria.
	//
	CriterionUserListMembershipStatusCLOSED CriterionUserListMembershipStatus = "CLOSED"
)

//
// Specified by user to pause or unpause a criterion.
//
type UserStatus string

const (

	//
	// Default state of a criterion (e.g. not paused).
	//
	UserStatusENABLED UserStatus = "ENABLED"

	//
	// Criterion is removed.
	//
	UserStatusREMOVED UserStatus = "REMOVED"

	//
	// Criterion is paused. Also used to pause a criterion.
	//
	UserStatusPAUSED UserStatus = "PAUSED"
)

//
// Operand value of {@link WebpageCondition}.
//
type WebpageConditionOperand string

const (

	//
	// Operand denoting a webpage URL targeting condition.
	// The operator {@link StringConditionOperator#CONTAINS} will be used for
	// such conditions.
	//
	WebpageConditionOperandURL WebpageConditionOperand = "URL"

	//
	// Operand denoting a webpage category targeting condition.
	// The operator {@link StringConditionOperator#EQUALS} will be used for
	// such conditions.
	//
	WebpageConditionOperandCATEGORY WebpageConditionOperand = "CATEGORY"

	//
	// Operand denoting a webpage title targeting condition.
	// The operator {@link StringConditionOperator#CONTAINS} will be used for
	// such conditions.
	//
	WebpageConditionOperandPAGE_TITLE WebpageConditionOperand = "PAGE_TITLE"

	//
	// Operand denoting a webpage content targeting condition.
	// The operator {@link StringConditionOperator#CONTAINS} will be used for
	// such conditions.
	//
	WebpageConditionOperandPAGE_CONTENT WebpageConditionOperand = "PAGE_CONTENT"

	//
	// Operand denoting a webpage custom label targeting condition.<br>
	// The operator {@link StringConditionOperator#EQUALS} will be used for such conditions.
	//
	WebpageConditionOperandCUSTOM_LABEL WebpageConditionOperand = "CUSTOM_LABEL"

	//
	// <span class="constraint Rejected">Used for return value only. An enumeration could not be processed, typically due to incompatibility with your WSDL version.</span>
	//
	WebpageConditionOperandUNKNOWN WebpageConditionOperand = "UNKNOWN"
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

	Rval *AdGroupCriterionPage `xml:"rval,omitempty"`
}

type Mutate struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 mutate"`

	//
	// <span class="constraint ContentsNotNull">This field must not contain {@code null} elements.</span>
	// <span class="constraint DistinctIds">Elements in this field must have distinct IDs for following {@link Operator}s : SET, REMOVE.</span>
	// <span class="constraint NotEmpty">This field must contain at least one element.</span>
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	// <span class="constraint SupportedOperators">The following {@link Operator}s are supported: ADD, SET, REMOVE.</span>
	//
	Operations []*AdGroupCriterionOperation `xml:"operations,omitempty"`
}

type MutateResponse struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 mutateResponse"`

	Rval *AdGroupCriterionReturnValue `xml:"rval,omitempty"`
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
	Operations []*AdGroupCriterionLabelOperation `xml:"operations,omitempty"`
}

type MutateLabelResponse struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 mutateLabelResponse"`

	Rval *AdGroupCriterionLabelReturnValue `xml:"rval,omitempty"`
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

	Rval *AdGroupCriterionPage `xml:"rval,omitempty"`
}

type AdGroupCriterion struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 AdGroupCriterion"`

	//
	// The ad group this criterion is in.
	// <span class="constraint Selectable">This field can be selected using the value "AdGroupId".</span><span class="constraint Filterable">This field can be filtered on.</span>
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	//
	AdGroupId int64 `xml:"adGroupId,omitempty"`

	//
	// <span class="constraint Selectable">This field can be selected using the value "CriterionUse".</span><span class="constraint Filterable">This field can be filtered on.</span>
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	CriterionUse *CriterionUse `xml:"criterionUse,omitempty"`

	//
	// The criterion part of the ad group criterion.
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	//
	Criterion *Criterion `xml:"criterion,omitempty"`

	//
	// Labels that are attached to the AdGroupCriterion. To associate an existing {@link Label} to an
	// existing {@link AdGroupCriterion}, use {@link AdGroupCriterionService#mutateLabel} with ADD
	// operator. To remove an associated {@link Label} from the {@link AdGroupCriterion}, use
	// {@link AdGroupCriterionService#mutateLabel} with REMOVE operator. To filter on {@link Label}s,
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
	// ID of the base campaign from which this draft/trial ad group criterion was created.
	// This field is only returned on get requests.
	// <span class="constraint Selectable">This field can be selected using the value "BaseCampaignId".</span><span class="constraint Filterable">This field can be filtered on.</span>
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	BaseCampaignId int64 `xml:"baseCampaignId,omitempty"`

	//
	// ID of the base ad group from which this draft/trial ad group criterion was created. For
	// base ad groups this is equal to the ad group ID.  If the ad group was created
	// in the draft or trial and has no corresponding base ad group, this field is null.
	// This field is only returned on get requests.
	// <span class="constraint Selectable">This field can be selected using the value "BaseAdGroupId".</span><span class="constraint Filterable">This field can be filtered on.</span>
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	BaseAdGroupId int64 `xml:"baseAdGroupId,omitempty"`

	//
	// Indicates that this instance is a subtype of AdGroupCriterion.
	// Although this field is returned in the response, it is ignored on input
	// and cannot be selected. Specify xsi:type instead.
	//
	AdGroupCriterionType string `xml:"AdGroupCriterion.Type,omitempty"`
}

type AdGroupCriterionError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 AdGroupCriterionError"`

	*ApiError

	//
	// The error reason represented by an enum.
	//
	Reason *AdGroupCriterionErrorReason `xml:"reason,omitempty"`
}

type AdGroupCriterionLabel struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 AdGroupCriterionLabel"`

	//
	// The id of the adgroup containing the criterion that the label is applied to.
	// <span class="constraint Required">This field is required and should not be {@code null} when it is contained within {@link Operator}s : ADD, REMOVE.</span>
	//
	AdGroupId int64 `xml:"adGroupId,omitempty"`

	//
	// The id of the criterion that the label is applied to.
	// <span class="constraint Required">This field is required and should not be {@code null} when it is contained within {@link Operator}s : ADD, REMOVE.</span>
	//
	CriterionId int64 `xml:"criterionId,omitempty"`

	//
	// The id of an existing label to be applied to the adgroup criterion.
	// <span class="constraint Required">This field is required and should not be {@code null} when it is contained within {@link Operator}s : ADD, REMOVE.</span>
	//
	LabelId int64 `xml:"labelId,omitempty"`
}

type AdGroupCriterionLabelOperation struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 AdGroupCriterionLabelOperation"`

	*Operation

	//
	// AdGroupCriterionLabel to operate on.
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	//
	Operand *AdGroupCriterionLabel `xml:"operand,omitempty"`
}

type AdGroupCriterionLabelReturnValue struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 AdGroupCriterionLabelReturnValue"`

	*ListReturnValue

	Value []*AdGroupCriterionLabel `xml:"value,omitempty"`

	PartialFailureErrors []*ApiError `xml:"partialFailureErrors,omitempty"`
}

type AdGroupCriterionLimitExceeded struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 AdGroupCriterionLimitExceeded"`

	*EntityCountLimitExceeded

	LimitType *AdGroupCriterionLimitExceededCriteriaLimitType `xml:"limitType,omitempty"`
}

type AdGroupCriterionOperation struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 AdGroupCriterionOperation"`

	*Operation

	//
	// The adgroup criterion being operated on.
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	//
	Operand *AdGroupCriterion `xml:"operand,omitempty"`

	//
	// List of exemption requests for policy violations flagged by this criterion.
	//
	// <p>Only set this field when adding criteria that trigger policy violations
	// for which you wish to get exemptions for.
	//
	ExemptionRequests []*ExemptionRequest `xml:"exemptionRequests,omitempty"`
}

type AdGroupCriterionPage struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 AdGroupCriterionPage"`

	*Page

	//
	// The result entries in this page.
	//
	Entries []*AdGroupCriterion `xml:"entries,omitempty"`
}

type AdGroupCriterionReturnValue struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 AdGroupCriterionReturnValue"`

	*ListReturnValue

	//
	// List of adgroup criteria.
	//
	Value []*AdGroupCriterion `xml:"value,omitempty"`

	//
	// List of partial failure errors.
	//
	PartialFailureErrors []*ApiError `xml:"partialFailureErrors,omitempty"`
}

type AdxError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 AdxError"`

	*ApiError

	//
	// The error reason represented by an enum.
	//
	Reason *AdxErrorReason `xml:"reason,omitempty"`
}

type AgeRange struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 AgeRange"`

	*Criterion

	//
	// <span class="constraint Selectable">This field can be selected using the value "AgeRangeType".</span>
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	AgeRangeType *AgeRangeAgeRangeType `xml:"ageRangeType,omitempty"`
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

type AppPaymentModel struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 AppPaymentModel"`

	*Criterion

	//
	// <span class="constraint Selectable">This field can be selected using the value "AppPaymentModelType".</span>
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	AppPaymentModelType *AppPaymentModelAppPaymentModelType `xml:"appPaymentModelType,omitempty"`
}

type AppUrl struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 AppUrl"`

	//
	// The app deep link url. E.g. "android-app://com.my.App"
	//
	Url string `xml:"url,omitempty"`

	//
	// The operating system targeted by this url.
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	//
	OsType *AppUrlOsType `xml:"osType,omitempty"`
}

type AppUrlList struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 AppUrlList"`

	//
	// List of URLs. On SET operation, empty list indicates to clear the list.
	// <span class="constraint CollectionSize">The maximum size of this collection is 10.</span>
	//
	AppUrls []*AppUrl `xml:"appUrls,omitempty"`
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

type Bid struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 Bid"`

	//
	// Bid amount.
	//
	Amount *Money `xml:"amount,omitempty"`
}

type BiddableAdGroupCriterion struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 BiddableAdGroupCriterion"`

	*AdGroupCriterion

	//
	// Current user-set state of criterion.
	// UserStatus may not be set to {@code REMOVED} and is not supported for ProductPartition
	// criterion. On add, defaults to {@code ENABLED} if unspecified.
	// <span class="constraint Selectable">This field can be selected using the value "Status".</span><span class="constraint Filterable">This field can be filtered on.</span>
	//
	UserStatus *UserStatus `xml:"userStatus,omitempty"`

	//
	// Serving status.
	// <span class="constraint Selectable">This field can be selected using the value "SystemServingStatus".</span><span class="constraint Filterable">This field can be filtered on.</span>
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	SystemServingStatus *SystemServingStatus `xml:"systemServingStatus,omitempty"`

	//
	// Approval status.
	// <span class="constraint Selectable">This field can be selected using the value "ApprovalStatus".</span><span class="constraint Filterable">This field can be filtered on.</span>
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	ApprovalStatus *ApprovalStatus `xml:"approvalStatus,omitempty"`

	//
	// List of disapproval reasons.
	// <span class="constraint Selectable">This field can be selected using the value "DisapprovalReasons".</span>
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	DisapprovalReasons []string `xml:"disapprovalReasons,omitempty"`

	//
	// Destination URL override when Ad is triggered by this criterion.
	//
	// <p>Some sample valid URLs are: "http://www.website.com",
	// "http://www.domain.com/somepath".
	// <p>Set to the empty string ("") to clear the destination URL.
	// <span class="constraint Selectable">This field can be selected using the value "DestinationUrl".</span><span class="constraint Filterable">This field can be filtered on.</span>
	//
	DestinationUrl string `xml:"destinationUrl,omitempty"`

	//
	// First page Cpc for this criterion.
	// <span class="constraint Selectable">This field can be selected using the value "FirstPageCpc".</span><span class="constraint Filterable">This field can be filtered on.</span>
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	FirstPageCpc *Bid `xml:"firstPageCpc,omitempty"`

	//
	// An estimate of the cpc bid needed for your ad to appear above the
	// first page of Google search results when a query matches the keywords exactly.
	// Note that meeting this estimate is not a guarantee of ad position,
	// which may depend on other factors.
	// <span class="constraint Selectable">This field can be selected using the value "TopOfPageCpc".</span><span class="constraint Filterable">This field can be filtered on.</span>
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	TopOfPageCpc *Bid `xml:"topOfPageCpc,omitempty"`

	//
	// An estimate of the cpc bid needed for your ad to regularly appear in the top position above
	// the search results on google.com when a query matches the keywords exactly.  Note that meeting
	// this estimate is not a guarantee of ad position, which may depend on other factors.
	// <span class="constraint Selectable">This field can be selected using the value "FirstPositionCpc".</span><span class="constraint Filterable">This field can be filtered on.</span>
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	FirstPositionCpc *Bid `xml:"firstPositionCpc,omitempty"`

	//
	// Contains quality information about the criterion.
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	QualityInfo *QualityInfo `xml:"qualityInfo,omitempty"`

	//
	// Bidding configuration for this ad group criterion. To set the bids on the ad groups
	// use {@link BiddingStrategyConfiguration#bids}. Multiple bids can be set on
	// ad group criterion at the same time. Only the bids that apply to the campaign's bidding
	// strategy {@linkplain Campaign#biddingStrategyConfiguration bidding strategy}
	// will be used.
	//
	BiddingStrategyConfiguration *BiddingStrategyConfiguration `xml:"biddingStrategyConfiguration,omitempty"`

	//
	// Bid modifier of the criterion which is used when the criterion is not in an absolute bidding
	// dimension.
	// <span class="constraint Selectable">This field can be selected using the value "BidModifier".</span><span class="constraint Filterable">This field can be filtered on.</span>
	//
	BidModifier float64 `xml:"bidModifier,omitempty"`

	//
	// A list of possible final URLs after all cross domain redirects.
	// <span class="constraint Selectable">This field can be selected using the value "FinalUrls".</span><span class="constraint Filterable">This field can be filtered on.</span>
	// <span class="constraint CampaignType">This field may not be set for campaign channel type SHOPPING with campaign channel subtype SHOPPING_UNIVERSAL_ADS.</span>
	//
	FinalUrls *UrlList `xml:"finalUrls,omitempty"`

	//
	// A list of possible final mobile URLs after all cross domain redirects.
	// <span class="constraint Selectable">This field can be selected using the value "FinalMobileUrls".</span><span class="constraint Filterable">This field can be filtered on.</span>
	// <span class="constraint CampaignType">This field may not be set for campaign channel type SHOPPING with campaign channel subtype SHOPPING_UNIVERSAL_ADS.</span>
	//
	FinalMobileUrls *UrlList `xml:"finalMobileUrls,omitempty"`

	//
	// A list of final app URLs that will be used on mobile if the user has the specific app
	// installed.
	// <span class="constraint Selectable">This field can be selected using the value "FinalAppUrls".</span><span class="constraint Filterable">This field can be filtered on.</span>
	// <span class="constraint CampaignType">This field may not be set for campaign channel type SHOPPING with campaign channel subtype SHOPPING_UNIVERSAL_ADS.</span>
	//
	FinalAppUrls *AppUrlList `xml:"finalAppUrls,omitempty"`

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
	//
	FinalUrlSuffix string `xml:"finalUrlSuffix,omitempty"`

	//
	// A list of mappings to be used for substituting URL custom parameter tags in the
	// trackingUrlTemplate, finalUrls, and/or finalMobileUrls.
	// <span class="constraint Selectable">This field can be selected using the value "UrlCustomParameters".</span><span class="constraint Filterable">This field can be filtered on.</span>
	//
	UrlCustomParameters *CustomParameters `xml:"urlCustomParameters,omitempty"`
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
	// <span class="constraint CampaignType">This field may not be set for campaign channel type SHOPPING.</span>
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
	// <span class="constraint CampaignType">This field may only be set to NONE for campaign channel type SHOPPING with ad group type SHOPPING_SHOWCASE_ADS.</span>
	// <span class="constraint CampaignType">This field may only be set to the values: MANUAL_CPC, ENHANCED_CPC, TARGET_ROAS, TARGET_SPEND, NONE for campaign channel type SHOPPING.</span>
	// <span class="constraint CampaignType">This field may only be set to NONE for campaign channel type SHOPPING with campaign channel subtype SHOPPING_UNIVERSAL_ADS.</span>
	// <span class="constraint CampaignType">This field may only be set to NONE for campaign channel type DISPLAY.</span>
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
	// <span class="constraint CampaignType">This field may not be set.</span>
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

type CollectionSizeError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 CollectionSizeError"`

	*ApiError

	//
	// The error reason represented by an enum.
	//
	Reason *CollectionSizeErrorReason `xml:"reason,omitempty"`
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
	//
	Bid *Money `xml:"bid,omitempty"`

	//
	// The level (ad group, ad group strategy, or campaign strategy) at which the bid was set.
	// This is applicable only at the ad group level.
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
	// <span class="constraint Selectable">This field can be selected using the value "CpcBidSource".</span><span class="constraint Filterable">This field can be filtered on.</span>
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
	// <span class="constraint Selectable">This field can be selected using the value "CpmBidSource".</span><span class="constraint Filterable">This field can be filtered on.</span>
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	CpmBidSource *BidSource `xml:"cpmBidSource,omitempty"`
}

type Criterion struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 Criterion"`

	//
	// ID of this criterion.
	// <span class="constraint Selectable">This field can be selected using the value "Id".</span><span class="constraint Filterable">This field can be filtered on.</span>
	// <span class="constraint Required">This field is required and should not be {@code null} when it is contained within {@link Operator}s : SET, REMOVE.</span>
	//
	Id int64 `xml:"id,omitempty"`

	//
	// <span class="constraint Selectable">This field can be selected using the value "CriteriaType".</span><span class="constraint Filterable">This field can be filtered on.</span>
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	Type_ *CriterionType `xml:"type,omitempty"`

	//
	// Indicates that this instance is a subtype of Criterion.
	// Although this field is returned in the response, it is ignored on input
	// and cannot be selected. Specify xsi:type instead.
	//
	CriterionType string `xml:"Criterion.Type,omitempty"`
}

type CriterionError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 CriterionError"`

	*ApiError

	Reason *CriterionErrorReason `xml:"reason,omitempty"`
}

type CriterionParameter struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 CriterionParameter"`

	//
	// Indicates that this instance is a subtype of CriterionParameter.
	// Although this field is returned in the response, it is ignored on input
	// and cannot be selected. Specify xsi:type instead.
	//
	CriterionParameterType string `xml:"CriterionParameter.Type,omitempty"`
}

type CriterionPolicyError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 CriterionPolicyError"`

	*PolicyViolationError
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

type ExemptionRequest struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 ExemptionRequest"`

	//
	// Identifies the violation to request an exemption for.
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	//
	Key *PolicyViolationKey `xml:"key,omitempty"`
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

type Gender struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 Gender"`

	*Criterion

	//
	// <span class="constraint Selectable">This field can be selected using the value "GenderType".</span>
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	GenderType *GenderGenderType `xml:"genderType,omitempty"`
}

type IdError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 IdError"`

	*ApiError

	//
	// The error reason represented by an enum.
	//
	Reason *IdErrorReason `xml:"reason,omitempty"`
}

type IncomeRange struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 IncomeRange"`

	*Criterion

	//
	// <span class="constraint Selectable">This field can be selected using the value "IncomeRangeType".</span>
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	IncomeRangeType *IncomeRangeIncomeRangeType `xml:"incomeRangeType,omitempty"`
}

type InternalApiError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 InternalApiError"`

	*ApiError

	//
	// The error reason represented by an enum.
	//
	Reason *InternalApiErrorReason `xml:"reason,omitempty"`
}

type Keyword struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 Keyword"`

	*Criterion

	//
	// Text of this keyword (at most 80 characters and ten words).
	// <span class="constraint Selectable">This field can be selected using the value "KeywordText".</span><span class="constraint Filterable">This field can be filtered on.</span>
	// <span class="constraint MatchesRegex">Keyword text must not contain NUL (code point 0x0) characters. This is checked by the regular expression '[^\x00]*'.</span>
	// <span class="constraint Required">This field is required and should not be {@code null} when it is contained within {@link Operator}s : ADD.</span>
	//
	Text string `xml:"text,omitempty"`

	//
	// Match type of this keyword.
	// <span class="constraint Selectable">This field can be selected using the value "KeywordMatchType".</span><span class="constraint Filterable">This field can be filtered on.</span>
	// <span class="constraint Required">This field is required and should not be {@code null} when it is contained within {@link Operator}s : ADD.</span>
	//
	MatchType *KeywordMatchType `xml:"matchType,omitempty"`
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

type MobileAppCategory struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 MobileAppCategory"`

	*Criterion

	//
	// ID of this mobile app category. A complete list of the available mobile app categories is
	// available <a href="/adwords/api/docs/appendix/mobileappcategories">here</a>.
	// <span class="constraint Selectable">This field can be selected using the value "MobileAppCategoryId".</span>
	// <span class="constraint Required">This field is required and should not be {@code null} when it is contained within {@link Operator}s : ADD.</span>
	//
	MobileAppCategoryId int32 `xml:"mobileAppCategoryId,omitempty"`

	//
	// Name of this mobile app category.
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	DisplayName string `xml:"displayName,omitempty"`
}

type MobileApplication struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 MobileApplication"`

	*Criterion

	//
	// A string that uniquely identifies a mobile application to AdWords API. The format of this
	// string is "<code>{platform}-{platform_native_id}</code>", where <code>platform</code> is "1"
	// for iOS apps and "2" for Android apps, and where <code>platform_native_id</code> is the mobile
	// application identifier native to the corresponding platform.
	// For iOS, this native identifier is the 9 digit string that appears at the end of an App Store
	// URL (e.g., "476943146" for "Flood-It! 2" whose App Store link is
	// http://itunes.apple.com/us/app/flood-it!-2/id476943146).
	// For Android, this native identifier is the application's package name (e.g.,
	// "com.labpixies.colordrips" for "Color Drips" given Google Play link
	// https://play.google.com/store/apps/details?id=com.labpixies.colordrips).
	// A well formed app id for AdWords API would thus be "1-476943146" for iOS and
	// "2-com.labpixies.colordrips" for Android.
	// <span class="constraint Selectable">This field can be selected using the value "AppId".</span><span class="constraint Filterable">This field can be filtered on.</span>
	// <span class="constraint Required">This field is required and should not be {@code null} when it is contained within {@link Operator}s : ADD.</span>
	//
	AppId string `xml:"appId,omitempty"`

	//
	// Title of this mobile application.
	// <span class="constraint Selectable">This field can be selected using the value "DisplayName".</span><span class="constraint Filterable">This field can be filtered on.</span>
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	DisplayName string `xml:"displayName,omitempty"`
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

type NegativeAdGroupCriterion struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 NegativeAdGroupCriterion"`

	*AdGroupCriterion
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

type PagingError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 PagingError"`

	*ApiError

	//
	// The error reason represented by an enum.
	//
	Reason *PagingErrorReason `xml:"reason,omitempty"`
}

type Parent struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 Parent"`

	*Criterion

	//
	// <span class="constraint Selectable">This field can be selected using the value "ParentType".</span>
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	ParentType *ParentParentType `xml:"parentType,omitempty"`
}

type Placement struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 Placement"`

	*Criterion

	//
	// Url of the placement.
	//
	// <p>For example, "http://www.domain.com".
	// <span class="constraint Selectable">This field can be selected using the value "PlacementUrl".</span><span class="constraint Filterable">This field can be filtered on.</span>
	// <span class="constraint Required">This field is required and should not be {@code null} when it is contained within {@link Operator}s : ADD.</span>
	//
	Url string `xml:"url,omitempty"`
}

type PolicyViolationError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 PolicyViolationError"`

	*ApiError

	//
	// Unique identifier for the violation.
	//
	Key *PolicyViolationKey `xml:"key,omitempty"`

	//
	// Name of policy suitable for display to users. In the user's preferred
	// language.
	//
	ExternalPolicyName string `xml:"externalPolicyName,omitempty"`

	//
	// Url with writeup about the policy.
	//
	ExternalPolicyUrl string `xml:"externalPolicyUrl,omitempty"`

	//
	// Localized description of the violation.
	//
	ExternalPolicyDescription string `xml:"externalPolicyDescription,omitempty"`

	//
	// Whether user can file an exemption request for this violation.
	//
	IsExemptable bool `xml:"isExemptable,omitempty"`

	//
	// Lists the parts that violate the policy.
	//
	ViolatingParts []*PolicyViolationErrorPart `xml:"violatingParts,omitempty"`
}

type PolicyViolationErrorPart struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 PolicyViolationError.Part"`

	//
	// Index of the starting position of the violating text within the line.
	//
	Index int32 `xml:"index,omitempty"`

	//
	// The length of the violating text.
	//
	Length int32 `xml:"length,omitempty"`
}

type PolicyViolationKey struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 PolicyViolationKey"`

	//
	// Unique id of the violated policy.
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	//
	PolicyName string `xml:"policyName,omitempty"`

	//
	// The text that violates the policy if specified. Otherwise, refers to the
	// policy in general (e.g. when requesting to be exempt from the whole
	// policy).
	//
	// May be null for criterion exemptions, in which case this refers to the
	// whole policy. Must be specified for ad exemptions.
	//
	ViolatingText string `xml:"violatingText,omitempty"`
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

type ProductAdwordsGrouping struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 ProductAdwordsGrouping"`

	*ProductDimension

	//
	// <span class="constraint StringLength">This string must not be empty.</span>
	//
	Value string `xml:"value,omitempty"`
}

type ProductAdwordsLabels struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 ProductAdwordsLabels"`

	*ProductDimension

	//
	// <span class="constraint StringLength">This string must not be empty.</span>
	//
	Value string `xml:"value,omitempty"`
}

type ProductBiddingCategory struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 ProductBiddingCategory"`

	*ProductDimension

	//
	// Dimension type of the category. Indicates the level of the category in the taxonomy.
	// <span class="constraint Filterable">This field can be filtered on using the value "ParentDimensionType".</span>
	// <span class="constraint OneOf">The value must be one of {BIDDING_CATEGORY_L1, BIDDING_CATEGORY_L2, BIDDING_CATEGORY_L3, BIDDING_CATEGORY_L4, BIDDING_CATEGORY_L5}.</span>
	// <span class="constraint Required">This field is required and should not be {@code null} when it is contained within {@link Operator}s : ADD.</span>
	//
	Type_ *ProductDimensionType `xml:"type,omitempty"`

	//
	// ID of the product category.
	// <span class="constraint Filterable">This field can be filtered on using the value "ParentDimensionId".</span>
	//
	Value int64 `xml:"value,omitempty"`
}

type ProductBrand struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 ProductBrand"`

	*ProductDimension

	//
	// <span class="constraint StringLength">This string must not be empty, (trimmed).</span>
	//
	Value string `xml:"value,omitempty"`
}

type ProductCanonicalCondition struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 ProductCanonicalCondition"`

	*ProductDimension

	Condition *ProductCanonicalConditionCondition `xml:"condition,omitempty"`
}

type ProductChannel struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 ProductChannel"`

	*ProductDimension

	Channel *ShoppingProductChannel `xml:"channel,omitempty"`
}

type ProductChannelExclusivity struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 ProductChannelExclusivity"`

	*ProductDimension

	ChannelExclusivity *ShoppingProductChannelExclusivity `xml:"channelExclusivity,omitempty"`
}

type ProductLegacyCondition struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 ProductLegacyCondition"`

	*ProductDimension

	Value string `xml:"value,omitempty"`
}

type ProductCustomAttribute struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 ProductCustomAttribute"`

	*ProductDimension

	//
	// Dimension type of the custom attribute. Indicates the index of the custom attribute.
	// <span class="constraint OneOf">The value must be one of {CUSTOM_ATTRIBUTE_0, CUSTOM_ATTRIBUTE_1, CUSTOM_ATTRIBUTE_2, CUSTOM_ATTRIBUTE_3, CUSTOM_ATTRIBUTE_4}.</span>
	// <span class="constraint Required">This field is required and should not be {@code null} when it is contained within {@link Operator}s : ADD.</span>
	//
	Type_ *ProductDimensionType `xml:"type,omitempty"`

	//
	// <span class="constraint StringLength">This string must not be empty, (trimmed).</span>
	//
	Value string `xml:"value,omitempty"`
}

type ProductDimension struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 ProductDimension"`

	//
	// Indicates that this instance is a subtype of ProductDimension.
	// Although this field is returned in the response, it is ignored on input
	// and cannot be selected. Specify xsi:type instead.
	//
	ProductDimensionType string `xml:"ProductDimension.Type,omitempty"`
}

type ProductOfferId struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 ProductOfferId"`

	*ProductDimension

	//
	// <span class="constraint StringLength">This string must not be empty, (trimmed).</span>
	//
	Value string `xml:"value,omitempty"`
}

type ProductPartition struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 ProductPartition"`

	*Criterion

	//
	// Type of the product partition.
	// <span class="constraint Selectable">This field can be selected using the value "PartitionType".</span>
	// <span class="constraint Required">This field is required and should not be {@code null} when it is contained within {@link Operator}s : ADD.</span>
	//
	PartitionType *ProductPartitionType `xml:"partitionType,omitempty"`

	//
	// ID of the parent product partition subdivision. Undefined for the root partition.
	// <span class="constraint Selectable">This field can be selected using the value "ParentCriterionId".</span>
	//
	ParentCriterionId int64 `xml:"parentCriterionId,omitempty"`

	//
	// Dimension value with which this product partition is refining its parent. Undefined for the
	// root partition.
	// <span class="constraint Selectable">This field can be selected using the value "CaseValue".</span>
	//
	CaseValue *ProductDimension `xml:"caseValue,omitempty"`
}

type ProductType struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 ProductType"`

	*ProductDimension

	//
	// Dimension type of the product type. Indicates the level of the product type.
	// <span class="constraint OneOf">The value must be one of {PRODUCT_TYPE_L1, PRODUCT_TYPE_L2, PRODUCT_TYPE_L3, PRODUCT_TYPE_L4, PRODUCT_TYPE_L5}.</span>
	// <span class="constraint Required">This field is required and should not be {@code null} when it is contained within {@link Operator}s : ADD.</span>
	//
	Type_ *ProductDimensionType `xml:"type,omitempty"`

	//
	// <span class="constraint StringLength">This string must not be empty, (trimmed).</span>
	//
	Value string `xml:"value,omitempty"`
}

type ProductTypeFull struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 ProductTypeFull"`

	*ProductDimension

	//
	// <span class="constraint StringLength">This string must not be empty.</span>
	//
	Value string `xml:"value,omitempty"`
}

type QualityInfo struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 QualityInfo"`

	//
	// The keyword quality score ranges from 1 (lowest) to 10 (highest).
	// <p>If there aren't enough impressions or clicks to determine an appropriate
	// quality score value, the QualityInfo object is not returned. For reports,
	// this field will return null (designated by "--").
	// <span class="constraint Selectable">This field can be selected using the value "QualityScore".</span><span class="constraint Filterable">This field can be filtered on.</span>
	//
	QualityScore int32 `xml:"qualityScore,omitempty"`
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

type UnknownProductDimension struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 UnknownProductDimension"`

	*ProductDimension
}

type UrlError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 UrlError"`

	*ApiError

	//
	// The error reason represented by an enum.
	//
	Reason *UrlErrorReason `xml:"reason,omitempty"`
}

type UrlList struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 UrlList"`

	//
	// List of URLs.  On SET operation, empty list indicates to clear the list.
	// <span class="constraint CollectionSize">The maximum size of this collection is 10.</span>
	// <span class="constraint ContentsStringLength">Strings in this field must be non-empty (trimmed).</span>
	//
	Urls []string `xml:"urls,omitempty"`
}

type CriterionUserInterest struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 CriterionUserInterest"`

	*Criterion

	//
	// Id of this user interest. This is a required field.
	// <span class="constraint Selectable">This field can be selected using the value "UserInterestId".</span>
	//
	UserInterestId int64 `xml:"userInterestId,omitempty"`

	//
	// Parent Id of this user interest.
	// <span class="constraint Selectable">This field can be selected using the value "UserInterestParentId".</span>
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	UserInterestParentId int64 `xml:"userInterestParentId,omitempty"`

	//
	// Name of this user interest.
	// <span class="constraint Selectable">This field can be selected using the value "UserInterestName".</span>
	//
	UserInterestName string `xml:"userInterestName,omitempty"`
}

type CriterionUserList struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 CriterionUserList"`

	*Criterion

	//
	// Id of this user list. This is a required field.
	// <span class="constraint Selectable">This field can be selected using the value "UserListId".</span>
	//
	UserListId int64 `xml:"userListId,omitempty"`

	//
	// <span class="constraint Selectable">This field can be selected using the value "UserListName".</span>
	//
	UserListName string `xml:"userListName,omitempty"`

	//
	// <span class="constraint Selectable">This field can be selected using the value "UserListMembershipStatus".</span><span class="constraint Filterable">This field can be filtered on.</span>
	//
	UserListMembershipStatus *CriterionUserListMembershipStatus `xml:"userListMembershipStatus,omitempty"`

	//
	// Determines whether a user list is eligible for targeting in the google.com
	// (search) network.
	// <span class="constraint Selectable">This field can be selected using the value "UserListEligibleForSearch".</span><span class="constraint Filterable">This field can be filtered on.</span>
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	UserListEligibleForSearch bool `xml:"userListEligibleForSearch,omitempty"`

	//
	// Determines whether a user list is eligible for targeting in the display network.
	// <span class="constraint Selectable">This field can be selected using the value "UserListEligibleForDisplay".</span><span class="constraint Filterable">This field can be filtered on.</span>
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	UserListEligibleForDisplay bool `xml:"userListEligibleForDisplay,omitempty"`
}

type Vertical struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 Vertical"`

	*Criterion

	//
	// Id of this vertical.
	// <span class="constraint Selectable">This field can be selected using the value "VerticalId".</span>
	//
	VerticalId int64 `xml:"verticalId,omitempty"`

	//
	// Id of the parent of this vertical.
	// <span class="constraint Selectable">This field can be selected using the value "VerticalParentId".</span>
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	VerticalParentId int64 `xml:"verticalParentId,omitempty"`

	//
	// The category to target or exclude. Each subsequent element in the array
	// describes a more specific sub-category. For example,
	// <code>{"Pets &amp; Animals", "Pets", "Dogs"}</code> represents the "Pets &amp;
	// Animals/Pets/Dogs" category. A complete list of available vertical categories
	// is available <a href="/adwords/api/docs/appendix/verticals">here</a>
	// This field is required and must not be empty.
	// <span class="constraint Selectable">This field can be selected using the value "Path".</span>
	//
	Path []string `xml:"path,omitempty"`
}

type Webpage struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 Webpage"`

	*Criterion

	//
	// The webpage criterion parameter.
	// <span class="constraint Selectable">This field can be selected using the value "Parameter".</span>
	// <span class="constraint Required">This field is required and should not be {@code null} when it is contained within {@link Operator}s : ADD.</span>
	//
	Parameter *WebpageParameter `xml:"parameter,omitempty"`

	//
	// Keywordless criteria coverage - Computed percentage of website coverage based on the
	// website target, negative website targets and negative keywords in the ad group and campaign.
	// <span class="constraint Selectable">This field can be selected using the value "CriteriaCoverage".</span>
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	CriteriaCoverage float64 `xml:"criteriaCoverage,omitempty"`

	//
	// Keywordless criteria samples - List of sample urls that matches with the website target.
	// <span class="constraint Selectable">This field can be selected using the value "CriteriaSamples".</span>
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	CriteriaSamples []string `xml:"criteriaSamples,omitempty"`
}

type WebpageCondition struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 WebpageCondition"`

	//
	// Operand of webpage targeting condition.
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	//
	Operand *WebpageConditionOperand `xml:"operand,omitempty"`

	//
	// Argument of the webpage targeting condition.
	// <span class="constraint MustNotContain">This string must not contain a substring that matches the regular expression '\*|\>\>|\=\=|\&\+'</span>
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	// <span class="constraint StringLength">The length of this string should be between 1 and 2048, inclusive.</span>
	//
	Argument string `xml:"argument,omitempty"`
}

type WebpageParameter struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 WebpageParameter"`

	*CriterionParameter

	//
	// The name of the criterion that is defined by this parameter.
	//
	// <p>This name value will be used for identifying, sorting and filtering
	// criteria with this type of parameters. For criteria with simpler
	// parameters, such as keywords and placements, the parameter value itself
	// is used for identification, sorting and filtering.
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	// <span class="constraint StringLength">The length of this string should be between 1 and 2048, inclusive.</span>
	//
	CriterionName string `xml:"criterionName,omitempty"`

	//
	// Conditions, or logical expressions, for webpage targeting.
	//
	// <p>The list of webpage targeting conditions are {@code and}-ed together
	// when evaluated for targeting. A {@code null} list of conditions means that
	// all webpages of the campaign's website are targeted.</p>
	// <span class="constraint CollectionSize">The maximum size of this collection is 3.</span>
	// <span class="constraint ContentsDistinct">This field must contain distinct elements.</span>
	//
	Conditions []*WebpageCondition `xml:"conditions,omitempty"`
}

type YouTubeChannel struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 YouTubeChannel"`

	*Criterion

	//
	// The YouTube uploader channel id or the channel code of a YouTube content channel.
	// <p>The uploader channel id can be obtained from the YouTube id-based URL. For example, in
	// <code>https://www.youtube.com/channel/UCEN58iXQg82TXgsDCjWqIkg</code> the channel id is
	// <code>UCEN58iXQg82TXgsDCjWqIkg</code>
	// <p>For more information see: https://support.google.com/youtube/answer/6180214
	// <span class="constraint Selectable">This field can be selected using the value "ChannelId".</span>
	//
	ChannelId string `xml:"channelId,omitempty"`

	//
	// The public name for a YouTube user channel.
	// <span class="constraint Selectable">This field can be selected using the value "ChannelName".</span>
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	ChannelName string `xml:"channelName,omitempty"`
}

type YouTubeVideo struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 YouTubeVideo"`

	*Criterion

	//
	// YouTube video id as it appears on the YouTube watch page.
	// <span class="constraint Selectable">This field can be selected using the value "VideoId".</span>
	//
	VideoId string `xml:"videoId,omitempty"`

	//
	// Name of the video.
	// <span class="constraint Selectable">This field can be selected using the value "VideoName".</span>
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	VideoName string `xml:"videoName,omitempty"`
}

type AdGroupCriterionServiceInterface struct {
	client *SOAPClient
}

func NewAdGroupCriterionServiceInterface(url string, tls bool, auth *BasicAuth) *AdGroupCriterionServiceInterface {
	if url == "" {
		url = ""
	}
	client := NewSOAPClient(url, tls, auth)

	return &AdGroupCriterionServiceInterface{
		client: client,
	}
}

func NewAdGroupCriterionServiceInterfaceWithTLSConfig(url string, tlsCfg *tls.Config, auth *BasicAuth) *AdGroupCriterionServiceInterface {
	if url == "" {
		url = ""
	}
	client := NewSOAPClientWithTLSConfig(url, tlsCfg, auth)

	return &AdGroupCriterionServiceInterface{
		client: client,
	}
}

func (service *AdGroupCriterionServiceInterface) AddHeader(header interface{}) {
	service.client.AddHeader(header)
}

// Backwards-compatible function: use AddHeader instead
func (service *AdGroupCriterionServiceInterface) SetHeader(header interface{}) {
	service.client.AddHeader(header)
}

// Error can be either of the following types:
//
//   - ApiException
/*
   Gets adgroup criteria.

   @param serviceSelector filters the adgroup criteria to be returned.
   @return a page (subset) view of the criteria selected
   @throws ApiException when there is at least one error with the request
*/
func (service *AdGroupCriterionServiceInterface) Get(request *Get) (*GetResponse, error) {
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
   Adds, removes or updates adgroup criteria.

   @param operations operations to do
   during checks on keywords to be added.
   @return added and updated adgroup criteria (without optional parts)
   @throws ApiException when there is at least one error with the request
*/
func (service *AdGroupCriterionServiceInterface) Mutate(request *Mutate) (*MutateResponse, error) {
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
   Adds labels to the AdGroupCriterion or removes labels from the AdGroupCriterion
   <p>Add - Apply an existing label to an existing
   {@linkplain AdGroupCriterion ad group criterion}. The {@code adGroupId} and
   {@code criterionId}
   must reference an existing {@linkplain AdGroupCriterion ad group criterion}. The
   {@code labelId} must reference an existing {@linkplain Label label}.
   <p>Remove - Removes the link between the specified
   {@linkplain AdGroupCriterion ad group criterion} and {@linkplain Label label}.</p>
   @param operations the operations to apply
   @return a list of AdGroupCriterionLabel where each entry in the list is the result of
   applying the operation in the input list with the same index. For an
   add operation, the returned AdGroupCriterionLabel contains the AdGroupId, CriterionId and the
   LabelId. In the case of a remove operation, the returned AdGroupCriterionLabel will only have
   AdGroupId and CriterionId.
   @throws ApiException when there are one or more errors with the request
*/
func (service *AdGroupCriterionServiceInterface) MutateLabel(request *MutateLabel) (*MutateLabelResponse, error) {
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
   Returns the list of AdGroupCriterion that match the query.

   @param query The SQL-like AWQL query string
   @returns A list of AdGroupCriterion
   @throws ApiException when the query is invalid or there are errors processing the request.
*/
func (service *AdGroupCriterionServiceInterface) Query(request *Query) (*QueryResponse, error) {
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
