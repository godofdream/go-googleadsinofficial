package AdwordsUserListService

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

type EntityNotFoundReason string

const (

	//
	// The specified id refered to an entity which either doesn't exist or is not accessible to the
	// customer. e.g. campaign belongs to another customer.
	//
	EntityNotFoundReasonINVALID_ID EntityNotFoundReason = "INVALID_ID"
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

//
// This indicates the way the entity such as UserList is related to a user.
//
type AccessReason string

const (

	//
	// The entity is owned by the user.
	//
	AccessReasonOWNED AccessReason = "OWNED"

	//
	// The entity is shared to the user.
	//
	AccessReasonSHARED AccessReason = "SHARED"

	//
	// The entity is licensed to the user.
	//
	AccessReasonLICENSED AccessReason = "LICENSED"

	//
	// The user subscribed to the entity.
	//
	AccessReasonSUBSCRIBED AccessReason = "SUBSCRIBED"
)

//
// Status in the AccountUserListStatus table. This indicates if the user list share or
// the licensing of the userlist is still active.
//
type AccountUserListStatus string

const (
	AccountUserListStatusACTIVE AccountUserListStatus = "ACTIVE"

	AccountUserListStatusINACTIVE AccountUserListStatus = "INACTIVE"
)

//
// Logical operator connecting two rules.
//
type CombinedRuleUserListRuleOperator string

const (

	//
	// <span class="constraint Rejected">Used for return value only. An enumeration could not be processed, typically due to incompatibility with your WSDL version.</span>
	//
	CombinedRuleUserListRuleOperatorUNKNOWN CombinedRuleUserListRuleOperator = "UNKNOWN"

	CombinedRuleUserListRuleOperatorAND CombinedRuleUserListRuleOperator = "AND"

	CombinedRuleUserListRuleOperatorAND_NOT CombinedRuleUserListRuleOperator = "AND_NOT"
)

//
// User can create only BOOMERANG_EVENT conversion types. For all other types
// UserListService service will return OTHER.
//
type UserListConversionTypeCategory string

const (
	UserListConversionTypeCategoryBOOMERANG_EVENT UserListConversionTypeCategory = "BOOMERANG_EVENT"

	UserListConversionTypeCategoryOTHER UserListConversionTypeCategory = "OTHER"
)

//
// Enum to indicate source of CRM upload data.
//
type CrmDataSourceType string

const (

	//
	// <span class="constraint Rejected">Used for return value only. An enumeration could not be processed, typically due to incompatibility with your WSDL version.</span>
	//
	CrmDataSourceTypeUNKNOWN CrmDataSourceType = "UNKNOWN"

	//
	// The uploaded data is first party data.
	//
	CrmDataSourceTypeFIRST_PARTY CrmDataSourceType = "FIRST_PARTY"

	//
	// The uploaded data is from third party credit bureau.
	//
	CrmDataSourceTypeTHIRD_PARTY_CREDIT_BUREAU CrmDataSourceType = "THIRD_PARTY_CREDIT_BUREAU"

	//
	// The uploaded data is from third party voter file.
	//
	CrmDataSourceTypeTHIRD_PARTY_VOTER_FILE CrmDataSourceType = "THIRD_PARTY_VOTER_FILE"
)

//
// Enum to indicate what type of data are the user list's members matched from.
//
type CustomerMatchUploadKeyType string

const (

	//
	// <span class="constraint Rejected">Used for return value only. An enumeration could not be processed, typically due to incompatibility with your WSDL version.</span>
	//
	CustomerMatchUploadKeyTypeUNKNOWN CustomerMatchUploadKeyType = "UNKNOWN"

	//
	// Members are matched from customer info such as email address, phone number or
	// physical address.
	//
	CustomerMatchUploadKeyTypeCONTACT_INFO CustomerMatchUploadKeyType = "CONTACT_INFO"

	//
	// Members are matched from advertiser generated and assigned user id.
	//
	CustomerMatchUploadKeyTypeCRM_ID CustomerMatchUploadKeyType = "CRM_ID"

	//
	// Members are matched from mobile advertising ids.
	//
	CustomerMatchUploadKeyTypeMOBILE_ADVERTISING_ID CustomerMatchUploadKeyType = "MOBILE_ADVERTISING_ID"
)

//
// Supported rule operator for date type.
//
type DateRuleItemDateOperator string

const (

	//
	// <span class="constraint Rejected">Used for return value only. An enumeration could not be processed, typically due to incompatibility with your WSDL version.</span>
	//
	DateRuleItemDateOperatorUNKNOWN DateRuleItemDateOperator = "UNKNOWN"

	DateRuleItemDateOperatorEQUALS DateRuleItemDateOperator = "EQUALS"

	DateRuleItemDateOperatorNOT_EQUAL DateRuleItemDateOperator = "NOT_EQUAL"

	DateRuleItemDateOperatorBEFORE DateRuleItemDateOperator = "BEFORE"

	DateRuleItemDateOperatorAFTER DateRuleItemDateOperator = "AFTER"
)

//
// Reasons
//
type MutateMembersErrorReason string

const (
	MutateMembersErrorReasonUNKNOWN MutateMembersErrorReason = "UNKNOWN"

	MutateMembersErrorReasonUNSUPPORTED_METHOD MutateMembersErrorReason = "UNSUPPORTED_METHOD"

	MutateMembersErrorReasonINVALID_USER_LIST_ID MutateMembersErrorReason = "INVALID_USER_LIST_ID"

	MutateMembersErrorReasonINVALID_USER_LIST_TYPE MutateMembersErrorReason = "INVALID_USER_LIST_TYPE"

	MutateMembersErrorReasonINVALID_DATA_TYPE MutateMembersErrorReason = "INVALID_DATA_TYPE"

	MutateMembersErrorReasonINVALID_SHA256_FORMAT MutateMembersErrorReason = "INVALID_SHA256_FORMAT"

	MutateMembersErrorReasonOPERATOR_CONFLICT_FOR_SAME_USER_LIST_ID MutateMembersErrorReason = "OPERATOR_CONFLICT_FOR_SAME_USER_LIST_ID"

	MutateMembersErrorReasonINVALID_REMOVEALL_OPERATION MutateMembersErrorReason = "INVALID_REMOVEALL_OPERATION"

	MutateMembersErrorReasonINVALID_OPERATION_ORDER MutateMembersErrorReason = "INVALID_OPERATION_ORDER"

	MutateMembersErrorReasonMISSING_MEMBER_IDENTIFIER MutateMembersErrorReason = "MISSING_MEMBER_IDENTIFIER"

	MutateMembersErrorReasonINCOMPATIBLE_UPLOAD_KEY_TYPE MutateMembersErrorReason = "INCOMPATIBLE_UPLOAD_KEY_TYPE"
)

//
// Supported operator for numbers.
//
type NumberRuleItemNumberOperator string

const (

	//
	// <span class="constraint Rejected">Used for return value only. An enumeration could not be processed, typically due to incompatibility with your WSDL version.</span>
	//
	NumberRuleItemNumberOperatorUNKNOWN NumberRuleItemNumberOperator = "UNKNOWN"

	NumberRuleItemNumberOperatorGREATER_THAN NumberRuleItemNumberOperator = "GREATER_THAN"

	NumberRuleItemNumberOperatorGREATER_THAN_OR_EQUAL NumberRuleItemNumberOperator = "GREATER_THAN_OR_EQUAL"

	NumberRuleItemNumberOperatorEQUALS NumberRuleItemNumberOperator = "EQUALS"

	NumberRuleItemNumberOperatorNOT_EQUAL NumberRuleItemNumberOperator = "NOT_EQUAL"

	NumberRuleItemNumberOperatorLESS_THAN NumberRuleItemNumberOperator = "LESS_THAN"

	NumberRuleItemNumberOperatorLESS_THAN_OR_EQUAL NumberRuleItemNumberOperator = "LESS_THAN_OR_EQUAL"
)

//
// The status of pre-population
//
type RuleBasedUserListPrepopulationStatus string

const (
	RuleBasedUserListPrepopulationStatusNONE RuleBasedUserListPrepopulationStatus = "NONE"

	RuleBasedUserListPrepopulationStatusREQUESTED RuleBasedUserListPrepopulationStatus = "REQUESTED"

	RuleBasedUserListPrepopulationStatusFINISHED RuleBasedUserListPrepopulationStatus = "FINISHED"

	RuleBasedUserListPrepopulationStatusFAILED RuleBasedUserListPrepopulationStatus = "FAILED"
)

//
// Size range in terms of number of users of a UserList/UserInterest.
//
type SizeRange string

const (
	SizeRangeLESS_THAN_FIVE_HUNDRED SizeRange = "LESS_THAN_FIVE_HUNDRED"

	SizeRangeLESS_THAN_ONE_THOUSAND SizeRange = "LESS_THAN_ONE_THOUSAND"

	SizeRangeONE_THOUSAND_TO_TEN_THOUSAND SizeRange = "ONE_THOUSAND_TO_TEN_THOUSAND"

	SizeRangeTEN_THOUSAND_TO_FIFTY_THOUSAND SizeRange = "TEN_THOUSAND_TO_FIFTY_THOUSAND"

	SizeRangeFIFTY_THOUSAND_TO_ONE_HUNDRED_THOUSAND SizeRange = "FIFTY_THOUSAND_TO_ONE_HUNDRED_THOUSAND"

	SizeRangeONE_HUNDRED_THOUSAND_TO_THREE_HUNDRED_THOUSAND SizeRange = "ONE_HUNDRED_THOUSAND_TO_THREE_HUNDRED_THOUSAND"

	SizeRangeTHREE_HUNDRED_THOUSAND_TO_FIVE_HUNDRED_THOUSAND SizeRange = "THREE_HUNDRED_THOUSAND_TO_FIVE_HUNDRED_THOUSAND"

	SizeRangeFIVE_HUNDRED_THOUSAND_TO_ONE_MILLION SizeRange = "FIVE_HUNDRED_THOUSAND_TO_ONE_MILLION"

	SizeRangeONE_MILLION_TO_TWO_MILLION SizeRange = "ONE_MILLION_TO_TWO_MILLION"

	SizeRangeTWO_MILLION_TO_THREE_MILLION SizeRange = "TWO_MILLION_TO_THREE_MILLION"

	SizeRangeTHREE_MILLION_TO_FIVE_MILLION SizeRange = "THREE_MILLION_TO_FIVE_MILLION"

	SizeRangeFIVE_MILLION_TO_TEN_MILLION SizeRange = "FIVE_MILLION_TO_TEN_MILLION"

	SizeRangeTEN_MILLION_TO_TWENTY_MILLION SizeRange = "TEN_MILLION_TO_TWENTY_MILLION"

	SizeRangeTWENTY_MILLION_TO_THIRTY_MILLION SizeRange = "TWENTY_MILLION_TO_THIRTY_MILLION"

	SizeRangeTHIRTY_MILLION_TO_FIFTY_MILLION SizeRange = "THIRTY_MILLION_TO_FIFTY_MILLION"

	SizeRangeOVER_FIFTY_MILLION SizeRange = "OVER_FIFTY_MILLION"
)

//
// Supported operators for strings.
//
type StringRuleItemStringOperator string

const (

	//
	// <span class="constraint Rejected">Used for return value only. An enumeration could not be processed, typically due to incompatibility with your WSDL version.</span>
	//
	StringRuleItemStringOperatorUNKNOWN StringRuleItemStringOperator = "UNKNOWN"

	StringRuleItemStringOperatorCONTAINS StringRuleItemStringOperator = "CONTAINS"

	StringRuleItemStringOperatorEQUALS StringRuleItemStringOperator = "EQUALS"

	StringRuleItemStringOperatorSTARTS_WITH StringRuleItemStringOperator = "STARTS_WITH"

	StringRuleItemStringOperatorENDS_WITH StringRuleItemStringOperator = "ENDS_WITH"

	StringRuleItemStringOperatorNOT_EQUAL StringRuleItemStringOperator = "NOT_EQUAL"

	StringRuleItemStringOperatorNOT_CONTAIN StringRuleItemStringOperator = "NOT_CONTAIN"

	StringRuleItemStringOperatorNOT_START_WITH StringRuleItemStringOperator = "NOT_START_WITH"

	StringRuleItemStringOperatorNOT_END_WITH StringRuleItemStringOperator = "NOT_END_WITH"
)

//
// Indicates the reason why the userlist was closed.
//
type UserListClosingReason string

const (

	//
	// <span class="constraint Rejected">Used for return value only. An enumeration could not be processed, typically due to incompatibility with your WSDL version.</span>
	//
	UserListClosingReasonUNKNOWN UserListClosingReason = "UNKNOWN"

	//
	// The userlist was closed because of not being used for over one year.
	//
	UserListClosingReasonUNUSED_LIST UserListClosingReason = "UNUSED_LIST"
)

type UserListErrorReason string

const (

	//
	// Creating and updating external remarketing user lists is not supported.
	//
	UserListErrorReasonEXTERNAL_REMARKETING_USER_LIST_MUTATE_NOT_SUPPORTED UserListErrorReason = "EXTERNAL_REMARKETING_USER_LIST_MUTATE_NOT_SUPPORTED"

	//
	// Concrete type of user list (logical v.s. remarketing) is required for
	// ADD and SET operations.
	//
	UserListErrorReasonCONCRETE_TYPE_REQUIRED UserListErrorReason = "CONCRETE_TYPE_REQUIRED"

	//
	// Adding/updating user list conversion types requires specifying the conversion
	// type id.
	//
	UserListErrorReasonCONVERSION_TYPE_ID_REQUIRED UserListErrorReason = "CONVERSION_TYPE_ID_REQUIRED"

	//
	// Remarketing user list cannot have duplicate conversion types.
	//
	UserListErrorReasonDUPLICATE_CONVERSION_TYPES UserListErrorReason = "DUPLICATE_CONVERSION_TYPES"

	//
	// Conversion type is invalid/unknown.
	//
	UserListErrorReasonINVALID_CONVERSION_TYPE UserListErrorReason = "INVALID_CONVERSION_TYPE"

	//
	// User list description is empty or invalid
	//
	UserListErrorReasonINVALID_DESCRIPTION UserListErrorReason = "INVALID_DESCRIPTION"

	//
	// User list name is empty or invalid.
	//
	UserListErrorReasonINVALID_NAME UserListErrorReason = "INVALID_NAME"

	//
	// Type of the UserList does not match.
	//
	UserListErrorReasonINVALID_TYPE UserListErrorReason = "INVALID_TYPE"

	//
	// Can't use similar list in logical user list rule when operator is NONE.
	//
	UserListErrorReasonCAN_NOT_ADD_SIMILAR_LIST_AS_LOGICAL_LIST_NONE_OPERAND UserListErrorReason = "CAN_NOT_ADD_SIMILAR_LIST_AS_LOGICAL_LIST_NONE_OPERAND"

	//
	// Embedded logical user lists are not allowed.
	//
	UserListErrorReasonCAN_NOT_ADD_LOGICAL_LIST_AS_LOGICAL_LIST_OPERAND UserListErrorReason = "CAN_NOT_ADD_LOGICAL_LIST_AS_LOGICAL_LIST_OPERAND"

	//
	// User list rule operand is invalid.
	//
	UserListErrorReasonINVALID_USER_LIST_LOGICAL_RULE_OPERAND UserListErrorReason = "INVALID_USER_LIST_LOGICAL_RULE_OPERAND"

	//
	// Name is already being used for another user list for the account.
	//
	UserListErrorReasonNAME_ALREADY_USED UserListErrorReason = "NAME_ALREADY_USED"

	//
	// Name is required when creating a new conversion type.
	//
	UserListErrorReasonNEW_CONVERSION_TYPE_NAME_REQUIRED UserListErrorReason = "NEW_CONVERSION_TYPE_NAME_REQUIRED"

	//
	// The given conversion type name has been used.
	//
	UserListErrorReasonCONVERSION_TYPE_NAME_ALREADY_USED UserListErrorReason = "CONVERSION_TYPE_NAME_ALREADY_USED"

	//
	// Only an owner account may edit a user list.
	//
	UserListErrorReasonOWNERSHIP_REQUIRED_FOR_SET UserListErrorReason = "OWNERSHIP_REQUIRED_FOR_SET"

	//
	// Removing user lists is not supported.
	//
	UserListErrorReasonREMOVE_NOT_SUPPORTED UserListErrorReason = "REMOVE_NOT_SUPPORTED"

	//
	// The user list of the type is not mutable
	//
	UserListErrorReasonUSER_LIST_MUTATE_NOT_SUPPORTED UserListErrorReason = "USER_LIST_MUTATE_NOT_SUPPORTED"

	//
	// Rule is invalid.
	//
	UserListErrorReasonINVALID_RULE UserListErrorReason = "INVALID_RULE"

	//
	// The specified date range is empty.
	//
	UserListErrorReasonINVALID_DATE_RANGE UserListErrorReason = "INVALID_DATE_RANGE"

	//
	// A userlist which is privacy sensitive or legal rejected cannot be mutated by external users.
	//
	UserListErrorReasonCAN_NOT_MUTATE_SENSITIVE_USERLIST UserListErrorReason = "CAN_NOT_MUTATE_SENSITIVE_USERLIST"

	//
	// Maximum number of rulebased user lists a customer can have.
	//
	UserListErrorReasonMAX_NUM_RULEBASED_USERLISTS UserListErrorReason = "MAX_NUM_RULEBASED_USERLISTS"

	//
	// BasicUserList's billable record field cannot be modified once it is set.
	//
	UserListErrorReasonCANNOT_MODIFY_BILLABLE_RECORD_COUNT UserListErrorReason = "CANNOT_MODIFY_BILLABLE_RECORD_COUNT"

	//
	// appId field can only be set when uploadKeyType is MOBILE_ADVERTISING_ID.
	//
	UserListErrorReasonAPP_ID_NOT_ALLOWED UserListErrorReason = "APP_ID_NOT_ALLOWED"

	//
	// appId field must be set when uploadKeyType is MOBILE_ADVERTISING_ID.
	//
	UserListErrorReasonAPP_ID_NOT_SET UserListErrorReason = "APP_ID_NOT_SET"

	//
	// Default generic error.
	//
	UserListErrorReasonUSER_LIST_SERVICE_ERROR UserListErrorReason = "USER_LIST_SERVICE_ERROR"
)

type UserListLogicalRuleOperator string

const (

	//
	// And - all of the operands.
	//
	UserListLogicalRuleOperatorALL UserListLogicalRuleOperator = "ALL"

	//
	// Or - at least one of the operands.
	//
	UserListLogicalRuleOperatorANY UserListLogicalRuleOperator = "ANY"

	//
	// Not - none of the operands.
	//
	UserListLogicalRuleOperatorNONE UserListLogicalRuleOperator = "NONE"

	//
	// <span class="constraint Rejected">Used for return value only. An enumeration could not be processed, typically due to incompatibility with your WSDL version.</span>
	//
	UserListLogicalRuleOperatorUNKNOWN UserListLogicalRuleOperator = "UNKNOWN"
)

//
// Membership status of the user list. This status indicates whether a user list
// can accumulate more users and may be targeted to.
//
type UserListMembershipStatus string

const (

	//
	// Open status - list is accruing members and can be targeted to.
	//
	UserListMembershipStatusOPEN UserListMembershipStatus = "OPEN"

	//
	// Closed status - No new members being added. Can not be used for targeting.
	//
	UserListMembershipStatusCLOSED UserListMembershipStatus = "CLOSED"
)

//
// Rule based userlist rule type.
//
type UserListRuleTypeEnumsEnum string

const (

	//
	// <span class="constraint Rejected">Used for return value only. An enumeration could not be processed, typically due to incompatibility with your WSDL version.</span>
	//
	UserListRuleTypeEnumsEnumUNKNOWN UserListRuleTypeEnumsEnum = "UNKNOWN"

	//
	// AND of ORs: Conjunctive normal form.
	//
	UserListRuleTypeEnumsEnumCNF UserListRuleTypeEnumsEnum = "CNF"

	//
	// OR of ANDs: Disjunctive normal form.
	//
	UserListRuleTypeEnumsEnumDNF UserListRuleTypeEnumsEnum = "DNF"
)

//
// The user list types
//
type UserListType string

const (

	//
	// UNKNOWN value can not be passed as input.
	// <span class="constraint Rejected">Used for return value only. An enumeration could not be processed, typically due to incompatibility with your WSDL version.</span>
	//
	UserListTypeUNKNOWN UserListType = "UNKNOWN"

	//
	// UserList represented as a collection of conversion types.
	//
	UserListTypeREMARKETING UserListType = "REMARKETING"

	//
	// UserList represented as a combination of other user lists/interests.
	//
	UserListTypeLOGICAL UserListType = "LOGICAL"

	//
	// UserList created in the DoubleClick platform.
	//
	UserListTypeEXTERNAL_REMARKETING UserListType = "EXTERNAL_REMARKETING"

	//
	// UserList associated with a rule.
	//
	UserListTypeRULE_BASED UserListType = "RULE_BASED"

	//
	// UserList with users similar to users of another UserList.
	//
	UserListTypeSIMILAR UserListType = "SIMILAR"

	//
	// UserList of first party CRM data provided by advertiser in the form of emails or
	// other formats.
	//
	UserListTypeCRM_BASED UserListType = "CRM_BASED"
)

//
// The status of the upload/remove-all operation on a CRM based UserList.
//
type UserListUploadStatus string

const (

	//
	// <span class="constraint Rejected">Used for return value only. An enumeration could not be processed, typically due to incompatibility with your WSDL version.</span>
	//
	UserListUploadStatusUNKNOWN UserListUploadStatus = "UNKNOWN"

	//
	// The upload/remove-all operation of this UserList is still in process.
	//
	UserListUploadStatusIN_PROCESS UserListUploadStatus = "IN_PROCESS"

	//
	// The upload/remove-all operation of this UserList has succeeded.
	//
	UserListUploadStatusSUCCESS UserListUploadStatus = "SUCCESS"

	//
	// The upload/remove-all operation of this UserList has failed.
	//
	UserListUploadStatusFAILURE UserListUploadStatus = "FAILURE"
)

type Get struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/rm/v201802 get"`

	//
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	//
	ServiceSelector *Selector `xml:"serviceSelector,omitempty"`
}

