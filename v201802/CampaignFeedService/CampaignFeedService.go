package CampaignFeedService

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
// Status of the CampaignFeed.
//
type CampaignFeedStatus string

const (

	//
	// This CampaignFeed's data is currently being used.
	//
	CampaignFeedStatusENABLED CampaignFeedStatus = "ENABLED"

	//
	// This CampaignFeed's data is not used anymore.
	//
	CampaignFeedStatusREMOVED CampaignFeedStatus = "REMOVED"

	//
	// Unknown status.
	// <span class="constraint Rejected">Used for return value only. An enumeration could not be processed, typically due to incompatibility with your WSDL version.</span>
	//
	CampaignFeedStatusUNKNOWN CampaignFeedStatus = "UNKNOWN"
)

//
// Error reasons.
//
type CampaignFeedErrorReason string

const (

	//
	// An active feed already exists for this campaign and place holder type.
	//
	CampaignFeedErrorReasonFEED_ALREADY_EXISTS_FOR_PLACEHOLDER_TYPE CampaignFeedErrorReason = "FEED_ALREADY_EXISTS_FOR_PLACEHOLDER_TYPE"

	//
	// The specified id does not exist.
	//
	CampaignFeedErrorReasonINVALID_ID CampaignFeedErrorReason = "INVALID_ID"

	//
	// The specified feed is deleted.
	//
	CampaignFeedErrorReasonCANNOT_ADD_FOR_DELETED_FEED CampaignFeedErrorReason = "CANNOT_ADD_FOR_DELETED_FEED"

	//
	// The CampaignFeed already exists. SET should be used to modify the existing CampaignFeed.
	//
	CampaignFeedErrorReasonCANNOT_ADD_ALREADY_EXISTING_CAMPAIGN_FEED CampaignFeedErrorReason = "CANNOT_ADD_ALREADY_EXISTING_CAMPAIGN_FEED"

	//
	// Cannot operate on deleted campaign feed.
	//
	CampaignFeedErrorReasonCANNOT_OPERATE_ON_REMOVED_CAMPAIGN_FEED CampaignFeedErrorReason = "CANNOT_OPERATE_ON_REMOVED_CAMPAIGN_FEED"

	//
	// Invalid placeholder type ids.
	//
	CampaignFeedErrorReasonINVALID_PLACEHOLDER_TYPES CampaignFeedErrorReason = "INVALID_PLACEHOLDER_TYPES"

	//
	// Feed mapping for this placeholder type does not exist.
	//
	CampaignFeedErrorReasonMISSING_FEEDMAPPING_FOR_PLACEHOLDER_TYPE CampaignFeedErrorReason = "MISSING_FEEDMAPPING_FOR_PLACEHOLDER_TYPE"

	//
	// Location CampaignFeeds cannot be created unless there is a location CustomerFeed
	// for the specified feed.
	//
	CampaignFeedErrorReasonNO_EXISTING_LOCATION_CUSTOMER_FEED CampaignFeedErrorReason = "NO_EXISTING_LOCATION_CUSTOMER_FEED"

	CampaignFeedErrorReasonUNKNOWN CampaignFeedErrorReason = "UNKNOWN"
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
// The types of constant operands.
//
type ConstantOperandConstantType string

const (

	//
	// Boolean constant type. booleanValue should be set for this type.
	//
	ConstantOperandConstantTypeBOOLEAN ConstantOperandConstantType = "BOOLEAN"

	//
	// Double constant type. doubleValue should be set for this type.
	//
	ConstantOperandConstantTypeDOUBLE ConstantOperandConstantType = "DOUBLE"

	//
	// Long constant type. longValue should be set for this type.
	//
	ConstantOperandConstantTypeLONG ConstantOperandConstantType = "LONG"

	//
	// String constant type. stringValue should be set for this type.
	//
	ConstantOperandConstantTypeSTRING ConstantOperandConstantType = "STRING"
)

//
// The units of constant operands, if applicable.
//
type ConstantOperandUnit string

const (

	//
	// Meters.
	//
	ConstantOperandUnitMETERS ConstantOperandUnit = "METERS"

	//
	// Miles.
	//
	ConstantOperandUnitMILES ConstantOperandUnit = "MILES"

	ConstantOperandUnitNONE ConstantOperandUnit = "NONE"
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
// Operators that can be used in functions.
//
type FunctionOperator string

const (

	//
	// The IN operator.
	//
	FunctionOperatorIN FunctionOperator = "IN"

	//
	// The IDENTITY operator.
	//
	FunctionOperatorIDENTITY FunctionOperator = "IDENTITY"

	//
	// The EQUALS operator
	//
	FunctionOperatorEQUALS FunctionOperator = "EQUALS"

	//
	// Operator that takes two or more operands that are of type FunctionOperand
	// and checks that all the operands evaluate to true.
	// For functions related to ad formats, all the operands must be in lhsOperand.
	// Return ConstantOperand with Bool type.
	//
	FunctionOperatorAND FunctionOperator = "AND"

	//
	// Operator that returns true if the elements in lhsOperand contains any of the elements
	// in rhsOperands. Otherwise, return false.
	//
	FunctionOperatorCONTAINS_ANY FunctionOperator = "CONTAINS_ANY"

	//
	// <span class="constraint Rejected">Used for return value only. An enumeration could not be processed, typically due to incompatibility with your WSDL version.</span>
	//
	FunctionOperatorUNKNOWN FunctionOperator = "UNKNOWN"
)

//
// The reasons for the target error.
//
type FunctionErrorReason string

const (

	//
	// The format of the function is not recognized as a supported function format.
	//
	FunctionErrorReasonINVALID_FUNCTION_FORMAT FunctionErrorReason = "INVALID_FUNCTION_FORMAT"

	//
	// Operand data types do not match.
	//
	FunctionErrorReasonDATA_TYPE_MISMATCH FunctionErrorReason = "DATA_TYPE_MISMATCH"

	//
	// The operands cannot be used together in a conjunction.
	//
	FunctionErrorReasonINVALID_CONJUNCTION_OPERANDS FunctionErrorReason = "INVALID_CONJUNCTION_OPERANDS"

	//
	// Invalid numer of Operands.
	//
	FunctionErrorReasonINVALID_NUMBER_OF_OPERANDS FunctionErrorReason = "INVALID_NUMBER_OF_OPERANDS"

	//
	// Operand Type not supported.
	//
	FunctionErrorReasonINVALID_OPERAND_TYPE FunctionErrorReason = "INVALID_OPERAND_TYPE"

	//
	// Operator not supported.
	//
	FunctionErrorReasonINVALID_OPERATOR FunctionErrorReason = "INVALID_OPERATOR"

	//
	// Request context type not supported.
	//
	FunctionErrorReasonINVALID_REQUEST_CONTEXT_TYPE FunctionErrorReason = "INVALID_REQUEST_CONTEXT_TYPE"

	//
	// The matching function is not allowed for call placeholders
	//
	FunctionErrorReasonINVALID_FUNCTION_FOR_CALL_PLACEHOLDER FunctionErrorReason = "INVALID_FUNCTION_FOR_CALL_PLACEHOLDER"

	//
	// The matching function is not allowed for the specified placeholder
	//
	FunctionErrorReasonINVALID_FUNCTION_FOR_PLACEHOLDER FunctionErrorReason = "INVALID_FUNCTION_FOR_PLACEHOLDER"

	//
	// Invalid operand.
	//
	FunctionErrorReasonINVALID_OPERAND FunctionErrorReason = "INVALID_OPERAND"

	//
	// Missing value for the constant operand.
	//
	FunctionErrorReasonMISSING_CONSTANT_OPERAND_VALUE FunctionErrorReason = "MISSING_CONSTANT_OPERAND_VALUE"

	//
	// The value of the constant operand is invalid.
	//
	FunctionErrorReasonINVALID_CONSTANT_OPERAND_VALUE FunctionErrorReason = "INVALID_CONSTANT_OPERAND_VALUE"

	//
	// Invalid function nesting.
	//
	FunctionErrorReasonINVALID_NESTING FunctionErrorReason = "INVALID_NESTING"

	//
	// The Feed ID was different from another Feed ID in the same function.
	//
	FunctionErrorReasonMULTIPLE_FEED_IDS_NOT_SUPPORTED FunctionErrorReason = "MULTIPLE_FEED_IDS_NOT_SUPPORTED"

	//
	// The matching function is invalid for use with a feed with a fixed schema.
	//
	FunctionErrorReasonINVALID_FUNCTION_FOR_FEED_WITH_FIXED_SCHEMA FunctionErrorReason = "INVALID_FUNCTION_FOR_FEED_WITH_FIXED_SCHEMA"

	//
	// Invalid attribute name.
	//
	FunctionErrorReasonINVALID_ATTRIBUTE_NAME FunctionErrorReason = "INVALID_ATTRIBUTE_NAME"

	FunctionErrorReasonUNKNOWN FunctionErrorReason = "UNKNOWN"
)

//
// Function parsing error reason.
//
type FunctionParsingErrorReason string

const (

	//
	// Unexpected end of function string.
	//
	FunctionParsingErrorReasonNO_MORE_INPUT FunctionParsingErrorReason = "NO_MORE_INPUT"

	//
	// Could not find an expected character.
	//
	FunctionParsingErrorReasonEXPECTED_CHARACTER FunctionParsingErrorReason = "EXPECTED_CHARACTER"

	//
	// Unexpected separator character.
	//
	FunctionParsingErrorReasonUNEXPECTED_SEPARATOR FunctionParsingErrorReason = "UNEXPECTED_SEPARATOR"

	//
	// Unmatched left bracket or parenthesis.
	//
	FunctionParsingErrorReasonUNMATCHED_LEFT_BRACKET FunctionParsingErrorReason = "UNMATCHED_LEFT_BRACKET"

	//
	// Unmatched right bracket or parenthesis.
	//
	FunctionParsingErrorReasonUNMATCHED_RIGHT_BRACKET FunctionParsingErrorReason = "UNMATCHED_RIGHT_BRACKET"

	//
	// Functions are nested too deeply.
	//
	FunctionParsingErrorReasonTOO_MANY_NESTED_FUNCTIONS FunctionParsingErrorReason = "TOO_MANY_NESTED_FUNCTIONS"

	//
	// Missing right-hand-side operand.
	//
	FunctionParsingErrorReasonMISSING_RIGHT_HAND_OPERAND FunctionParsingErrorReason = "MISSING_RIGHT_HAND_OPERAND"

	//
	// Invalid operator/function name.
	//
	FunctionParsingErrorReasonINVALID_OPERATOR_NAME FunctionParsingErrorReason = "INVALID_OPERATOR_NAME"

	//
	// Feed attribute operand's argument is not an integer.
	//
	FunctionParsingErrorReasonFEED_ATTRIBUTE_OPERAND_ARGUMENT_NOT_INTEGER FunctionParsingErrorReason = "FEED_ATTRIBUTE_OPERAND_ARGUMENT_NOT_INTEGER"

	//
	// Missing function operands.
	//
	FunctionParsingErrorReasonNO_OPERANDS FunctionParsingErrorReason = "NO_OPERANDS"

	//
	// Function had too many operands.
	//
	FunctionParsingErrorReasonTOO_MANY_OPERANDS FunctionParsingErrorReason = "TOO_MANY_OPERANDS"

	FunctionParsingErrorReasonUNKNOWN FunctionParsingErrorReason = "UNKNOWN"
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

type RequestContextOperandContextType string

const (

	//
	// Feed item id in the request context.
	//
	RequestContextOperandContextTypeFEED_ITEM_ID RequestContextOperandContextType = "FEED_ITEM_ID"

	//
	// The device's platform (possible values are 'Desktop' or 'Mobile').
	//
	RequestContextOperandContextTypeDEVICE_PLATFORM RequestContextOperandContextType = "DEVICE_PLATFORM"

	//
	// <span class="constraint Rejected">Used for return value only. An enumeration could not be processed, typically due to incompatibility with your WSDL version.</span>
	//
	RequestContextOperandContextTypeUNKNOWN RequestContextOperandContextType = "UNKNOWN"
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

type Get struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 get"`

	//
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	//
	Selector *Selector `xml:"selector,omitempty"`
}

type GetResponse struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 getResponse"`

	Rval *CampaignFeedPage `xml:"rval,omitempty"`
}

type Mutate struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 mutate"`

	//
	// <span class="constraint ContentsNotNull">This field must not contain {@code null} elements.</span>
	// <span class="constraint NotEmpty">This field must contain at least one element.</span>
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	// <span class="constraint SupportedOperators">The following {@link Operator}s are supported: ADD, SET, REMOVE.</span>
	//
	Operations []*CampaignFeedOperation `xml:"operations,omitempty"`
}

type MutateResponse struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 mutateResponse"`

	Rval *CampaignFeedReturnValue `xml:"rval,omitempty"`
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

	Rval *CampaignFeedPage `xml:"rval,omitempty"`
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

type CampaignFeed struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 CampaignFeed"`

	//
	// Id of the Feed associated with the CampaignFeed.
	// <span class="constraint Selectable">This field can be selected using the value "FeedId".</span><span class="constraint Filterable">This field can be filtered on.</span>
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	//
	FeedId int64 `xml:"feedId,omitempty"`

	//
	// Id of the Campaign associated with the CampaignFeed.
	// <span class="constraint Selectable">This field can be selected using the value "CampaignId".</span><span class="constraint Filterable">This field can be filtered on.</span>
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	//
	CampaignId int64 `xml:"campaignId,omitempty"`

	//
	// Matching function associated with the CampaignFeed.
	// The matching function will return true/false indicating
	// which feed items may serve.
	// <span class="constraint Selectable">This field can be selected using the value "MatchingFunction".</span>
	// <span class="constraint Required">This field is required and should not be {@code null} when it is contained within {@link Operator}s : ADD.</span>
	//
	MatchingFunction *Function `xml:"matchingFunction,omitempty"`

	//
	// Indicates which <a href="/adwords/api/docs/appendix/placeholders">
	// placeholder types</a> the feed may populate under the
	// connected Campaign.
	// <span class="constraint Selectable">This field can be selected using the value "PlaceholderTypes".</span><span class="constraint Filterable">This field can be filtered on.</span>
	// <span class="constraint Required">This field is required and should not be {@code null} when it is contained within {@link Operator}s : ADD.</span>
	//
	PlaceholderTypes []int32 `xml:"placeholderTypes,omitempty"`

	//
	// Status of the CampaignFeed.
	// <span class="constraint Selectable">This field can be selected using the value "Status".</span><span class="constraint Filterable">This field can be filtered on.</span>
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	Status *CampaignFeedStatus `xml:"status,omitempty"`

	//
	// ID of the base campaign from which this draft/trial feed was created.
	// This field is only returned on get requests.
	// <span class="constraint Selectable">This field can be selected using the value "BaseCampaignId".</span><span class="constraint Filterable">This field can be filtered on.</span>
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	BaseCampaignId int64 `xml:"baseCampaignId,omitempty"`
}

