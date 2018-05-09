package ReportDefinitionService

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
// The single reason for the whitelist error.
//
type NotWhitelistedErrorReason string

const (

	//
	// Customer is not whitelisted for accessing the API.
	//
	NotWhitelistedErrorReasonCUSTOMER_NOT_WHITELISTED_FOR_API NotWhitelistedErrorReason = "CUSTOMER_NOT_WHITELISTED_FOR_API"
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
// The reasons for the target error.
//
type RejectedErrorReason string

const (

	//
	// Unknown value encountered
	//
	RejectedErrorReasonUNKNOWN_VALUE RejectedErrorReason = "UNKNOWN_VALUE"
)

//
// Enums for report types.
//
type ReportDefinitionReportType string

const (

	//
	// Reports performance data for your keywords.
	//
	ReportDefinitionReportTypeKEYWORDS_PERFORMANCE_REPORT ReportDefinitionReportType = "KEYWORDS_PERFORMANCE_REPORT"

	//
	// Reports performance data for your ads.
	//
	ReportDefinitionReportTypeAD_PERFORMANCE_REPORT ReportDefinitionReportType = "AD_PERFORMANCE_REPORT"

	//
	// Reports performance data for URLs which triggered your ad and
	// received clicks.
	//
	ReportDefinitionReportTypeURL_PERFORMANCE_REPORT ReportDefinitionReportType = "URL_PERFORMANCE_REPORT"

	//
	// Reports ad group performance data for one or more of your campaigns.
	//
	ReportDefinitionReportTypeADGROUP_PERFORMANCE_REPORT ReportDefinitionReportType = "ADGROUP_PERFORMANCE_REPORT"

	//
	// Reports performance data for your campaigns.
	//
	ReportDefinitionReportTypeCAMPAIGN_PERFORMANCE_REPORT ReportDefinitionReportType = "CAMPAIGN_PERFORMANCE_REPORT"

	//
	// Reports performance data for your entire account.
	//
	ReportDefinitionReportTypeACCOUNT_PERFORMANCE_REPORT ReportDefinitionReportType = "ACCOUNT_PERFORMANCE_REPORT"

	//
	// Reports performance data by geographic origin.
	//
	ReportDefinitionReportTypeGEO_PERFORMANCE_REPORT ReportDefinitionReportType = "GEO_PERFORMANCE_REPORT"

	//
	// Reports performance data for search queries which triggered your ad and
	// received clicks.
	//
	ReportDefinitionReportTypeSEARCH_QUERY_PERFORMANCE_REPORT ReportDefinitionReportType = "SEARCH_QUERY_PERFORMANCE_REPORT"

	//
	// Reports performance data for automatic placements on the content network.
	//
	ReportDefinitionReportTypeAUTOMATIC_PLACEMENTS_PERFORMANCE_REPORT ReportDefinitionReportType = "AUTOMATIC_PLACEMENTS_PERFORMANCE_REPORT"

	//
	// Reports performance data for negative keywords at the campaign level.
	//
	ReportDefinitionReportTypeCAMPAIGN_NEGATIVE_KEYWORDS_PERFORMANCE_REPORT ReportDefinitionReportType = "CAMPAIGN_NEGATIVE_KEYWORDS_PERFORMANCE_REPORT"

	//
	// Reports performance data for the negative placements at the campaign
	// level.
	//
	ReportDefinitionReportTypeCAMPAIGN_NEGATIVE_PLACEMENTS_PERFORMANCE_REPORT ReportDefinitionReportType = "CAMPAIGN_NEGATIVE_PLACEMENTS_PERFORMANCE_REPORT"

	//
	// Reports performance data for destination urls.
	//
	ReportDefinitionReportTypeDESTINATION_URL_REPORT ReportDefinitionReportType = "DESTINATION_URL_REPORT"

	//
	// Reports data for shared sets.
	//
	ReportDefinitionReportTypeSHARED_SET_REPORT ReportDefinitionReportType = "SHARED_SET_REPORT"

	//
	// Reports data for campaign shared sets.
	//
	ReportDefinitionReportTypeCAMPAIGN_SHARED_SET_REPORT ReportDefinitionReportType = "CAMPAIGN_SHARED_SET_REPORT"

	//
	// Provides a downloadable snapshot of shared set criteria.
	//
	ReportDefinitionReportTypeSHARED_SET_CRITERIA_REPORT ReportDefinitionReportType = "SHARED_SET_CRITERIA_REPORT"

	//
	// Reports performance data for creative conversions (e.g. free clicks).
	//
	ReportDefinitionReportTypeCREATIVE_CONVERSION_REPORT ReportDefinitionReportType = "CREATIVE_CONVERSION_REPORT"

	//
	// Reports per-phone-call details for calls tracked using call metrics.
	//
	ReportDefinitionReportTypeCALL_METRICS_CALL_DETAILS_REPORT ReportDefinitionReportType = "CALL_METRICS_CALL_DETAILS_REPORT"

	//
	// Reports performance data for keywordless ads.
	//
	ReportDefinitionReportTypeKEYWORDLESS_QUERY_REPORT ReportDefinitionReportType = "KEYWORDLESS_QUERY_REPORT"

	//
	// Reports performance data for keywordless ads.
	//
	ReportDefinitionReportTypeKEYWORDLESS_CATEGORY_REPORT ReportDefinitionReportType = "KEYWORDLESS_CATEGORY_REPORT"

	//
	// Reports performance data for all published criteria types including keywords,
	// placements, topics, user-lists in a single report.
	//
	ReportDefinitionReportTypeCRITERIA_PERFORMANCE_REPORT ReportDefinitionReportType = "CRITERIA_PERFORMANCE_REPORT"

	//
	// Reports performance data for clicks.
	//
	ReportDefinitionReportTypeCLICK_PERFORMANCE_REPORT ReportDefinitionReportType = "CLICK_PERFORMANCE_REPORT"

	//
	// Reports performance data for budgets.
	//
	ReportDefinitionReportTypeBUDGET_PERFORMANCE_REPORT ReportDefinitionReportType = "BUDGET_PERFORMANCE_REPORT"

	//
	// Reports performance data for your (shared) bid strategies.
	//
	ReportDefinitionReportTypeBID_GOAL_PERFORMANCE_REPORT ReportDefinitionReportType = "BID_GOAL_PERFORMANCE_REPORT"

	//
	// Reports performance data for your display keywords.
	//
	ReportDefinitionReportTypeDISPLAY_KEYWORD_PERFORMANCE_REPORT ReportDefinitionReportType = "DISPLAY_KEYWORD_PERFORMANCE_REPORT"

	//
	// Reports performance data for your placeholder feed items
	//
	ReportDefinitionReportTypePLACEHOLDER_FEED_ITEM_REPORT ReportDefinitionReportType = "PLACEHOLDER_FEED_ITEM_REPORT"

	//
	// Reports performance data for your placements.
	//
	ReportDefinitionReportTypePLACEMENT_PERFORMANCE_REPORT ReportDefinitionReportType = "PLACEMENT_PERFORMANCE_REPORT"

	//
	// Reports performance data for negative location targets at campaign level.
	//
	ReportDefinitionReportTypeCAMPAIGN_NEGATIVE_LOCATIONS_REPORT ReportDefinitionReportType = "CAMPAIGN_NEGATIVE_LOCATIONS_REPORT"

	//
	// Reports performance data for managed and automatic genders in a combined report.
	//
	ReportDefinitionReportTypeGENDER_PERFORMANCE_REPORT ReportDefinitionReportType = "GENDER_PERFORMANCE_REPORT"

	//
	// Reports performance data for managed and automatic age ranges in a combined report.
	//
	ReportDefinitionReportTypeAGE_RANGE_PERFORMANCE_REPORT ReportDefinitionReportType = "AGE_RANGE_PERFORMANCE_REPORT"

	//
	// Reports performance data for campaign level location targets.
	//
	ReportDefinitionReportTypeCAMPAIGN_LOCATION_TARGET_REPORT ReportDefinitionReportType = "CAMPAIGN_LOCATION_TARGET_REPORT"

	//
	// Reports performance data for campaign level ad schedule targets.
	//
	ReportDefinitionReportTypeCAMPAIGN_AD_SCHEDULE_TARGET_REPORT ReportDefinitionReportType = "CAMPAIGN_AD_SCHEDULE_TARGET_REPORT"

	//
	// Paid & organic report
	//
	ReportDefinitionReportTypePAID_ORGANIC_QUERY_REPORT ReportDefinitionReportType = "PAID_ORGANIC_QUERY_REPORT"

	//
	// Reports performance data for your audience criteria.
	//
	ReportDefinitionReportTypeAUDIENCE_PERFORMANCE_REPORT ReportDefinitionReportType = "AUDIENCE_PERFORMANCE_REPORT"

	//
	// Reports performance data for your topic criteria.
	//
	ReportDefinitionReportTypeDISPLAY_TOPICS_PERFORMANCE_REPORT ReportDefinitionReportType = "DISPLAY_TOPICS_PERFORMANCE_REPORT"

	//
	// Distance report
	//
	ReportDefinitionReportTypeUSER_AD_DISTANCE_REPORT ReportDefinitionReportType = "USER_AD_DISTANCE_REPORT"

	//
	// Performance data for shopping campaigns.
	//
	ReportDefinitionReportTypeSHOPPING_PERFORMANCE_REPORT ReportDefinitionReportType = "SHOPPING_PERFORMANCE_REPORT"

	//
	// Performance data for product partitions in shopping campaigns.
	//
	ReportDefinitionReportTypePRODUCT_PARTITION_REPORT ReportDefinitionReportType = "PRODUCT_PARTITION_REPORT"

	//
	// Reports performance data for managed and automatic parental statuses in a combined report.
	//
	ReportDefinitionReportTypePARENTAL_STATUS_PERFORMANCE_REPORT ReportDefinitionReportType = "PARENTAL_STATUS_PERFORMANCE_REPORT"

	//
	// Performance data for Extension placeholders
	//
	ReportDefinitionReportTypePLACEHOLDER_REPORT ReportDefinitionReportType = "PLACEHOLDER_REPORT"

	//
	// Reports performance of ad placeholders when instantiated with specific FeedItems.
	//
	ReportDefinitionReportTypeAD_CUSTOMIZERS_FEED_ITEM_REPORT ReportDefinitionReportType = "AD_CUSTOMIZERS_FEED_ITEM_REPORT"

	//
	// Reports stats and settings details for labels.
	//
	ReportDefinitionReportTypeLABEL_REPORT ReportDefinitionReportType = "LABEL_REPORT"

	//
	// Reports performance data for final urls.
	//
	ReportDefinitionReportTypeFINAL_URL_REPORT ReportDefinitionReportType = "FINAL_URL_REPORT"

	//
	// Video performance report.
	//
	ReportDefinitionReportTypeVIDEO_PERFORMANCE_REPORT ReportDefinitionReportType = "VIDEO_PERFORMANCE_REPORT"

	//
	// Reports performance data for top content bid modifier criteria.
	//
	ReportDefinitionReportTypeTOP_CONTENT_PERFORMANCE_REPORT ReportDefinitionReportType = "TOP_CONTENT_PERFORMANCE_REPORT"

	//
	// Report to show campaign criteria structure.
	//
	ReportDefinitionReportTypeCAMPAIGN_CRITERIA_REPORT ReportDefinitionReportType = "CAMPAIGN_CRITERIA_REPORT"

	//
	// Report performance data for Campaign Groups.
	//
	ReportDefinitionReportTypeCAMPAIGN_GROUP_PERFORMANCE_REPORT ReportDefinitionReportType = "CAMPAIGN_GROUP_PERFORMANCE_REPORT"

	//
	// Reports performance data for landing pages on unexpanded and expanded final url levels.
	//
	ReportDefinitionReportTypeLANDING_PAGE_REPORT ReportDefinitionReportType = "LANDING_PAGE_REPORT"

	//
	// Report performance data for Marketplace Ads Clients.
	//
	ReportDefinitionReportTypeMARKETPLACE_PERFORMANCE_REPORT ReportDefinitionReportType = "MARKETPLACE_PERFORMANCE_REPORT"

	//
	// <span class="constraint Rejected">Used for return value only. An enumeration could not be processed, typically due to incompatibility with your WSDL version.</span>
	//
	ReportDefinitionReportTypeUNKNOWN ReportDefinitionReportType = "UNKNOWN"
)

//
// Enums for all the reasons an error can be thrown to the user during
// a {@link ReportDefinitionService#mutate(java.util.List)} operation.
//
type ReportDefinitionErrorReason string

const (

	//
	// Customer passed in invalid date range for a report type.
	//
	ReportDefinitionErrorReasonINVALID_DATE_RANGE_FOR_REPORT ReportDefinitionErrorReason = "INVALID_DATE_RANGE_FOR_REPORT"

	//
	// Customer passed in invalid field name for a report type
	//
	ReportDefinitionErrorReasonINVALID_FIELD_NAME_FOR_REPORT ReportDefinitionErrorReason = "INVALID_FIELD_NAME_FOR_REPORT"

	//
	// Unable to locate a field mapping for this report type.
	//
	ReportDefinitionErrorReasonUNABLE_TO_FIND_MAPPING_FOR_THIS_REPORT ReportDefinitionErrorReason = "UNABLE_TO_FIND_MAPPING_FOR_THIS_REPORT"

	//
	// Customer passed in invalid column name for a report type
	//
	ReportDefinitionErrorReasonINVALID_COLUMN_NAME_FOR_REPORT ReportDefinitionErrorReason = "INVALID_COLUMN_NAME_FOR_REPORT"

	//
	// Customer passed in invalid report definition id.
	//
	ReportDefinitionErrorReasonINVALID_REPORT_DEFINITION_ID ReportDefinitionErrorReason = "INVALID_REPORT_DEFINITION_ID"

	//
	// Report selector cannot be null.
	//
	ReportDefinitionErrorReasonREPORT_SELECTOR_CANNOT_BE_NULL ReportDefinitionErrorReason = "REPORT_SELECTOR_CANNOT_BE_NULL"

	//
	// No Enums exist for this column name.
	//
	ReportDefinitionErrorReasonNO_ENUMS_FOR_THIS_COLUMN_NAME ReportDefinitionErrorReason = "NO_ENUMS_FOR_THIS_COLUMN_NAME"

	//
	// Invalid view name.
	//
	ReportDefinitionErrorReasonINVALID_VIEW ReportDefinitionErrorReason = "INVALID_VIEW"

	//
	// Sorting is not supported for reports.
	//
	ReportDefinitionErrorReasonSORTING_NOT_SUPPORTED ReportDefinitionErrorReason = "SORTING_NOT_SUPPORTED"

	//
	// Paging is not supported for reports.
	//
	ReportDefinitionErrorReasonPAGING_NOT_SUPPORTED ReportDefinitionErrorReason = "PAGING_NOT_SUPPORTED"

	//
	// Customer can not create report of a selected type.
	//
	ReportDefinitionErrorReasonCUSTOMER_SERVING_TYPE_REPORT_MISMATCH ReportDefinitionErrorReason = "CUSTOMER_SERVING_TYPE_REPORT_MISMATCH"

	//
	// Cross client report has an client selector without any valid identifier
	// for a customer.
	//
	ReportDefinitionErrorReasonCLIENT_SELECTOR_NO_CUSTOMER_IDENTIFIER ReportDefinitionErrorReason = "CLIENT_SELECTOR_NO_CUSTOMER_IDENTIFIER"

	//
	// Cross client report has an invalid external customer ID in the client
	// selector.
	//
	ReportDefinitionErrorReasonCLIENT_SELECTOR_INVALID_CUSTOMER_ID ReportDefinitionErrorReason = "CLIENT_SELECTOR_INVALID_CUSTOMER_ID"

	ReportDefinitionErrorReasonREPORT_DEFINITION_ERROR ReportDefinitionErrorReason = "REPORT_DEFINITION_ERROR"
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

type GetReportFields struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 getReportFields"`

	//
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	//
	ReportType *ReportDefinitionReportType `xml:"reportType,omitempty"`
}

type GetReportFieldsResponse struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 getReportFieldsResponse"`

	Rval []*ReportDefinitionField `xml:"rval,omitempty"`
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

type EnumValuePair struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 EnumValuePair"`

	//
	// The api enum value.
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	EnumValue string `xml:"enumValue,omitempty"`

	//
	// The enum value displayed in the downloaded report.
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	EnumDisplayValue string `xml:"enumDisplayValue,omitempty"`
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

type NotEmptyError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 NotEmptyError"`

	*ApiError

	//
	// The error reason represented by an enum.
	//
	Reason *NotEmptyErrorReason `xml:"reason,omitempty"`
}

