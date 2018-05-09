package TargetingIdeaService

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

//
// Enum that represents the different Targeting Status values for a Location criterion.
//
type LocationTargetingStatus string

const (

	//
	// The location is active.
	//
	LocationTargetingStatusACTIVE LocationTargetingStatus = "ACTIVE"

	//
	// The location is not available for targeting.
	//
	LocationTargetingStatusOBSOLETE LocationTargetingStatus = "OBSOLETE"

	//
	// The location is phasing out, it will marked obsolete soon.
	//
	LocationTargetingStatusPHASING_OUT LocationTargetingStatus = "PHASING_OUT"
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
// The reasons for the validation error.
//
type OperatorErrorReason string

const (
	OperatorErrorReasonOPERATOR_NOT_SUPPORTED OperatorErrorReason = "OPERATOR_NOT_SUPPORTED"
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
// The reasons for the validation error.
//
type RegionCodeErrorReason string

const (
	RegionCodeErrorReasonINVALID_REGION_CODE RegionCodeErrorReason = "INVALID_REGION_CODE"
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

type AdGroupCriterionError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 AdGroupCriterionError"`

	*ApiError

	//
	// The error reason represented by an enum.
	//
	Reason *AdGroupCriterionErrorReason `xml:"reason,omitempty"`
}

type AdGroupCriterionLimitExceeded struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 AdGroupCriterionLimitExceeded"`

	*EntityCountLimitExceeded

	LimitType *AdGroupCriterionLimitExceededCriteriaLimitType `xml:"limitType,omitempty"`
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

type BudgetError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 BudgetError"`

	*ApiError

	//
	// The error reason represented by an enum.
	//
	Reason *BudgetErrorReason `xml:"reason,omitempty"`
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

type CriterionPolicyError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 CriterionPolicyError"`

	*PolicyViolationError
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

type Language struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 Language"`

	*Criterion

	//
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	Code string `xml:"code,omitempty"`

	//
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	Name string `xml:"name,omitempty"`
}

type Location struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 Location"`

	*Criterion

	//
	// Name of the location criterion. <b> Note:</b> This field is filterable only in
	// LocationCriterionService. If used as a filter, a location name cannot be greater than 300
	// characters.
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	LocationName string `xml:"locationName,omitempty"`

	//
	// Display type of the location criterion.
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	DisplayType string `xml:"displayType,omitempty"`

	//
	// The targeting status of the location criterion.
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	TargetingStatus *LocationTargetingStatus `xml:"targetingStatus,omitempty"`

	//
	// Ordered list of parents of the location criterion.
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	ParentLocations []*Location `xml:"parentLocations,omitempty"`
}

type LongValue struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 LongValue"`

	*NumberValue

	//
	// the underlying long value.
	//
	Number int64 `xml:"number,omitempty"`
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

type NetworkSetting struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 NetworkSetting"`

	//
	// Ads will be served with Google.com search results.
	// <span class="constraint AdxEnabled">This is disabled for AdX.</span>
	// <span class="constraint CampaignType">This field may only be set to true for campaign channel type SEARCH.</span>
	// <span class="constraint CampaignType">This field may only be set to true for campaign channel type MULTI_CHANNEL.</span>
	// <span class="constraint CampaignType">This field may only be set to false for campaign channel type DISPLAY.</span>
	// <span class="constraint CampaignType">This field may only be set to true for campaign channel subtype SHOPPING_UNIVERSAL_ADS.</span>
	//
	TargetGoogleSearch bool `xml:"targetGoogleSearch,omitempty"`

	//
	// Ads will be served on partner sites in the Google Search Network
	// (requires {@code GOOGLE_SEARCH}).
	// <span class="constraint AdxEnabled">This is disabled for AdX.</span>
	// <span class="constraint CampaignType">This field may only be set to true for campaign channel type MULTI_CHANNEL.</span>
	// <span class="constraint CampaignType">This field may only be set to true for campaign channel subtype SHOPPING_UNIVERSAL_ADS.</span>
	//
	TargetSearchNetwork bool `xml:"targetSearchNetwork,omitempty"`

	//
	// Ads will be served on specified placements in the Google Display Network.
	// Placements are specified using {@code Placement} criteria.
	// <span class="constraint CampaignType">This field may only be set to true for campaign channel type MULTI_CHANNEL.</span>
	// <span class="constraint CampaignType">This field may only be set to false for campaign channel subtype SEARCH_MOBILE_APP.</span>
	// <span class="constraint CampaignType">This field may only be set to true for campaign channel subtype SHOPPING_UNIVERSAL_ADS.</span>
	//
	TargetContentNetwork bool `xml:"targetContentNetwork,omitempty"`

	//
	// Ads will be served on the Google Partner Network. This is available to
	// only some specific Google partner accounts.
	// <span class="constraint AdxEnabled">This is disabled for AdX.</span>
	// <span class="constraint CampaignType">This field may only be set to false for campaign channel type MULTI_CHANNEL.</span>
	//
	TargetPartnerSearchNetwork bool `xml:"targetPartnerSearchNetwork,omitempty"`
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