type CampaignFeedError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 CampaignFeedError"`

	*ApiError

	//
	// Error reason.
	//
	Reason *CampaignFeedErrorReason `xml:"reason,omitempty"`
}

type CampaignFeedOperation struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 CampaignFeedOperation"`

	*Operation

	//
	// CampaignFeed operand.
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	//
	Operand *CampaignFeed `xml:"operand,omitempty"`
}

type CampaignFeedPage struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 CampaignFeedPage"`

	*NullStatsPage

	//
	// The resulting CampaignFeeds.
	//
	Entries []*CampaignFeed `xml:"entries,omitempty"`
}

type CampaignFeedReturnValue struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 CampaignFeedReturnValue"`

	*ListReturnValue

	//
	// The resulting CampaignFeeds.
	//
	Value []*CampaignFeed `xml:"value,omitempty"`

	PartialFailureErrors []*ApiError `xml:"partialFailureErrors,omitempty"`
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

type ConstantOperand struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 ConstantOperand"`

	*FunctionArgumentOperand

	//
	// Type of constant in this operand.
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	//
	Type_ *ConstantOperandConstantType `xml:"type,omitempty"`

	//
	// Units of constant in this operand.
	//
	Unit *ConstantOperandUnit `xml:"unit,omitempty"`

	//
	// Long value of the operand if it is a long type.
	//
	LongValue int64 `xml:"longValue,omitempty"`

	//
	// Boolean value of the operand if it is a boolean type.
	//
	BooleanValue bool `xml:"booleanValue,omitempty"`

	//
	// Double value of the operand if it is a double type.
	//
	DoubleValue float64 `xml:"doubleValue,omitempty"`

	//
	// String value of the operand if it is a string type.
	//
	StringValue string `xml:"stringValue,omitempty"`
}