type GetResponse struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/rm/v201802 getResponse"`

	Rval *UserListPage `xml:"rval,omitempty"`
}

type Mutate struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/rm/v201802 mutate"`

	//
	// <span class="constraint CollectionSize">The minimum size of this collection is 1. The maximum size of this collection is 10000.</span>
	// <span class="constraint ContentsNotNull">This field must not contain {@code null} elements.</span>
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	// <span class="constraint SupportedOperators">The following {@link Operator}s are supported: ADD, SET.</span>
	//
	Operations []*UserListOperation `xml:"operations,omitempty"`
}

type MutateResponse struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/rm/v201802 mutateResponse"`

	Rval *UserListReturnValue `xml:"rval,omitempty"`
}

type MutateMembers struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/rm/v201802 mutateMembers"`

	//
	// <span class="constraint CollectionSize">The minimum size of this collection is 1. The maximum size of this collection is 10000.</span>
	// <span class="constraint ContentsNotNull">This field must not contain {@code null} elements.</span>
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	// <span class="constraint SupportedOperators">The following {@link Operator}s are supported: ADD, REMOVE.</span>
	//
	Operations []*MutateMembersOperation `xml:"operations,omitempty"`
}

type MutateMembersResponse struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/rm/v201802 mutateMembersResponse"`

	Rval *MutateMembersReturnValue `xml:"rval,omitempty"`
}