type NotWhitelistedError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 NotWhitelistedError"`

	*ApiError

	//
	// The error reason represented by an enum.
	//
	Reason *NotWhitelistedErrorReason `xml:"reason,omitempty"`
}

type NullError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 NullError"`

	*ApiError

	//
	// The error reason represented by an enum.
	//
	Reason *NullErrorReason `xml:"reason,omitempty"`
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

type ReportDefinitionError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 ReportDefinitionError"`

	*ApiError

	Reason *ReportDefinitionErrorReason `xml:"reason,omitempty"`
}

type ReportDefinitionField struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 ReportDefinitionField"`

	//
	// The field name.
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	FieldName string `xml:"fieldName,omitempty"`

	//
	// The name that is displayed in the downloaded report.
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	DisplayFieldName string `xml:"displayFieldName,omitempty"`

	//
	// The XML attribute in the downloaded report.
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	XmlAttributeName string `xml:"xmlAttributeName,omitempty"`

	//
	// The type of field. Useful for knowing what operation type to pass in for
	// a given field in a predicate.
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	FieldType string `xml:"fieldType,omitempty"`

	//
	// The behavior of this field. Possible values are: "ATTRIBUTE", "METRIC"
	// and "SEGMENT".
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	FieldBehavior string `xml:"fieldBehavior,omitempty"`

	//
	// List of enum values for the corresponding field if and only if the
	// field is an enum type.
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	EnumValues []string `xml:"enumValues,omitempty"`

	//
	// Indicates if the user can select this field.
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	CanSelect bool `xml:"canSelect,omitempty"`

	//
	// Indicates if the user can filter on this field.
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	CanFilter bool `xml:"canFilter,omitempty"`

	//
	// Indicates that the field is an enum type.
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	IsEnumType bool `xml:"isEnumType,omitempty"`

	//
	// Indicates that the field is only available with beta access.
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	IsBeta bool `xml:"isBeta,omitempty"`

	//
	// Indicates if the field can be selected in queries that explicitly request zero rows.
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	IsZeroRowCompatible bool `xml:"isZeroRowCompatible,omitempty"`

	//
	// List of enum values in api to their corresponding display values in the
	// downloaded reports.
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	EnumValuePairs []*EnumValuePair `xml:"enumValuePairs,omitempty"`

	//
	// List of mutually exclusive fields of this field. This field cannot be selected or used in
	// a predicate together with any of the mutually exclusive fields in this list.
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	ExclusiveFields []string `xml:"exclusiveFields,omitempty"`
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