type Platform struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 Platform"`

	*Criterion

	//
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	PlatformName string `xml:"platformName,omitempty"`
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

type RegionCodeError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 RegionCodeError"`

	*ApiError

	//
	// The error reason represented by an enum.
	//
	Reason *RegionCodeErrorReason `xml:"reason,omitempty"`
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

//
// Represents the type of {@link Attribute}.
// <p><b>{@link IdeaType} KEYWORD supports the following {@link AttributeType}s:</b><br/>
// <ul><li>{@link #AVERAGE_CPC}</li>
// <li>{@link #CATEGORY_PRODUCTS_AND_SERVICES}</li>
// <li>{@link #COMPETITION}</li>
// <li>{@link #EXTRACTED_FROM_WEBPAGE}</li>
// <li>{@link #IDEA_TYPE}</li>
// <li>{@link #KEYWORD_TEXT}</li>
// <li>{@link #SEARCH_VOLUME}</li>
// <li>{@link #TARGETED_MONTHLY_SEARCHES}</li>
// </ul>
//
type AttributeType string

const (

	//
	// Value substituted in when the actual value is not available in the Web API
	// version being used.  (Please upgrade to the latest published WSDL.)
	// <p>This element is not supported directly by any {@link IdeaType}.
	//
	AttributeTypeUNKNOWN AttributeType = "UNKNOWN"

	//
	// Represents a category ID in the "Products and Services" taxonomy.
	//
	// <p>Resulting attribute is {@link IntegerSetAttribute}.
	// <p>This element is supported by following {@link IdeaType}s: KEYWORD.
	//
	AttributeTypeCATEGORY_PRODUCTS_AND_SERVICES AttributeType = "CATEGORY_PRODUCTS_AND_SERVICES"

	//
	// Represents the relative amount of competition associated with the given keyword idea,
	// relative to other keywords. This value will be between 0 and 1 (inclusive).
	//
	// <p>Resulting attribute is {@link DoubleAttribute}.
	// <p>This element is supported by following {@link IdeaType}s: KEYWORD.
	//
	AttributeTypeCOMPETITION AttributeType = "COMPETITION"

	//
	// Represents the webpage from which this keyword idea was extracted (if applicable.)
	//
	// <p>Resulting attribute is {@link WebpageDescriptorAttribute}.
	// <p>This element is supported by following {@link IdeaType}s: KEYWORD.
	//
	AttributeTypeEXTRACTED_FROM_WEBPAGE AttributeType = "EXTRACTED_FROM_WEBPAGE"

	//
	// Represents the type of the given idea.
	//
	// <p>Resulting attribute is {@link IdeaTypeAttribute}.
	// <p>This element is supported by following {@link IdeaType}s: KEYWORD.
	//
	AttributeTypeIDEA_TYPE AttributeType = "IDEA_TYPE"

	//
	// Represents the keyword text for the given keyword idea.
	//
	// <p>Resulting attribute is {@link StringAttribute}.
	// <p>This element is supported by following {@link IdeaType}s: KEYWORD.
	//
	AttributeTypeKEYWORD_TEXT AttributeType = "KEYWORD_TEXT"

	//
	// Represents either the (approximate) number of searches for the given keyword idea on google.com
	// or google.com and partners, depending on the user's targeting.
	//
	// <p>Resulting attribute is {@link LongAttribute}.
	// <p>This element is supported by following {@link IdeaType}s: KEYWORD.
	//
	AttributeTypeSEARCH_VOLUME AttributeType = "SEARCH_VOLUME"

	//
	// Represents the average cost per click historically paid for the keyword.
	//
	// <p>Resulting attribute is {@link MoneyAttribute}.
	// <p>This element is supported by following {@link IdeaType}s: KEYWORD.
	//
	AttributeTypeAVERAGE_CPC AttributeType = "AVERAGE_CPC"

	//
	// Represents the (approximated) number of searches on this keyword idea (as available for the
	// past twelve months), targeted to the specified geographies.
	//
	// <p>Resulting attribute is {@link MonthlySearchVolumeAttribute}.
	// <p>This element is supported by following {@link IdeaType}s: KEYWORD.
	//
	AttributeTypeTARGETED_MONTHLY_SEARCHES AttributeType = "TARGETED_MONTHLY_SEARCHES"
)