type DatabaseError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 DatabaseError"`

	*ApiError

	//
	// The error reason represented by an enum.
	//
	Reason *DatabaseErrorReason `xml:"reason,omitempty"`
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

type FeedAttributeOperand struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 FeedAttributeOperand"`

	*FunctionArgumentOperand

	//
	// Id of associated feed.
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	//
	FeedId int64 `xml:"feedId,omitempty"`

	//
	// Id of the referenced feed attribute.
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	//
	FeedAttributeId int64 `xml:"feedAttributeId,omitempty"`
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

type Function struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 Function"`

	//
	// Operator for a function.
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	//
	Operator *FunctionOperator `xml:"operator,omitempty"`

	//
	// Operand on the LHS in the equation. This is also the operand to be used for
	// single operand expressions such as NOT.
	// <span class="constraint CollectionSize">The minimum size of this collection is 1.</span>
	//
	LhsOperand []*FunctionArgumentOperand `xml:"lhsOperand,omitempty"`

	//
	// Operand on the RHS of the equation.
	//
	RhsOperand []*FunctionArgumentOperand `xml:"rhsOperand,omitempty"`

	//
	// String representation of the {@code Function}.
	//
	// <p>For mutate actions, this field can be set instead of the {@code operator},
	// {@code lhsOperand}, and {@code rhsOperand} fields. This field will be parsed and used to
	// populate the other fields.
	//
	// <p>When {@code Function} objects are returned from get or mutate calls, this field contains the
	// string representation of the {@code Function}. Note that because multiple strings may map to
	// the same {@code Function} (whitespace and single versus double quotation marks, for example),
	// the value returned may not be identical to the string sent in the request.
	//
	FunctionString string `xml:"functionString,omitempty"`
}