type Query struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/rm/v201802 query"`

	//
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	//
	Query string `xml:"query,omitempty"`
}

type QueryResponse struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/rm/v201802 queryResponse"`

	Rval *UserListPage `xml:"rval,omitempty"`
}

type AddressInfo struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/rm/v201802 AddressInfo"`

	//
	// First name of the member, which is hashed as SHA-256 after normalized (Lowercase all
	// characters; Remove any extra spaces before, after, and in between).
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	//
	HashedFirstName string `xml:"hashedFirstName,omitempty"`

	//
	// Last name of the member, which is hashed as SHA-256 after normalized (lower case only and no
	// punctuation).
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	//
	HashedLastName string `xml:"hashedLastName,omitempty"`

	//
	// 2-letter country code in ISO-3166-1 alpha-2 of the member's address.
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	// <span class="constraint StringLength">The length of this string should be between 2 and 2, inclusive.</span>
	//
	CountryCode string `xml:"countryCode,omitempty"`

	//
	// Zip code of the member's address.
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	// <span class="constraint StringLength">This string must not be empty.</span>
	//
	ZipCode string `xml:"zipCode,omitempty"`
}

type CombinedRuleUserList struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/rm/v201802 CombinedRuleUserList"`

	*RuleBasedUserList

	//
	// <span class="constraint Required">This field is required and should not be {@code null} when it is contained within {@link Operator}s : ADD.</span>
	//
	LeftOperand *Rule `xml:"leftOperand,omitempty"`

	//
	// <span class="constraint Required">This field is required and should not be {@code null} when it is contained within {@link Operator}s : ADD.</span>
	//
	RightOperand *Rule `xml:"rightOperand,omitempty"`

	//
	// <span class="constraint Required">This field is required and should not be {@code null} when it is contained within {@link Operator}s : ADD.</span>
	//
	RuleOperator *CombinedRuleUserListRuleOperator `xml:"ruleOperator,omitempty"`
}