//
// An enumeration of possible values to be used in conjunction with the
// {@link CompetitionSearchParameter} to specify the granularity of
// competition to be filtered.
//
type CompetitionSearchParameterLevel string

const (

	//
	// Low - competition rate [0.0000, 0.3333]
	//
	CompetitionSearchParameterLevelLOW CompetitionSearchParameterLevel = "LOW"

	//
	// Medium - competition rate (0.3333, 0.6667]
	//
	CompetitionSearchParameterLevelMEDIUM CompetitionSearchParameterLevel = "MEDIUM"

	//
	// High - competition rate (0.6667, 1.0000]
	//
	CompetitionSearchParameterLevelHIGH CompetitionSearchParameterLevel = "HIGH"
)

//
// Encodes the reason (cause) of a particular {@link CurrencyCodeError}.
//
type CurrencyCodeErrorReason string

const (
	CurrencyCodeErrorReasonUNSUPPORTED_CURRENCY_CODE CurrencyCodeErrorReason = "UNSUPPORTED_CURRENCY_CODE"
)

//
// Represents the type of idea.
// <span class="constraint AdxEnabled">This is disabled for AdX.</span>
//
type IdeaType string

const (

	//
	// Keyword idea.
	//
	IdeaTypeKEYWORD IdeaType = "KEYWORD"
)

//
// Represents the type of the request.
//
type RequestType string

const (

	//
	// Request for new ideas based on other entries in selector.
	// This {@link RequestType} can be used to request other ideas
	// using keyword/placements that the user already has.
	//
	RequestTypeIDEAS RequestType = "IDEAS"

	//
	// Request for stats for entries in selector.
	// This {@link RequestType} can be used to request
	// the stats for keywords/placements that the user already has.
	//
	// <p>Stats are generated once a month (typically on the last week of the
	// month) from the historical data of previous months.</p>
	//
	RequestTypeSTATS RequestType = "STATS"
)

//
// An enumeration of {@link TargetingIdeaService} specific errors.
//
type TargetingIdeaErrorReason string