type ReportDefinitionServiceInterface struct {
	client *SOAPClient
}

func NewReportDefinitionServiceInterface(url string, tls bool, auth *BasicAuth) *ReportDefinitionServiceInterface {
	if url == "" {
		url = ""
	}
	client := NewSOAPClient(url, tls, auth)

	return &ReportDefinitionServiceInterface{
		client: client,
	}
}

func NewReportDefinitionServiceInterfaceWithTLSConfig(url string, tlsCfg *tls.Config, auth *BasicAuth) *ReportDefinitionServiceInterface {
	if url == "" {
		url = ""
	}
	client := NewSOAPClientWithTLSConfig(url, tlsCfg, auth)

	return &ReportDefinitionServiceInterface{
		client: client,
	}
}

func (service *ReportDefinitionServiceInterface) AddHeader(header interface{}) {
	service.client.AddHeader(header)
}

// Backwards-compatible function: use AddHeader instead
func (service *ReportDefinitionServiceInterface) SetHeader(header interface{}) {
	service.client.AddHeader(header)
}

// Error can be either of the following types:
//
//   - ApiException
/*
   Returns the available report fields for a given report type.
   When using this method the {@code clientCustomerId} header field is
   optional. Callers are discouraged from setting the clientCustomerId
   header field in calls to this method as its presence will trigger an
   authorization error if the caller does not have access to the customer
   with the included ID.

   @param reportType The type of report.
   @return The list of available report fields. Each
   {@link ReportDefinitionField} encapsulates the field name, the
   field data type, and the enum values (if the field's type is
   {@code enum}).
   @throws ApiException if a problem occurred while fetching the
   ReportDefinitionField information.
*/
func (service *ReportDefinitionServiceInterface) GetReportFields(request *GetReportFields) (*GetReportFieldsResponse, error) {
	response := new(GetReportFieldsResponse)
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