type UserListConversionType struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/rm/v201802 UserListConversionType"`

	//
	// Conversion type id
	//
	Id int64 `xml:"id,omitempty"`

	//
	// Name of this conversion type
	//
	Name string `xml:"name,omitempty"`

	//
	// The category of the ConversionType based on the location where the
	// conversion event was generated (from a user's perspective).
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	Category *UserListConversionTypeCategory `xml:"category,omitempty"`
}

type CrmBasedUserList struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/rm/v201802 CrmBasedUserList"`

	*UserList

	//
	// A string that uniquely identifies a mobile application from which the data was
	// collected to AdWords API.
	// For iOS, the ID string is the 9 digit string that appears at the end of an App Store URL
	// (e.g., "476943146" for "Flood-It! 2" whose App Store link is
	// http://itunes.apple.com/us/app/flood-it!-2/id476943146).
	// For Android, the ID string is the application's package name
	// (e.g., "com.labpixies.colordrips" for "Color Drips" given Google Play link
	// https://play.google.com/store/apps/details?id=com.labpixies.colordrips).
	//
	// Required when creating CrmBasedUserList for uploading mobile advertising IDs.
	// <span class="constraint Selectable">This field can be selected using the value "AppId".</span>
	//
	AppId string `xml:"appId,omitempty"`

	//
	// Matching key type of the list.
	// This field is read only and set on the first upload by API.
	// Mixed data types are not allowed on the same list.
	// From v201802, this field will be required for an ADD operation.
	// <span class="constraint Selectable">This field can be selected using the value "UploadKeyType".</span>
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API for the following {@link Operator}s: SET.</span>
	//
	UploadKeyType *CustomerMatchUploadKeyType `xml:"uploadKeyType,omitempty"`

	//
	// Data source of the list.
	// Default value is FIRST_PARTY. Only whitelisted customers can create third party sourced crm
	// lists.
	// <span class="constraint Selectable">This field can be selected using the value "DataSourceType".</span><span class="constraint Filterable">This field can be filtered on.</span>
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API for the following {@link Operator}s: SET.</span>
	//
	DataSourceType *CrmDataSourceType `xml:"dataSourceType,omitempty"`

	//
	// <span class="constraint Selectable">This field can be selected using the value "DataUploadResult".</span>
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API for the following {@link Operator}s: ADD.</span>
	//
	DataUploadResult *DataUploadResult `xml:"dataUploadResult,omitempty"`
}