const (

	//
	// Error returned when there are multiple instance of same type of {@link SearchParameter}s.
	//
	TargetingIdeaErrorReasonDUPLICATE_SEARCH_FILTER_TYPES_PRESENT TargetingIdeaErrorReason = "DUPLICATE_SEARCH_FILTER_TYPES_PRESENT"

	//
	// Error returned when the {@link TargetingIdeaSelector} doesn't have enough
	// {@link SearchParameter}s to execute request.
	//
	TargetingIdeaErrorReasonINSUFFICIENT_SEARCH_PARAMETERS TargetingIdeaErrorReason = "INSUFFICIENT_SEARCH_PARAMETERS"

	//
	// Error returned when an {@link AttributeType} doesn't match the {@link IdeaType} specified in
	// the {@link TargetingIdeaSelector}. For example, if the {@code KEYWORD} {@code IDEAS} selector
	// contains an {@code STATS} only AttributeType, this error will be returned.
	//
	TargetingIdeaErrorReasonINVALID_ATTRIBUTE_TYPE TargetingIdeaErrorReason = "INVALID_ATTRIBUTE_TYPE"

	//
	// Error returned when a {@link SearchParameter} doesn't match the {@link IdeaType} specified in
	// the {@link TargetingIdeaSelector} or is otherwise invalid.  Error trigger usually contains
	// the parameter name, and error details contain a more detailed explanation.
	//
	TargetingIdeaErrorReasonINVALID_SEARCH_PARAMETERS TargetingIdeaErrorReason = "INVALID_SEARCH_PARAMETERS"

	//
	// Error returned when the {@link TargetingIdeaSelector} contains a
	// {@link DomainSuffixSearchParameter}s that contains an invalid domain suffix.
	//
	TargetingIdeaErrorReasonINVALID_DOMAIN_SUFFIX TargetingIdeaErrorReason = "INVALID_DOMAIN_SUFFIX"

	//
	// Error returned when a selector contains mutually exclusive parameters.
	//
	TargetingIdeaErrorReasonMUTUALLY_EXCLUSIVE_SEARCH_PARAMETERS_IN_QUERY TargetingIdeaErrorReason = "MUTUALLY_EXCLUSIVE_SEARCH_PARAMETERS_IN_QUERY"

	//
	// Error returned when the {@link TargetingIdeaService} is not available.
	//
	TargetingIdeaErrorReasonSERVICE_UNAVAILABLE TargetingIdeaErrorReason = "SERVICE_UNAVAILABLE"

	//
	// Error returned when the URL value specified in the {@link TargetingIdeaSelector}, such as
	// {@link RelatedToUrlSearchParameter}, is not a valid URL.
	//
	TargetingIdeaErrorReasonINVALID_URL_IN_SEARCH_PARAMETER TargetingIdeaErrorReason = "INVALID_URL_IN_SEARCH_PARAMETER"

	//
	// Error returned when the requested number of entries in {@link TargetingIdeaSelector}'s
	// {@link Paging} is greater than the maximum allowed.
	//
	TargetingIdeaErrorReasonTOO_MANY_RESULTS_REQUESTED TargetingIdeaErrorReason = "TOO_MANY_RESULTS_REQUESTED"

	//
	// Error returned when the requested {@link Paging} is missing from the
	// {@link TargetingIdeaSelector} when required.
	//
	TargetingIdeaErrorReasonNO_PAGING_IN_SELECTOR TargetingIdeaErrorReason = "NO_PAGING_IN_SELECTOR"

	//
	// Error returned when included keywords and excluded keywords in
	// {@link IdeaTextFilterSearchParameter}, {@link IdeaTextMatchesSearchParameter}
	// or {@link ExcludedKeywordSearchParameter} are overlapped.
	//
	TargetingIdeaErrorReasonINVALID_INCLUDED_EXCLUDED_KEYWORDS TargetingIdeaErrorReason = "INVALID_INCLUDED_EXCLUDED_KEYWORDS"
)

type TrafficEstimatorErrorReason string

const (

	//
	// When the request with {@code null} campaign ID in {@link CampaignEstimateRequest} contains an
	// {@link AdGroupEstimateRequest} with an ID.
	//
	TrafficEstimatorErrorReasonNO_CAMPAIGN_FOR_AD_GROUP_ESTIMATE_REQUEST TrafficEstimatorErrorReason = "NO_CAMPAIGN_FOR_AD_GROUP_ESTIMATE_REQUEST"

	//
	// When the request with {@code null} adgroup ID in {@link AdGroupEstimateRequest} contains a
	// {@link KeywordEstimateRequest} with an ID.
	//
	TrafficEstimatorErrorReasonNO_AD_GROUP_FOR_KEYWORD_ESTIMATE_REQUEST TrafficEstimatorErrorReason = "NO_AD_GROUP_FOR_KEYWORD_ESTIMATE_REQUEST"

	//
	// All {@link KeywordEstimateRequest} items should have maxCpc associated with them.
	//
	TrafficEstimatorErrorReasonNO_MAX_CPC_FOR_KEYWORD_ESTIMATE_REQUEST TrafficEstimatorErrorReason = "NO_MAX_CPC_FOR_KEYWORD_ESTIMATE_REQUEST"

	//
	// When there are more {@link KeywordEstimateRequest}s in the request than
	// TrafficEstimatorService allows.
	//
	TrafficEstimatorErrorReasonTOO_MANY_KEYWORD_ESTIMATE_REQUESTS TrafficEstimatorErrorReason = "TOO_MANY_KEYWORD_ESTIMATE_REQUESTS"

	//
	// When there are more {@link CampaignEstimateRequest}s in the request than
	// TrafficEstimatorService allows.
	//
	TrafficEstimatorErrorReasonTOO_MANY_CAMPAIGN_ESTIMATE_REQUESTS TrafficEstimatorErrorReason = "TOO_MANY_CAMPAIGN_ESTIMATE_REQUESTS"

	//
	// When there are more {@link AdGroupEstimateRequest}s in the request than
	// TrafficEstimatorService allows.
	//
	TrafficEstimatorErrorReasonTOO_MANY_ADGROUP_ESTIMATE_REQUESTS TrafficEstimatorErrorReason = "TOO_MANY_ADGROUP_ESTIMATE_REQUESTS"

	//
	// When there are more targets in the request than TrafficEstimatorService allows. See
	// documentation on {@link CampaignEstimateRequest} for more information about this error.
	//
	TrafficEstimatorErrorReasonTOO_MANY_TARGETS TrafficEstimatorErrorReason = "TOO_MANY_TARGETS"

	//
	// Request contains a keyword that is too long for backends to handle.
	//
	TrafficEstimatorErrorReasonKEYWORD_TOO_LONG TrafficEstimatorErrorReason = "KEYWORD_TOO_LONG"

	//
	// Request contains a keyword that contains broad match modifiers.
	//
	TrafficEstimatorErrorReasonKEYWORD_CONTAINS_BROAD_MATCH_MODIFIERS TrafficEstimatorErrorReason = "KEYWORD_CONTAINS_BROAD_MATCH_MODIFIERS"

	//
	// When an unexpected error occurs.
	//
	TrafficEstimatorErrorReasonINVALID_INPUT TrafficEstimatorErrorReason = "INVALID_INPUT"

	//
	// When backend service calls fail.
	//
	TrafficEstimatorErrorReasonSERVICE_UNAVAILABLE TrafficEstimatorErrorReason = "SERVICE_UNAVAILABLE"
)

