package ConstantDataService

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
// An enumeration of possible user interest taxonomy types.
//
type ConstantDataServiceUserInterestTaxonomyType string

const (

	//
	// The brand for this user interest.
	//
	ConstantDataServiceUserInterestTaxonomyTypeBRAND ConstantDataServiceUserInterestTaxonomyType = "BRAND"

	//
	// The market for this user interest.
	//
	ConstantDataServiceUserInterestTaxonomyTypeIN_MARKET ConstantDataServiceUserInterestTaxonomyType = "IN_MARKET"

	//
	// Users known to have installed applications in the specified categories.
	//
	ConstantDataServiceUserInterestTaxonomyTypeMOBILE_APP_INSTALL_USER ConstantDataServiceUserInterestTaxonomyType = "MOBILE_APP_INSTALL_USER"

	//
	// The geographical location of the interest-based vertical.
	//
	ConstantDataServiceUserInterestTaxonomyTypeVERTICAL_GEO ConstantDataServiceUserInterestTaxonomyType = "VERTICAL_GEO"

	//
	// User interest criteria for new smart phone users.
	//
	ConstantDataServiceUserInterestTaxonomyTypeNEW_SMART_PHONE_USER ConstantDataServiceUserInterestTaxonomyType = "NEW_SMART_PHONE_USER"
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

type MobileDeviceDeviceType string

const (
	MobileDeviceDeviceTypeDEVICE_TYPE_MOBILE MobileDeviceDeviceType = "DEVICE_TYPE_MOBILE"

	MobileDeviceDeviceTypeDEVICE_TYPE_TABLET MobileDeviceDeviceType = "DEVICE_TYPE_TABLET"
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
// The operator type.
//
type OperatingSystemVersionOperatorType string

const (
	OperatingSystemVersionOperatorTypeGREATER_THAN_EQUAL_TO OperatingSystemVersionOperatorType = "GREATER_THAN_EQUAL_TO"

	OperatingSystemVersionOperatorTypeEQUAL_TO OperatingSystemVersionOperatorType = "EQUAL_TO"

	OperatingSystemVersionOperatorTypeUNKNOWN OperatingSystemVersionOperatorType = "UNKNOWN"
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
// Status of a bidding dimension (category) in a bidding taxonomy.
//
type ShoppingBiddingDimensionStatus string

const (

	//
	// Default status. Should not be used.
	// <span class="constraint Rejected">Used for return value only. An enumeration could not be processed, typically due to incompatibility with your WSDL version.</span>
	//
	ShoppingBiddingDimensionStatusUNKNOWN ShoppingBiddingDimensionStatus = "UNKNOWN"

	//
	// The dimension is active and it can be used for bidding.
	//
	ShoppingBiddingDimensionStatusACTIVE ShoppingBiddingDimensionStatus = "ACTIVE"

	//
	// The dimension is deprecated and should not be used for bidding.
	//
	ShoppingBiddingDimensionStatusOBSOLETE ShoppingBiddingDimensionStatus = "OBSOLETE"
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

type GetAgeRangeCriterion struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 getAgeRangeCriterion"`
}

type GetAgeRangeCriterionResponse struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 getAgeRangeCriterionResponse"`

	Rval []*AgeRange `xml:"rval,omitempty"`
}

type GetCarrierCriterion struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 getCarrierCriterion"`
}

type GetCarrierCriterionResponse struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 getCarrierCriterionResponse"`

	Rval []*Carrier `xml:"rval,omitempty"`
}

type GetGenderCriterion struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 getGenderCriterion"`
}

type GetGenderCriterionResponse struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 getGenderCriterionResponse"`

	Rval []*Gender `xml:"rval,omitempty"`
}

type GetLanguageCriterion struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 getLanguageCriterion"`
}

type GetLanguageCriterionResponse struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 getLanguageCriterionResponse"`

	Rval []*Language `xml:"rval,omitempty"`
}

type GetMobileAppCategoryCriterion struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 getMobileAppCategoryCriterion"`
}

type GetMobileAppCategoryCriterionResponse struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 getMobileAppCategoryCriterionResponse"`

	Rval []*MobileAppCategory `xml:"rval,omitempty"`
}