type DataUploadResult struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/rm/v201802 DataUploadResult"`

	//
	// Indicates status of the upload operation.
	// Upload operation is triggered when {@link MutateMembersOperand#removeAll removeAll} is not set
	// to true and {@link Operator operator} is "ADD" or "REMOVE".
	//
	UploadStatus *UserListUploadStatus `xml:"uploadStatus,omitempty"`

	//
	// Indicates status of the remove-all operation.
	// Remove-all operation is triggered when {@link MutateMembersOperand#removeAll removeAll} is set
	// to true and {@link Operator operator} is "REMOVE".
	//
	RemoveAllStatus *UserListUploadStatus `xml:"removeAllStatus,omitempty"`
}

type DateKey struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/rm/v201802 DateKey"`

	//
	// <span class="constraint MatchesRegex">A name must begin with US-ascii letters or underscore or UTF8 code that is greater than 127 and consist of US-ascii letters or digits or underscore or UTF8 code that is greater than 127. This is checked by the regular expression '^[a-zA-Z_?-?][a-zA-Z0-9_?-?]*$'.</span>
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	// <span class="constraint StringLength">This string must not be empty, (trimmed).</span>
	//
	Name string `xml:"name,omitempty"`
}

type DateRuleItem struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/rm/v201802 DateRuleItem"`

	//
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	//
	Key *DateKey `xml:"key,omitempty"`

	//
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	//
	Op *DateRuleItemDateOperator `xml:"op,omitempty"`

	//
	// The right hand side of date rule item. The date's format should be YYYYMMDD.
	//
	Value string `xml:"value,omitempty"`

	//
	// The relative date value of the right hand side. The {@code value} field will
	// override this field when both are present.
	//
	RelativeValue *RelativeDate `xml:"relativeValue,omitempty"`
}

type DateSpecificRuleUserList struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/rm/v201802 DateSpecificRuleUserList"`

	*RuleBasedUserList

	//
	// Boolean rule that defines visitor of a page. This field is selected by default.
	// <span class="constraint Selectable">This field can be selected using the value "DateSpecificListRule".</span>
	// <span class="constraint Required">This field is required and should not be {@code null} when it is contained within {@link Operator}s : ADD.</span>
	//
	Rule *Rule `xml:"rule,omitempty"`

	//
	// Start date of users visit. If set to <code>20000101</code>, then includes
	// all users before <code>endDate</code>. The date's format should be YYYYMMDD.
	// This field is selected by default.
	// <span class="constraint Selectable">This field can be selected using the value "DateSpecificListStartDate".</span>
	// <span class="constraint Required">This field is required and should not be {@code null} when it is contained within {@link Operator}s : ADD.</span>
	//
	StartDate string `xml:"startDate,omitempty"`

	//
	// End date of users visit. If set to <code>20371230</code>, then includes
	// all users after <code>startDate</code>. The date's format should be YYYYMMDD.
	// This field is selected by default.
	// <span class="constraint Selectable">This field can be selected using the value "DateSpecificListEndDate".</span>
	// <span class="constraint Required">This field is required and should not be {@code null} when it is contained within {@link Operator}s : ADD.</span>
	//
	EndDate string `xml:"endDate,omitempty"`
}

type ExpressionRuleUserList struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/rm/v201802 ExpressionRuleUserList"`

	*RuleBasedUserList

	//
	// Boolean rule that defines this user list. The rule consists of a list of rule item groups and
	// each rule item group consists of a list of rule items.
	// All the rule item groups are ORed together for evaluation before version V201705.
	// Starting from version V201705, the group operator is based on {@link Rule#getRuleType()}.
	// This field is selected by default.
	// <span class="constraint Selectable">This field can be selected using the value "ExpressionListRule".</span>
	// <span class="constraint Required">This field is required and should not be {@code null} when it is contained within {@link Operator}s : ADD.</span>
	//
	Rule *Rule `xml:"rule,omitempty"`
}