type Get struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/o/v201802 get"`

	//
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	//
	Selector *TargetingIdeaSelector `xml:"selector,omitempty"`
}

type GetResponse struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/o/v201802 getResponse"`

	Rval *TargetingIdeaPage `xml:"rval,omitempty"`
}

type Attribute struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/o/v201802 Attribute"`

	//
	// Indicates that this instance is a subtype of Attribute.
	// Although this field is returned in the response, it is ignored on input
	// and cannot be selected. Specify xsi:type instead.
	//
	AttributeType string `xml:"Attribute.Type,omitempty"`
}

type BooleanAttribute struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/o/v201802 BooleanAttribute"`

	*Attribute

	//
	// Boolean value contained by this {@link Attribute}.
	//
	Value bool `xml:"value,omitempty"`
}

type CategoryProductsAndServicesSearchParameter struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/o/v201802 CategoryProductsAndServicesSearchParameter"`

	*SearchParameter

	//
	// A keyword category ID in the "Products and Services" taxonomy that all
	// search results should belong to.
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	//
	CategoryId int32 `xml:"categoryId,omitempty"`
}

type CompetitionSearchParameter struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/o/v201802 CompetitionSearchParameter"`

	*SearchParameter

	//
	// A set of {@link Level}s indicating a relative amount of competition that
	// {@code KEYWORD} {@link IdeaType}s should have in the  results.
	// <span class="constraint ContentsDistinct">This field must contain distinct elements.</span>
	// <span class="constraint ContentsNotNull">This field must not contain {@code null} elements.</span>
	// <span class="constraint NotEmpty">This field must contain at least one element.</span>
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	//
	Levels []*CompetitionSearchParameterLevel `xml:"levels,omitempty"`
}

type CriterionAttribute struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/o/v201802 CriterionAttribute"`

	*Attribute

	//
	// Criterion value contained by this {@link Attribute}.
	//
	Value *Criterion `xml:"value,omitempty"`
}

type CurrencyCodeError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/o/v201802 CurrencyCodeError"`

	*ApiError

	//
	// The error reason represented by an enum.
	//
	Reason *CurrencyCodeErrorReason `xml:"reason,omitempty"`
}

type DoubleAttribute struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/o/v201802 DoubleAttribute"`

	*Attribute

	//
	// Double value contained by this {@link Attribute}.
	//
	Value float64 `xml:"value,omitempty"`
}

type IdeaTextFilterSearchParameter struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/o/v201802 IdeaTextFilterSearchParameter"`

	*SearchParameter

	//
	// A set of strings specifying which ideas should be included in the results.
	// <span class="constraint ContentsDistinct">This field must contain distinct elements.</span>
	// <span class="constraint ContentsNotNull">This field must not contain {@code null} elements.</span>
	// <span class="constraint ContentsStringLength">Strings in this field must be non-empty (trimmed).</span>
	//
	Included []string `xml:"included,omitempty"`

	//
	// A set of strings specifying which ideas should be excluded from the results.
	// <span class="constraint ContentsDistinct">This field must contain distinct elements.</span>
	// <span class="constraint ContentsNotNull">This field must not contain {@code null} elements.</span>
	// <span class="constraint ContentsStringLength">Strings in this field must be non-empty (trimmed).</span>
	//
	Excluded []string `xml:"excluded,omitempty"`
}