type GetMobileDeviceCriterion struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 getMobileDeviceCriterion"`
}

type GetMobileDeviceCriterionResponse struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 getMobileDeviceCriterionResponse"`

	Rval []*MobileDevice `xml:"rval,omitempty"`
}

type GetOperatingSystemVersionCriterion struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 getOperatingSystemVersionCriterion"`
}

type GetOperatingSystemVersionCriterionResponse struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 getOperatingSystemVersionCriterionResponse"`

	Rval []*OperatingSystemVersion `xml:"rval,omitempty"`
}

type GetProductBiddingCategoryData struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 getProductBiddingCategoryData"`

	//
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	//
	Selector *Selector `xml:"selector,omitempty"`
}

type GetProductBiddingCategoryDataResponse struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 getProductBiddingCategoryDataResponse"`

	Rval []*ProductBiddingCategoryData `xml:"rval,omitempty"`
}

type GetUserInterestCriterion struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 getUserInterestCriterion"`

	UserInterestTaxonomyType *ConstantDataServiceUserInterestTaxonomyType `xml:"userInterestTaxonomyType,omitempty"`
}

type GetUserInterestCriterionResponse struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 getUserInterestCriterionResponse"`

	Rval []*CriterionUserInterest `xml:"rval,omitempty"`
}

type GetVerticalCriterion struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 getVerticalCriterion"`
}

type GetVerticalCriterionResponse struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 getVerticalCriterionResponse"`

	Rval []*Vertical `xml:"rval,omitempty"`
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

type Carrier struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 Carrier"`

	*Criterion

	//
	// Name of the carrier.
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	Name string `xml:"name,omitempty"`

	//
	// Country code of the carrier.
	// Can be {@code null} if not applicable, e.g., for Carrier "Wifi".
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	CountryCode string `xml:"countryCode,omitempty"`
}

type ClientTermsError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 ClientTermsError"`

	*ApiError

	Reason *ClientTermsErrorReason `xml:"reason,omitempty"`
}

type ConstantData struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 ConstantData"`

	//
	// Indicates that this instance is a subtype of ConstantData.
	// Although this field is returned in the response, it is ignored on input
	// and cannot be selected. Specify xsi:type instead.
	//
	ConstantDataType string `xml:"ConstantData.Type,omitempty"`
}

type Criterion struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 Criterion"`

	//
	// ID of this criterion.
	// <span class="constraint Required">This field is required and should not be {@code null} when it is contained within {@link Operator}s : SET, REMOVE.</span>
	//
	Id int64 `xml:"id,omitempty"`

	//
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

type Gender struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 Gender"`

	*Criterion

	//
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
	// <span class="constraint MatchesRegex">Keyword text must not contain NUL (code point 0x0) characters. This is checked by the regular expression '[^\x00]*'.</span>
	// <span class="constraint Required">This field is required and should not be {@code null} when it is contained within {@link Operator}s : ADD.</span>
	//
	Text string `xml:"text,omitempty"`

	//
	// Match type of this keyword.
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

type MobileAppCategory struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 MobileAppCategory"`

	*Criterion

	//
	// ID of this mobile app category. A complete list of the available mobile app categories is
	// available <a href="/adwords/api/docs/appendix/mobileappcategories">here</a>.
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
	// <span class="constraint Required">This field is required and should not be {@code null} when it is contained within {@link Operator}s : ADD.</span>
	//
	AppId string `xml:"appId,omitempty"`

	//
	// Title of this mobile application.
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	DisplayName string `xml:"displayName,omitempty"`
}

type MobileDevice struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 MobileDevice"`

	*Criterion

	//
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	DeviceName string `xml:"deviceName,omitempty"`

	//
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	ManufacturerName string `xml:"manufacturerName,omitempty"`

	//
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	DeviceType *MobileDeviceDeviceType `xml:"deviceType,omitempty"`

	//
	// <span class="constraint Selectable">This field can be selected using the value "OperatingSystemName".</span>
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	OperatingSystemName string `xml:"operatingSystemName,omitempty"`
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