type LogicalUserList struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/rm/v201802 LogicalUserList"`

	*UserList

	//
	// Logical list rules that define this user list.  The rules are defined as
	// logical operator (ALL/ANY/NONE) and a list of user lists. All the rules are
	// anded for the evaluation. Required for ADD operation.
	// <span class="constraint Selectable">This field can be selected using the value "Rules".</span>
	//
	Rules []*UserListLogicalRule `xml:"rules,omitempty"`
}

type LogicalUserListOperand struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/rm/v201802 LogicalUserListOperand"`

	UserList *UserList `xml:"UserList,omitempty"`
}

type Member struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/rm/v201802 Member"`

	//
	// Hashed email address using SHA-256 hash function after normalization.
	//
	HashedEmail string `xml:"hashedEmail,omitempty"`

	//
	// Mobile device ID (advertising ID/IDFA).
	//
	MobileId string `xml:"mobileId,omitempty"`

	//
	// Hashed phone number using SHA-256 hash function after normalization.
	//
	HashedPhoneNumber string `xml:"hashedPhoneNumber,omitempty"`

	//
	// Address info.
	//
	AddressInfo *AddressInfo `xml:"addressInfo,omitempty"`

	//
	// Advertiser generated and assigned user ID. Accessible to whitelisted US customers only.
	// <span class="constraint StringLength">The length of this string should be between 1 and 512, inclusive.</span>
	//
	UserId string `xml:"userId,omitempty"`
}

type MutateMembersError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/rm/v201802 MutateMembersError"`

	*ApiError

	Reason *MutateMembersErrorReason `xml:"reason,omitempty"`
}

type MutateMembersOperand struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/rm/v201802 MutateMembersOperand"`

	//
	// The id of the user list.
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	//
	UserListId int64 `xml:"userListId,omitempty"`

	//
	// Set to indicate a remove-all operation which will remove all members from the user list.
	// Can only be set with {@code Operator#REMOVE} and
	// when set to true {@link #members} must be null or empty.
	//
	RemoveAll bool `xml:"removeAll,omitempty"`

	//
	// A list of members to be added or removed.
	//
	// <p>If {@link #removeAll} is {@code true}, this list must be {@code null} or empty. Otherwise,
	// this field is required and there must be at least one member.
	// <span class="constraint CollectionSize">The maximum size of this collection is 1000000.</span>
	// <span class="constraint ContentsNotNull">This field must not contain {@code null} elements.</span>
	//
	MembersList []*Member `xml:"membersList,omitempty"`
}

type MutateMembersOperation struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/rm/v201802 MutateMembersOperation"`

	*Operation

	//
	// The mutate members operand to operate on.
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	//
	Operand *MutateMembersOperand `xml:"operand,omitempty"`
}

type MutateMembersReturnValue struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/rm/v201802 MutateMembersReturnValue"`

	//
	// The user lists associated in mutate members operations.
	//
	UserLists []*UserList `xml:"userLists,omitempty"`
}

type NumberKey struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/rm/v201802 NumberKey"`

	//
	// <span class="constraint MatchesRegex">A name must begin with US-ascii letters or underscore or UTF8 code that is greater than 127 and consist of US-ascii letters or digits or underscore or UTF8 code that is greater than 127. This is checked by the regular expression '^[a-zA-Z_?-?][a-zA-Z0-9_?-?]*$'.</span>
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	// <span class="constraint StringLength">This string must not be empty, (trimmed).</span>
	//
	Name string `xml:"name,omitempty"`
}

type NumberRuleItem struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/rm/v201802 NumberRuleItem"`

	//
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	//
	Key *NumberKey `xml:"key,omitempty"`

	//
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	//
	Op *NumberRuleItemNumberOperator `xml:"op,omitempty"`

	//
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	//
	Value float64 `xml:"value,omitempty"`
}

type RelativeDate struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/rm/v201802 RelativeDate"`

	//
	// Number of days offset from current date.
	// <span class="constraint InRange">This field must be greater than or equal to 0.</span>
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	//
	OffsetInDays int32 `xml:"offsetInDays,omitempty"`
}

type BasicUserList struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/rm/v201802 BasicUserList"`

	*UserList

	//
	// Conversion types associated with this user list.
	// <span class="constraint Selectable">This field can be selected using the value "ConversionTypes".</span>
	//
	ConversionTypes []*UserListConversionType `xml:"conversionTypes,omitempty"`
}

type Rule struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/rm/v201802 Rule"`

	//
	// List of rule item groups that defines this rule.
	// Rule item groups are ORed together for evaluation before version V201705.
	// Starting from version V201705, rule item groups are grouped together based on
	// {@link #getRuleType()} for evaluation.
	// <span class="constraint CollectionSize">The minimum size of this collection is 1.</span>
	// <span class="constraint ContentsDistinct">This field must contain distinct elements.</span>
	// <span class="constraint ContentsNotNull">This field must not contain {@code null} elements.</span>
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	//
	Groups []*RuleItemGroup `xml:"groups,omitempty"`

	//
	// Rule type is used to determine how to group rule item groups and rule items inside rule item
	// group. Currently, conjunctive normal form (AND of ORs) is only supported for
	// ExpressionRuleUserList. If no ruleType is specified, it will be treated as disjunctive normal
	// form (OR of ANDs), namely rule item groups are ORed together and inside each rule item group,
	// rule items are ANDed together.
	//
	RuleType *UserListRuleTypeEnumsEnum `xml:"ruleType,omitempty"`
}

type RuleBasedUserList struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/rm/v201802 RuleBasedUserList"`

	*UserList

	//
	// Status of pre-population. The field is default to NONE if not set which means the previous
	// users will not be considered. If set to REQUESTED, past site visitors or app users who match
	// the list definition will be included in the list (works on the Display Network only). This will
	// only pre-populate past users within up to the last 30 days, depending on the list's membership
	// duration and the date when the remarketing tag is added. The status will be updated to FINISHED
	// once request is processed, or FAILED if the request fails.
	// <span class="constraint Selectable">This field can be selected using the value "PrepopulationStatus".</span>
	//
	PrepopulationStatus *RuleBasedUserListPrepopulationStatus `xml:"prepopulationStatus,omitempty"`
}

type RuleItem struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/rm/v201802 RuleItem"`

	DateRuleItem *DateRuleItem `xml:"DateRuleItem,omitempty"`

	NumberRuleItem *NumberRuleItem `xml:"NumberRuleItem,omitempty"`

	StringRuleItem *StringRuleItem `xml:"StringRuleItem,omitempty"`
}