type IdeaTypeAttribute struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/o/v201802 IdeaTypeAttribute"`

	*Attribute

	//
	// {@link IdeaType} value contained by this {@link Attribute}.
	//
	Value *IdeaType `xml:"value,omitempty"`
}

type IncludeAdultContentSearchParameter struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/o/v201802 IncludeAdultContentSearchParameter"`

	*SearchParameter
}

type IntegerAttribute struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/o/v201802 IntegerAttribute"`

	*Attribute

	//
	// Integer value contained by this {@link Attribute}.
	//
	Value int32 `xml:"value,omitempty"`
}

type IntegerSetAttribute struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/o/v201802 IntegerSetAttribute"`

	*Attribute

	//
	// Set of integer values contained by this {@link Attribute}.
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	//
	Value []int32 `xml:"value,omitempty"`
}

type KeywordAttribute struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/o/v201802 KeywordAttribute"`

	*Attribute

	//
	// {@link Keyword} value contained by this {@link Attribute}.
	//
	Value *Keyword `xml:"value,omitempty"`
}

type LanguageSearchParameter struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/o/v201802 LanguageSearchParameter"`

	*SearchParameter

	//
	// A list of {@link Language}s indicating the desired languages being targeted in the results.
	// <span class="constraint ContentsDistinct">This field must contain distinct elements.</span>
	// <span class="constraint ContentsNotNull">This field must not contain {@code null} elements.</span>
	// <span class="constraint NotEmpty">This field must contain at least one element.</span>
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	//
	Languages []*Language `xml:"languages,omitempty"`
}

type LocationSearchParameter struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/o/v201802 LocationSearchParameter"`

	*SearchParameter

	//
	// A list of {@link Location}s indicating the desired locations (e.g countries) being targeted
	// in the results.
	// <span class="constraint ContentsDistinct">This field must contain distinct elements.</span>
	// <span class="constraint ContentsNotNull">This field must not contain {@code null} elements.</span>
	// <span class="constraint NotEmpty">This field must contain at least one element.</span>
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	//
	Locations []*Location `xml:"locations,omitempty"`
}

type LongAttribute struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/o/v201802 LongAttribute"`

	*Attribute

	//
	// Long value contained by this {@link Attribute}.
	//
	Value int64 `xml:"value,omitempty"`
}

type LongComparisonOperation struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/o/v201802 LongComparisonOperation"`

	//
	// The minimum value of elements returned by this operation (inclusive).
	//
	Minimum int64 `xml:"minimum,omitempty"`

	//
	// The maximum value of elements returned by this operation (inclusive).
	//
	Maximum int64 `xml:"maximum,omitempty"`
}

type LongRangeAttribute struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/o/v201802 LongRangeAttribute"`

	*Attribute

	//
	// {@link Range} of {@link LongValue} values contained by this
	// {@link Attribute}.
	//
	Value *Range `xml:"value,omitempty"`
}

type MoneyAttribute struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/o/v201802 MoneyAttribute"`

	*Attribute

	//
	// {@link Money} value contained by this {@link Attribute}.
	//
	Value *Money `xml:"value,omitempty"`
}

type MonthlySearchVolume struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/o/v201802 MonthlySearchVolume"`

	//
	// The year this search volume occurred in. (i.e. 2009)
	//
	Year int32 `xml:"year,omitempty"`

	//
	// The month this search volume occurred in. Month is 1 indexed,
	// such that 1=January and 12=December.
	//
	Month int32 `xml:"month,omitempty"`

	//
	// The approximate number of searches in this year/month. A {@code null} count
	// means that data is unavailable or unknown.
	//
	Count int64 `xml:"count,omitempty"`
}

type MonthlySearchVolumeAttribute struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/o/v201802 MonthlySearchVolumeAttribute"`

	*Attribute

	//
	// List of {@link MonthlySearchVolume} values contained by this
	// {@link Attribute}. The list contains the data for the past 12 months
	// (excluding the current month) in sorted order started with the most recent
	// month.
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	//
	Value []*MonthlySearchVolume `xml:"value,omitempty"`
}