type OperatingSystemVersion struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 OperatingSystemVersion"`

	*Criterion

	//
	// The name of the operating system.
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	Name string `xml:"name,omitempty"`

	//
	// The OS Major Version number.
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	OsMajorVersion int32 `xml:"osMajorVersion,omitempty"`

	//
	// The OS Minor Version number.
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	OsMinorVersion int32 `xml:"osMinorVersion,omitempty"`

	//
	// The operator type.
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	OperatorType *OperatingSystemVersionOperatorType `xml:"operatorType,omitempty"`
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
	// <span class="constraint Required">This field is required and should not be {@code null} when it is contained within {@link Operator}s : ADD.</span>
	//
	Url string `xml:"url,omitempty"`
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

type ProductBiddingCategoryData struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 ProductBiddingCategoryData"`

	*ConstantData

	//
	// The dimension value that corresponds to this category.
	//
	DimensionValue *ProductBiddingCategory `xml:"dimensionValue,omitempty"`

	//
	// The dimension value that corresponds to parent category.
	//
	ParentDimensionValue *ProductBiddingCategory `xml:"parentDimensionValue,omitempty"`

	//
	// The country of the taxonomy. It applies to all categories from the taxonomy.
	// <span class="constraint Selectable">This field can be selected using the value "Country".</span><span class="constraint Filterable">This field can be filtered on.</span>
	//
	Country string `xml:"country,omitempty"`

	//
	// The status of the taxonomy. It applies to all categories from the taxonomy.
	// <span class="constraint Selectable">This field can be selected using the value "BiddingCategoryStatus".</span><span class="constraint Filterable">This field can be filtered on.</span>
	//
	Status *ShoppingBiddingDimensionStatus `xml:"status,omitempty"`

	//
	// A map of displayValues by their language code. The language code is a 2-letter string
	// conforming with ISO 639-1 standard.
	//
	DisplayValue []*String_StringMapEntry `xml:"displayValue,omitempty"`
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

type String_StringMapEntry struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 String_StringMapEntry"`

	Key string `xml:"key,omitempty"`

	Value string `xml:"value,omitempty"`
}

type UnknownProductDimension struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 UnknownProductDimension"`

	*ProductDimension
}

type CriterionUserInterest struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 CriterionUserInterest"`

	*Criterion

	//
	// Id of this user interest. This is a required field.
	//
	UserInterestId int64 `xml:"userInterestId,omitempty"`

	//
	// Parent Id of this user interest.
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	UserInterestParentId int64 `xml:"userInterestParentId,omitempty"`

	//
	// Name of this user interest.
	//
	UserInterestName string `xml:"userInterestName,omitempty"`
}

type CriterionUserList struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 CriterionUserList"`

	*Criterion

	//
	// Id of this user list. This is a required field.
	//
	UserListId int64 `xml:"userListId,omitempty"`

	UserListName string `xml:"userListName,omitempty"`

	UserListMembershipStatus *CriterionUserListMembershipStatus `xml:"userListMembershipStatus,omitempty"`

	//
	// Determines whether a user list is eligible for targeting in the google.com
	// (search) network.
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	UserListEligibleForSearch bool `xml:"userListEligibleForSearch,omitempty"`

	//
	// Determines whether a user list is eligible for targeting in the display network.
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	UserListEligibleForDisplay bool `xml:"userListEligibleForDisplay,omitempty"`
}

type Vertical struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 Vertical"`

	*Criterion

	//
	// Id of this vertical.
	//
	VerticalId int64 `xml:"verticalId,omitempty"`

	//
	// Id of the parent of this vertical.
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
	//
	Path []string `xml:"path,omitempty"`
}

type ConstantDataServiceInterface struct {
	client *SOAPClient
}

func NewConstantDataServiceInterface(url string, tls bool, auth *BasicAuth) *ConstantDataServiceInterface {
	if url == "" {
		url = ""
	}
	client := NewSOAPClient(url, tls, auth)

	return &ConstantDataServiceInterface{
		client: client,
	}
}

func NewConstantDataServiceInterfaceWithTLSConfig(url string, tlsCfg *tls.Config, auth *BasicAuth) *ConstantDataServiceInterface {
	if url == "" {
		url = ""
	}
	client := NewSOAPClientWithTLSConfig(url, tlsCfg, auth)

	return &ConstantDataServiceInterface{
		client: client,
	}
}