type RuleItemGroup struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/rm/v201802 RuleItemGroup"`

	//
	// Before version V201705, rule items are ANDed together.
	// Starting from version V201705, rule items will be grouped together based on
	// {@link Rule#getRuleType()}.
	// <span class="constraint CollectionSize">The minimum size of this collection is 1. The maximum size of this collection is 1000.</span>
	// <span class="constraint ContentsDistinct">This field must contain distinct elements.</span>
	// <span class="constraint ContentsNotNull">This field must not contain {@code null} elements.</span>
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	//
	Items []*RuleItem `xml:"items,omitempty"`
}

type SimilarUserList struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/rm/v201802 SimilarUserList"`

	*UserList

	//
	// Seed UserListId from which this list is derived.
	// <span class="constraint Selectable">This field can be selected using the value "SeedUserListId".</span><span class="constraint Filterable">This field can be filtered on.</span>
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API for the following {@link Operator}s: SET.</span>
	// <span class="constraint Required">This field is required and should not be {@code null} when it is contained within {@link Operator}s : ADD.</span>
	//
	SeedUserListId int64 `xml:"seedUserListId,omitempty"`

	//
	// Name of the seed user list.
	// <span class="constraint Selectable">This field can be selected using the value "SeedUserListName".</span>
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	SeedUserListName string `xml:"seedUserListName,omitempty"`

	//
	// Description of this seed user list.
	// <span class="constraint Selectable">This field can be selected using the value "SeedUserListDescription".</span>
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	SeedUserListDescription string `xml:"seedUserListDescription,omitempty"`

	//
	// Membership status of this seed user list.
	// <span class="constraint Selectable">This field can be selected using the value "SeedUserListStatus".</span>
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	SeedUserListStatus *UserListMembershipStatus `xml:"seedUserListStatus,omitempty"`

	//
	// Estimated number of users in this seed user list.
	// This value is null if the number of users has not yet been determined.
	// <span class="constraint Selectable">This field can be selected using the value "SeedListSize".</span><span class="constraint Filterable">This field can be filtered on.</span>
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	SeedListSize int64 `xml:"seedListSize,omitempty"`
}

type StringKey struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/rm/v201802 StringKey"`

	//
	// <span class="constraint MatchesRegex">A name must begin with US-ascii letters or underscore or UTF8 code that is greater than 127 and consist of US-ascii letters or digits or underscore or UTF8 code that is greater than 127. This is checked by the regular expression '^[a-zA-Z_?-?][a-zA-Z0-9_?-?]*$'.</span>
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	// <span class="constraint StringLength">This string must not be empty, (trimmed).</span>
	//
	Name string `xml:"name,omitempty"`
}

type StringRuleItem struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/rm/v201802 StringRuleItem"`

	//
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	//
	Key *StringKey `xml:"key,omitempty"`

	//
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	//
	Op *StringRuleItemStringOperator `xml:"op,omitempty"`

	//
	// The right hand side of the string rule item. For URL/Referrer URL,
	// <code>value</code> can not contain illegal URL chars such as: <code>"()'\"\t"</code>.
	// <span class="constraint MatchesRegex">String value can not contain newline (
	// ) or both single quote and double quote. This is checked by the regular expression '[^
	// ']*|[^
	// "]*'.</span>
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	//
	Value string `xml:"value,omitempty"`
}

type UserList struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/rm/v201802 UserList"`

	//
	// Id of this user list.
	// <span class="constraint Selectable">This field can be selected using the value "Id".</span><span class="constraint Filterable">This field can be filtered on.</span>
	// <span class="constraint Required">This field is required and should not be {@code null} when it is contained within {@link Operator}s : SET.</span>
	//
	Id int64 `xml:"id,omitempty"`

	//
	// A flag that indicates if a user may edit a list. Depends on the list ownership
	// and list type. For example, external remarketing user lists are not editable.
	// <span class="constraint Selectable">This field can be selected using the value "IsReadOnly".</span>
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	IsReadOnly bool `xml:"isReadOnly,omitempty"`

	//
	// Name of this user list. Depending on its AccessReason, the user list name
	// may not be unique (e.g. if {@code AccessReason=SHARED}).
	// <span class="constraint Selectable">This field can be selected using the value "Name".</span><span class="constraint Filterable">This field can be filtered on.</span>
	//
	Name string `xml:"name,omitempty"`

	//
	// Description of this user list.
	// <span class="constraint Selectable">This field can be selected using the value "Description".</span>
	//
	Description string `xml:"description,omitempty"`

	//
	// Membership status of this user list. Indicates whether a user list is open
	// or active. Only open user lists can accumulate more users and can be targeted to.
	// <span class="constraint Selectable">This field can be selected using the value "Status".</span><span class="constraint Filterable">This field can be filtered on.</span>
	//
	Status *UserListMembershipStatus `xml:"status,omitempty"`

	//
	// An Id from external system. It is used by user list sellers to correlate ids on their
	// systems.
	// <span class="constraint Selectable">This field can be selected using the value "IntegrationCode".</span><span class="constraint Filterable">This field can be filtered on.</span>
	//
	IntegrationCode string `xml:"integrationCode,omitempty"`

	//
	// Indicates the reason this account has been granted access to the list. The reason can be
	// Shared, Owned, Licensed or Subscribed.
	// <span class="constraint Selectable">This field can be selected using the value "AccessReason".</span><span class="constraint Filterable">This field can be filtered on.</span>
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	AccessReason *AccessReason `xml:"accessReason,omitempty"`

	//
	// Indicates if this share is still active. When a UserList is shared with the user
	// this field is set to Active. Later the userList owner can decide to revoke the
	// share and make it Inactive. The default value of this field is set to Active.
	// <span class="constraint Selectable">This field can be selected using the value "AccountUserListStatus".</span><span class="constraint Filterable">This field can be filtered on.</span>
	//
	AccountUserListStatus *AccountUserListStatus `xml:"accountUserListStatus,omitempty"`

	//
	// Number of days a user's cookie stays on your list since its most recent addition to the list.
	// This field must be between 0 and 540 inclusive. However, for CRM based userlists, this field
	// can be set to 10000 which means no expiration.
	//
	// <p>It'll be ignored for {@link LogicalUserList}.
	// <span class="constraint Selectable">This field can be selected using the value "MembershipLifeSpan".</span><span class="constraint Filterable">This field can be filtered on.</span>
	//
	MembershipLifeSpan int64 `xml:"membershipLifeSpan,omitempty"`

	//
	// Estimated number of users in this user list, on the Google Display Network.
	// This value is null if the number of users has not yet been determined.
	// <span class="constraint Selectable">This field can be selected using the value "Size".</span><span class="constraint Filterable">This field can be filtered on.</span>
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	Size int64 `xml:"size,omitempty"`

	//
	// Size range in terms of number of users of the UserList.
	// <span class="constraint Selectable">This field can be selected using the value "SizeRange".</span>
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	SizeRange *SizeRange `xml:"sizeRange,omitempty"`

	//
	// Estimated number of users in this user list in the google.com domain.
	// These are the users available for targeting in search campaigns.
	// This value is null if the number of users has not yet been determined.
	// <span class="constraint Selectable">This field can be selected using the value "SizeForSearch".</span><span class="constraint Filterable">This field can be filtered on.</span>
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	SizeForSearch int64 `xml:"sizeForSearch,omitempty"`

	//
	// Size range in terms of number of users of the UserList, for Search ads.
	// <span class="constraint Selectable">This field can be selected using the value "SizeRangeForSearch".</span>
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	SizeRangeForSearch *SizeRange `xml:"sizeRangeForSearch,omitempty"`

	//
	// Type of this list: remarketing/logical/external remarketing.
	// <span class="constraint Selectable">This field can be selected using the value "ListType".</span><span class="constraint Filterable">This field can be filtered on.</span>
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	ListType *UserListType `xml:"listType,omitempty"`

	//
	// A flag that indicates this user list is eligible for Google Search Network.
	// <span class="constraint Selectable">This field can be selected using the value "IsEligibleForSearch".</span><span class="constraint Filterable">This field can be filtered on.</span>
	//
	IsEligibleForSearch bool `xml:"isEligibleForSearch,omitempty"`

	//
	// A flag that indicates this user list is eligible for Display Network.
	// <span class="constraint Selectable">This field can be selected using the value "IsEligibleForDisplay".</span><span class="constraint Filterable">This field can be filtered on.</span>
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	IsEligibleForDisplay bool `xml:"isEligibleForDisplay,omitempty"`

	//
	// Indicating the reason why this user list membership status is closed. It is only populated on
	// lists that were automatically closed due to inactivity, and will be cleared once the list
	// membership status becomes open.
	// <span class="constraint Selectable">This field can be selected using the value "ClosingReason".</span><span class="constraint Filterable">This field can be filtered on.</span>
	//
	ClosingReason *UserListClosingReason `xml:"closingReason,omitempty"`

	//
	// Indicates that this instance is a subtype of UserList.
	// Although this field is returned in the response, it is ignored on input
	// and cannot be selected. Specify xsi:type instead.
	//
	UserListType string `xml:"UserList.Type,omitempty"`
}