type NetworkSearchParameter struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/o/v201802 NetworkSearchParameter"`

	*SearchParameter

	//
	// The network targeted for this search.
	//
	// <p>Currently we can support two options:
	// <ul>
	// <li>number of google search impressions</li>
	// <li>number of search impressions on the google search network(AFS)</li>
	// </ul>
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	//
	NetworkSetting *NetworkSetting `xml:"networkSetting,omitempty"`
}

type Range struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/o/v201802 Range"`

	Min *ComparableValue `xml:"min,omitempty"`

	Max *ComparableValue `xml:"max,omitempty"`
}

type RelatedToQuerySearchParameter struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/o/v201802 RelatedToQuerySearchParameter"`

	*SearchParameter

	//
	// A list of exact keyword match query {@link String}s that the search result
	// should be related to.
	// <span class="constraint ContentsDistinct">This field must contain distinct elements.</span>
	// <span class="constraint ContentsNotNull">This field must not contain {@code null} elements.</span>
	// <span class="constraint NotEmpty">This field must contain at least one element.</span>
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	//
	Queries []string `xml:"queries,omitempty"`
}

type RelatedToUrlSearchParameter struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/o/v201802 RelatedToUrlSearchParameter"`

	*SearchParameter

	//
	// A set of URL strings to which search results should be related.
	// For {@code KEYWORD} queries, only one URL may be submitted.
	// <span class="constraint ContentsDistinct">This field must contain distinct elements.</span>
	// <span class="constraint ContentsNotNull">This field must not contain {@code null} elements.</span>
	// <span class="constraint NotEmpty">This field must contain at least one element.</span>
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	//
	Urls []string `xml:"urls,omitempty"`

	//
	// Whether to crawl links off of the {@code urls} for the same domain.
	// Default is {@code false}.
	//
	IncludeSubUrls bool `xml:"includeSubUrls,omitempty"`
}

type SearchParameter struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/o/v201802 SearchParameter"`

	//
	// Indicates that this instance is a subtype of SearchParameter.
	// Although this field is returned in the response, it is ignored on input
	// and cannot be selected. Specify xsi:type instead.
	//
	SearchParameterType string `xml:"SearchParameter.Type,omitempty"`
}

type SearchVolumeSearchParameter struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/o/v201802 SearchVolumeSearchParameter"`

	*SearchParameter

	//
	// Used to specify the range of search volume.
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	//
	Operation *LongComparisonOperation `xml:"operation,omitempty"`
}

type SeedAdGroupIdSearchParameter struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/o/v201802 SeedAdGroupIdSearchParameter"`

	*SearchParameter

	//
	// The id for the ad group that should be used as a seed for generating new ideas.
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	//
	AdGroupId int64 `xml:"adGroupId,omitempty"`
}

type StringAttribute struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/o/v201802 StringAttribute"`

	*Attribute

	//
	// String value contained by this {@link Attribute}.
	//
	Value string `xml:"value,omitempty"`
}

type TargetingIdea struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/o/v201802 TargetingIdea"`

	//
	// Map of {@link AttributeType} to {@link Attribute}. Stores all data retrieved for each key
	// {@code AttributeType}.
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	//
	Data []*Type_AttributeMapEntry `xml:"data,omitempty"`
}

type TargetingIdeaError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/o/v201802 TargetingIdeaError"`

	*ApiError

	//
	// The error reason represented by an enum.
	//
	Reason *TargetingIdeaErrorReason `xml:"reason,omitempty"`
}

type TargetingIdeaPage struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/o/v201802 TargetingIdeaPage"`

	//
	// Total number of entries that can be retrieved using the get method.
	//
	TotalNumEntries int32 `xml:"totalNumEntries,omitempty"`

	//
	// The result entries in this page, as list of {@link TargetingIdea}s.
	//
	Entries []*TargetingIdea `xml:"entries,omitempty"`
}