func (service *ConstantDataServiceInterface) AddHeader(header interface{}) {
	service.client.AddHeader(header)
}

// Backwards-compatible function: use AddHeader instead
func (service *ConstantDataServiceInterface) SetHeader(header interface{}) {
	service.client.AddHeader(header)
}

// Error can be either of the following types:
//
//   - ApiException
/*
   Returns a list of all age range criteria.

   @return A list of age ranges.
   @throws ApiException when there is at least one error with the request.
*/
func (service *ConstantDataServiceInterface) GetAgeRangeCriterion(request *GetAgeRangeCriterion) (*GetAgeRangeCriterionResponse, error) {
	response := new(GetAgeRangeCriterionResponse)
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
   Returns a list of all carrier criteria.

   @return A list of carriers.
   @throws ApiException when there is at least one error with the request.
*/
func (service *ConstantDataServiceInterface) GetCarrierCriterion(request *GetCarrierCriterion) (*GetCarrierCriterionResponse, error) {
	response := new(GetCarrierCriterionResponse)
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
   Returns a list of all gender criteria.

   @return A list of genders.
   @throws ApiException when there is at least one error with the request.
*/
func (service *ConstantDataServiceInterface) GetGenderCriterion(request *GetGenderCriterion) (*GetGenderCriterionResponse, error) {
	response := new(GetGenderCriterionResponse)
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
   Returns a list of all language criteria.

   @return A list of languages.
   @throws ApiException when there is at least one error with the request.
*/
func (service *ConstantDataServiceInterface) GetLanguageCriterion(request *GetLanguageCriterion) (*GetLanguageCriterionResponse, error) {
	response := new(GetLanguageCriterionResponse)
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
   Returns a list of all mobile app category criteria.

   @return A list of mobile app categories.
   @throws ApiException when there is at least one error with the request.
*/
func (service *ConstantDataServiceInterface) GetMobileAppCategoryCriterion(request *GetMobileAppCategoryCriterion) (*GetMobileAppCategoryCriterionResponse, error) {
	response := new(GetMobileAppCategoryCriterionResponse)
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
   Returns a list of all mobile devices.

   @return A list of mobile devices.
   @throws ApiException when there is at least one error with the request.
*/
func (service *ConstantDataServiceInterface) GetMobileDeviceCriterion(request *GetMobileDeviceCriterion) (*GetMobileDeviceCriterionResponse, error) {
	response := new(GetMobileDeviceCriterionResponse)
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
   Returns a list of all operating system version criteria.

   @return A list of operating system versions.
   @throws ApiException when there is at least one error with the request.
*/
func (service *ConstantDataServiceInterface) GetOperatingSystemVersionCriterion(request *GetOperatingSystemVersionCriterion) (*GetOperatingSystemVersionCriterionResponse, error) {
	response := new(GetOperatingSystemVersionCriterionResponse)
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
   Returns a list of shopping bidding categories.

   A country predicate must be included in the selector, only {@link Predicate.Operator#EQUALS}
   and {@link Predicate.Operator#IN} with a single value are supported in the country predicate.
   An empty parentDimensionType predicate will filter for root categories.

   @return A list of shopping bidding categories.
   @throws ApiException when there is at least one error with the request.
*/
func (service *ConstantDataServiceInterface) GetProductBiddingCategoryData(request *GetProductBiddingCategoryData) (*GetProductBiddingCategoryDataResponse, error) {
	response := new(GetProductBiddingCategoryDataResponse)
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
   Returns a list of user interests.

   @param userInterestTaxonomyType The type of taxonomy to use when requesting user interests.
   @return A list of user interests.
   @throws ApiException when there is at least one error with the request.
*/
func (service *ConstantDataServiceInterface) GetUserInterestCriterion(request *GetUserInterestCriterion) (*GetUserInterestCriterionResponse, error) {
	response := new(GetUserInterestCriterionResponse)
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
   Returns a list of content verticals.

   @return A list of verticals.
   @throws ApiException when there is at least one error with the request.
*/
func (service *ConstantDataServiceInterface) GetVerticalCriterion(request *GetVerticalCriterion) (*GetVerticalCriterionResponse, error) {
	response := new(GetVerticalCriterionResponse)
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