type UserListError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/rm/v201802 UserListError"`

	*ApiError

	//
	// The error reason represented by an enum.
	//
	Reason *UserListErrorReason `xml:"reason,omitempty"`
}

type UserListLogicalRule struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/rm/v201802 UserListLogicalRule"`

	//
	// The logical operator of the rule.
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	//
	Operator *UserListLogicalRuleOperator `xml:"operator,omitempty"`

	//
	// The list of operands of the rule.
	// <span class="constraint ContentsNotNull">This field must not contain {@code null} elements.</span>
	// <span class="constraint NotEmpty">This field must contain at least one element.</span>
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	//
	RuleOperands []*LogicalUserListOperand `xml:"ruleOperands,omitempty"`
}

type UserListOperation struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/rm/v201802 UserListOperation"`

	*Operation

	//
	// UserList to operate on
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	//
	Operand *UserList `xml:"operand,omitempty"`
}

type UserListPage struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/rm/v201802 UserListPage"`

	*Page

	//
	// The result entries in this page.
	//
	Entries []*UserList `xml:"entries,omitempty"`
}

type UserListReturnValue struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/rm/v201802 UserListReturnValue"`

	*ListReturnValue

	Value []*UserList `xml:"value,omitempty"`
}

type AdwordsUserListServiceInterface struct {
	client *SOAPClient
}

func NewAdwordsUserListServiceInterface(url string, tls bool, auth *BasicAuth) *AdwordsUserListServiceInterface {
	if url == "" {
		url = ""
	}
	client := NewSOAPClient(url, tls, auth)

	return &AdwordsUserListServiceInterface{
		client: client,
	}
}

func NewAdwordsUserListServiceInterfaceWithTLSConfig(url string, tlsCfg *tls.Config, auth *BasicAuth) *AdwordsUserListServiceInterface {
	if url == "" {
		url = ""
	}
	client := NewSOAPClientWithTLSConfig(url, tlsCfg, auth)

	return &AdwordsUserListServiceInterface{
		client: client,
	}
}

func (service *AdwordsUserListServiceInterface) AddHeader(header interface{}) {
	service.client.AddHeader(header)
}

// Backwards-compatible function: use AddHeader instead
func (service *AdwordsUserListServiceInterface) SetHeader(header interface{}) {
	service.client.AddHeader(header)
}

// Error can be either of the following types:
//
//   - ApiException
/*
   Returns the list of user lists that meet the selector criteria.

   @param serviceSelector the selector specifying the {@link UserList}s to return.
   @return a list of UserList entities which meet the selector criteria.
   @throws ApiException if problems occurred while fetching UserList information.
*/
func (service *AdwordsUserListServiceInterface) Get(request *Get) (*GetResponse, error) {
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
   Applies a list of mutate operations (i.e. add, set):

   Add - creates a set of user lists
   Set - updates a set of user lists
   Remove - not supported

   @param operations the operations to apply
   @return a list of UserList objects
*/
func (service *AdwordsUserListServiceInterface) Mutate(request *Mutate) (*MutateResponse, error) {
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
   Mutate members of user lists by either adding or removing their lists of members.
   The following {@link Operator}s are supported: ADD and REMOVE. The SET operator
   is not supported.

   <p>Note that operations cannot have same user list id but different operators.

   @param operations the mutate members operations to apply
   @return a list of UserList objects
   @throws ApiException when there are one or more errors with the request
*/
func (service *AdwordsUserListServiceInterface) MutateMembers(request *MutateMembers) (*MutateMembersResponse, error) {
	response := new(MutateMembersResponse)
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
   Returns the list of user lists that match the query.

   @param query The SQL-like AWQL query string
   @return A list of UserList
   @throws ApiException when the query is invalid or there are errors processing the request.
*/
func (service *AdwordsUserListServiceInterface) Query(request *Query) (*QueryResponse, error) {
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