type TargetingIdeaSelector struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/o/v201802 TargetingIdeaSelector"`

	//
	// Search for targeting ideas based on these search rules.
	//
	// <p>An empty set indicates this selector is valid for selecting metadata
	// with default parameters.
	// <span class="constraint ContentsNotNull">This field must not contain {@code null} elements.</span>
	// <span class="constraint DistinctTypes">Elements in this field must have distinct types.</span>
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	//
	SearchParameters []*SearchParameter `xml:"searchParameters,omitempty"`

	//
	// Limits the request to responses of this {@link IdeaType}, e.g. {@code KEYWORDS}.
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	//
	IdeaType *IdeaType `xml:"ideaType,omitempty"`

	//
	// Specifies the {@link RequestType}, e.g. {@code IDEAS} or {@code STATS}.
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	//
	RequestType *RequestType `xml:"requestType,omitempty"`

	//
	// Request {@link Attribute}s and associated data for this set of {@link Type}s.
	//
	// <p>An empty set indicates a request for {@link KeywordAttribute}, {@link PlacementAttribute},
	// and {@link IdeaType}.
	// <span class="constraint ContentsDistinct">This field must contain distinct elements.</span>
	// <span class="constraint ContentsNotNull">This field must not contain {@code null} elements.</span>
	//
	RequestedAttributeTypes []*AttributeType `xml:"requestedAttributeTypes,omitempty"`

	//
	// A {@link Paging} object that specifies the desired starting index and
	// number of results to return.
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	//
	Paging *Paging `xml:"paging,omitempty"`

	//
	// The locale code (for example "en_US") used for localizing strings,
	// controlling numeric formatting, and the like.  See RFC 3066 for
	// information on the format used.
	//
	LocaleCode string `xml:"localeCode,omitempty"`

	//
	// The currency code to be used for all monetary values returned in results in
	// ISO format (see
	// https://developers.google.com/adwords/api/docs/appendix/currencycodes
	// for supported currencies). The default is "USD" (US dollars).
	//
	CurrencyCode string `xml:"currencyCode,omitempty"`
}

type TrafficEstimatorError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/o/v201802 TrafficEstimatorError"`

	*ApiError

	Reason *TrafficEstimatorErrorReason `xml:"reason,omitempty"`
}

type Type_AttributeMapEntry struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/o/v201802 Type_AttributeMapEntry"`

	Key *AttributeType `xml:"key,omitempty"`

	Value *Attribute `xml:"value,omitempty"`
}

type WebpageDescriptor struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/o/v201802 WebpageDescriptor"`

	//
	// The URL of the webpage.
	//
	Url string `xml:"url,omitempty"`

	//
	// The title of the webpage.
	//
	Title string `xml:"title,omitempty"`
}

type WebpageDescriptorAttribute struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/o/v201802 WebpageDescriptorAttribute"`

	*Attribute

	//
	// {@link WebpageDescriptor} value contained by this {@link Attribute}.
	//
	Value *WebpageDescriptor `xml:"value,omitempty"`
}

type TargetingIdeaServiceInterface struct {
	client *SOAPClient
}

func NewTargetingIdeaServiceInterface(url string, tls bool, auth *BasicAuth) *TargetingIdeaServiceInterface {
	if url == "" {
		url = ""
	}
	client := NewSOAPClient(url, tls, auth)

	return &TargetingIdeaServiceInterface{
		client: client,
	}
}

func NewTargetingIdeaServiceInterfaceWithTLSConfig(url string, tlsCfg *tls.Config, auth *BasicAuth) *TargetingIdeaServiceInterface {
	if url == "" {
		url = ""
	}
	client := NewSOAPClientWithTLSConfig(url, tlsCfg, auth)

	return &TargetingIdeaServiceInterface{
		client: client,
	}
}

func (service *TargetingIdeaServiceInterface) AddHeader(header interface{}) {
	service.client.AddHeader(header)
}

// Backwards-compatible function: use AddHeader instead
func (service *TargetingIdeaServiceInterface) SetHeader(header interface{}) {
	service.client.AddHeader(header)
}

// Error can be either of the following types:
//
//   - ApiException
/*
   Returns a page of ideas that match the query described by the specified
   {@link TargetingIdeaSelector}.

   <p>The selector must specify a {@code paging} value, with {@code numberResults} set to 700 or
   less.  Large result sets must be composed through multiple calls to this method, advancing the
   paging {@code startIndex} value by {@code numberResults} with each call.</p>

   @param selector Query describing the types of results to return when
   finding matches (similar keyword ideas).
   @return A {@link TargetingIdeaPage} of results, that is a subset of the
   list of possible ideas meeting the criteria of the
   {@link TargetingIdeaSelector}.
   @throws ApiException If problems occurred while querying for ideas.
*/
func (service *TargetingIdeaServiceInterface) Get(request *Get) (*GetResponse, error) {
	response := new(GetResponse)
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