type FunctionError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 FunctionError"`

	*ApiError

	//
	// The error reason represented by an enum
	//
	Reason *FunctionErrorReason `xml:"reason,omitempty"`
}

type FunctionOperand struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 FunctionOperand"`

	*FunctionArgumentOperand

	//
	// The function held in this operand.
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	//
	Value *Function `xml:"value,omitempty"`
}

type FunctionParsingError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 FunctionParsingError"`

	*ApiError

	Reason *FunctionParsingErrorReason `xml:"reason,omitempty"`

	OffendingText string `xml:"offendingText,omitempty"`

	OffendingTextIndex int32 `xml:"offendingTextIndex,omitempty"`
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

type NullStatsPage struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 NullStatsPage"`

	*Page
}

type FunctionArgumentOperand struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 FunctionArgumentOperand"`

	//
	// Indicates that this instance is a subtype of FunctionArgumentOperand.
	// Although this field is returned in the response, it is ignored on input
	// and cannot be selected. Specify xsi:type instead.
	//
	FunctionArgumentOperandType string `xml:"FunctionArgumentOperand.Type,omitempty"`
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

type RequestContextOperand struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 RequestContextOperand"`

	*FunctionArgumentOperand

	//
	// Type of value to be referred in the request context.
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	//
	ContextType *RequestContextOperandContextType `xml:"contextType,omitempty"`
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

type CampaignFeedServiceInterface struct {
	client *SOAPClient
}

func NewCampaignFeedServiceInterface(url string, tls bool, auth *BasicAuth) *CampaignFeedServiceInterface {
	if url == "" {
		url = ""
	}
	client := NewSOAPClient(url, tls, auth)

	return &CampaignFeedServiceInterface{
		client: client,
	}
}

func NewCampaignFeedServiceInterfaceWithTLSConfig(url string, tlsCfg *tls.Config, auth *BasicAuth) *CampaignFeedServiceInterface {
	if url == "" {
		url = ""
	}
	client := NewSOAPClientWithTLSConfig(url, tlsCfg, auth)

	return &CampaignFeedServiceInterface{
		client: client,
	}
}

func (service *CampaignFeedServiceInterface) AddHeader(header interface{}) {
	service.client.AddHeader(header)
}

// Backwards-compatible function: use AddHeader instead
func (service *CampaignFeedServiceInterface) SetHeader(header interface{}) {
	service.client.AddHeader(header)
}

// Error can be either of the following types:
//
//   - ApiException
/*
   Returns a list of CampaignFeeds that meet the selector criteria.

   @param selector Determines which CampaignFeeds to return. If empty all
   Campaign feeds are returned.
   @return The list of CampaignFeeds.
   @throws ApiException Indicates a problem with the request.
*/
func (service *CampaignFeedServiceInterface) Get(request *Get) (*GetResponse, error) {
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
   Adds, sets or removes CampaignFeeds.

   @param operations The operations to apply.
   @return The resulting Feeds.
   @throws ApiException Indicates a problem with the request.
*/
func (service *CampaignFeedServiceInterface) Mutate(request *Mutate) (*MutateResponse, error) {
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
   Returns a list of {@link CampaignFeed}s inside a {@link CampaignFeedPage} that matches
   the query.

   @param query The SQL-like AWQL query string.
   @throws ApiException when there are one or more errors with the request.
*/
func (service *CampaignFeedServiceInterface) Query(request *Query) (*QueryResponse, error) {
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
